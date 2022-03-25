package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
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

var (
	threadsFlag      = flag.Int("threads", core.DefaultThreads, "Number of concurrent threads (default number of logical CPUs)")
	silentFlag       = flag.Bool("silent", core.DefaultSilent, "Suppress all output except for errors")
	debugFlag        = flag.Bool("debug", core.DefaultDebug, "Print debugging information")
	maxRepoFlag      = flag.Uint("maximum-repository-size", core.DefaultMaximumRepositorySize, "Maximum repository size to process in KB")
	maxFileSizeFlag  = flag.Uint("maximum-file-size", core.DefaultMaximumFileSize, "Maximum file size to process in KB")
	cloneTimeoutFlag = flag.Uint("clone-repository-timeout", core.DefaultCloneRepositoryTimeout, "Maximum time it should take to clone a repository in seconds. Increase this if you have a slower connection")
	entropyFlag      = flag.Float64("entropy-threshold", core.DefaultEntropy, "Set to 0 to disable entropy checks")
	minStarsFlag     = flag.Uint("minimum-stars", core.DefaultMinimumStars, "Only process repositories with this many stars. Default 0 will ignore star count")
	pathChecksFlag   = flag.Bool("path-checks", core.DefaultPathChecks, "Set to false to disable checking of filepaths, i.e. just match regex patterns of file contents")
	gistsFlag        = flag.Bool("process-gists", core.DefaultProcessGists, "Will watch and process Gists. Set to false to disable.")
	tempDirFlag      = flag.String("temp-directory", core.DefaultTempDirectory, "Directory to process and store repositories/matches")
	csvFlag          = flag.String("csv-path", "", "CSV file path to log found secrets to. Leave blank to disable")
	searchQueryFlag  = flag.String("search-query", "", "Specify a search string to ignore signatures and filter on files containing this string (regex compatible)")
	localFlag        = flag.String("local", "", "Specify local directory (absolute path) which to scan. Scans only given directory recursively. No need to have GitHub tokens with local run.")
	liveFlag         = flag.String("live", "", "Your shhgit live endpoint")
	configPathFlag   = flag.String("config-path", "", "Searches for config.yaml from given directory. If not set, tries to find if from shhgit binary's and current directory")
	configNameFlag   = flag.String("config-name", "config.yaml", "filename to search for")
)

type MatchEvent struct {
	Url       string
	Matches   []string
	Signature string
	File      string
	Stars     int
	Source    core.GitResourceType
}

var session *core.Session

func ProcessRepositories(s *core.Session) {
	threadNum := *s.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			for {
				timeout := time.Duration(*s.Options.CloneRepositoryTimeout) * time.Second
				_, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				repository := <-s.Repositories

				repo, err := core.GetRepository(session, repository.Id)

				if err != nil {
					s.Log.Warn("Failed to retrieve repository %d: %s", repository.Id, err)
					continue
				}

				if repo.GetPermissions()["pull"] &&
					uint(repo.GetStargazersCount()) >= *s.Options.MinimumStars &&
					uint(repo.GetSize()) < *s.Options.MaximumRepositorySize {

					processRepositoryOrGist(s, repo.GetCloneURL(), repository.Ref, repo.GetStargazersCount(), core.GITHUB_SOURCE)
				}
			}
		}(i)
	}
}

func ProcessGists(s *core.Session) {
	threadNum := *s.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			for {
				gistUrl := <-s.Gists
				processRepositoryOrGist(s, gistUrl, "", -1, core.GIST_SOURCE)
			}
		}(i)
	}
}

func ProcessComments(s *core.Session) {
	threadNum := *s.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			for {
				commentBody := <-s.Comments
				dir := core.GetTempDir(s, core.GetHash(commentBody))
				ioutil.WriteFile(filepath.Join(dir, "comment.ignore"), []byte(commentBody), 0644)

				if !checkSignatures(s, dir, "ISSUE", 0, core.GITHUB_COMMENT) {
					os.RemoveAll(dir)
				}
			}
		}(i)
	}
}

func processRepositoryOrGist(s *core.Session, url string, ref string, stars int, source core.GitResourceType) {
	var (
		matchedAny bool = false
	)

	dir := core.GetTempDir(s, core.GetHash(url))
	_, err := core.CloneRepository(session, url, ref, dir)

	if err != nil {
		s.Log.Debug("[%s] Cloning failed: %s", url, err.Error())
		os.RemoveAll(dir)
		return
	}

	s.Log.Debug("[%s] Cloning %s in to %s", url, ref, strings.Replace(dir, *s.Options.TempDirectory, "", -1))
	matchedAny = checkSignatures(s, dir, url, stars, source)
	if !matchedAny {
		os.RemoveAll(dir)
	}
}

