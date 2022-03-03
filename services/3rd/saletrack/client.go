package saletrack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

type Client struct {
}

func (c *Client) doWithAuth(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func (c *Client) postJSON(apiURL string, headers map[string]string, jsonObject interface{}, result interface{}) error {
	bodyBytes, _ := json.Marshal(jsonObject)
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := c.doWithAuth(req)
	if err != nil {
		return fmt.Errorf("failed request: %v", err)
	}
	if resp.StatusCode >= 300 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("http response bad status %d %s", resp.StatusCode, err.Error())
		}
		return fmt.Errorf("http response bad status %d %s", resp.StatusCode, string(bodyBytes))
	}
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) getJSON(url string, headers map[string]string, result interface{}) (int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := c.doWithAuth(req)
	if err != nil {
		return 0, fmt.Errorf("failed request: %v", err)
	}
	if resp.StatusCode >= 300 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, fmt.Errorf("http response bad status %d %s", resp.StatusCode, err.Error())
		}
		return resp.StatusCode, fmt.Errorf("http response bad status %d %s", resp.StatusCode, string(bodyBytes))
	}
	if result != nil {
		return resp.StatusCode, json.NewDecoder(resp.Body).Decode(result)
	}
	return resp.StatusCode, nil
}

type MagicEdenSaleResp struct {
	TxType            string `json:"txType"`
	TransactionID     string `json:"transaction_id"`
	BlockTime         int64  `json:"blockTime"`
	Mint              string `json:"mint"`
	BuyerAddress      string `json:"buyer_address"`
	SellerAddress     string `json:"seller_address"`
	ParsedTransaction struct {
		BuyerAddress  string `json:"buyer_address"`
		SellerAddress string `json:"seller_address"`
		TotalAmount   uint64 `json:"total_amount"`
	} `json:"parsedTransaction"`
}

func (c *Client) GetMagicEdenSaleHistories(mint string) ([]*MagicEdenSaleResp, error) {
	q, err := json.Marshal(map[string]interface{}{
		"$match": map[string]interface{}{
			"mint": mint,
		},
		"$sort": map[string]interface{}{
			"blockTime": -1,
			"createdAt": -1,
		},
		"$skip": 0,
	})
	if err != nil {
		return nil, err
	}
	uri, err := url.Parse("https://api-mainnet.magiceden.io/rpc/getGlobalActivitiesByQuery")
	if err != nil {
		return nil, err
	}
	query := url.Values{}
	query.Add("q", string(q))
	uri.RawQuery = query.Encode()
	var rs struct {
		Results []*MagicEdenSaleResp `json:"results"`
	}
	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
	resp, err := client.Get(uri.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("http response bad status %d %s", resp.StatusCode, err.Error())
		}
		return nil, fmt.Errorf("http response bad status %d %s", resp.StatusCode, string(bodyBytes))
	}
	err = json.NewDecoder(resp.Body).Decode(&rs)
	if err != nil {
		return nil, err
	}
	return rs.Results, nil
}
