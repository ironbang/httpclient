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
	KeepAlives	bool
	client      *http.Client
}

func (client *HttpClient) NewClient() (*HttpClient, error) {
	transport := &http.Transport{}

	if len(client.ProxyIp) > 0 {
		urli := url.URL{}
		urlproxy, err := urli.Parse(fmt.Sprintf("%s://%s", client.ProxyScheme, client.ProxyIp))
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(urlproxy)
	}

	transport.DisableKeepAlives = client.KeepAlives
	transport.IdleConnTimeout = (client.DialTimeout+client.ReadTimeout) * 2 * time.Second
	if client.DialTimeout > 0 && client.ReadTimeout > 0 {
		transport.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			c, err := net.DialTimeout(network, addr, client.DialTimeout)
			if err != nil {
				// 连接不上代理服务器
				return nil, err
			}
			c.SetDeadline(time.Now().Add(client.ReadTimeout))
			return c, nil
		}
	}
	client.client = &http.Client{
		Transport: transport,
	}
	return client, nil
}

func (client *HttpClient) Get(url string,headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil{
		return nil,err
	}
	for key,val := range headers{
		req.Header.Add(key,val)
	}
	return client.client.Do(req)
}
