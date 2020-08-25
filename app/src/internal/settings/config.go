package settings

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sources    ConfigSources    `yaml:"sources"`
	BlackLists ConfigBlacklists `yaml:"blacklists"`
}

type ConfigSources struct {
	GitHub    ConfigSource `yaml:"github"`
	Gist      ConfigSource `yaml:"gist"`
	GitLab    ConfigSource `yaml:"gitlab"`
	BitBucket ConfigSource `yaml:"bitbucket"`
}

type ConfigSource struct {
	Enabled       bool     `yaml:"enabled"`
	Endpoint      string   `yaml:"endpoint"`
	PerPage       uint     `yaml:"per_page"`
	CheckInterval uint     `yaml:"check_interval"`
	Tokens        []string `yaml:"tokens"`
}

type ConfigBlacklists struct {
	Strings           []string `yaml:"strings"`
	Extensions        []string `yaml:"extensions"`
	EntropyExtensions []string `yaml:"entropy_extensions"`
	Paths             []string `yaml:"paths"`
}

func ParseConfig(options *Options) (*Config, error) {
	config := &Config{}
	var (
		data []byte
		err  error
	)

	if len(*options.ConfigPath) > 0 {
		data, err = ioutil.ReadFile(path.Join(*options.ConfigPath, "config.yaml"))
		if err != nil {
			return config, err
		}
	} else {
		ex, err := os.Executable()
		dir := filepath.Dir(ex)
		data, err = ioutil.ReadFile(path.Join(dir, "config.yaml"))

		if err != nil {
			dir, _ = os.Getwd()
			data, err = ioutil.ReadFile(path.Join(dir, "config.yaml"))

			if err != nil {
				return config, err
			}
		}
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = Config{}
	type plain Config

	err := unmarshal((*plain)(c))

	if err != nil {
		return err
	}

	return nil
}