func checkSignatures(s *core.Session, dir string, url string, stars int, source core.GitResourceType) (matchedAny bool) {

	for _, file := range core.GetMatchingFiles(s, dir) {
		var (
			matches          []string
			relativeFileName string
		)
		if strings.Contains(dir, *s.Options.TempDirectory) {
			relativeFileName = strings.Replace(file.Path, *s.Options.TempDirectory, "", -1)
		} else {
			relativeFileName = strings.Replace(file.Path, dir, "", -1)
		}

		if *s.Options.SearchQuery != "" {
			queryRegex := regexp.MustCompile(*s.Options.SearchQuery)
			for _, match := range queryRegex.FindAllSubmatch(file.Contents, -1) {
				matches = append(matches, string(match[0]))
			}

			if matches != nil {
				count := len(matches)
				m := strings.Join(matches, ", ")
				s.Log.Important("[%s] %d %s for %s in file %s: %s", url, count, core.Pluralize(count, "match", "matches"), color.GreenString("Search Query"), relativeFileName, color.YellowString(m))
				s.WriteToCsv([]string{url, "Search Query", relativeFileName, m})
			}
		} else {
			for _, signature := range s.Signatures {
				if matched, part := signature.Match(file); matched {
					if part == core.PartContents {
						if matches = signature.GetContentsMatches(s, file.Contents); len(matches) > 0 {
							count := len(matches)
							m := strings.Join(matches, ", ")
							publish(s, &MatchEvent{Source: source, Url: url, Matches: matches, Signature: signature.Name(), File: relativeFileName, Stars: stars})
							matchedAny = true

							s.Log.Important("[%s] %d %s for %s in file %s: %s", url, count, core.Pluralize(count, "match", "matches"), color.GreenString(signature.Name()), relativeFileName, color.YellowString(m))
							s.WriteToCsv([]string{url, signature.Name(), relativeFileName, m})
						}
					} else {
						if *s.Options.PathChecks {
							publish(s, &MatchEvent{Source: source, Url: url, Matches: matches, Signature: signature.Name(), File: relativeFileName, Stars: stars})
							matchedAny = true

							s.Log.Important("[%s] Matching file %s for %s", url, color.YellowString(relativeFileName), color.GreenString(signature.Name()))
							s.WriteToCsv([]string{url, signature.Name(), relativeFileName, ""})
						}

						if *s.Options.EntropyThreshold > 0 && file.CanCheckEntropy(s) {
							scanner := bufio.NewScanner(bytes.NewReader(file.Contents))

							for scanner.Scan() {
								line := scanner.Text()

								if len(line) > 6 && len(line) < 100 {
									entropy := core.GetEntropy(line)

									if entropy >= *s.Options.EntropyThreshold {
										blacklistedMatch := false

										for _, blacklistedString := range s.Config.BlacklistedStrings {
											if strings.Contains(strings.ToLower(line), strings.ToLower(blacklistedString)) {
												blacklistedMatch = true
											}
										}

										if !blacklistedMatch {
											publish(s, &MatchEvent{Source: source, Url: url, Matches: []string{line}, Signature: "High entropy string", File: relativeFileName, Stars: stars})
											matchedAny = true

											s.Log.Important("[%s] Potential secret in %s = %s", url, color.YellowString(relativeFileName), color.GreenString(line))
											s.WriteToCsv([]string{url, "High entropy string", relativeFileName, line})
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if !matchedAny && len(*s.Options.Local) <= 0 {
			os.Remove(file.Path)
		}
	}
	return
}

func publish(s *core.Session, event *MatchEvent) {
	// todo: implement a modular plugin system to handle the various outputs (console, live, csv, webhooks, etc)
	if len(*s.Options.Live) > 0 {
		data, _ := json.Marshal(event)
		http.Post(*s.Options.Live, "application/json", bytes.NewBuffer(data))
	}
}

func main() {
	flag.Parse()
	ctx := context.Background()

	s, err := core.NewSession(ctx, &core.Options{
		Threads:                threadsFlag,
		Silent:                 silentFlag,
		Debug:                  debugFlag,
		MaximumRepositorySize:  maxRepoFlag,
		MaximumFileSize:        maxFileSizeFlag,
		CloneRepositoryTimeout: cloneTimeoutFlag,
		EntropyThreshold:       entropyFlag,
		MinimumStars:           minStarsFlag,
		PathChecks:             pathChecksFlag,
		ProcessGists:           gistsFlag,
		TempDirectory:          tempDirFlag,
		CSVPath:                csvFlag,
		SearchQuery:            searchQueryFlag,
		Local:                  localFlag,
		Live:                   liveFlag,
		ConfigPath:             configPathFlag,
		ConfigName:             configNameFlag,
	})

	if err != nil {
		fmt.Printf("Failed to create session: %v\n", err)
		os.Exit(1)
	}

	s.Log.Info(color.HiBlueString(core.Banner))
	s.Log.Info("\t%s\n", color.HiCyanString(core.Author))
	s.Log.Info("[*] Loaded %s signatures. Using %s worker threads. Temp work dir: %s\n", color.BlueString("%d", len(s.Signatures)), color.BlueString("%d", *s.Options.Threads), color.BlueString(*s.Options.TempDirectory))

	if s.Options.Local != nil {
		s.Log.Info("[*] Scanning local directory: %s - skipping public repository checks...", color.BlueString(*s.Options.Local))
		rc := 0
		if checkSignatures(s, *s.Options.Local, *s.Options.Local, -1, core.LOCAL_SOURCE) {
			rc = 1
		} else {
			s.Log.Info("[*] No matching secrets found in %s!", color.BlueString(*s.Options.Local))
		}
		os.Exit(rc)
	}
	if *s.Options.SearchQuery != "" {
		s.Log.Important("Search Query '%s' given. Only returning matching results.", *s.Options.SearchQuery)
	}

	go core.GetRepositories(s)
	go ProcessRepositories(s)
	go ProcessComments(s)

	if *s.Options.ProcessGists {
		go core.GetGists(s)
		go ProcessGists(s)
	}

	spinny := core.ShowSpinner()
	select {}
	spinny()
}
