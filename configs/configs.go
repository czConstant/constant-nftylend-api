package configs

import (
	"encoding/json"
	"os"

	"github.com/czConstant/blockchain-api/bcclient"
	"github.com/czConstant/constant-evn/client"
)

var config *Config

func init() {
	file, err := os.Open("configs/config.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	v := Config{}
	err = decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	config = &v
	evnClient := &client.Client{
		URL: config.EvnURL,
	}

	dbURL, _, err := evnClient.GetSecret("DB-DEMO-NFTLEND-URL")
	if err != nil {
		panic(err)
	}
	config.DbURL = client.ParseDBURL(dbURL)
}

func GetConfig() *Config {
	return config
}

type Config struct {
	Env               string `json:"env"`
	EvnURL            string `json:"evn_url"`
	RavenDNS          string `json:"raven_dns"`
	RavenENV          string `json:"raven_env"`
	Port              int    `json:"port"`
	LogPath           string `json:"log_path"`
	DbURL             string `json:"db_url"`
	Debug             bool   `json:"debug"`
	RecaptchaV3Serect string `json:"recaptcha_v3_serect"`
	JobToken          string `json:"job_token"`
	Datadog           struct {
		Env     string `json:"env"`
		Service string `json:"service"`
		Version string `json:"version"`
	} `json:"datadog"`
	Blockchain bcclient.Config `json:"blockchain"`
}
