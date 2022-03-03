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

	dbURL, _, err := evnClient.GetSecret("DB-NFTLEND-URL")
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
	Mailer struct {
		URL string `json:"url"`
	} `json:"mailer"`
	Backend struct {
		URL string `json:"url"`
	} `json:"backend"`
	Core struct {
		URL string `json:"url"`
	} `json:"core"`
	Blockchain bcclient.Config `json:"blockchain"`
	Contract   struct {
		Ronin struct {
			AdminAddress           string `json:"admin_address"`
			AxieAddress            string `json:"axie_address"`
			AxieMarketplaceAddress string `json:"axie_marketplace_address"`
			AxieEtherTokenAddress  string `json:"axie_ether_token_address"`
		} `json:"ronin"`
	} `json:"contract"`
	AxieMarketplace struct {
		URL string `json:"url"`
	} `json:"axie_marketplace"`
	AxieGameapi struct {
		URL string `json:"url"`
	} `json:"axie_gameapi"`
	AxieTracker struct {
		URL  string `json:"url"`
		Host string `json:"host"`
		Key  string `json:"key"`
	} `json:"axie_tracker"`

	Liquidation struct {
		URL string `json:"url"`
	} `json:"liquidation"`

	Chat struct {
		URL string `json:"url"`
	} `json:"chat"`
}
