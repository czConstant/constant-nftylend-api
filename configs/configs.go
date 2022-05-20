package configs

import (
	"encoding/json"
	"os"

	"github.com/czConstant/blockchain-api/bcclient"
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
	WebUrl            string `json:"web_url"`
	Datadog           struct {
		Env     string `json:"env"`
		Service string `json:"service"`
		Version string `json:"version"`
	} `json:"datadog"`
	Moralis struct {
		APIKey string `json:"api_key"`
	} `json:"moralis"`
	SaleTrack struct {
		NftbankKey string `json:"nftbank_key"`
	} `json:"sale_track"`
	Contract struct {
		ProgramID            string `json:"program_id"`
		MaticNftypawnAddress string `json:"matic_nftypawn_address"`
		AvaxNftypawnAddress  string `json:"avax_nftypawn_address"`
		BscNftypawnAddress   string `json:"bsc_nftypawn_address"`
		BobaNftypawnAddress  string `json:"boba_nftypawn_address"`
		NearNftypawnAddress  string `json:"near_nftypawn_address"`
		OneNftypawnAddress   string `json:"one_nftypawn_address"`
	} `json:"contract"`
	Blockchain bcclient.Config `json:"blockchain"`
	Mailer     struct {
		URL string `json:"url"`
	} `json:"mailer"`
}
