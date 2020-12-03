package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStripIgnoredLinesMultiLineWithoutPragma(t *testing.T) {
	multiLineString := `line1
	line2
	line3`
	multiLineByteArray := []byte(multiLineString)
	assert.Equal(t, StripIgnoredLines(multiLineByteArray, []string{"PRAGMA"}), multiLineByteArray)
}

func TestStripIgnoredLinesMultiLineWithPragma(t *testing.T) {
	multiLineString := `line1
	line2 #PRAGMA
	line3`
	multiLineByteArray := []byte(multiLineString)

	expectedString := `line1
	line3`
	expectedByteArray := []byte(expectedString)

	assert.Equal(t, StripIgnoredLines(multiLineByteArray, []string{"PRAGMA"}), expectedByteArray)
}

func TestStripIgnoredLinesSingleLineWithoutPragma(t *testing.T) {
	multiLineString := `line1`
	multiLineByteArray := []byte(multiLineString)

	assert.Equal(t, StripIgnoredLines(multiLineByteArray, []string{"PRAGMA"}), multiLineByteArray)
}

func TestStripIgnoredLinesSingleLineWithPragma(t *testing.T) {
	multiLineString := `line1 #PRAGMA`
	multiLineByteArray := []byte(multiLineString)
	assert.Equal(t, StripIgnoredLines(multiLineByteArray, []string{"PRAGMA"}), []byte{})
}

func TestStripIgnoredLinesEmptyText(t *testing.T) {
	multiLineString := ``
	multiLineByteArray := []byte(multiLineString)
	assert.Equal(t, StripIgnoredLines(multiLineByteArray, []string{"PRAGMA"}), []byte(nil))
}

func TestStripIgnoredLinesMultiLineWithMultiPragma(t *testing.T) {
	multiLineString := `line1 #PRAGMA
line2
line3 #PRAGMA`
	multiLineByteArray := []byte(multiLineString)

	expectedString := `line2`
	expectedByteArray := []byte(expectedString)

	assert.Equal(t, StripIgnoredLines(multiLineByteArray, []string{"PRAGMA"}), expectedByteArray)
}

func TestStripIgnoredLinesMultiLineWithMultiPragmas(t *testing.T) {
	multiLineString := `line1 #PRAGMA
line2
line3 #AMGARP`
	multiLineByteArray := []byte(multiLineString)

	expectedString := `line2`
	expectedByteArray := []byte(expectedString)

	assert.Equal(t, StripIgnoredLines(multiLineByteArray, []string{"PRAGMA", "AMGARP"}), expectedByteArray)
}
