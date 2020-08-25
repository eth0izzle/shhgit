package settings

import (
	"flag"
	"os"
	"path/filepath"
)

type Options struct {
	Workers                *int
	Silent                 *bool
	Debug                  *bool
	MaximumRepositorySize  *int64
	MaximumFileSize        *int64
	CloneRepositoryTimeout *int
	ScanFileTimeout        *int
	EntropyThreshold       *float64
	MinimumStars           *int
	WorkingDirectory       *string
	LocalDirectory         *string
	ConfigPath             *string
	RulesPath              *string
}

func ParseOptions() (*Options, error) {
	options := &Options{
		Workers:                flag.Int("workers", 10, "Number of workers to process repositories. More workers = faster processing but higher memory useage."),
		Silent:                 flag.Bool("silent", false, "Suppress all output except for errors"),
		Debug:                  flag.Bool("debug", false, "Print debugging information"),
		MaximumRepositorySize:  flag.Int64("maximum-repository-size", 5120, "Maximum repository size to process in KB"),
		MaximumFileSize:        flag.Int64("maximum-file-size", 256, "Maximum file size to process in KB"),
		CloneRepositoryTimeout: flag.Int("clone-repository-timeout", 10, "Maximum time it should take to clone a repository in seconds. Increase this if you have a slower connection"),
		ScanFileTimeout:        flag.Int("scan-file-timeout", 5, "Maximum seconds to spend scanning a single file for secrets. Increase this if you are targeting larger files"),
		EntropyThreshold:       flag.Float64("entropy-threshold", 5.0, "Set to 0 to disable entropy checks"),
		MinimumStars:           flag.Int("minimum-stars", 0, "Only process repositories with this many stars. Default 0 will ignore star count"),
		WorkingDirectory:       flag.String("working-directory", filepath.Join(os.TempDir(), "shhgit"), "Directory to process and store repositories/matches"),
		LocalDirectory:         flag.String("local-dir", "", "Specify a local directory to recursively scan."),
		ConfigPath:             flag.String("config-path", "", "Searches for config.yaml from given directory. If not set, tries to find if from shhgit binary's and current directory"),
		RulesPath:              flag.String("rules-path", "rules/", "Directory to load and process yara rules."),
	}

	flag.Parse()

	return options, nil
}
