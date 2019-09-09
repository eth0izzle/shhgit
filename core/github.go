package core

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

type GitHubClientWrapper struct {
	*github.Client
	Token            string
	RateLimitedUntil time.Duration
}

const (
	baseUrl = "https://www.github.com"
	perPage = 300
	sleep   = 30 * time.Second
)

func ReadEvents(session *Session) {
	localCtx, cancel := context.WithCancel(session.Context)
	observedKeys := map[int64]bool{}

	for c := time.Tick(sleep); ; {
		opt := &github.ListOptions{PerPage: perPage}
		client := session.GetClient()

		for {
			events, resp, err := client.Activity.ListEvents(localCtx, opt)

			if err != nil {
				if _, ok := err.(*github.RateLimitError); ok {
					session.Log.Warn("Token %s rate limited. Reset at %s", client.Token, resp.Rate.Reset)
					client.RateLimitedUntil = time.Until(resp.Rate.Reset.Time)
					break
				}

				if _, ok := err.(*github.AbuseRateLimitError); ok {
					GetSession().Log.Fatal("GitHub API abused detected. Quitting...")
				}

				GetSession().Log.Important("Error getting GitHub events... trying again", err)
			}

			newEvents := make([]*github.Event, 0, len(events))
			for _, e := range events {
				if observedKeys[e.GetRepo().GetID()] {
					continue
				}

				newEvents = append(newEvents, e)
			}

			for _, e := range newEvents {
				if *e.Type == "PushEvent" {
					observedKeys[e.GetRepo().GetID()] = true
					session.Repositories <- e.GetRepo().GetName()
				}
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
			time.Sleep(5 * time.Second)
		}

		select {
		case <-c:
			continue
		case <-localCtx.Done():
			cancel()
			return
		}
	}
}

func GetRepository(session *Session, name string) *github.Repository {
	ownerRepo := strings.Split(name, "/")
	client := session.GetClient()
	repo, resp, _ := client.Repositories.Get(session.Context, ownerRepo[0], ownerRepo[1])

	if resp.Rate.Remaining <= 1 {
		session.Log.Warn("Token %s rate limited. Reset at %s", client.Token, resp.Rate.Reset)
		client.RateLimitedUntil = time.Until(resp.Rate.Reset.Time)
	}

	return repo
}

func GetRepositoryUrl(name string) string {
	return baseUrl + "/" + name
}
