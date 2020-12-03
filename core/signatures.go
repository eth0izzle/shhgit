package core

import (
	"bytes"
	"regexp"
	"regexp/syntax"
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
	Name() string
	Match(file MatchFile) (bool, string)
	GetContentsMatches(contents []byte) []string
}

type SimpleSignature struct {
	part  string
	match string
	name  string
}

type PatternSignature struct {
	part  string
	match *regexp.Regexp
	name  string
}

func (s SimpleSignature) Match(file MatchFile) (bool, string) {
	var (
		haystack  *string
		matchPart = ""
	)

	switch s.part {
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

	return (s.match == *haystack), matchPart
}

func (s SimpleSignature) GetContentsMatches(contents []byte) []string {
	return nil
}

func (s SimpleSignature) Name() string {
	return s.name
}

func (s PatternSignature) Match(file MatchFile) (bool, string) {
	var (
		haystack  *string
		matchPart = ""
	)

	switch s.part {
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
		return s.match.Match(file.Contents), PartContents
	default:
		return false, matchPart
	}

	return s.match.MatchString(*haystack), matchPart
}

func (s PatternSignature) GetContentsMatches(contents []byte) []string {
	matches := make([]string, 0)
	if len(session.Config.IgnoreLinePragmas) > 0 {
		contents = StripIgnoredLines(contents, session.Config.IgnoreLinePragmas)
	}
	for _, match := range s.match.FindAllSubmatch(contents, -1) {
		match := string(match[0])
		blacklistedMatch := false

		for _, blacklistedString := range session.Config.BlacklistedStrings {
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

func (s PatternSignature) Name() string {
	return s.name
}

func StripIgnoredLines(contents []byte, ignoreLinePragmas []string) []byte {
	processedLines := make([][]byte, 0)
	newlineByte := []byte("\n")

	lines := strings.Split(string(contents[:]), "\n")
	for _, line := range lines {
		ignoreCurrentLine := false
		for _, ignoreLinePragma := range ignoreLinePragmas {
			if strings.Contains(strings.ToLower(line), strings.ToLower(ignoreLinePragma)) {
				ignoreCurrentLine = true
			}
		}
		if !ignoreCurrentLine {
			processedLines = append(processedLines, []byte(line))
		}
	}
	return bytes.Join(processedLines, newlineByte)
}

func GetSignatures(s *Session) []Signature {
	var signatures []Signature
	for _, signature := range s.Config.Signatures {
		if signature.Match != "" {
			signatures = append(signatures, SimpleSignature{
				name:  signature.Name,
				part:  signature.Part,
				match: signature.Match,
			})
		} else {
			if _, err := syntax.Parse(signature.Match, syntax.FoldCase); err == nil {
				signatures = append(signatures, PatternSignature{
					name:  signature.Name,
					part:  signature.Part,
					match: regexp.MustCompile(signature.Regex),
				})
			}
		}
	}

	return signatures
}
