package config

import (
	merlinsdk "github.com/openmerlin/merlin-sdk/httpclient"
	"os"

	"github.com/openmerlin/cronjob/utils"
)

func LoadConfig(path string, remove bool) (cfg Config, err error) {
	if remove {
		defer os.Remove(path)
	}

	if err = utils.LoadFromYaml(path, &cfg); err != nil {
		return
	}

	SetDefault(&cfg)

	err = Validate(&cfg)

	return
}

type DownloadConfig struct {
	Spec            string `json:"spec"     required:"true"`
	OriginalDataUrl string `json:"original_data_url" required:"true"`
}

type Config struct {
	Merlin        merlinsdk.Config `json:"merlin"`
	DownloadCount DownloadConfig   `json:"download_count"`
}

type MerlinConfig struct {
	Token string `json:"token"     required:"true"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.DownloadCount,
	}
}

// SetDefault sets default values for the Config struct.
func (cfg *Config) SetDefault() {

}

// Validate validates the configuration.
func (cfg *Config) Validate() error {
	return utils.CheckConfig(cfg, "")
}
