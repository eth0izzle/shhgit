package core

import (
	"fmt"
	"net/http"
	"net"
	"os"
	"regexp"
	"strings"
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

func (l *Logger) Log(level int, format string, args ...interface{}) {
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

	if level > WARN && session.Config.Webhook != "" {
		text := colorStrip(fmt.Sprintf(format, args...))
		payload := fmt.Sprintf(session.Config.WebhookPayload, text)
		http.Post(session.Config.Webhook, "application/json", strings.NewReader(payload))
	}
	
	if session.Config.Logstash != "" {
		text := colorStrip(fmt.Sprintf(format, args...))
		payload := fmt.Sprintf(session.Config.Logstash, text)				
  		pc, err := net.ListenPacket("udp4", ":" + session.Config.LogstashPort)
  		if err != nil {
  			fmt.Printf("Logstash: Unknown Port\n")
  		}
  		defer pc.Close()

  		addr, err := net.ResolveUDPAddr("udp4", session.Config.Logstash)
  		if err != nil {
  			fmt.Printf("Logstash: No Such Host\n")	
  		}
  		
  		pc.WriteTo([]byte(payload), addr)	
	}
	
	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Log(FATAL, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(ERROR, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.Log(WARN, format, args...)
}

func (l *Logger) Important(format string, args ...interface{}) {
	l.Log(IMPORTANT, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(INFO, format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(DEBUG, format, args...)
}

func colorStrip(str string) string {
	ansi := "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	re := regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}
