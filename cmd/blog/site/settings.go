package site

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dusted-go/config/env"
	"github.com/dusted-go/diagnostic/v2/log"
	"github.com/dusted-go/http/v3/request"
)

type Assets struct {
	CSSPath string
	JSPath  string
}

type Settings struct {
	EnvironmentName    string
	LogLevel           log.Level
	ApplicationName    string
	ApplicationVersion string
	HTTPPort           int
	ProxyCount         int
	PublicHost         string
	BaseURL            string
	RedirectWWW        bool
	CDN                string
	MaxRequestSize     int64
}

func (s *Settings) MinLogLevel() log.Level {
	return s.LogLevel
}

func (s *Settings) IsProduction() bool {
	return strings.ToLower(s.EnvironmentName) == "production"
}

func (s *Settings) HotReload() bool {
	return !s.IsProduction()
}

func (s *Settings) ServerAddress() string {
	return ":" + strconv.Itoa(s.HTTPPort)
}

func (s *Settings) PublicHosts() []string {
	return []string{s.PublicHost}
}

func (s *Settings) DomainRedirects() map[string]string {
	redirects := map[string]string{}
	if s.RedirectWWW {
		for _, dest := range s.PublicHosts() {
			from := fmt.Sprintf("www.%s", dest)
			redirects[from] = dest
		}
	}
	return redirects
}

func (s *Settings) URLs(r *http.Request) *URLs {
	return &URLs{
		RequestURL: request.FullURL(r),
		BaseURL:    s.BaseURL,
		CDN:        s.CDN,
	}
}

func InitSettings() *Settings {
	return &Settings{
		EnvironmentName:    env.GetOrDefault("ENV_NAME", "Development"),
		LogLevel:           log.ParseLevel(env.GetOrDefault("LOG_LEVEL", "Debug")),
		ApplicationName:    env.GetOrDefault("APP_NAME", "dustedcodes"),
		ApplicationVersion: env.GetOrDefault("APP_VERSION", "0.1.0"),
		HTTPPort:           env.GetIntOrDefault("HTTP_PORT", 3000),
		ProxyCount:         env.GetIntOrDefault("PROXY_COUNT", 0),
		PublicHost:         env.GetOrDefault("PUBLIC_HOST", "dusted.codes"),
		BaseURL:            env.GetOrDefault("BASE_URL", "https://dusted.codes"),
		RedirectWWW:        env.GetBoolOrDefault("REDIRECT_WWW", false),
		CDN:                env.GetOrDefault("CDN", "https://cdn.dusted.codes"),
		MaxRequestSize:     int64(env.GetIntOrDefault("MAX_REQUEST_SIZE", 500000)),
	}
}
