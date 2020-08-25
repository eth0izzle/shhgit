package helpers

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func GetTempDir(prefix string, suffix string) string {
	dir := filepath.Join(prefix, suffix)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	} else {
		os.RemoveAll(dir)
	}

	return dir
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func GetHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}

func Pluralize(count int, singular string, plural string) string {
	if count == 1 {
		return singular
	}

	return plural
}

func GetEntropy(data string) (entropy float64) {
	if data == "" {
		return 0
	}

	for i := 0; i < 256; i++ {
		px := float64(strings.Count(data, string(byte(i)))) / float64(len(data))
		if px > 0 {
			entropy += -px * math.Log2(px)
		}
	}

	return entropy
}

func ReverseSlice(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if strings.Contains(n, x) {
			return true
		}
	}
	return false
}

func GetFilesInPath(dir string, ext string) []string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get cwd?")
	}

	path := filepath.Join(cwd, dir)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("Couldn't read files in: ", path)
	}

	var matches []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ext {
			matches = append(matches, filepath.Join(path, file.Name()))
		}
	}

	return matches
}

func FetchUrlAs(url string, auth string, v interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if len(auth) > 0 {
		req.Header.Add("Authorization", auth)
	}

	if resp, err := http.DefaultClient.Do(req); err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			return errors.New("rate limited")
		} else if resp.StatusCode == http.StatusInternalServerError {
			return errors.New("internal server error")
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("got %s, wanted 200 OK", resp.Status)
		}

		if contents, err := ioutil.ReadAll(resp.Body); err == nil {
			return json.Unmarshal(contents, v)
		}
	}

	return err
}

func CloneGitRepository(url string, dir string, timeout int) error {
	timeoutSecs := time.Duration(timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecs)
	defer cancel()

	cloneCmd := exec.CommandContext(ctx, "git", "clone", url, dir, "--quiet", "--no-tags", "--single-branch", "--depth=1")
	if err := cloneCmd.Run(); err != nil {
		return err
	}

	return nil
}

func CloneMercurialRepository(url string, dir string, timeout int) error {
	timeoutSecs := time.Duration(timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecs)
	defer cancel()

	cloneCmd := exec.CommandContext(ctx, "hg", "clone", url, dir, "--stream")
	if err := cloneCmd.Run(); err != nil {
		return err
	}

	return nil
}

func GetRandomToken(tokens []string) string {
	numberOfTokens := len(tokens)
	return tokens[rand.Intn(numberOfTokens)]
}

func GetDirectorySize(dir string) (int64, error) {
	var size int64

	err := filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return err
	})

	return size, err
}
