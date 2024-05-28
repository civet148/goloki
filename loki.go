package goloki

import (
	"encoding/json"
	"fmt"
	"github.com/civet148/httpc"
	"github.com/civet148/log"
	"strings"
	"sync"
	"time"
)

const (
	lokiPushUri = "/loki/api/v1/push"
)

type LokiValue []string

type LokiClient struct {
	mu      sync.RWMutex
	lokiUrl string
	labels  map[string]string
	hc      *httpc.Client
}

type LokiStream struct {
	Stream map[string]string `json:"stream"` //=> labels map
	Values []LokiValue       `json:"values"` //=> [[timestamp, log]]
}

type LokiBody struct {
	Streams []LokiStream `json:"streams"`
}

func NewLokiClient(strUrl string, labels map[string]string) *LokiClient {
	if strUrl == "" {
		panic("loki server url must not be empty")
	}
	strUrl = strings.TrimSuffix(strUrl, "/")
	lokiUrl := fmt.Sprintf("%s%s", strUrl, lokiPushUri)
	opt := httpc.Option{
		Timeout: 3,
	}
	return &LokiClient{
		lokiUrl: lokiUrl,
		labels:  labels,
		hc:      httpc.NewHttpClient(&opt),
	}
}

func (l *LokiClient) Println(strOutput string, labels ...map[string]string) (err error) {
	var label = make(map[string]string)
	if len(labels) > 0 {
		label = labels[0]
		l.mu.RLock()
		for k, v := range l.labels {
			label[k] = v
		}
		l.mu.RUnlock()
	}
	var timestamp = fmt.Sprintf("%v", time.Now().UnixNano())
	var body = LokiBody{
		Streams: []LokiStream{
			{
				Stream: label,
				Values: []LokiValue{
					{timestamp, strOutput},
				},
			},
		},
	}
	rsp, err := l.hc.PostJson(l.lokiUrl, body)
	if err != nil {
		return log.Errorf("LOKI PUSH ERROR: %s", err.Error())
	}
	if rsp.StatusCode != 200 && rsp.StatusCode != 204 {
		data, _ := json.MarshalIndent(body, "", "\t")
		return log.Errorf("POST url [%s] body [%s] response code [%v] message [%s]", l.lokiUrl, data, rsp.StatusCode, rsp.Body)
	}
	return nil
}
