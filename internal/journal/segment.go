package journal

import (
	"bytes"
	"fmt"
	"github.com/kataras/golog"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type Segment struct {
	pathname string
	file     *os.File
	index    uint64
	size     uint64
}

func getPathname(basedir string, index uint64) string {
	filename := fmt.Sprintf("%016x.jsonl", index)
	pathname := filepath.Join(basedir, filename)
	return pathname
}

func (segment *Segment) ensureOpen() error {
	if segment.file != nil {
		return nil
	}

	file, err := os.OpenFile(segment.pathname, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		return fmt.Errorf("Opening segment file index=%x, pathname=%s failed: %v", segment.index, segment.pathname, err.Error())
	}
	segment.file = file

	return nil
}

func (segment *Segment) mustSwitch() bool {
	return segment.size > 1024
}

func (segment *Segment) Append(json []byte) error {

	err := segment.ensureOpen()
	if err != nil {
		return err
	}

	lineEnd := []byte("\n")
	buf := bytes.Join([][]byte{json, lineEnd}, nil)

	if _, err := (*segment.file).Write(buf); err != nil {
		return err
	}

	segment.size = segment.size + uint64(len(buf))

	return nil
}

func (segment *Segment) Close() error {

	var err error
	if segment.file != nil {
		err = segment.file.Close()
	}
	segment.file = nil

	return err
}

type SegmentManager struct {
	basedir  string
	segments []*Segment
}

func (manager *SegmentManager) collectFiles() error {

	baseDirStats, err := os.Stat(manager.basedir)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(manager.basedir, 0750)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err

	} else {
		if !baseDirStats.IsDir() {
			return fmt.Errorf("Cannot create %s. File is in the way.", manager.basedir)
		}
	}

	matcher, err := regexp.Compile("^([0-9a-f]{16})\\.jsonl$")
	if err != nil {
		return err
	}

	segments := make([]*Segment, 0)
	if err := filepath.Walk(manager.basedir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		matchResult := matcher.FindStringSubmatch(info.Name())
		if len(matchResult) > 1 {
			index, err := strconv.ParseUint(matchResult[1], 16, 64)
			if err != nil {
				return err
			}
			segments = append(segments, &Segment{index: index, size: uint64(info.Size()), pathname: path})
		}

		sort.Sort(SegmentCollection(segments))

		manager.segments = segments

		return nil
	}); err != nil {
		return err
	}

	golog.Debugf("Found %d segment file(s).", len(manager.segments))

	return nil
}

func (manager *SegmentManager) getLatestSegment() (*Segment, error) {
	segmentCount := len(manager.segments)
	if segmentCount == 0 {
		return manager.switchSegment()
	}

	lastSegment := manager.segments[segmentCount-1]
	if lastSegment.mustSwitch() {
		return manager.switchSegment()
	}

	return manager.segments[segmentCount-1], nil
}

func (manager *SegmentManager) switchSegment() (*Segment, error) {
	segmentCount := len(manager.segments)
	var newSegment Segment
	if segmentCount == 0 {
		newPathname := getPathname(manager.basedir, 0)
		newSegment = Segment{size: 0, index: 0, pathname: newPathname}
		manager.segments = []*Segment{&newSegment}
	} else {
		newIndex := manager.segments[segmentCount-1].index + 1
		newPathname := getPathname(manager.basedir, newIndex)
		newSegment = Segment{size: 0, index: newIndex, pathname: newPathname}
		_ = manager.segments[segmentCount-1].Close()
		manager.segments = append(manager.segments, &newSegment)
	}
	golog.Infof("Switched to segment #%d, pathname=%s.", newSegment.index, newSegment.pathname)
	return &newSegment, nil
}

func (manager *SegmentManager) Close() error {
	var err error
	if manager.segments != nil {
		for _, segment := range manager.segments {
			err = segment.Close()
		}
	}
	manager.segments = nil
	return err
}

type SegmentCollection []*Segment

func (s SegmentCollection) Len() int {
	return len(s)
}

func (s SegmentCollection) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SegmentCollection) Less(i, j int) bool {
	return s[i].index < s[j].index
}
