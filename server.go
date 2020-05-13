// Copyright (c) 2020 SwitchBit, Inc.
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
	"container/heap"
	"context"
	"crypto/tls"
	"github.com/switch-bit/orlop/log"
	syslog "log"
	"net"
	"net/http"
)

type HandlerHeap []*handlerPair

func (h *HandlerHeap) Handle(pattern string, handler http.Handler) {
	heap.Push(h, &handlerPair{
		pattern: pattern,
		handler: handler,
	})
}
func (h HandlerHeap) Len() int           { return len(h) }
func (h HandlerHeap) Less(i, j int) bool { return len(h[i].pattern) > len(h[j].pattern) }
func (h HandlerHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *HandlerHeap) Push(x interface{}) {
	p := x.(*handlerPair)
	*h = append(*h, &handlerPair{
		pattern: p.pattern,
		handler: p.handler,
	})
}

func (h *HandlerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Serve sets up the server and listens for requests
func Serve(ctx context.Context, serviceName string, options ...ServerOption) error {
	var err error

	// Setup the server options
	serverOptions := &serverOptions{
		serviceName: serviceName,
		tlsProvider: NewSimpleTLSProvider(),
		log:         log.WithField("service", serviceName),
	}

	// Set default server config
	err = WithServerConfig(ServerConfig{
		Bind:   "0.0.0.0",
		Listen: 5000,
		TLS:    TLSConfig{},
	}).apply(ctx, serverOptions)
	if err != nil {
		return err
	}

	// Add default health check
	err = WithHealthCheck(nil).apply(ctx, serverOptions)
	if err != nil {
		return err
	}

	// Add default metrics handler
	err = WithPrometheusMetrics().apply(ctx, serverOptions)
	if err != nil {
		return err
	}

	// Process all server options (which may override any of the above)
	for _, option := range options {
		err = option.apply(ctx, serverOptions)
		if err != nil {
			return err
		}
	}

	handlers := &HandlerHeap{}
	heap.Init(handlers)

	for _, option := range options {
		err = option.addHandler(ctx, serverOptions, handlers)
		if err != nil {
			return err
		}
	}

	// Create the HTTP server
	mux := http.NewServeMux()

	for handlers.Len() > 0 {
		pair := heap.Pop(handlers).(*handlerPair)
		mux.Handle(pair.pattern, pair.handler)
	}

	w := log.Writer()
	defer w.Close()

	// Start listening
	serverOptions.log.Info("listening")
	ln, err := net.Listen("tcp", serverOptions.addr)
	if err != nil {
		return err
	}

	// Serve requests
	if serverOptions.config.GetTLS().GetEnabled() {
		serverOptions.log.Trace("loading server tls certs")
		config, err := serverOptions.tlsProvider.NewServerTLSConfig(serverOptions.config.GetTLS())
		if err != nil {
			ln.Close()

			return err
		}

		ln = tls.NewListener(ln, config)
	}

	defer ln.Close()

	serverOptions.log.Info("serving")
	srv := &http.Server{
		Addr:     serverOptions.addr,
		Handler:  mux,
		ErrorLog: syslog.New(w, "[http]", 0),
	}

	return srv.Serve(ln)
}
