package core

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GitHubAccessTokens           []string          `yaml:"github_access_tokens"`
	SlackWebhook                 string            `yaml:"slack_webhook,omitempty"`
	BlacklistedExtensions        []string          `yaml:"blacklisted_extensions"`
	BlacklistedPaths             []string          `yaml:"blacklisted_paths"`
	BlacklistedEntropyExtensions []string          `yaml:"blacklisted_entropy_extensions"`
	Signatures                   []ConfigSignature `yaml:"signatures"`
}

type ConfigSignature struct {
	Name     string `yaml:"name"`
	Part     string `yaml:"part"`
	Match    string `yaml:"match,omitempty"`
	Regex    string `yaml:"regex,omitempty"`
	Verifier string `yaml:"verifier,omitempty"`
}

func ParseConfig() (*Config, error) {
	config := &Config{}

	dir, _ := os.Getwd()
	data, err := ioutil.ReadFile(path.Join(dir, "config.yaml"))
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config, err
	}

	for i := 0; i < len(config.GitHubAccessTokens); i++ {
		config.GitHubAccessTokens[i] = os.ExpandEnv(config.GitHubAccessTokens[i])
	}

	if len(config.SlackWebhook) > 0 {
		config.SlackWebhook = os.ExpandEnv(config.SlackWebhook)
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