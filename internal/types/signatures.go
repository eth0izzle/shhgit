package types

import (
	"regexp"
	"strings"
)

const (
	TypeSimple  = "simple"
	TypePattern = "pattern"

	PartExtension = "extension"
	PartFilename  = "filename"
	PartPath      = "path"
	PartContents  = "contents"
)

type Signature interface {
	Match(file MatchFile) (bool, string)
	GetContentsMatches(contents []byte, blacklistedStrings []string) []string
}

type SimpleSignature struct {
	Part    string
	MatchOn string
	Name    string
}

type PatternSignature struct {
	Part    string
	MatchOn *regexp.Regexp
	Name    string
}

func (s SimpleSignature) Match(file MatchFile) (bool, string) {
	var (
		haystack  *string
		matchPart = ""
	)

	switch s.Part {
	case PartPath:
		haystack = &file.Path
		matchPart = PartPath
	case PartFilename:
		haystack = &file.Filename
		matchPart = PartPath
	case PartExtension:
		haystack = &file.Extension
		matchPart = PartPath
	default:
		return false, matchPart
	}

	return (s.MatchOn == *haystack), matchPart
}

func (s SimpleSignature) GetContentsMatches(contents []byte, blacklistedStrings []string) []string {
	return nil
}

func (s PatternSignature) Match(file MatchFile) (bool, string) {
	var (
		haystack  *string
		matchPart = ""
	)

	switch s.Part {
	case PartPath:
		haystack = &file.Path
		matchPart = PartPath
	case PartFilename:
		haystack = &file.Filename
		matchPart = PartFilename
	case PartExtension:
		haystack = &file.Extension
		matchPart = PartExtension
	case PartContents:
		return s.MatchOn.Match(file.Contents), PartContents
	default:
		return false, matchPart
	}

	return s.MatchOn.MatchString(*haystack), matchPart
}

func (s PatternSignature) GetContentsMatches(contents []byte, blacklistedStrings []string) []string {
	matches := make([]string, 0)

	for _, match := range s.MatchOn.FindAllSubmatch(contents, -1) {
		match := string(match[0])
		blacklistedMatch := false

		for _, blacklistedString := range blacklistedStrings {
			if strings.Contains(strings.ToLower(match), strings.ToLower(blacklistedString)) {
				blacklistedMatch = true
			}
		}

		if !blacklistedMatch {
			matches = append(matches, match)
		}
	}

	return matches
}
