package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

// A MatchEvent holds a list of matches for a particular Url and signature
type MatchEvent struct {
	Url       string
	Matches   []string
	Signature string
	File      string
	Stars     int
	Source    GitResourceType
}

// Line returns a string containing the following: Url, Signature, File, and
// Matches
func (m MatchEvent) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Url: %s  ", m.Url))
	b.WriteString(fmt.Sprintf("Signature: %s  ", m.Signature))
	b.WriteString(fmt.Sprintf("File: %s  ", m.File))
	b.WriteString(fmt.Sprintf("Matches: %s\n", strings.Join(m.Matches, ", ")))

	return b.String()
}

// Line returns a slice of strings containing the following: Url, Signature,
// File, and Matches
func (m MatchEvent) Line() []string {
	return []string{m.Url, m.Signature, m.File, strings.Join(m.Matches, ", ")}
}

// Json returns a JSON formatted string that includes all of the data in a
// MatchEvent.
func (m MatchEvent) Json() string {
	b, err := json.Marshal(m)
	if err != nil {
		LogIfError("unable to create JSON, %s", err)
		return ""
	}

	return string(b)
}
