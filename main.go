package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/eth0izzle/shhgit/core"
	"github.com/fatih/color"
)

var session = core.GetSession()

func ProcessEvents() {
	threadNum := *session.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			var (
				dir, url string
				matches  []string
			)

			for {
				repositoryName := <-session.Repositories
				repo := core.GetRepository(session, repositoryName)

				if repo.GetPermissions()["pull"] &&
					uint(repo.GetStargazersCount()) >= *session.Options.MinimumStars &&
					uint(repo.GetSize()) < *session.Options.MaximumRepositorySize {

					dir = core.GetTempDir(core.GetHash(repositoryName))
					url = core.GetRepositoryUrl(repositoryName)
					_, err := core.CloneRepository(session, url, dir)
					matchedAny := false

					if err != nil {
						session.Log.Debug("[%s] Cloning %s failed: %s", repositoryName, url, err.Error())
						os.RemoveAll(dir)
						continue
					}

					session.Log.Debug("[%s] Cloning %s to %s", repositoryName, url, dir)

					for _, file := range core.GetMatchingFiles(dir) {
						for _, signature := range session.Signatures {
							if matched, part := signature.Match(file); matched {
								matchedAny = true
								relativeFileName := strings.Replace(file.Path, *session.Options.TempDirectory, "", -1)

								if part == core.PartContents {
									if matches = signature.GetMatches(file); matches != nil {
										count := len(matches)
										m := strings.Join(matches, ", ")
										session.Log.Important("[%s] %d %s for %s in file %s: %s", repositoryName, count, core.Pluralize(count, "match", "matches"), color.GreenString(signature.Name()), relativeFileName, color.YellowString(m))
										session.WriteToCsv([]string{repositoryName, signature.Name(), relativeFileName, m})
									}
								} else {
									if *session.Options.PathChecks {
										session.Log.Important("[%s] Matching file %s for %s", repositoryName, color.YellowString(relativeFileName), color.GreenString(signature.Name()))
										session.WriteToCsv([]string{repositoryName, signature.Name(), relativeFileName, ""})
									}

									if *session.Options.EntropyThreshold > 0 && file.CanCheckEntropy() {
										scanner := bufio.NewScanner(bytes.NewReader(file.Contents))

										for scanner.Scan() {
											line := scanner.Text()

											if len(line) > 6 && len(line) < 100 {
												entropy := core.GetEntropy(scanner.Text())

												if entropy >= *session.Options.EntropyThreshold {
													session.Log.Important("[%s] Potential secret in %s = %s", repositoryName, color.YellowString(relativeFileName), color.GreenString(scanner.Text()))
													session.WriteToCsv([]string{repositoryName, signature.Name(), relativeFileName, scanner.Text()})
												}
											}
										}
									}
								}
							}
						}

						if !matchedAny {
							os.Remove(file.Path)
						}
					}

					if !matchedAny {
						os.RemoveAll(dir)
					}
				}
			}
		}(i)
	}
}

func main() {
	session.Log.Info("%s v%s started. Loaded %d signatures. Using %d threads. Work dir: %s", core.Name, core.Version, len(session.Signatures), *session.Options.Threads, *session.Options.TempDirectory)

	go core.ReadEvents(session)
	go ProcessEvents()

	session.Log.Info("Press Ctrl+C to stop and exit.\n")
	select {}
}
