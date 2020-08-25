package session

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/outputs"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/sources/bitbucket"
	"github.com/eth0izzle/shhgit/internal/sources/github"
	"github.com/eth0izzle/shhgit/internal/sources/gitlab"
	"github.com/eth0izzle/shhgit/internal/types"

	"github.com/hillu/go-yara/v4"
)

type Session struct {
	sync.Mutex

	Version      string
	Log          *outputs.Logger
	Options      *settings.Options
	Config       *settings.Config
	ScannerRules *yara.Rules
	Outputs      []outputs.Publisher
	Context      context.Context
}

var (
	session     *Session
	sessionSync sync.Once
	err         error
)

func (s *Session) initLogger() {
	s.Log = &outputs.Logger{}
	s.Log.SetDebug(*s.Options.Debug)
	s.Log.SetSilent(*s.Options.Silent)
}

func (s *Session) initScanner() {
	c, err := yara.NewCompiler()
	if err != nil {
		s.Log.Fatal("Failed to initialize YARA compiler: %s", err)
	}

	c.DisableIncludes()

	_ = c.DefineVariable("filename", "")
	_ = c.DefineVariable("filepath", "")
	_ = c.DefineVariable("extension", "")
	_ = c.DefineVariable("repository_name", "")
	_ = c.DefineVariable("repository_description", "")
	_ = c.DefineVariable("repository_owner", "")
	_ = c.DefineVariable("repository_size", 0)
	_ = c.DefineVariable("repository_stars", 0)

	for _, rule := range helpers.GetFilesInPath(*s.Options.RulesPath, ".yara") {
		f, err := os.Open(rule)
		defer f.Close()

		if err != nil {
			s.Log.Fatal("Could not open rule file %s: %s", rule, err)
		}

		err = c.AddFile(f, "default")
		if err != nil {
			s.Log.Warn("Could not parse rule file %s: %s - ignorning", rule, err)
		}
	}

	rules, err := c.GetRules()
	if err != nil {
		s.Log.Fatal("Failed to compile rules: %s", err)
	}

	s.Log.Info("Loading %d rules from %s", len(rules.GetRules()), *s.Options.RulesPath)
	s.ScannerRules = rules
}

func (s *Session) initOutputs() {

}

// clean this up to make it more dynamic - perhaps use reflection?
func (s *Session) FetchFromSources(repositories chan<- types.RepositoryResource) {
	enabledSources := make([]string, 0)

	if s.Config.Sources.GitHub.Enabled {
		enabledSources = append(enabledSources, "GitHub")
		go github.FetchRepositories(s.Config.Sources.GitHub, repositories)
	}

	if s.Config.Sources.Gist.Enabled {
		enabledSources = append(enabledSources, "Gist")
		go github.FetchGists(s.Config.Sources.Gist, repositories)
	}

	if s.Config.Sources.GitLab.Enabled {
		enabledSources = append(enabledSources, "GitLab")
		go gitlab.Fetch(s.Config.Sources.GitLab, repositories)
	}

	if s.Config.Sources.BitBucket.Enabled {
		enabledSources = append(enabledSources, "BitBucket")
		go bitbucket.Fetch(s.Config.Sources.BitBucket, repositories)
	}

	s.Log.Info("[*] Fetching events from %d sources: %s", len(enabledSources), strings.Join(enabledSources, ", "))
}

func Fetch() *Session {
	sessionSync.Do(func() {
		session = &Session{Context: context.Background()}

		if session.Options, err = settings.ParseOptions(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if session.Config, err = settings.ParseConfig(session.Options); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		session.initLogger()
		session.initScanner()
		session.initOutputs()

		session.Log.Debug("Started new session at %s", time.Now())
	})

	return session
}
