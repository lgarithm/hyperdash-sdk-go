package hyperdash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	_endpoint = `https://hyperdash.io/api/v1/sdk/http`
)

var (
	endpoint, _ = url.Parse(_endpoint)
)

type Client struct {
	endpoint url.URL
	config   Config
	client   *http.Client
	stop     chan struct{}
	uuid     string
	jobName  string
}

func NewClient(name string) (*Client, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return &Client{
		endpoint: *endpoint,
		config:   *config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		stop:    make(chan struct{}, 1),
		uuid:    genUUID(),
		jobName: name,
	}, nil
}

func (c *Client) Start() error {
	payload := &runStarted{
		JobName: c.jobName,
	}
	msg := c.newSDKMessage(kTypeRunStarted, payload)
	go func() {
		tk := time.NewTicker(10 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-tk.C:
				c.HeartBeat()
			case <-c.stop:
				return
			}
		}
	}()
	return c.send(msg)
}

func (c *Client) HeartBeat() error {
	msg := c.newSDKMessage(kTypeHeartbeat, nil)
	return c.send(msg)
}

func (c *Client) Log(s string) error {
	payload := &userLog{
		UUID:  genUUID(),
		Level: kLevelInfo,
		Body:  s,
	}
	msg := c.newSDKMessage(kTypeLog, payload)
	return c.send(msg)
}

func (c *Client) Metric(name string, value interface{}) error {
	payload := &metric{
		Name:      name,
		Timestamp: time.Now().Unix() * 1000,
		Value:     value,
	}
	msg := c.newSDKMessage(kTypeMetric, payload)
	return c.send(msg)
}

func (c *Client) Param(params map[string]interface{}) error {
	payload := &param{
		Params: params,
	}
	msg := c.newSDKMessage(kTypeParam, payload)
	return c.send(msg)
}

func (c *Client) Close() error {
	c.stop <- struct{}{}
	payload := &runEnded{
		FinalStatus: kOutcomeSuccess,
	}
	msg := c.newSDKMessage(kTypeRunEnded, payload)
	return c.send(msg)
}

func (c *Client) send(msg interface{}) error {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	if err := e.Encode(msg); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.endpoint.String(), &b)
	if err != nil {
		return err
	}
	req.Header.Set(`x-hyperdash-auth`, c.config.APIKey)
	req.Header.Set(`Content-Type`, `application/json`)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	fmt.Fprintf(os.Stdout, "\n")
	return nil
}

func (c *Client) newSDKMessage(t MsgType, payload interface{}) *SDKMessage {
	return &SDKMessage{
		Type:       t,
		Timestamp:  time.Now().Unix() * 1000,
		SDKRunUUID: c.uuid,
		Payload:    payload,
	}
}
