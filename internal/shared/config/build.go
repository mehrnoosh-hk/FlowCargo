package config

type Build struct {
	Version string `json:"version" mapstructure:"version" validate:"required"`
	Commit  string `json:"commit" mapstructure:"commit" validate:"required"`
	Date    string `json:"date" mapstructure:"date" validate:"required"`
}
