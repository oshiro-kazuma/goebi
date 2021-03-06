package goebi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/kyokomi/goebi/goebi/notice"
)

// TODO: とりあえずgobrake参考
var defaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 3*time.Second)
		},
		ResponseHeaderTimeout: 5 * time.Second,
	},
	Timeout: 10 * time.Second,
}

// Client is an errbit client.
type Client struct {
	client    *http.Client
	noticeURL string
	options   Options
}

// New errbitのClientを生成します
func New(opt Options) *Client {

	c := Client{}
	c.client = defaultHTTPClient
	c.noticeURL = opt.createNoticeBaseURL()
	c.options = opt

	return &c
}

// SendNotice エラー通知します
func (c Client) SendNotice(n notice.Notice) error {

	data, err := json.Marshal(n)
	if err != nil {
		return err
	}

	u := c.options.createNoticeBaseURL()

	res, err := c.client.Post(u, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if !isSuccess(res.StatusCode) {
		data, _ := ioutil.ReadAll(res.Body)

		return fmt.Errorf("error response code %d %s", res.StatusCode, string(data))
	}

	return nil
}

func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
