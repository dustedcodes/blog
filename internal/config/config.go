package config

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/dusted-go/config/env"
)

type Config struct {
	EnvironmentName    string
	LogLevel           slog.Leveler
	ApplicationName    string
	ApplicationVersion string
	HTTPPort           int
	ProxyCount         int
	PublicHost         string
	BaseURL            string
	RedirectWWW        bool
	CDN                string
	MaxRequestSize     int64
	DisqusShortname    string
}

func parseLogLevel(value string) slog.Leveler {
	if len(value) == 0 {
		return slog.LevelInfo
	}

	level := strings.ToLower(strings.Trim(value, " "))

	if level == "debug" {
		return slog.LevelDebug
	}
	if level == "info" {
		return slog.LevelInfo
	}
	if level == "warning" {
		return slog.LevelWarn
	}
	if level == "error" {
		return slog.LevelError
	}

	i, err := strconv.Atoi(value)
	if err == nil {
		return slog.Level(i)
	}

	return slog.LevelInfo
}

func (c *Config) MinLogLevel() slog.Leveler {
	return c.LogLevel
}

func (c *Config) IsProduction() bool {
	return strings.ToLower(c.EnvironmentName) == "production"
}

func (c *Config) HotReload() bool {
	return !c.IsProduction()
}

func (c *Config) ServerAddress() string {
	return ":" + strconv.Itoa(c.HTTPPort)
}

func (c *Config) PublicHosts() []string {
	return []string{c.PublicHost}
}

func (c *Config) DomainRedirects() map[string]string {
	redirects := map[string]string{}

	if c.RedirectWWW {
		for _, dest := range c.PublicHosts() {
			from := "www." + dest
			redirects[from] = dest
		}
	}

	return redirects
}

func Load() *Config {
	return &Config{
		EnvironmentName:    env.GetOrDefault("ENV_NAME", "Development"),
		LogLevel:           parseLogLevel(env.GetOrDefault("LOG_LEVEL", "Debug")),
		ApplicationName:    env.GetOrDefault("APP_NAME", "dustedcodes"),
		ApplicationVersion: env.GetOrDefault("APP_VERSION", "0.1.0"),
		HTTPPort:           env.GetIntOrDefault("HTTP_PORT", 3000),
		ProxyCount:         env.GetIntOrDefault("PROXY_COUNT", 0),
		PublicHost:         env.GetOrDefault("PUBLIC_HOST", "dusted.codes"),
		BaseURL:            env.GetOrDefault("BASE_URL", "https://dusted.codes"),
		RedirectWWW:        env.GetBoolOrDefault("REDIRECT_WWW", false),
		CDN:                env.GetOrDefault("CDN", "https://cdn.dusted.codes"),
		MaxRequestSize:     int64(env.GetIntOrDefault("MAX_REQUEST_SIZE", 500000)),
		DisqusShortname:    env.GetOrDefault("DISQUS_SHORTNAME", ""),
	}
}
