package configs

import "github.com/spf13/viper"

type conf struct {
	BrasilApiUrl string `mastructure:"BRASIL_API_URL"`
	ViaCepUrl    string `mastructure:"VIA_CEP_URL"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	brasilApiUrl, _ := viper.Get("BRASIL_API_URL").(string)
	viaCepUrl, _ := viper.Get("VIA_CEP_URL").(string)
	cfg.BrasilApiUrl = brasilApiUrl
	cfg.ViaCepUrl = viaCepUrl
	return cfg, err
}
