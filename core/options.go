package core

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
	SearchQuery            *string
	Local                  *string
	Live                   *string
	ConfigPath             *string
	ConfigName             *string
}

func (o *Options) Merge(new *Options) {
	if new.Threads != nil {
		o.Threads = new.Threads
	}
	if new.Silent != nil {
		o.Silent = new.Silent
	}
	if new.Debug != nil {
		o.Debug = new.Debug
	}
	if new.MaximumRepositorySize != nil {
		o.MaximumRepositorySize = new.MaximumRepositorySize
	}
	if new.MaximumFileSize != nil {
		o.MaximumFileSize = new.MaximumFileSize
	}
	if new.CloneRepositoryTimeout != nil {
		o.CloneRepositoryTimeout = new.CloneRepositoryTimeout
	}
	if new.EntropyThreshold != nil {
		o.EntropyThreshold = new.EntropyThreshold
	}
	if new.MinimumStars != nil {
		o.MinimumStars = new.MinimumStars
	}
	if new.PathChecks != nil {
		o.PathChecks = new.PathChecks
	}
	if new.ProcessGists != nil {
		o.ProcessGists = new.ProcessGists
	}
	if new.TempDirectory != nil {
		o.TempDirectory = new.TempDirectory
	}
	if new.CsvPath != nil {
		o.CsvPath = new.CsvPath
	}
	if new.SearchQuery != nil {
		o.SearchQuery = new.SearchQuery
	}
	if new.Local != nil {
		o.Local = new.Local
	}
	if new.Live != nil {
		o.Live = new.Live
	}
	if new.ConfigPath != nil {
		o.ConfigPath = new.ConfigPath
	}
	if new.ConfigName != nil {
		o.ConfigName = new.ConfigName
	}
}

var (
	// Defaults that don't represent Go struct defaults
	DefaultMaximumRepositorySize  = uint(5120)
	DefaultMaximumFileSize        = uint(256)
	DefaultCloneRepositoryTimeout = uint(10)
	DefaultEntropy                = float64(5.0)
	DefaultPathChecks             = true
	DefaultProcessGists           = true

	DefaultOptions = Options{
		MaximumRepositorySize:  &DefaultMaximumRepositorySize,
		MaximumFileSize:        &DefaultMaximumFileSize,
		CloneRepositoryTimeout: &DefaultCloneRepositoryTimeout,
		EntropyThreshold:       &DefaultEntropy,
		PathChecks:             &DefaultPathChecks,
		ProcessGists:           &DefaultProcessGists,
	}
)
