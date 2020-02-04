// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dapr/dapr/pkg/runtime"
	log "github.com/sirupsen/logrus"

	// HTTP Middleware
	middleware "github.com/dapr/components-contrib/middleware"
	http_middleware_loader "github.com/dapr/dapr/pkg/components/middleware/http"
	http_middleware "github.com/dapr/dapr/pkg/middleware/http"
	"github.com/valyala/fasthttp"
)

func main() {
	rt, err := runtime.FromFlags()
	if err != nil {
		log.Fatal(err)
	}

	err = rt.Run(
		runtime.WithHTTPMiddleware(
			http_middleware_loader.New("uppercase", func(metadata middleware.Metadata) http_middleware.Middleware {
				return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
					return func(ctx *fasthttp.RequestCtx) {
						body := string(ctx.PostBody())
						ctx.Request.SetBody([]byte(strings.ToUpper(body)))
						h(ctx)
					}
				}
			}),
		),
	)
	if err != nil {
		log.Fatalf("fatal error from runtime: %s", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	<-stop
	gracefulShutdownDuration := 5 * time.Second
	log.Info("dapr shutting down. Waiting 5 seconds to finish outstanding operations")
	rt.Stop()
	<-time.After(gracefulShutdownDuration)
}
