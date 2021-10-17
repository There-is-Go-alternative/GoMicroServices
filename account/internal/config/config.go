package config

import (
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/utils"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
	"github.com/pkg/errors"
	"os"
)

const (
	HTTPPortEnvVar   = "HTTP_PORT"
	APIKeyEnvVar     = "API_KEY"
	AuthURLEnvVar    = "AUTH_URL"
	BalanceURLEnvVar = "BALANCE_URL"
)

type Config struct {
	Host            string `json:"host,omitempty"`
	HTTPPort        string `json:"port,omitempty"`
	APIKey          string `json:"api_key,omitempty"`
	AccountEndpoint string `json:"account_endpoint,omitempty"`
	AuthURL         string `json:"auth_url,omitempty"`
	BalanceURL      string `json:"balance_url,omitempty"`
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
		fallback: func(c *Config) bool { return c.AuthURL == "" },
		opt: func(c *Config) error {
			if c.AuthURL = os.Getenv(AuthURLEnvVar); c.AuthURL == "" {
				return errors.Wrapf(xerrors.Config, "%v not found in config or env", AuthURLEnvVar)
			}
			return nil
		},
	},
	{
		fallback: func(c *Config) bool { return c.BalanceURL == "" },
		opt: func(c *Config) error {
			if c.BalanceURL = os.Getenv(BalanceURLEnvVar); c.BalanceURL == "" {
				return errors.Wrapf(xerrors.Config, "%v not found in config or env", BalanceURLEnvVar)
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
