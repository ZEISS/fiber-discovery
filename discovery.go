// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package discovery

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

const defaultWellknownURL = ".well-known"

const (
	// FormatJSON returns a JSON response.
	FormatJSON int = iota
	// FormatXML returns an XML response.
	FormatXML
)

// WellknownURL returns the well-known URL for the given keys.
func WellknownURL(keys ...string) string {
	return strings.Join(append([]string{defaultWellknownURL}, keys...), "/")
}

// Config caputes the configuration for running the goth middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	Next func(c fiber.Ctx) bool
	// ResponseFormat defines the response format for the health check.
	ResponseFormat int
}

// DefaultConfig returns the default configuration for the healthz middleware.
var DefaultConfig = Config{
	ResponseFormat: FormatJSON,
}

// New returns a new instance of the goth middleware with the given configuration.
func New(response any, config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	return func(c fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		if c.Method() != fiber.MethodGet && c.Method() != fiber.MethodHead {
			return c.Next()
		}

		switch cfg.ResponseFormat {
		case FormatJSON:
			return c.JSON(response)
		case FormatXML:
			return c.XML(response)
		default: // FormatText
			return c.SendString("")
		}
	}
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return DefaultConfig
	}

	// Override default config
	cfg := config[0]

	if cfg.Next == nil {
		cfg.Next = DefaultConfig.Next
	}

	if cfg.ResponseFormat == 0 {
		cfg.ResponseFormat = DefaultConfig.ResponseFormat
	}

	return cfg
}
