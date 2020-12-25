package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/eth0izzle/shhgit/core"
	"github.com/fatih/color"
)

type MatchEvent struct {
	Url       string
	Matches   []string
	Signature string
	File      string
	Stars     int
	Source    core.GitResourceType
}

var session = core.GetSession()

func ProcessRepositories() {
	threadNum := *session.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			for {
				timeout := time.Duration(*session.Options.CloneRepositoryTimeout) * time.Second
				_, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				repository := <-session.Repositories

				repo, err := core.GetRepository(session, repository.Id)

				if err != nil {
					session.Log.Warn("Failed to retrieve repository %d: %s", repository.Id, err)
					continue
				}

				if repo.GetPermissions()["pull"] &&
					uint(repo.GetStargazersCount()) >= *session.Options.MinimumStars &&
					uint(repo.GetSize()) < *session.Options.MaximumRepositorySize {

					processRepositoryOrGist(repo.GetCloneURL(), repository.Ref, repo.GetStargazersCount(), core.GITHUB_SOURCE)
				}
			}
		}(i)
	}
}

func ProcessGists() {
	threadNum := *session.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			for {
				gistUrl := <-session.Gists
				processRepositoryOrGist(gistUrl, "", -1, core.GIST_SOURCE)
			}
		}(i)
	}
}

func ProcessComments() {
	threadNum := *session.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			for {
				commentBody := <-session.Comments
				dir := core.GetTempDir(core.GetHash(commentBody))
				ioutil.WriteFile(filepath.Join(dir, "comment.ignore"), []byte(commentBody), 0644)

				if !checkSignatures(dir, "ISSUE", 0, core.GITHUB_COMMENT) {
					os.RemoveAll(dir)
				}
			}
		}(i)
	}
}

func processRepositoryOrGist(url string, ref string, stars int, source core.GitResourceType) {
	var (
		matchedAny bool = false
	)

	dir := core.GetTempDir(core.GetHash(url))
	_, err := core.CloneRepository(session, url, ref, dir)

	if err != nil {
		session.Log.Debug("[%s] Cloning failed: %s", url, err.Error())
		os.RemoveAll(dir)
		return
	}

	session.Log.Debug("[%s] Cloning %s in to %s", url, ref, strings.Replace(dir, *session.Options.TempDirectory, "", -1))
	matchedAny = checkSignatures(dir, url, stars, source)
	if !matchedAny {
		os.RemoveAll(dir)
	}
}

func checkSignatures(dir string, url string, stars int, source core.GitResourceType) (matchedAny bool) {
	for _, file := range core.GetMatchingFiles(dir) {
		var (
			matches          []string
			relativeFileName string
		)
		if strings.Contains(dir, *session.Options.TempDirectory) {
			relativeFileName = strings.Replace(file.Path, *session.Options.TempDirectory, "", -1)
		} else {
			relativeFileName = strings.Replace(file.Path, dir, "", -1)
		}

		if *session.Options.SearchQuery != "" {
			queryRegex := regexp.MustCompile(*session.Options.SearchQuery)
			for _, match := range queryRegex.FindAllSubmatch(file.Contents, -1) {
				matches = append(matches, string(match[0]))
			}

			if matches != nil {
				count := len(matches)
				m := strings.Join(matches, ", ")
				session.Log.Important("[%s] %d %s for %s in file %s: %s", url, count, core.Pluralize(count, "match", "matches"), color.GreenString("Search Query"), relativeFileName, color.YellowString(m))
				session.WriteToCsv([]string{url, "Search Query", relativeFileName, m})
			}
		} else {
			for _, signature := range session.Signatures {
				if matched, part := signature.Match(file); matched {
					matchedAny = true

					if part == core.PartContents {
						if matches = signature.GetContentsMatches(file.Contents); matches != nil {
							count := len(matches)
							m := strings.Join(matches, ", ")
							publish(&MatchEvent{Source: source, Url: url, Matches: matches, Signature: signature.Name(), File: relativeFileName, Stars: stars})
							session.Log.Important("[%s] %d %s for %s in file %s: %s", url, count, core.Pluralize(count, "match", "matches"), color.GreenString(signature.Name()), relativeFileName, color.YellowString(m))
							session.WriteToCsv([]string{url, signature.Name(), relativeFileName, m})
						}
					} else {
						if *session.Options.PathChecks {
							publish(&MatchEvent{Source: source, Url: url, Matches: matches, Signature: signature.Name(), File: relativeFileName, Stars: stars})
							session.Log.Important("[%s] Matching file %s for %s", url, color.YellowString(relativeFileName), color.GreenString(signature.Name()))
							session.WriteToCsv([]string{url, signature.Name(), relativeFileName, ""})
						}

						if *session.Options.EntropyThreshold > 0 && file.CanCheckEntropy() {
							scanner := bufio.NewScanner(bytes.NewReader(file.Contents))

							for scanner.Scan() {
								line := scanner.Text()

								if len(line) > 6 && len(line) < 100 {
									entropy := core.GetEntropy(line)

									if entropy >= *session.Options.EntropyThreshold {
										blacklistedMatch := false

										for _, blacklistedString := range session.Config.BlacklistedStrings {
											if strings.Contains(strings.ToLower(line), strings.ToLower(blacklistedString)) {
												blacklistedMatch = true
											}
										}

										if !blacklistedMatch {
											publish(&MatchEvent{Source: source, Url: url, Matches: []string{line}, Signature: "High entropy string", File: relativeFileName, Stars: stars})
											session.Log.Important("[%s] Potential secret in %s = %s", url, color.YellowString(relativeFileName), color.GreenString(line))
											session.WriteToCsv([]string{url, "High entropy string", relativeFileName, line})
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if !matchedAny && len(*session.Options.Local) <= 0 {
			os.Remove(file.Path)
		}
	}
	return
}

func publish(event *MatchEvent) {
	// todo: implement a modular plugin system to handle the various outputs (console, live, csv, webhooks, etc)
	if len(*session.Options.Live) > 0 {
		data, _ := json.Marshal(event)
		http.Post(*session.Options.Live, "application/json", bytes.NewBuffer(data))
	}
}

func main() {
	session.Log.Info(color.HiBlueString(core.Banner))
	session.Log.Info("\t%s\n", color.HiCyanString(core.Author))
	session.Log.Info("[*] Loaded %s signatures. Using %s worker threads. Temp work dir: %s\n", color.BlueString("%d", len(session.Signatures)), color.BlueString("%d", *session.Options.Threads), color.BlueString(*session.Options.TempDirectory))

	if len(*session.Options.Local) > 0 {
		session.Log.Info("[*] Scanning local directory: %s - skipping public repository checks...", color.BlueString(*session.Options.Local))
		rc := 0
		if checkSignatures(*session.Options.Local, *session.Options.Local, -1, core.LOCAL_SOURCE) {
			rc = 1
		} else {
			session.Log.Info("[*] No matching secrets found in %s!", color.BlueString(*session.Options.Local))
		}
		os.Exit(rc)
	} else {
		if *session.Options.SearchQuery != "" {
			session.Log.Important("Search Query '%s' given. Only returning matching results.", *session.Options.SearchQuery)
		}

		go core.GetRepositories(session)
		go ProcessRepositories()
		go ProcessComments()

		if *session.Options.ProcessGists {
			go core.GetGists(session)
			go ProcessGists()
		}

		spinny := core.ShowSpinner()
		select {}
		spinny()
	}
}
