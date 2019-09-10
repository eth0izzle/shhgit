package core

import (
	"regexp"
	"regexp/syntax"
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
	GetContentsMatches(file MatchFile) []string
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

func (s SimpleSignature) GetContentsMatches(file MatchFile) []string {
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

func (s PatternSignature) GetContentsMatches(file MatchFile) []string {
	matches := make([]string, 0)

	for _, match := range s.match.FindAllSubmatch(file.Contents, -1) {
		matches = append(matches, string(match[0]))
	}

	return matches
}

func (s PatternSignature) Name() string {
	return s.name
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
