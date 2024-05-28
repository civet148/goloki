package goloki

import (
	"github.com/civet148/log"
	"testing"
)

func TestLoki(t *testing.T) {
	log.SetLevel("debug")
	loki := NewLokiClient("http://127.0.0.1:3100", map[string]string{
		"app": "gotest",
	})
	loki.Println("hello world", map[string]string{
		"level": "INFO",
		"tag":   "code farmer",
	})
}
