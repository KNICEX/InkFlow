package bff

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/spf13/viper"
)

func initCloudinary() *cloudinary.Cloudinary {
	type Config struct {
		Key       string `mapstructure:"key"`
		Secret    string `mapstructure:"secret"`
		CloudName string `mapstructure:"cloud_name"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("file.cloudinary", &cfg); err != nil {
		panic(err)
	}
	cld, err := cloudinary.NewFromURL(fmt.Sprintf("cloudinary://%s:%s@%s", cfg.Key, cfg.Secret, cfg.CloudName))
	if err != nil {
		panic(err)
	}
	return cld
}
