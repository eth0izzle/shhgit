package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// A MatchEvent holds a list of matches for a particular Url and signature
type Match struct {
	Rule                string
	Matches             []string
	Filename            string
	FileExtension       string
	FileSize            int64
	RepositoryName      string
	RepositoryURL       string
	RepositoryLocalPath string
	RepositorySize      int64
	RepositoryStars     int
	RepositoryOwner     string
	Found               time.Time
	RepositorySource    RepositoryResourceType
}

// Line returns a string containing the following: Url, Signature, File, and
// Matches
func (m Match) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Rule: %s  ", m.Rule))
	b.WriteString(fmt.Sprintf("File: %s  ", m.Filename))
	b.WriteString(fmt.Sprintf("Source: %s  ", m.RepositorySource))
	b.WriteString(fmt.Sprintf("Matches:\n\t%s", strings.Join(m.Matches, "\n\t")))

	return b.String()
}

// Json returns a JSON formatted string that includes all of the data in a
// MatchEvent.
func (m Match) JSON() string {
	b, _ := json.Marshal(m)

	return string(b)
}
