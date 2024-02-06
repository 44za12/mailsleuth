package requestor

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/44za12/mailsleuth/internal/utils"
	"github.com/andybalholm/brotli"
)

type CookieJar struct {
	cookies map[string][]*http.Cookie
}

func NewCookieJar() *CookieJar {
	return &CookieJar{
		cookies: make(map[string][]*http.Cookie),
	}
}

func (jar *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies[u.Host] = cookies
}

func (jar *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func (jar *CookieJar) EmptyCookies() {
	jar.cookies = make(map[string][]*http.Cookie)
}

type Request struct {
	Client     *http.Client
	Parameters map[string]string
}

func NewRequest(client *http.Client) *Request {
	return &Request{
		Client: client,
	}
}

type Response struct {
	StatusCode int
	Body       string
}

func NewResponse(statusCode int, body string) *Response {
	return &Response{
		StatusCode: statusCode,
		Body:       body,
	}
}

type Requestor struct {
	Input           string
	Request         *Request
	Response        *Response
	AddExtraHeaders bool
	Headers         map[string]string
}

func NewRequestor(email string, proxyURL string) (*Requestor, error) {
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		return &Requestor{
			Input: email,
			Request: NewRequest(&http.Client{
				Jar:     NewCookieJar(),
				Timeout: time.Duration(time.Second * 10),
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxy),
				},
			}),
			AddExtraHeaders: true,
			Headers:         utils.StandardHeaders(),
		}, nil
	}
	return &Requestor{
		Input: email,
		Request: NewRequest(&http.Client{
			Jar:     NewCookieJar(),
			Timeout: time.Duration(time.Second * 10),
		}),
		AddExtraHeaders: true,
		Headers:         utils.StandardHeaders(),
	}, nil
}

func (r *Requestor) GET(url *url.URL) error {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}
	response, err := r.Request.Client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := r.decodeResponseBody(response)
	if err != nil {
		return err
	}
	r.Response = NewResponse(response.StatusCode, body)
	return nil
}

func (r *Requestor) POST(url *url.URL) error {
	formData := url.Query()
	for key, value := range r.Request.Parameters {
		formData.Add(key, value)
	}
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.addHeaders(req)
	response, err := r.Request.Client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := r.decodeResponseBody(response)
	if err != nil {
		return err
	}
	r.Response = NewResponse(response.StatusCode, body)
	return nil
}

func (r *Requestor) addHeaders(req *http.Request) {
	if r.AddExtraHeaders {
		for key, value := range r.Headers {
			req.Header.Set(key, value)
		}
	}
}

func (r *Requestor) decodeResponseBody(resp *http.Response) (string, error) {
	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
	case "deflate":
		reader = flate.NewReader(resp.Body)
	case "br":
		reader = brotli.NewReader(resp.Body)
	default:
		reader = resp.Body
	}

	decodedBody, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decodedBody), nil
}
