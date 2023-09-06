package main

import (
	"context"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/dusted-go/config/dotenv"
	"github.com/dusted-go/http/v6/middleware/assets"
	"github.com/dusted-go/http/v6/middleware/firewall"
	"github.com/dusted-go/http/v6/middleware/headers"
	"github.com/dusted-go/http/v6/middleware/healthz"
	"github.com/dusted-go/http/v6/middleware/mware"
	"github.com/dusted-go/http/v6/middleware/proxy"
	"github.com/dusted-go/http/v6/middleware/recoverer"
	"github.com/dusted-go/http/v6/middleware/redirect"
	"github.com/dusted-go/logging/prettylog"
	"github.com/dusted-go/logging/stackdriver"

	"github.com/dustedcodes/blog/cmd/blog/model"
	"github.com/dustedcodes/blog/cmd/blog/web"
	"github.com/dustedcodes/blog/internal/blog"
	"github.com/dustedcodes/blog/internal/cloudtrace"
	"github.com/dustedcodes/blog/internal/config"
)

func main() {
	// -----------------------------
	// Load config
	// -----------------------------
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := dotenv.Load(".env", true)
	if err != nil {
		panic(err)
	}
	config := config.Load()

	// -----------------------------
	// Init default logger
	// -----------------------------
	var logHandler slog.Handler
	var loggingMiddleware func(http.Handler) http.Handler
	if config.IsProduction() {
		logHandlerOptions := &stackdriver.HandlerOptions{
			ServiceName:    config.ApplicationName,
			ServiceVersion: config.ApplicationVersion,
			MinLevel:       config.MinLogLevel(),
			AddSource:      config.IsProduction(),
		}
		logMiddlewareOptions := &stackdriver.MiddlewareOptions{
			GCPProjectID:   config.GoogleCloudProjectID,
			AddTrace:       config.IsProduction(),
			AddHTTPRequest: config.IsProduction(),
		}
		logHandler = stackdriver.NewHandler(logHandlerOptions)
		loggingMiddleware = stackdriver.Logging(logHandlerOptions, logMiddlewareOptions)
	} else {
		logHandler = prettylog.NewHandler(&slog.HandlerOptions{
			Level:       config.MinLogLevel(),
			AddSource:   config.IsProduction(),
			ReplaceAttr: stackdriver.ReplaceLogLevel,
		})
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// ----------------------------------------
	// Bootstrap:
	// ----------------------------------------
	assetMiddleware, err :=
		assets.NewMiddleware(
			"dist/assets/",
			"public, max-age=15552000",
			!config.IsProduction(),
			false)
	if err != nil {
		panic(err)
	}
	siteAssets := &model.Assets{
		CSSPath: assetMiddleware.CSS.VirtualFileName,
		JSPath:  assetMiddleware.JS.VirtualFileName,
	}
	blogPosts, err := blog.ReadPosts(ctx, blog.DefaultBlogPostPath)
	if err != nil {
		panic(err)
	}
	// Sort blog posts by date (newest first)
	sort.Slice(blogPosts, func(i, j int) bool {
		return blogPosts[i].PublishDate.After(blogPosts[j].PublishDate)
	})
	webHandler := web.NewHandler(
		config,
		siteAssets,
		blogPosts)

	// ----------------------------------------
	// Web Server:
	// ----------------------------------------
	middleware := mware.Bind(
		recoverer.HandlePanics(webHandler.Recover),
		healthz.LivenessProbe,
		cloudtrace.Middleware,
		loggingMiddleware,
		proxy.ForwardedHeaders(config.ProxyCount),
		redirect.TrailingSlash,
		redirect.Hosts(config.DomainRedirects(), true),
		redirect.ForceHTTPS(
			config.IsProduction(),
			config.PublicHosts()...),
		assetMiddleware.ServeFiles,
		firewall.LimitRequestSize(config.MaxRequestSize),
		headers.Security(60*60*24*30),
	)
	webApp := middleware(webHandler)

	// -----------------------------
	// Launch web server
	// -----------------------------
	logger.Info("Starting server...",
		"web-server-address", config.ServerAddress(),
		"base-url", config.BaseURL)
	httpServer := &http.Server{
		Addr:              config.ServerAddress(),
		ReadHeaderTimeout: 3 * time.Second,
		// Google Cloud LoadBalancer has a keepalive timeout of 600s
		// and recommends to set the backend server to have one of 620s
		// https://cloud.google.com/load-balancing/docs/https#timeouts_and_retries
		IdleTimeout:    620 * time.Second,
		Handler:        webApp,
		MaxHeaderBytes: int(config.MaxRequestSize),
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
