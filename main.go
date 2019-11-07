package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/eth0izzle/shhgit/core"
	"github.com/fatih/color"
	//_ "net/http/pprof"
)

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

					processRepositoryOrGist(repo.GetCloneURL())
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
				processRepositoryOrGist(gistUrl)
			}
		}(i)
	}
}

func processRepositoryOrGist(url string) {
	var (
		matches []string
		//matchedAny bool = false
	)

	dir := core.GetTempDir(core.GetHash(url))
	_, err := core.CloneRepository(session, url, dir)

	if err != nil {
		session.Log.Debug("[%s] Cloning failed: %s", url, err.Error())
		os.RemoveAll(dir)
		return
	}

	session.Log.Debug("[%s] Cloning in to %s", url, strings.Replace(dir, *session.Options.TempDirectory, "", -1))

	maxFileSize := *session.Options.MaximumFileSize * 1024
	defer os.RemoveAll(dir)

	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil || f.IsDir() || uint(f.Size()) > maxFileSize || core.IsSkippableFile(path) {
			return nil
		}

		file := core.NewMatchFile(path)
		defer os.Remove(file.Path)
		relativeFileName := strings.Replace(file.Path, *session.Options.TempDirectory, "", -1)

		if *session.Options.SearchQuery != "" {
			queryRegex := regexp.MustCompile(*session.Options.SearchQuery)
			for _, match := range queryRegex.FindAllSubmatch(file.Contents, -1) {
				matches = append(matches, string(match[0]))
			}

			if matches != nil {
				count := len(matches)
				m := strings.Join(matches, ", ")
				session.Log.ImportantFile("[%s] %d %s for %s in file %s: %s", &file, url, count, core.Pluralize(count, "match", "matches"), color.GreenString("Search Query"), relativeFileName, color.YellowString(m))
				session.WriteToCsv([]string{url, "Search Query", relativeFileName, m})
			}
		} else {
			for _, signature := range session.Signatures {
				if matched, part := signature.Match(file); matched {

					if part == core.PartContents {
						if matches = signature.GetContentsMatches(file); matches != nil {
							count := len(matches)
							m := strings.Join(matches, ", ")
							session.Log.ImportantFile("[%s] %d %s for %s in file %s: %s", &file, url, count, core.Pluralize(count, "match", "matches"), color.GreenString(signature.Name()), relativeFileName, color.YellowString(m))
							session.WriteToCsv([]string{url, signature.Name(), relativeFileName, m})
						}
					} else {
						if *session.Options.PathChecks {
							session.Log.ImportantFile("[%s] Matching file %s for %s", &file, url, color.YellowString(relativeFileName), color.GreenString(signature.Name()))
							session.WriteToCsv([]string{url, signature.Name(), relativeFileName, ""})
						}

						if *session.Options.EntropyThreshold > 0 && file.CanCheckEntropy() {
							scanner := bufio.NewScanner(bytes.NewReader(file.Contents))

							for scanner.Scan() {
								line := scanner.Text()

								if len(line) > 6 && len(line) < 100 {
									entropy := core.GetEntropy(scanner.Text())

									if entropy >= *session.Options.EntropyThreshold {
										session.Log.ImportantFile("[%s] Potential secret in %s = %s", &file, url, color.YellowString(relativeFileName), color.GreenString(scanner.Text()))
										session.WriteToCsv([]string{url, signature.Name(), relativeFileName, scanner.Text()})
									}
								}
							}
						}
					}
				}
			}
		}

		return nil
	})
}

func main() {
	session.Log.Info("%s v%s started. Loaded %d signatures. Using %d GitHub tokens and %d threads. Work dir: %s", core.Name, core.Version, len(session.Signatures), len(session.Clients), *session.Options.Threads, *session.Options.TempDirectory)

	if *session.Options.SearchQuery != "" {
		session.Log.Important("Search Query '%s' given. Only returning matching results.", *session.Options.SearchQuery)
	}

	if *session.Options.NoColor == true {
		color.NoColor = true
	}

	go core.GetRepositories(session)
	go ProcessRepositories()

	if *session.Options.ProcessGists {
		go core.GetGists(session)
		go ProcessGists()
	}

	session.Log.Info("Press Ctrl+C to stop and exit.\n")
	//log.Fatal(http.ListenAndServe(":8080", nil))
	select {}
}
