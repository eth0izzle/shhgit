package settings

import (
	"flag"
	"os"
	"path/filepath"
)

type Options struct {
	Threads                *int
	Silent                 *bool
	Debug                  *bool
	MaximumRepositorySize  *int64
	MaximumFileSize        *int64
	CloneRepositoryTimeout *uint
	ScanFileTimeout        *uint
	EntropyThreshold       *float64
	MinimumStars           *uint
	PathChecks             *bool
	WorkingDirectory       *string
	LocalDirectory         *string
	Live                   *string
	ConfigPath             *string
	RulesPath              *string
}

func ParseOptions() (*Options, error) {
	options := &Options{
		Threads:                flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
		Silent:                 flag.Bool("silent", false, "Suppress all output except for errors"),
		Debug:                  flag.Bool("debug", false, "Print debugging information"),
		MaximumRepositorySize:  flag.Int64("maximum-repository-size", 5120, "Maximum repository size to process in KB"),
		MaximumFileSize:        flag.Int64("maximum-file-size", 256, "Maximum file size to process in KB"),
		CloneRepositoryTimeout: flag.Uint("clone-repository-timeout", 10, "Maximum time it should take to clone a repository in seconds. Increase this if you have a slower connection"),
		ScanFileTimeout:        flag.Uint("scan-file-timeout", 5, "Maximum seconds to spend scanning a single file for secrets. Increase this if you are targeting larger files"),
		EntropyThreshold:       flag.Float64("entropy-threshold", 5.0, "Set to 0 to disable entropy checks"),
		MinimumStars:           flag.Uint("minimum-stars", 0, "Only process repositories with this many stars. Default 0 will ignore star count"),
		PathChecks:             flag.Bool("path-checks", true, "Set to false to disable checking of filepaths, i.e. just match regex patterns of file contents"),
		WorkingDirectory:       flag.String("working-directory", filepath.Join(os.TempDir(), "shhgit"), "Directory to process and store repositories/matches"),
		LocalDirectory:         flag.String("local-dir", "", "Specify a local directory to recursively scan."),
		Live:                   flag.String("live", "", "Your shhgit live endpoint"),
		ConfigPath:             flag.String("config-path", "", "Searches for config.yaml from given directory. If not set, tries to find if from shhgit binary's and current directory"),
		RulesPath:              flag.String("rules-path", "rules/", "Directory to load and process yara rules."),
	}

	flag.Parse()

	return options, nil
}
