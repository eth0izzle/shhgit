package core

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GitHubAccessTokens           []string          `yaml:"github_access_tokens"`
	Webhook                      string            `yaml:"webhook,omitempty"`
	WebhookPayload               string            `yaml:"webhook_payload,omitempty"`
	BlacklistedStrings           []string          `yaml:"blacklisted_strings"`
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
		// Trying to first find the configuration next to executable
		// Helps e.g. with Drone where workdir is different than shhgit dir
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

	if len(*options.Local) <= 0 && (len(config.GitHubAccessTokens) < 1 || strings.TrimSpace(strings.Join(config.GitHubAccessTokens, "")) == "") {
		return config, errors.New("You need to provide at least one GitHub Access Token. See https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line")
	}

	for i := 0; i < len(config.GitHubAccessTokens); i++ {
		config.GitHubAccessTokens[i] = os.ExpandEnv(config.GitHubAccessTokens[i])
	}

	if len(config.Webhook) > 0 {
		config.Webhook = os.ExpandEnv(config.Webhook)
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
