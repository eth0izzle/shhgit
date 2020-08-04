package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"

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
				repositoryId := <-session.Repositories
				repo, err := core.GetRepository(session, repositoryId)

				if err != nil {
					session.Log.Warn("Failed to retrieve repository %d: %s", repositoryId, err)
					continue
				}

				if repo.GetPermissions()["pull"] &&
					uint(repo.GetStargazersCount()) >= *session.Options.MinimumStars &&
					uint(repo.GetSize()) < *session.Options.MaximumRepositorySize {

					processRepositoryOrGist(repo.GetCloneURL(), repo.GetStargazersCount(), core.GITHUB_SOURCE)
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
				processRepositoryOrGist(gistUrl, -1, core.GIST_SOURCE)
			}
		}(i)
	}
}

func processRepositoryOrGist(url string, stars int, source core.GitResourceType) {
	var (
		matchedAny bool = false
	)

	dir := core.GetTempDir(core.GetHash(url))
	_, err := core.CloneRepository(session, url, dir)

	if err != nil {
		session.Log.Debug("[%s] Cloning failed: %s", url, err.Error())
		os.RemoveAll(dir)
		return
	}

	session.Log.Debug("[%s] Cloning in to %s", url, strings.Replace(dir, *session.Options.TempDirectory, "", -1))
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
						if matches = signature.GetContentsMatches(file); matches != nil {
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
											session.WriteToCsv([]string{url, signature.Name(), relativeFileName, line})
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
	if len(*session.Options.Local) > 0 {
		session.Log.Info("Scanning local dir %s with %s v%s. Loaded %d signatures.", session.Options.Local, core.Name, core.Version, len(session.Signatures))
		rc := 0
		if checkSignatures(*session.Options.Local, *session.Options.Local, -1, core.LOCAL_SOURCE) {
			rc = 1
		}
		os.Exit(rc)
	} else {
		session.Log.Info("%s v%s started. Loaded %d signatures. Using %d GitHub tokens and %d threads. Work dir: %s", core.Name, core.Version, len(session.Signatures), len(session.Clients), *session.Options.Threads, *session.Options.TempDirectory)

		if *session.Options.SearchQuery != "" {
			session.Log.Important("Search Query '%s' given. Only returning matching results.", *session.Options.SearchQuery)
		}

		go core.GetRepositories(session)
		go ProcessRepositories()

		if *session.Options.ProcessGists {
			go core.GetGists(session)
			go ProcessGists()
		}

		session.Log.Info("Press Ctrl+C to stop and exit.\n")
		select {}
	}
}
