package core

import (
	"context"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Session struct {
	sync.Mutex

	Version      string
	Log          *Logger
	Options      *Options
	Config       *Config
	Signatures   []Signature
	Repositories chan int64
	Gists        chan string
	Context      context.Context
	Clients      []*GitHubClientWrapper
	CsvWriter    *csv.Writer
}

var (
	session     *Session
	sessionSync sync.Once
	err         error
)

func (s *Session) Start() {
	rand.Seed(time.Now().Unix())

	s.InitLogger()
	s.InitThreads()
	s.InitSignatures()
	s.InitGitHubClients()
	s.InitCsvWriter()
}

func (s *Session) InitLogger() {
	s.Log = &Logger{}
	s.Log.SetDebug(*s.Options.Debug)
	s.Log.SetSilent(*s.Options.Silent)
}

func (s *Session) InitSignatures() {
	s.Signatures = GetSignatures(s)
}

func (s *Session) InitGitHubClients() {
	for _, token := range s.Config.GitHubAccessTokens {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(s.Context, ts)

		client := github.NewClient(tc)

		client.UserAgent = fmt.Sprintf("%s v%s", Name, Version)
		_, _, err := client.Users.Get(s.Context, "")

		if err != nil {
			if _, ok := err.(*github.ErrorResponse); ok {
				s.Log.Warn("Failed to validate token %s[..]: %s", token[:10], err)
				continue
			}
		}

		s.Clients = append(s.Clients, &GitHubClientWrapper{client, token, time.Now().Add(-1 * time.Second)})
	}

	if len(s.Clients) < 1 {
		s.Log.Fatal("No valid GitHub tokens provided. Quitting!")
	}
}

func (s *Session) GetClient() *GitHubClientWrapper {
	sleepTime := 0 * time.Second

	for _, client := range s.Clients {
		if client.RateLimitedUntil.After(time.Now()) {
			sleepTime = time.Until(client.RateLimitedUntil)
			continue
		}

		return client
	}

	s.Log.Warn("All GitHub tokens exchausted/rate limited. Sleeping for %s", sleepTime.String())
	time.Sleep(sleepTime)

	return s.GetClient()
}

func (s *Session) InitThreads() {
	if *s.Options.Threads == 0 {
		numCPUs := runtime.NumCPU()
		s.Options.Threads = &numCPUs
	}

	runtime.GOMAXPROCS(*s.Options.Threads + 1)
}

func (s *Session) InitCsvWriter() {
	if *s.Options.CsvPath == "" {
		return
	}

	writeHeader := false
	if !PathExists(*s.Options.CsvPath) {
		writeHeader = true
	}

	file, err := os.OpenFile(*s.Options.CsvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	LogIfError("Could not create/open CSV file", err)

	s.CsvWriter = csv.NewWriter(file)

	if writeHeader {
		s.WriteToCsv([]string{"Repository name", "Signature name", "Matching file", "Matches"})
	}
}

func (s *Session) WriteToCsv(line []string) {
	if *s.Options.CsvPath == "" {
		return
	}

	s.CsvWriter.Write(line)
	s.CsvWriter.Flush()
}

func GetSession() *Session {
	sessionSync.Do(func() {
		session = &Session{
			Context:      context.Background(),
			Repositories: make(chan int64, 1000),
			Gists:        make(chan string, 100),
		}

		if session.Options, err = ParseOptions(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if session.Config, err = ParseConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		session.Version = Version
		session.Start()
	})

	return session
}
