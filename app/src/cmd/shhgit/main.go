package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/session"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/types"
	"github.com/fatih/color"
	"github.com/hillu/go-yara/v4"
)

var activeSession = session.Fetch()

func processRepository(repository types.RepositoryResource) ([]types.Match, error) {
	var repositoryPath = repository.Location
	var isLocalRepository = repository.Type == types.LOCAL_SOURCE

	// no need to clone a local repository
	if !isLocalRepository {
		var err error = nil

		if (repository.Stars >= 0 && repository.Stars > *activeSession.Options.MinimumStars) || (repository.Size >= 0 && repository.Size > *activeSession.Options.MaximumRepositorySize) {
			return nil, fmt.Errorf("Skipping %s as it exceeeds star/size critera (%d stars, %dkb)", repository.Location, repository.Stars, repository.Size)
		}

		repositoryPath = helpers.GetTempDir(*activeSession.Options.WorkingDirectory, helpers.GetHash(string(repository.Location)))
		if repository.Type == types.GITHUB_COMMENT {
			// we write the comment itself to disk for easier processing, i.e. treat it as if it was a repo
			err = ioutil.WriteFile(filepath.Join(repositoryPath, "comment.ignore"), []byte(repository.Description), 0644)
		} else if repository.SCM == types.GIT_SCM {
			err = helpers.CloneGitRepository(repository.Location, repositoryPath, *activeSession.Options.CloneRepositoryTimeout)
		} else if repository.SCM == types.HG_SCM {
			err = helpers.CloneMercurialRepository(repository.Location, repositoryPath, *activeSession.Options.CloneRepositoryTimeout)
		}

		// clean up if the clone failed
		if err != nil {
			os.RemoveAll(repositoryPath)
			return nil, fmt.Errorf("Failed to process %s: %s", repository.Location, err)
		}
	}

	scanner, err := yara.NewScanner(activeSession.ScannerRules)
	if err != nil {
		return nil, fmt.Errorf("Failed to create rules scanner engine: %s", err)
	}

	scanner.SetTimeout(time.Duration(*activeSession.Options.ScanFileTimeout) * time.Second)
	//scanner.SetFlags(yara.ScanFlagsFastMode)
	scanner.DefineVariable("repository_name", repository.Name)
	scanner.DefineVariable("repository_description", repository.Description)
	scanner.DefineVariable("repository_owner", repository.Owner)
	scanner.DefineVariable("repository_size", repository.Size)
	scanner.DefineVariable("repository_stars", repository.Stars)

	var matches = make([]types.Match, 0)
	for _, file := range types.GetProcessableFiles(repositoryPath, *activeSession.Options.MaximumFileSize, activeSession.Config.BlackLists) {
		matchTemplate := types.Match{
			RepositoryLocalPath: repositoryPath,
			RepositoryURL:       repository.Location,
			Filename:            file.Filename,
			RepositoryName:      repository.Name,
			RepositoryStars:     repository.Stars,
			RepositoryOwner:     repository.Owner,
			RepositorySize:      repository.Size,
			FileSize:            file.Size,
			FileExtension:       file.Extension,
			RepositorySource:    repository.Type,
			Found:               time.Now(),
		}

		// manual entropy checks (todo: move to yara using the math module)
		if *activeSession.Options.EntropyThreshold > 0 && file.CanCheckEntropy(activeSession.Config.BlackLists.EntropyExtensions) {
			entropyLines, err := file.GetEntropy(*activeSession.Options.EntropyThreshold, activeSession.Config.BlackLists.Strings)

			if err != nil {
				activeSession.Log.Debug("Cannot read file %s: %s", file.Path, err)
				continue
			}

			if len(entropyLines) > 0 {
				matchTemplate.Rule = "High entropy string"
				matchTemplate.Matches = entropyLines
				matches = append(matches, matchTemplate)
			}
		}

		scanner.DefineVariable("filename", file.Filename)
		scanner.DefineVariable("filepath", file.Path)
		scanner.DefineVariable("extension", file.Extension)

		var yaraMatches yara.MatchRules
		scanner.SetCallback(&yaraMatches)
		scanner.ScanFile(file.Path)

		for _, match := range yaraMatches {
			matchStrings := make([]string, 0)

			for _, matchString := range match.Strings {
				data := string(matchString.Data)
				if len(data) > 0 && !helpers.Contains(matchStrings, data) {
					matchStrings = append(matchStrings, data)
				}
			}

			matchTemplate.Rule = match.Rule
			helpers.ReverseSlice(matchStrings)
			matchTemplate.Matches = matchStrings
			matches = append(matches, matchTemplate)
		}
	}

	if !isLocalRepository {
		os.RemoveAll(repositoryPath)
	}

	return matches, nil
}

func publish(match types.Match) {
	activeSession.Log.Info(match.String())

	for _, output := range activeSession.Outputs {
		output.Publish(match)
	}
}

func main() {
	activeSession.Log.Info(color.HiBlueString(settings.Banner))
	activeSession.Log.Info("\t%s\n", color.HiCyanString(settings.Author))

	if len(*activeSession.Options.LocalDirectory) > 0 {
		activeSession.Log.Info("[*] Scanning local directory: %s (skipping public repository checks)...", color.BlueString(*activeSession.Options.LocalDirectory))

		dirSize, _ := helpers.GetDirectorySize(*activeSession.Options.LocalDirectory)
		matches, _ := processRepository(types.RepositoryResource{
			Type:     types.LOCAL_SOURCE,
			Location: *activeSession.Options.LocalDirectory,
			Name:     "LOCAL",
			Size:     dirSize,
			Stars:    -1,
		})

		for _, match := range matches {
			publish(match)
		}
	} else {
		activeSession.Log.Info("[*] Using %s workers. Temp work dir: %s\n", color.BlueString("%d", *activeSession.Options.Workers), color.BlueString(*activeSession.Options.WorkingDirectory))

		repositories := make(chan types.RepositoryResource, 1000)
		activeSession.FetchFromSources(repositories)

		workers := *activeSession.Options.Workers
		for workerID := 0; workerID < workers; workerID++ {
			go func(workerID int) {
				for repository := range repositories {
					matches, _ := processRepository(repository)

					for _, match := range matches {
						publish(match)
					}
				}
			}(workerID)
		}

		// forces program to run indefinitely until it's killed
		select {}
	}
}
