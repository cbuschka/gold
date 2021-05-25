package journal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/karlseguin/jsonwriter"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"math"
	"strconv"
	"time"
)

type Message struct {
	Id               uint64
	Version          string
	Host             string
	SenderHost       string
	Short            string
	Full             string
	TimeUnix         time.Time
	ReceivedTimeUnix time.Time
	Level            int32
	Facility         string
	Extra            map[string]interface{}
}

func FromGelfMessage(gelfMessage *gelf.Message, senderHost string) *Message {
	var message Message
	message.Id = 0
	message.Version = gelfMessage.Version
	message.Host = gelfMessage.Host
	message.Short = gelfMessage.Short
	message.Full = gelfMessage.Full
	message.TimeUnix = float64ToTime(gelfMessage.TimeUnix)
	message.Level = gelfMessage.Level
	message.Facility = gelfMessage.Facility
	message.Extra = gelfMessage.Extra
	message.SenderHost = senderHost
	message.ReceivedTimeUnix = time.Now()
	message.Extra = make(map[string]interface{}, 0)
	for k, v := range gelfMessage.Extra {
		message.Extra[k] = v
	}
	return &message
}

func (m *Message) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))

	writer := jsonwriter.New(buf)
	writer.RootObject(func() {
		writer.KeyValue("_id", uint64ToString(m.Id))
		writer.KeyValue("version", m.Version)
		writer.KeyValue("host", m.Host)
		writer.KeyValue("short_message", m.Short)
		writer.KeyValue("full_message", m.Full)
		writer.KeyValue("timestamp", timeToString(m.TimeUnix))
		writer.KeyValue("level", float64(m.Level))
		writer.KeyValue("facility", m.Facility)
		writer.KeyValue("_sender_host", m.SenderHost)
		writer.KeyValue("_received_timestamp", timeToString(m.ReceivedTimeUnix))
		for k, v := range m.Extra {
			writer.KeyValue(k, v)
		}
	})
	return buf.Bytes(), nil
}

func (m *Message) UnmarshalJSON(data []byte) error {
	i := make(map[string]interface{}, 16)
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	for k, v := range i {
		ok := true
		switch k {
		case "_id":
			m.Id, ok = stringToUint64(v.(string))
		case "version":
			m.Version, ok = v.(string)
		case "host":
			m.Host, ok = v.(string)
		case "short_message":
			m.Short, ok = v.(string)
		case "full_message":
			m.Full, ok = v.(string)
		case "timestamp":
			m.TimeUnix, ok = stringToTime(v.(string))
		case "level":
			var level float64
			level, ok = v.(float64)
			m.Level = int32(level)
		case "facility":
			m.Facility, ok = v.(string)
		case "_sender_host":
			m.SenderHost, ok = v.(string)
		case "_received_timestamp":
			m.ReceivedTimeUnix, ok = stringToTime(v.(string))
		default:
			if k[0] == '_' {
				if m.Extra == nil {
					m.Extra = make(map[string]interface{}, 1)
				}
				m.Extra[k] = v
				ok = true
			}
		}

		if !ok {
			return fmt.Errorf("invalid type for field %s", k)
		}
	}
	return nil
}

func uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func stringToUint64(v string) (uint64, bool) {
	x, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, false
	}
	return x, true
}

func stringToTime(s string) (time.Time, bool) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Unix(0, 0), false
	}

	return t, true
}

func timeToString(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func float64ToTime(timeFloat float64) time.Time {
	sec, dec := math.Modf(timeFloat)
	return time.Unix(int64(sec), int64(dec*(1e9))).UTC()
}
