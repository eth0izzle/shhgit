package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/fatih/color"
)

const (
	FATAL     = 5
	ERROR     = 4
	IMPORTANT = 3
	WARN      = 2
	INFO      = 1
	DEBUG     = 0
)

var LogColors = map[int]*color.Color{
	FATAL:     color.New(color.FgRed).Add(color.Bold),
	ERROR:     color.New(color.FgRed),
	WARN:      color.New(color.FgYellow),
	IMPORTANT: color.New(),
	DEBUG:     color.New(color.Faint),
}

type Logger struct {
	sync.Mutex

	debug  bool
	silent bool
}

func (l *Logger) SetDebug(d bool) {
	l.debug = d
}

func (l *Logger) SetSilent(d bool) {
	l.silent = d
}

func (l *Logger) Log(level int, format string, file *MatchFile, args ...interface{}) {
	l.Lock()
	defer l.Unlock()

	if level == DEBUG && !l.debug {
		return
	}

	if l.silent && level < IMPORTANT {
		return
	}

	if c, ok := LogColors[level]; ok {
		c.Printf(format+"\n", args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}

	if level > INFO && session.Config.SlackWebhook != "" {
		values := map[string]string{"text": fmt.Sprintf(format+"\n", args...)}
		jsonValue, _ := json.Marshal(values)
		http.Post(session.Config.SlackWebhook, "application/json", bytes.NewBuffer(jsonValue))
	}

	if session.Config.Telegram.Token != "" && session.Config.Telegram.ChatID != "" {
		caption := fmt.Sprintf(format+"\n", args...)
		rcpt := session.Config.Telegram.ChatID
		if level != IMPORTANT && session.Config.Telegram.AdminID != "" {
			rcpt = session.Config.Telegram.AdminID
		}

		if file != nil {
			if len(caption) > 1023 {
				caption = caption[0:1023]
			}

			values := map[string]string{
				"caption":    caption,
				"chat_id":    rcpt,
				"parse_mode": "Markdown",
			}
			requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", session.Config.Telegram.Token)
			request, err := NewfileUploadRequest(requestURL, values, "document", file)
			if err != nil {
				log.Fatal(err)
			}
			client := &http.Client{}

			client.Do(request)
		} else {
			values := map[string]string{
				"text":       caption,
				"chat_id":    rcpt,
				"parse_mode": "Markdown",
			}
			jsonValue, _ := json.Marshal(values)
			requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", session.Config.Telegram.Token)
			http.Post(requestURL, "application/json", bytes.NewBuffer(jsonValue))
		}

	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Log(FATAL, format, nil, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(ERROR, format, nil, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.Log(WARN, format, nil, args...)
}

func (l *Logger) Important(format string, args ...interface{}) {
	l.Log(IMPORTANT, format, nil, args...)
}

func (l *Logger) ImportantFile(format string, file *MatchFile, args ...interface{}) {
	l.Log(IMPORTANT, format, file, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(INFO, format, nil, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(DEBUG, format, nil, args...)
}
