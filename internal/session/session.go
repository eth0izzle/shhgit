package session

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"regexp/syntax"
	"runtime"
	"sync"
	"time"

	"github.com/Velocidex/go-yara"
	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/outputs"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/types"
)

type Session struct {
	sync.Mutex

	Version      string
	Log          *outputs.Logger
	Options      *settings.Options
	Config       *settings.Config
	Signatures   []types.Signature
	Repositories chan types.RepositoryResource
	Rules        *yara.Rules
	Comments     chan string
	Context      context.Context
}

var (
	session     *Session
	sessionSync sync.Once
	err         error
)

func (s *Session) start() {
	rand.Seed(time.Now().Unix())

	s.initLogger()
	s.initThreads()
	s.initSignatures()
	s.initRules()
}

func (s *Session) initLogger() {
	s.Log = &outputs.Logger{}
	s.Log.SetDebug(*s.Options.Debug)
	s.Log.SetSilent(*s.Options.Silent)
}

func (s *Session) initSignatures() {
	var signatures []types.Signature

	for _, signature := range s.Config.Signatures {
		if signature.Match != "" {
			signatures = append(signatures, types.SimpleSignature{
				Name:    signature.Name,
				Part:    signature.Part,
				MatchOn: signature.Match,
			})
		} else {
			if _, err := syntax.Parse(signature.Match, syntax.FoldCase); err == nil {
				signatures = append(signatures, types.PatternSignature{
					Name:    signature.Name,
					Part:    signature.Part,
					MatchOn: regexp.MustCompile(signature.Regex),
				})
			}
		}
	}

	s.Signatures = signatures
}

func (s *Session) initRules() {
	c, err := yara.NewCompiler()
	if err != nil {
		s.Log.Fatal("Failed to initialize YARA compiler: %s", err)
	}

	_ = c.DefineVariable("filename", "filename.test")
	_ = c.DefineVariable("filepath", "test/path")
	_ = c.DefineVariable("extension", ".test")
	_ = c.DefineVariable("repository_name", "repository/name")
	_ = c.DefineVariable("repository_description", "Repository description")
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

	s.Rules = rules
}

func (s *Session) initThreads() {
	if *s.Options.Threads == 0 {
		numCPUs := runtime.NumCPU()
		s.Options.Threads = &numCPUs
	}

	runtime.GOMAXPROCS(*s.Options.Threads + 1)
}

func Start() *Session {
	sessionSync.Do(func() {
		session = &Session{
			Context:      context.Background(),
			Repositories: make(chan types.RepositoryResource, 10000),
			Comments:     make(chan string, 10000),
		}

		if session.Options, err = settings.ParseOptions(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if session.Config, err = settings.ParseConfig(session.Options); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		session.start()
	})

	return session
}
