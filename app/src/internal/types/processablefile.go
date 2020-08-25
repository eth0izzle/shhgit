package types

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/settings"
)

type ProcessableFile struct {
	Path      string
	Filename  string
	Extension string
	Size      int64
}

func NewProcessableFile(path string, size int64) ProcessableFile {
	path = filepath.ToSlash(path)
	_, filename := filepath.Split(path)
	extension := filepath.Ext(path)

	return ProcessableFile{
		Path:      path,
		Filename:  filename,
		Extension: extension,
		Size:      size,
	}
}

func (match ProcessableFile) CanCheckEntropy(blacklistedEntropyExtensions []string) bool {
	if match.Filename == "id_rsa" {
		return false
	}

	for _, skippableExt := range blacklistedEntropyExtensions {
		if match.Extension == skippableExt {
			return false
		}
	}

	return true
}

func (match ProcessableFile) GetEntropy(entropyThreshold float64, blacklistedEntropyStrings []string) (entropy []string, err error) {
	file, err := os.Open(match.Path)
	defer file.Close()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot read file %s: %s", match.Path, err))
	}

	scanner := bufio.NewScanner(file)
	matches := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		lineLen := len(line)

		if lineLen > 10 && lineLen < 100 {
			entropy := helpers.GetEntropy(line)

			if entropy > entropyThreshold {
				blacklistedMatch := false

				for _, blacklistedString := range blacklistedEntropyStrings {
					blacklistedMatch = strings.Contains(strings.ToLower(line), strings.ToLower(blacklistedString))
				}

				if !blacklistedMatch {
					matches = append(matches, line)
				}
			}
		}
	}

	return matches, nil
}

func GetProcessableFiles(dir string, maximumFileSize int64, blacklists settings.ConfigBlacklists) []ProcessableFile {
	fileList := make([]ProcessableFile, 0)
	maximumFileSizeKb := maximumFileSize * 1024

	filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil || file.IsDir() || file.Size() > maximumFileSizeKb {
			return nil
		}

		extension := strings.ToLower(filepath.Ext(path))
		for _, skippableExt := range blacklists.Extensions {
			if extension == skippableExt {
				return nil
			}
		}

		for _, skippablePathIndicator := range blacklists.Paths {
			if strings.Contains(filepath.ToSlash(path), skippablePathIndicator) {
				return nil
			}
		}

		fileList = append(fileList, NewProcessableFile(path, file.Size()))

		return nil
	})

	return fileList
}
