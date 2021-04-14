// Copyright (c) 2020 Ketch Kloud, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package orlop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.ketch.com/lib/orlop/errors"
	"io"
	"net/http"
	"strings"
)

// HttpClient provides a wrapper around http.Client, providing automatic TLS setup and header management
type HttpClient struct {
	cfg   HasClientConfig
	cli   *http.Client
	token []byte
}

// NewHttpClient creates a new HttpClient
func NewHttpClient(ctx context.Context, cfg HasClientConfig, vault HasVaultConfig) (*HttpClient, error) {
	ct, err := NewClientTLSConfigContext(ctx, cfg.GetTLS(), vault)
	if err != nil {
		return nil, err
	}

	httpCli := &HttpClient{
		cfg: cfg,
		cli: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: ct,
				IdleConnTimeout: cfg.GetConnTimeout(),
				WriteBufferSize: cfg.GetWriteBufferSize(),
				ReadBufferSize:  cfg.GetReadBufferSize(),
			},
		},
	}

	if cfg.GetToken() != nil && cfg.GetToken().GetShared() != nil {
		if httpCli.token, err = LoadKeyContext(ctx, cfg.GetToken().GetShared(), vault, "shared"); err != nil {
			return nil, err
		}
	}

	return httpCli, nil
}

// Do executes the request
func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	for k, v := range c.cfg.GetHeaders() {
		req.Header.Add(k, v)
	}

	req.Header.Set("User-Agent", c.cfg.GetUserAgent())
	if c.token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(c.token)))
	}

	return c.cli.Do(req)
}

// Get performs a GET against the given relative url
func (c *HttpClient) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.join(c.cfg.GetURL(), url), nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// GetJSON performs a GET against the given relative url and returns the results unmarshalled from JSON
func (c *HttpClient) GetJSON(ctx context.Context, url string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.join(c.cfg.GetURL(), url), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = c.handleError(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

// Head performs a HEAD against the given relative url
func (c *HttpClient) Head(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, c.join(c.cfg.GetURL(), url), nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Post performs a POST against the given relative url
func (c *HttpClient) Post(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.join(c.cfg.GetURL(), url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

// PostJSON performs a POST against the given relative url using the JSON body and returns JSON
func (c *HttpClient) PostJSON(ctx context.Context, url string, in interface{}, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.join(c.cfg.GetURL(), url), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = c.handleError(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

// Put performs a PUT against the given relative url
func (c *HttpClient) Put(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.join(c.cfg.GetURL(), url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

// PutJSON performs a PUT against the given relative url using the JSON body and returns JSON
func (c *HttpClient) PutJSON(ctx context.Context, url string, in interface{}, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.join(c.cfg.GetURL(), url), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = c.handleError(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

// Patch performs a PATCH against the given relative url
func (c *HttpClient) Patch(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.join(c.cfg.GetURL(), url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

// PatchJSON performs a PATCH against the given relative url using the JSON body and returns JSON
func (c *HttpClient) PatchJSON(ctx context.Context, url string, in interface{}, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.join(c.cfg.GetURL(), url), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = c.handleError(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

// Delete performs a DELETE against the given relative url
func (c *HttpClient) Delete(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.join(c.cfg.GetURL(), url), nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *HttpClient) handleError(resp *http.Response) error {
	if resp.StatusCode < 300 {
		return nil
	}

	var v map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return errors.New(resp.Status)
	}

	if errMsg, ok := v["error"]; ok {
		return errors.Errorf("%v", errMsg)
	}

	return errors.New(resp.Status)
}

func (c *HttpClient) join(root, suffix string) string {
	root = strings.TrimSuffix(root, "/")
	suffix = strings.TrimPrefix(suffix, "/")
	return root + "/" + suffix
}
