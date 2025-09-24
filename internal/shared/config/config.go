package config


type Config struct {
	Server Server `json:"server"`
	Database Database `json:"database"`
	Logger Logger `json:"logger"`
}

type Server struct {
	Port string `json:"port"`
}

type Database struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Logger struct {
	Level string `json:"level"`
	Source bool `json:"source"`
}

