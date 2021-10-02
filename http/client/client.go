// Copyright (c) 2021 Ketch Kloud, Inc.
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

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/secret"
	"go.ketch.com/lib/orlop/tls"
	"go.uber.org/fx"
	"io"
	"net/http"
	"strings"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Get(ctx context.Context, url string) (resp *http.Response, err error)
	GetJSON(ctx context.Context, url string, out interface{}) error
	Head(ctx context.Context, url string) (resp *http.Response, err error)
	Post(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error)
	PostJSON(ctx context.Context, url string, in interface{}, out interface{}) error
	Put(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error)
	PutJSON(ctx context.Context, url string, in interface{}, out interface{}) error
	Patch(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error)
	PatchJSON(ctx context.Context, url string, in interface{}, out interface{}) error
	Delete(ctx context.Context, url string) (resp *http.Response, err error)
}

type Provider interface {
	New(ctx context.Context, cfg Config) (Client, error)
}

type Params struct {
	fx.In

	Config         Config
	TLSProvider    tls.ClientProvider
	SecretProvider secret.Provider
}

func New(params Params) Provider {
	return &providerImpl{
		params: params,
	}
}

type providerImpl struct {
	params Params
}

// New creates a new Client
func (p providerImpl) New(ctx context.Context, cfg Config) (Client, error) {
	ct, err := p.params.TLSProvider.NewConfig(ctx, cfg.TLS)
	if err != nil {
		return nil, err
	}

	httpCli := &clientImpl{
		cfg: cfg,
		cli: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: ct,
				IdleConnTimeout: cfg.ConnTimeout,
				WriteBufferSize: cfg.WriteBufferSize,
				ReadBufferSize:  cfg.ReadBufferSize,
			},
		},
	}

	if httpCli.token, err = p.params.SecretProvider.Load(ctx, cfg.Token); err != nil {
		return nil, err
	}

	return httpCli, nil
}

// clientImpl provides a wrapper around http.Client, providing automatic TLS setup and header management
type clientImpl struct {
	cfg   Config
	cli   *http.Client
	token string
}

// Do executes the request
func (c *clientImpl) Do(req *http.Request) (*http.Response, error) {
	for k, v := range c.cfg.Headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("User-Agent", c.cfg.GetUserAgent())
	if len(c.token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	return c.cli.Do(req)
}

// Get performs a GET against the given relative url
func (c *clientImpl) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.join(c.cfg.GetURL(), url), nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// GetJSON performs a GET against the given relative url and returns the results unmarshalled from JSON
func (c *clientImpl) GetJSON(ctx context.Context, url string, out interface{}) error {
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
func (c *clientImpl) Head(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, c.join(c.cfg.GetURL(), url), nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Post performs a POST against the given relative url
func (c *clientImpl) Post(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.join(c.cfg.GetURL(), url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

// PostJSON performs a POST against the given relative url using the JSON body and returns JSON
func (c *clientImpl) PostJSON(ctx context.Context, url string, in interface{}, out interface{}) error {
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
func (c *clientImpl) Put(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.join(c.cfg.GetURL(), url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

// PutJSON performs a PUT against the given relative url using the JSON body and returns JSON
func (c *clientImpl) PutJSON(ctx context.Context, url string, in interface{}, out interface{}) error {
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
func (c *clientImpl) Patch(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.join(c.cfg.GetURL(), url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

// PatchJSON performs a PATCH against the given relative url using the JSON body and returns JSON
func (c *clientImpl) PatchJSON(ctx context.Context, url string, in interface{}, out interface{}) error {
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
func (c *clientImpl) Delete(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.join(c.cfg.GetURL(), url), nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *clientImpl) handleError(resp *http.Response) error {
	if resp.StatusCode < 300 {
		return nil
	}

	var v map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return errors.New(resp.Status)
	}

	if errMsg, ok := v["error"]; ok {
		return errors.Errorf("%v", errMsg)
	}

	return errors.New(resp.Status)
}

func (c *clientImpl) join(root, suffix string) string {
	root = strings.TrimSuffix(root, "/")
	suffix = strings.TrimPrefix(suffix, "/")
	return root + "/" + suffix
}
