package core

import (
	"context"
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

	Version          string
	Log              *Logger
	Options          *Options
	Config           *Config
	Signatures       []Signature
	Repositories     chan int64
	Gists            chan string
	Comments         chan string
	Context          context.Context
	Clients          chan *GitHubClientWrapper
	ExhaustedClients chan *GitHubClientWrapper
	Publishers       []Publisher
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
	s.InitPublishers()
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
	if len(*s.Options.Local) <= 0 {
		chanSize := *s.Options.Threads * (len(s.Config.GitHubAccessTokens) + 1)
		s.Clients = make(chan *GitHubClientWrapper, chanSize)
		s.ExhaustedClients = make(chan *GitHubClientWrapper, chanSize)
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

			for i := 0; i <= *s.Options.Threads; i++ {
				s.Clients <- &GitHubClientWrapper{client, token, time.Now().Add(-1 * time.Second)}
			}
		}

		if len(s.Clients) < 1 {
			s.Log.Fatal("No valid GitHub tokens provided. Quitting!")
		}
	}
}

func (s *Session) GetClient() *GitHubClientWrapper {
	for {
		select {

		case client := <-s.Clients:
			s.Log.Debug("Using client with token: %s", client.Token[:10])
			return client

		case client := <-s.ExhaustedClients:
			sleepTime := time.Until(client.RateLimitedUntil)
			s.Log.Warn("All GitHub tokens exhausted/rate limited. Sleeping for %s", sleepTime.String())
			time.Sleep(sleepTime)
			s.Log.Debug("Returning client %s to pool", client.Token[:10])
			s.FreeClient(client)

		default:
			s.Log.Debug("Available Clients: %d", len(s.Clients))
			s.Log.Debug("Exhausted Clients: %d", len(s.ExhaustedClients))
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

// FreeClient returns the GitHub Client to the pool of available,
// non-rate-limited channel of clients in the session
func (s *Session) FreeClient(client *GitHubClientWrapper) {
	if client.RateLimitedUntil.After(time.Now()) {
		s.ExhaustedClients <- client
	} else {
		s.Clients <- client
	}
}

func (s *Session) InitThreads() {
	if *s.Options.Threads == 0 {
		numCPUs := runtime.NumCPU()
		s.Options.Threads = &numCPUs
	}

	runtime.GOMAXPROCS(*s.Options.Threads + 1)
}

func (s *Session) InitPublishers() {
	header := []string{"Repository name", "Signature name", "Matching file", "Matches"}

	// Setup our delimited publisher
	if *s.Options.CsvPath != "" {
		// Convert delimiter to a rune slice. We only need the first rune
		delimiter := []rune(*s.Options.Delimiter)
		publisher, err := NewDelimitedSource(*s.Options.CsvPath, delimiter[0], header)
		if err != nil {
			s.Log.Error("Cannot create CSV publisher: %s", err)
		}
		if err == nil {
			s.Publishers = append(s.Publishers, &publisher)
		}
	}

	// Setup our Live publisher
	if *session.Options.Live != "" {
		publisher, err := NewWebSource(*session.Options.Live, "POST", "application/json")
		if err != nil {
			s.Log.Error("Cannot create Live publisher: %s", err)
		}
		if err == nil {
			s.Publishers = append(s.Publishers, &publisher)
		}
	}
}

func GetSession() *Session {
	sessionSync.Do(func() {
		session = &Session{
			Context:      context.Background(),
			Repositories: make(chan int64, 1000),
			Gists:        make(chan string, 100),
			Comments:     make(chan string, 1000),
		}

		if session.Options, err = ParseOptions(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if session.Config, err = ParseConfig(session.Options); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		session.Start()
	})

	return session
}
