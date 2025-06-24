package models

type Target struct {
	Account      string `mapstructure:"account"`
	Region       string `mapstructure:"region"`
	Environment  string `mapstructure:"environment"`
	Version      string `mapstructure:"version"`
	Organization string `mapstructure:"organization"`
	Name         string `mapstructure:"name"`
	Alias        string `mapstructure:"alias"`
	Domain       string `mapstructure:"domain"`
}

type DynamoDb struct {
	Table string `mapstructure:"table"`
}
