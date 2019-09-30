package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/google/go-github/github"
)

type GitHubClientWrapper struct {
	*github.Client
	Token            string
	RateLimitedUntil time.Time
}

const (
	baseUrl = "https://www.github.com"
	perPage = 300
	sleep   = 30 * time.Second
)

func GetRepositories(session *Session) {
	localCtx, cancel := context.WithCancel(session.Context)
	defer cancel()
	observedKeys := map[int64]bool{}
	for c := time.Tick(sleep); ; {
		opt := &github.ListOptions{PerPage: perPage}
		client := session.GetClient()

		for {
			events, resp, err := client.Activity.ListEvents(localCtx, opt)

			if err != nil {
				if _, ok := err.(*github.RateLimitError); ok {
					session.Log.Warn("Token %s[..] rate limited. Reset at %s", client.Token[:10], resp.Rate.Reset)
					client.RateLimitedUntil = resp.Rate.Reset.Time
					break
				}

				if _, ok := err.(*github.AbuseRateLimitError); ok {
					GetSession().Log.Fatal("GitHub API abused detected. Quitting...")
				}

				GetSession().Log.Important("Error getting GitHub events... trying again", err)
			}

			if opt.Page == 0 {
				session.Log.Warn("Token %s[..] has %d/%d calls remaining.", client.Token[:10], resp.Rate.Remaining, resp.Rate.Limit)
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
					session.Repositories <- e.GetRepo().GetID()
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

func GetGists(session *Session) {
	localCtx, cancel := context.WithCancel(session.Context)
	defer cancel()

	observedKeys := map[string]bool{}
	opt := &github.GistListOptions{}

	for c := time.Tick(sleep); ; {
		client := session.GetClient()
		gists, resp, err := client.Gists.ListAll(localCtx, opt)

		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				session.Log.Warn("Token %s[..] rate limited. Reset at %s", client.Token[:10], resp.Rate.Reset)
				client.RateLimitedUntil = resp.Rate.Reset.Time
				break
			}

			if _, ok := err.(*github.AbuseRateLimitError); ok {
				GetSession().Log.Fatal("GitHub API abused detected. Quitting...")
			}

			GetSession().Log.Important("Error getting GitHub Gists... trying again", err)
		}

		newGists := make([]*github.Gist, 0, len(gists))
		for _, e := range gists {
			if observedKeys[e.GetID()] {
				continue
			}

			newGists = append(newGists, e)
		}

		for _, e := range newGists {
			observedKeys[e.GetID()] = true
			session.Gists <- e.GetGitPullURL()
		}

		opt.Since = time.Now()

		select {
		case <-c:
			continue
		case <-localCtx.Done():
			cancel()
			return
		}
	}
}

func GetRepository(session *Session, id int64) (*github.Repository, error) {
	client := session.GetClient()
	repo, resp, err := client.Repositories.GetByID(session.Context, id)

	if err != nil {
		return nil, err
	}

	if resp.Rate.Remaining <= 1 {
		session.Log.Warn("Token %s[..] rate limited. Reset at %s", client.Token[:10], resp.Rate.Reset)
		client.RateLimitedUntil = resp.Rate.Reset.Time
	}

	return repo, nil
}

func GetUser(session *Session, id string) (*github.User, error) {
	client := session.GetClient()
	//println(id)
	user, resp, err := client.Users.Get(session.Context, id)

	if err != nil {
		return nil, err
	}

	if resp.Rate.Remaining <= 1 {
		session.Log.Warn("Token %s[..] rate limited. Reset at %s", client.Token[:10], resp.Rate.Reset)
		client.RateLimitedUntil = resp.Rate.Reset.Time
	}

	return user, nil
}
//The Google GitHub module didn't have support for this, so we have to make a web request and use some regex magic here.
func GetGistOwner(gistUrl string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	resp, err := client.Get(gistUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	r := regexp.MustCompile(`.com\/(.+)?\/.*.git`)
	res := r.FindStringSubmatch(bodyString)
	if len(res) != 0 {
		gistOwner := fmt.Sprint(res[1])
		return gistOwner, nil
	}
	return "err_getting_user", err
}
