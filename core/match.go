package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type MatchFile struct {
	Path      string
	Filename  string
	Extension string
	Contents  []byte
}

func NewMatchFile(path string) MatchFile {
	path = filepath.ToSlash(path)
	_, filename := filepath.Split(path)
	extension := filepath.Ext(path)
	contents, _ := ioutil.ReadFile(path)

	return MatchFile{
		Path:      path,
		Filename:  filename,
		Extension: extension,
		Contents:  contents,
	}
}

func IsSkippableFile(s *Session, path string) bool {
	extension := strings.ToLower(filepath.Ext(path))

	for _, skippableExt := range s.Config.BlacklistedExtensions {
		if extension == skippableExt {
			return true
		}
	}

	for _, skippablePathIndicator := range s.Config.BlacklistedPaths {
		skippablePathIndicator = strings.Replace(skippablePathIndicator, "{sep}", string(os.PathSeparator), -1)
		if strings.Contains(path, skippablePathIndicator) {
			return true
		}
	}

	return false
}

func (match MatchFile) CanCheckEntropy(s *Session) bool {
	if match.Filename == "id_rsa" {
		return false
	}

	for _, skippableExt := range s.Config.BlacklistedEntropyExtensions {
		if match.Extension == skippableExt {
			return false
		}
	}

	return true
}

func GetMatchingFiles(s *Session, dir string) []MatchFile {
	fileList := make([]MatchFile, 0)
	maxFileSize := *s.Options.MaximumFileSize * 1024

	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil || f.IsDir() || uint(f.Size()) > maxFileSize || IsSkippableFile(s, path) {
			return nil
		}
		fileList = append(fileList, NewMatchFile(path))
		return nil
	})

	return fileList
}
