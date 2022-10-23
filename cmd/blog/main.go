package main

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/dusted-go/config/dotenv"
	"github.com/dusted-go/diagnostic/v3/dlog"
	"github.com/dusted-go/http/v3/middleware/assets"
	"github.com/dusted-go/http/v3/middleware/httptrace"
	"github.com/dusted-go/http/v3/middleware/proxy"
	"github.com/dusted-go/http/v3/middleware/recoverer"
	"github.com/dusted-go/http/v3/middleware/redirect"
	"github.com/dusted-go/http/v3/server"
	"github.com/dustedcodes/blog/cmd/blog/site"
	"github.com/dustedcodes/blog/cmd/blog/web"
)

func createLogProvider(settings *site.Settings) httptrace.CreateLogProviderFunc {
	return func() *dlog.Provider {
		provider := dlog.
			NewProvider().
			SetFilter(assets.LogFilter).
			SetMinLogLevel(settings.MinLogLevel()).
			SetServiceName(settings.ApplicationName).
			SetServiceVersion(settings.ApplicationVersion).
			AddLabel("appName", settings.ApplicationName).
			AddLabel("appVersion", settings.ApplicationVersion)

		if settings.IsProduction() {
			provider.SetFormatter(dlog.NewStackdriverFormatter())
		}
		return provider
	}
}

func main() {
	// ----------------------------------------
	// Bootstrap:
	// ----------------------------------------

	// Init dotenv
	err := dotenv.Load(".env", true)
	if err != nil {
		panic(err)
	}

	// Init settings
	settings := site.InitSettings()

	// Init context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init logger
	defaultLogProvider := createLogProvider(settings)
	dlog.Context(ctx, defaultLogProvider())

	// Init public assets
	assetMiddleware, err :=
		assets.NewMiddleware(
			ctx,
			"dist/assets/",
			"public, max-age=15552000",
			!settings.IsProduction())
	if err != nil {
		panic(err)
	}

	siteAssets := &site.Assets{
		CSSPath: assetMiddleware.CSS.VirtualFileName,
		JSPath:  assetMiddleware.JS.VirtualFileName,
	}

	// Init all blog posts
	blogPosts, err := site.ReadBlogPosts(ctx, "dist/posts")
	if err != nil {
		panic(err)
	}

	// Sort blog posts by date (newest first)
	sort.Slice(blogPosts, func(i, j int) bool {
		return blogPosts[i].PublishDate.After(blogPosts[j].PublishDate)
	})

	// Init web handler
	webHandler := web.NewHandler(
		settings,
		siteAssets,
		blogPosts)

	// ----------------------------------------
	// Web Server:
	// ----------------------------------------

	middleware := server.CombineMiddlewares(
		recoverer.HandlePanics(webHandler.Recover),
		proxy.ForwardedHeaders(settings.ProxyCount),
		httptrace.GoogleCloudTrace(createLogProvider(settings)),
		redirect.TrailingSlash(),
		redirect.Hosts(settings.DomainRedirects(), true),
		redirect.ForceHTTPS(
			settings.IsProduction(),
			settings.PublicHosts()...),
		assetMiddleware,
	)
	webApp := middleware.Next(webHandler)

	dlog.New(ctx).Notice().Fmt("Starting web server on %s", settings.ServerAddress())

	httpServer := &http.Server{
		Addr:              settings.ServerAddress(),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		// Google Cloud LoadBalancer has a keepalive timeout of 600s
		// and recommends to set the backend server to have one of 620s
		// https://cloud.google.com/load-balancing/docs/https#timeouts_and_retries
		IdleTimeout:    620 * time.Second,
		Handler:        webApp,
		MaxHeaderBytes: int(settings.MaxRequestSize),
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
