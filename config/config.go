package config

type Config struct {
	Port string
	Path string
	Id   string
	Join string
	Addr string
}

func NewConfig(port, path, id, join, addr string) *Config {
	return &Config{
		Port: port,
		Path: path,
		Id:   id,
		Join: join,
		Addr: addr,
	}
}
