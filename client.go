package httpclient

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	ProxyScheme string        // 代理方式:http,https
	ProxyIp     string        // 代理IP
	DialTimeout time.Duration // 拨号超时
	ReadTimeout time.Duration // 读取超时
	client      *http.Client
}

func (client *HttpClient) NewClient() (*HttpClient, error) {
	urli := url.URL{}
	urlproxy, err := urli.Parse(fmt.Sprintf("%s://%s", client.ProxyScheme, client.ProxyIp))
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{}
	transport.DisableKeepAlives = true
	transport.IdleConnTimeout = time.Duration(10) * time.Second
	transport.Proxy = http.ProxyURL(urlproxy)
	transport.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
		c, err := net.DialTimeout(network, addr, client.DialTimeout)
		if err != nil {
			// 连接不上代理服务器
			return nil, err
		}
		c.SetDeadline(time.Now().Add(client.ReadTimeout))
		return c, nil
	}
	client.client = &http.Client{
		Transport: transport,
	}
	return client, err
}

func (client *HttpClient) Get(url string) (*http.Response, error) {
	return client.client.Get(url)
}
