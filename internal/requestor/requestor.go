package requestor

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
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
	Client        *http.Client
	Parameters    map[string]string
	RawParameters string
	Cookies       []*http.Cookie
}

func NewRequest(client *http.Client) *Request {
	return &Request{
		Client: client,
	}
}

type Response struct {
	StatusCode int
	Body       string
	Cookies    []*http.Cookie
	Headers    http.Header
}

func NewResponse(statusCode int, body string, cookies []*http.Cookie, headers http.Header) *Response {
	return &Response{
		StatusCode: statusCode,
		Body:       body,
		Cookies:    cookies,
		Headers:    headers,
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
	if len(r.Request.Cookies) > 1 {
		for _, c := range r.Request.Cookies {
			req.AddCookie(c)
		}
	}
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
	r.Response = NewResponse(response.StatusCode, body, response.Cookies(), response.Header)
	return nil
}

func (r *Requestor) POST(url *url.URL) error {
	isRaw := r.Request.RawParameters != ""
	if (len(r.Request.Parameters) == 0) && !isRaw {
		return errors.New("no formdata set")
	}
	switch isRaw {
	case true:
		{
			req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer([]byte(r.Request.RawParameters)))
			if err != nil {
				return err
			}
			if len(r.Request.Cookies) > 1 {
				for _, c := range r.Request.Cookies {
					req.AddCookie(c)
				}
			}
			req.Header.Set("Content-Type", "application/json")
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
			r.Response = NewResponse(response.StatusCode, body, response.Cookies(), response.Header)
			return nil
		}
	case false:
		{
			formData := url.Query()
			for key, value := range r.Request.Parameters {
				formData.Add(key, value)
			}
			req, err := http.NewRequest("POST", url.String(), strings.NewReader(formData.Encode()))
			if err != nil {
				return err
			}
			if len(r.Request.Cookies) > 1 {
				for _, c := range r.Request.Cookies {
					req.AddCookie(c)
				}
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
			r.Response = NewResponse(response.StatusCode, body, response.Cookies(), response.Header)
			return nil
		}
	}
	return nil
}

func (r *Requestor) GetCookie(cookieKey string) string {
	var specificCookie *http.Cookie
	for _, cookie := range r.Response.Cookies {
		if cookie.Name == cookieKey {
			specificCookie = cookie
			break
		}
	}
	return specificCookie.Value
}

func (r *Requestor) GetHeader(headerName string) string {
	return r.Response.Headers.Get(headerName)
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
