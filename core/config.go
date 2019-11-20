package core

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

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

	subst := []byte(os.ExpandEnv(string(data)))

	err = yaml.Unmarshal(subst, config)
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

	if len(c.GitHubAccessTokens) < 1 || strings.TrimSpace(strings.Join(c.GitHubAccessTokens, "")) == "" {
		return errors.New("You need to provide at least one GitHub Access Token. See https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line")
	}

	return nil
}
