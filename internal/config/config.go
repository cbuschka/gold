package config

type Config struct {
	DataDirPath             string
	CommandDomainSocketPath string
	GelfUdpListeners        []string
	GelfTcpListeners        []string
	GelfHttpListeners       []string
}

func GetConfig() *Config {
	return &Config{DataDirPath: "./tmp/db.leveldb", CommandDomainSocketPath: "./tmp/golfd.sock",
		GelfUdpListeners: []string{"127.0.0.1:12201"}, GelfTcpListeners: []string{"127.0.0.1:12201"},
		GelfHttpListeners: []string{"127.0.0.1:8080"}}
}
