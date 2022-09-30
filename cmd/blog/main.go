package main

import (
	"context"
	"net/http"
	"time"

	"github.com/dusted-go/config/dotenv"
	"github.com/dusted-go/diagnostic/v2/log"
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
	return func() *log.Provider {
		provider := log.
			NewProvider().
			SetFilter(assets.LogFilter).
			SetMinLogLevel(settings.MinLogLevel()).
			SetServiceName(settings.ApplicationName).
			SetServiceVersion(settings.ApplicationVersion).
			AddLabel("appName", settings.ApplicationName).
			AddLabel("appVersion", settings.ApplicationVersion)

		if settings.IsProduction() {
			provider.SetFormatter(log.NewStackdriverFormatter())
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

	// Init logger
	ctx := context.Background()
	defaultLogProvider := createLogProvider(settings)
	log.Context(ctx, defaultLogProvider())

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

	// Init web handler
	webHandler := web.NewHandler(settings, siteAssets)

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

	log.New(ctx).Notice().Fmt("Starting web server on %s", settings.ServerAddress())

	httpServer := &http.Server{
		Addr:              settings.ServerAddress(),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           webApp,
		MaxHeaderBytes:    int(settings.MaxRequestSize),
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
