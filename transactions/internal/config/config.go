package config

import (
	"os"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/internal/utils"
	"github.com/There-is-Go-alternative/GoMicroServices/transactions/internal/xerrors"
	"github.com/pkg/errors"
)

const (
	HTTPPortEnvVar   = "HTTP_PORT"
	APIKeyEnvVar     = "API_KEY"
	AccountURLEnvVar = "AUTH_URL"
	AdsURLEnvVar     = "ADS_URL"
	FundsURLEnvVar   = "FUNDS_URL"
)

type Config struct {
	Host       string `json:"host,omitempty"`
	HTTPPort   string `json:"port,omitempty"`
	APIKey     string `json:"api_key,omitempty"`
	AccountURL string `json:"account_url,omitempty"`
	FundsURL   string `json:"funds_url,omitempty"`
	AdsURL     string `json:"ads_url,omitempty"`
}

// Option is a Config option that can be given to the constructor
type Option func(*Config) error

func (c *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}
	return nil
}

// Complete list of default options and when to fallback on them.
var defaults = []struct {
	fallback func(*Config) bool
	opt      Option
}{
	{
		fallback: func(c *Config) bool { return c.HTTPPort == "" },
		opt: func(c *Config) error {
			if c.HTTPPort = os.Getenv(HTTPPortEnvVar); c.HTTPPort == "" {
				return errors.Wrapf(xerrors.Config, "%v not found in config or env", HTTPPortEnvVar)
			}
			return nil
		},
	},
	{
		fallback: func(c *Config) bool { return c.APIKey == "" },
		opt: func(c *Config) error {
			if c.APIKey = os.Getenv(APIKeyEnvVar); c.APIKey == "" {
				return errors.Wrapf(xerrors.Config, "%v not found in config or env", APIKeyEnvVar)
			}
			return nil
		},
	},
	{
		fallback: func(c *Config) bool { return c.AccountURL == "" },
		opt: func(c *Config) error {
			if c.AccountURL = os.Getenv(c.AccountURL); c.AccountURL == "" {
				return errors.Wrapf(xerrors.Config, "%v not found in config or env", AccountURLEnvVar)
			}
			return nil
		},
	},
	{
		fallback: func(c *Config) bool { return c.FundsURL == "" },
		opt: func(c *Config) error {
			if c.FundsURL = os.Getenv(FundsURLEnvVar); c.FundsURL == "" {
				return errors.Wrapf(xerrors.Config, "%v not found in config or env", FundsURLEnvVar)
			}
			return nil
		},
	},
	{
		fallback: func(c *Config) bool { return c.AdsURL == "" },
		opt: func(c *Config) error {
			if c.AdsURL = os.Getenv(AdsURLEnvVar); c.AdsURL == "" {
				return errors.Wrapf(xerrors.Config, "%v not found in config or env", AdsURLEnvVar)
			}
			return nil
		},
	},
}

func applyDefaults(c *Config) error {
	for _, def := range defaults {
		if !def.fallback(c) {
			continue
		}
		if err := def.opt(c); err != nil {
			return err
		}
	}
	return nil
}

func NewConfig(path string) (*Config, error) {
	var conf Config
	err := utils.DecodeJSONFromFile(path, &conf)
	if err != nil && !errors.Is(err, &os.PathError{}) {
		return nil, err
	}
	if err = conf.Apply(applyDefaults); err != nil {
		return nil, err
	}
	return &conf, nil
}

func ParseConfigFromPath(path string) (conf *Config, err error) {
	err = utils.DecodeJSONFromFile(path, &conf)
	if err != nil {
		return nil, err
	}
	if conf.HTTPPort == "" {
		os.Getenv("CONF_PATH")
	}
	return
}
