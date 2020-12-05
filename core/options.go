package core

import (
	"flag"
	"os"
	"path/filepath"
)

type Options struct {
	Threads                *int
	Silent                 *bool
	Debug                  *bool
	MaximumRepositorySize  *uint
	MaximumFileSize        *uint
	CloneRepositoryTimeout *uint
	EntropyThreshold       *float64
	MinimumStars           *uint
	PathChecks             *bool
	ProcessGists           *bool
	TempDirectory          *string
	CsvPath                *string
	Delimiter              *string
	SearchQuery            *string
	Local                  *string
	Live                   *string
	ConfigPath             *string
}

func ParseOptions() (*Options, error) {
	options := &Options{
		Threads:                flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
		Silent:                 flag.Bool("silent", false, "Suppress all output except for errors"),
		Debug:                  flag.Bool("debug", false, "Print debugging information"),
		MaximumRepositorySize:  flag.Uint("maximum-repository-size", 5120, "Maximum repository size to process in KB"),
		MaximumFileSize:        flag.Uint("maximum-file-size", 256, "Maximum file size to process in KB"),
		CloneRepositoryTimeout: flag.Uint("clone-repository-timeout", 10, "Maximum time it should take to clone a repository in seconds. Increase this if you have a slower connection"),
		EntropyThreshold:       flag.Float64("entropy-threshold", 5.0, "Set to 0 to disable entropy checks"),
		MinimumStars:           flag.Uint("minimum-stars", 0, "Only process repositories with this many stars. Default 0 will ignore star count"),
		PathChecks:             flag.Bool("path-checks", true, "Set to false to disable checking of filepaths, i.e. just match regex patterns of file contents"),
		ProcessGists:           flag.Bool("process-gists", true, "Will watch and process Gists. Set to false to disable."),
		TempDirectory:          flag.String("temp-directory", filepath.Join(os.TempDir(), Name), "Directory to process and store repositories/matches"),
		CsvPath:                flag.String("csv-path", "", "CSV file path to log found secrets to. Leave blank to disable"),
		Delimiter:              flag.String("delimiter", ",", "Delimiter for CSV file."),
		SearchQuery:            flag.String("search-query", "", "Specify a search string to ignore signatures and filter on files containing this string (regex compatible)"),
		Local:                  flag.String("local", "", "Specify local directory (absolute path) which to scan. Scans only given directory recursively. No need to have Githib tokens with local run."),
		Live:                   flag.String("live", "", "Your shhgit live endpoint"),
		ConfigPath:             flag.String("config-path", "", "Searches for config.yaml from given directory. If not set, tries to find if from shhgit binary's and current directory"),
	}

	flag.Parse()

	return options, nil
}
