package types

import (
	"io/ioutil"
	"path/filepath"
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

func (match MatchFile) CanCheckEntropy(blacklistedEntropyExtensions []string) bool {
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
