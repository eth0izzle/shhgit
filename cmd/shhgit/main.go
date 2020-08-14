package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Velocidex/go-yara"
	"github.com/fatih/color"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/session"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/sources"
	"github.com/eth0izzle/shhgit/internal/types"
)

type matchEvent struct {
	URL       string
	Matches   []string
	Signature string
	File      string
	Stars     int
	Source    types.RepositoryResourceType
}

var activeSession = session.Start()

func ProcessRepository(repository types.RepositoryResource) {
	var err error = nil
	var wg sync.WaitGroup
	wg.Add(1)

	var repositoryPath = repository.Location
	var isLocalRepository = repository.Type == types.LOCAL_SOURCE

	// no need to clone a local repository
	if !isLocalRepository {
		activeSession.Log.Debug("Processing %s from %v...", repository.Location, repository.Type)
		repositoryPath = helpers.GetTempDir(*activeSession.Options.WorkingDirectory, helpers.GetHash(string(repository.Location)))

		if repository.SCM == types.GIT_SCM {
			err = helpers.CloneGitRepository(repository.Location, repositoryPath, *activeSession.Options.CloneRepositoryTimeout)
		} else if repository.SCM == types.HG_SCM {
			err = helpers.CloneMercurialRepository(repository.Location, repositoryPath, *activeSession.Options.CloneRepositoryTimeout)
		}

		// clean up if the clone failed
		if err != nil {
			activeSession.Log.Debug("Failed to clone %s: %s", repository.Location, err)
			os.RemoveAll(repositoryPath)
			return
		}
	}

	for _, file := range helpers.GetCheckableFiles(repositoryPath, *activeSession.Options.MaximumFileSize, activeSession.Config.BlackLists) {
		relativeFilePath, _ := filepath.Rel(repositoryPath, file.Path)

		activeSession.Lock()
		activeSession.Rules.DefineVariable("filename", file.Filename)
		activeSession.Rules.DefineVariable("filepath", file.Path)
		activeSession.Rules.DefineVariable("extension", file.Extension)
		activeSession.Rules.DefineVariable("repository_name", repository.Name)
		activeSession.Rules.DefineVariable("repository_description", repository.Description)
		activeSession.Rules.DefineVariable("repository_size", repository.Size)
		activeSession.Rules.DefineVariable("repository_stars", repository.Stars)
		activeSession.Unlock()

		activeSession.Log.Debug("Scanning file %s", file.Path)
		matches, _ := activeSession.Rules.ScanFile(file.Path, yara.ScanFlagsFastMode, time.Duration(*activeSession.Options.ScanFileTimeout)*time.Second)

		for _, match := range matches {
			activeSession.Log.Important("[%v@%s->%s] matches %s (%s)", types.RepositoryResourceType(repository.Type), color.HiYellowString(repository.Name), color.YellowString(relativeFilePath), color.GreenString(match.Rule), (match.Meta["description"]))

			for _, matchString := range match.Strings {
				activeSession.Log.Important("\t[%s]", color.RedString(string(matchString.Data)))
			}
		}
	}

	if !isLocalRepository {
		os.RemoveAll(repositoryPath)
	}
}

// func processGists() {
// 	threadNum := *session.Options.Threads

// 	for i := 0; i < threadNum; i++ {
// 		go func(tid int) {
// 			for {
// 				gistUrl := <-activeSession.Gists
// 				processRepositoryOrGist(gistUrl, -1, core.GIST_SOURCE)
// 			}
// 		}(i)
// 	}
// }

// func processComments() {
// 	threadNum := *session.Options.Threads

// 	for i := 0; i < threadNum; i++ {
// 		go func(tid int) {
// 			for {
// 				commentBody := <-session.Comments
// 				dir := core.GetTempDir(core.GetHash(commentBody))
// 				ioutil.WriteFile(filepath.Join(dir, "comment.ignore"), []byte(commentBody), 0644)

// 				if !checkSignatures(dir, "ISSUE", 0, core.GITHUB_COMMENT) {
// 					os.RemoveAll(dir)
// 				}
// 			}
// 		}(i)
// 	}
// }

// func publish(event *MatchEvent) {
// 	// todo: implement a modular plugin system to handle the various outputs (console, live, csv, webhooks, etc)
// 	if len(*session.Options.Live) > 0 {
// 		data, _ := json.Marshal(event)
// 		http.Post(*session.Options.Live, "application/json", bytes.NewBuffer(data))
// 	}
// }

func main() {
	activeSession.Log.Info(color.HiBlueString(settings.Banner))
	activeSession.Log.Info("\t%s\n", color.HiCyanString(settings.Author))

	if len(*activeSession.Options.LocalDirectory) > 0 {
		activeSession.Log.Info("[*] Scanning local directory: %s against %s rules (skipping public repository checks)...", color.BlueString(*activeSession.Options.LocalDirectory), color.BlueString("%d", len(activeSession.Rules.GetRules())))
		helpers.ShowSpinner()

		dirSize, _ := helpers.GetDirectorySize(*activeSession.Options.LocalDirectory)
		ProcessRepository(types.RepositoryResource{
			Type:     types.LOCAL_SOURCE,
			Location: *activeSession.Options.LocalDirectory,
			Name:     "LOCAL",
			Size:     dirSize,
			Stars:    -1,
		})
	} else {
		activeSession.Log.Info("[*] Loaded %s rules. Using %s worker threads. Temp work dir: %s\n", color.BlueString("%d", len(activeSession.Rules.GetRules())), color.BlueString("%d", *activeSession.Options.Threads), color.BlueString(*activeSession.Options.WorkingDirectory))
		helpers.ShowSpinner()

		go sources.StartGitHub(activeSession.Config.Sources.GitHub, activeSession.Repositories)
		go sources.StartGitLab(activeSession.Config.Sources.GitLab, activeSession.Repositories)
		go sources.StartBitBucket(activeSession.Config.Sources.BitBucket, activeSession.Repositories)

		for i := 0; i < *activeSession.Options.Threads; i++ {
			go func() {
				for {
					timeout := time.Duration(*activeSession.Options.CloneRepositoryTimeout) * time.Second
					_, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					repository := <-activeSession.Repositories
					if (repository.Stars < 0 || repository.Stars >= int(*activeSession.Options.MinimumStars)) && (repository.Size < 0 || repository.Size <= *activeSession.Options.MaximumRepositorySize) {
						ProcessRepository(repository)
					}
				}
			}()
		}

		// forces program to run indefinitely until it's killed
		select {}
	}
}
