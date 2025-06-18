package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MongoURI       string   `yaml:"mongo_uri"`
	DatabaseName   string   `yaml:"database_name"`
	CollectionName string   `yaml:"collection_name"`
	BadWordsFile   string   `yaml:"badwords_file"`
	AllowedPlayers []string `yaml:"allowed_players"`
}


func LoadConfig(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
