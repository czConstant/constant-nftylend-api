package saletrack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/gorilla/websocket"
)

type Client struct {
	msgChain   chan string
	NftbankKey string
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

type SolnartSaleResp struct {
	Date          *time.Time `json:"date"`
	Mint          string     `json:"mint"`
	BuyerAdd      string     `json:"buyerAdd"`
	SellerAddress string     `json:"seller_address"`
	Price         float64    `json:"price"`
	Currency      string     `json:"currency"`
}

func (c *Client) GetSolnartSaleHistories(mint string) ([]*SolnartSaleResp, error) {
	var rs []*SolnartSaleResp
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("https://api.solanart.io/last_sales_token?address=%s", mint))
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
	return rs, nil
}

type OpenseaSaleResp struct {
	Data struct {
		AssetEvents struct {
			Edges []struct {
				Seller struct {
					Address string `json:"address"`
				} `json:"seller"`
				WinnerAccount struct {
					Address string `json:"address"`
				} `json:"winnerAccount"`
				Price struct {
					QuantityInEth string `json:"quantityInEth"`
				} `json:"price"`
				Transaction struct {
					BlockExplorerLink string `json:"blockExplorerLink"`
				} `json:"transaction"`
			} `json:"edges"`
		} `json:"assetEvents"`
	} `json:"data"`
}

func (c *Client) GetOpenseaSaleHistories(contractAddr string, tokenId string) ([]*OpenseaSaleResp, error) {
	var rs []*OpenseaSaleResp
	bodyBytes := fmt.Sprintf(
		`{
			"id": "EventHistoryQuery",
			"query": "query EventHistoryQuery(\n  $archetype: ArchetypeInputType\n  $bundle: BundleSlug\n  $collections: [CollectionSlug!]\n  $categories: [CollectionSlug!]\n  $chains: [ChainScalar!]\n  $eventTypes: [EventType!]\n  $cursor: String\n  $count: Int = 16\n  $showAll: Boolean = false\n  $identity: IdentityInputType\n) {\n  ...EventHistory_data_L1XK6\n}\n\nfragment AccountLink_data on AccountType {\n  address\n  config\n  isCompromised\n  user {\n    publicUsername\n    id\n  }\n  displayName\n  ...ProfileImage_data\n  ...wallet_accountKey\n  ...accounts_url\n}\n\nfragment AssetCell_asset on AssetType {\n  collection {\n    name\n    id\n  }\n  name\n  ...AssetMedia_asset\n  ...asset_url\n}\n\nfragment AssetCell_assetBundle on AssetBundleType {\n  assetQuantities(first: 2) {\n    edges {\n      node {\n        asset {\n          collection {\n            name\n            id\n          }\n          name\n          ...AssetMedia_asset\n          ...asset_url\n          id\n        }\n        relayId\n        id\n      }\n    }\n  }\n  name\n  ...bundle_url\n}\n\nfragment AssetMedia_asset on AssetType {\n  animationUrl\n  backgroundColor\n  collection {\n    displayData {\n      cardDisplayStyle\n    }\n    id\n  }\n  isDelisted\n  imageUrl\n  displayImageUrl\n}\n\nfragment AssetQuantity_data on AssetQuantityType {\n  asset {\n    ...Price_data\n    id\n  }\n  quantity\n}\n\nfragment CollectionLink_assetContract on AssetContractType {\n  address\n  blockExplorerLink\n}\n\nfragment CollectionLink_collection on CollectionType {\n  name\n  ...collection_url\n  ...verification_data\n}\n\nfragment EventHistory_data_L1XK6 on Query {\n  assetEvents(after: $cursor, bundle: $bundle, archetype: $archetype, first: $count, categories: $categories, collections: $collections, chains: $chains, eventTypes: $eventTypes, identity: $identity, includeHidden: true) {\n    edges {\n      node {\n        assetBundle @include(if: $showAll) {\n          relayId\n          ...AssetCell_assetBundle\n          ...bundle_url\n          id\n        }\n        assetQuantity {\n          asset @include(if: $showAll) {\n            relayId\n            assetContract {\n              ...CollectionLink_assetContract\n              id\n            }\n            ...AssetCell_asset\n            ...asset_url\n            collection {\n              ...CollectionLink_collection\n              id\n            }\n            id\n          }\n          ...quantity_data\n          id\n        }\n        relayId\n        eventTimestamp\n        eventType\n        offerExpired\n        customEventName\n        ...utilsAssetEventLabel\n        devFee {\n          asset {\n            assetContract {\n              chain\n              id\n            }\n            id\n          }\n          quantity\n          ...AssetQuantity_data\n          id\n        }\n        devFeePaymentEvent {\n          ...EventTimestamp_data\n          id\n        }\n        fromAccount {\n          address\n          ...AccountLink_data\n          id\n        }\n        price {\n          quantity\n          quantityInEth\n          ...AssetQuantity_data\n          id\n        }\n        endingPrice {\n          quantity\n          ...AssetQuantity_data\n          id\n        }\n        seller {\n          ...AccountLink_data\n          id\n        }\n        toAccount {\n          ...AccountLink_data\n          id\n        }\n        winnerAccount {\n          ...AccountLink_data\n          id\n        }\n        ...EventTimestamp_data\n        id\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n\nfragment EventTimestamp_data on AssetEventType {\n  eventTimestamp\n  transaction {\n    blockExplorerLink\n    id\n  }\n}\n\nfragment Price_data on AssetType {\n  decimals\n  imageUrl\n  symbol\n  usdSpotPrice\n  assetContract {\n    blockExplorerLink\n    chain\n    id\n  }\n}\n\nfragment ProfileImage_data on AccountType {\n  imageUrl\n}\n\nfragment accounts_url on AccountType {\n  address\n  user {\n    publicUsername\n    id\n  }\n}\n\nfragment asset_url on AssetType {\n  assetContract {\n    address\n    chain\n    id\n  }\n  tokenId\n}\n\nfragment bundle_url on AssetBundleType {\n  slug\n}\n\nfragment collection_url on CollectionType {\n  slug\n}\n\nfragment quantity_data on AssetQuantityType {\n  asset {\n    decimals\n    id\n  }\n  quantity\n}\n\nfragment utilsAssetEventLabel on AssetEventType {\n  isMint\n  eventType\n}\n\nfragment verification_data on CollectionType {\n  isMintable\n  isSafelisted\n  isVerified\n}\n\nfragment wallet_accountKey on AccountType {\n  address\n}\n",
			"variables": {
				"archetype": {
					"chain": "ETHEREUM",
					"tokenId": "%s",
					"assetContractAddress": "%s"
				},
				"bundle": null,
				"collections": null,
				"categories": null,
				"chains": null,
				"eventTypes": [
					"AUCTION_SUCCESSFUL"
				],
				"cursor": null,
				"count": 16,
				"showAll": false,
				"identity": null
			}
		}`,
		tokenId,
		contractAddr,
	)
	client := &http.Client{}
	resp, err := client.Post("https://api.opensea.io/graphql/", "application/json", bytes.NewBuffer([]byte(bodyBytes)))
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
	return rs, nil
}

func (c *Client) StartWssSolsea(msgReceivedFunc func(msg string)) {
	c.msgChain = make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	wc, _, err := websocket.DefaultDialer.Dial("wss://api.all.art/socket.io/?EIO=3&transport=websocket", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer wc.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := wc.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			if msgReceivedFunc != nil {
				msgReceivedFunc(string(message))
			}
		}
	}()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			{
				err := wc.WriteMessage(websocket.TextMessage, []byte("2"))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		case msg := <-c.msgChain:
			{
				log.Println("write:", msg)
				err := wc.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		case <-interrupt:
			log.Println("interrupt")
			err := wc.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (c *Client) PubSolseaMsg(msg string) {
	c.msgChain <- msg
}

type EvmNftMetaResp struct {
	Description  string `json:"description"`
	ExternalUrl  string `json:"external_url"`
	Image        string `json:"image"`
	Name         string `json:"name"`
	Collection   string `json:"collection"`
	CollectionID string `json:"collection_id"`
	CreatorID    string `json:"creator_id"`
	Attributes   []struct {
		TraitType string      `json:"trait_type"`
		Value     interface{} `json:"value"`
	} `json:"attributes"`
}

func (c *Client) GetEvmNftMetaResp(tokenURL string) (*EvmNftMetaResp, error) {
	tokenURL = strings.Replace(tokenURL, "https://ipfs.fleek.co/ipfs", "https://cloudflare-ipfs.com/ipfs", -1)
	var rs EvmNftMetaResp
	client := &http.Client{}
	resp, err := client.Get(tokenURL)
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
	return &rs, nil
}

type NearNftMetaResp struct {
	Description string `json:"description"`
	Collection  string `json:"collection"`
	Attributes  []struct {
		TraitType string      `json:"trait_type"`
		Value     interface{} `json:"value"`
	} `json:"attributes"`
}

func (c *Client) GetNearNftMetaResp(tokenURL string) (*NearNftMetaResp, error) {
	tokenURL = strings.Replace(tokenURL, "https://ipfs.fleek.co/ipfs", "https://ipfs.io/", -1)
	var rs NearNftMetaResp
	client := &http.Client{}
	resp, err := client.Get(tokenURL)
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
	return &rs, nil
}

type ParasSaleResp struct {
	IssuedAt        int64  `json:"issued_at"`
	ID              string `json:"_id"`
	From            string `json:"from"`
	To              string `json:"to"`
	FtTokenID       string `json:"ft_token_id"`
	TransactionHash string `json:"transaction_hash"`
	ContractID      string `json:"contract_id"`
	TokenID         string `json:"token_id"`
	Msg             struct {
		Params struct {
			Price numeric.BigInt `json:"price"`
		} `json:"params"`
	} `json:"msg"`
}

func (c *Client) GetParasSaleHistories(contractID string, tokenID string) ([]*ParasSaleResp, error) {
	var result struct {
		Data struct {
			Results []*ParasSaleResp `json:"results"`
		} `json:"data"`
	}
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("https://api-v2-mainnet.paras.id/activities?__limit=30&contract_id=%s&token_id=%s&type=resolve_purchase", contractID, tokenID))
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Data.Results, nil
}

type ParasProfileResp struct {
	ID        string `json:"_id"`
	AccountID string `json:"accountId"`
	IsCreator bool   `json:"isCreator"`
}

func (c *Client) GetParasProfile(contractID string) ([]*ParasProfileResp, error) {
	var result struct {
		Data struct {
			Results []*ParasProfileResp `json:"results"`
		} `json:"data"`
	}
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("https://api-v2-mainnet.paras.id/profiles?accountId=%s", contractID))
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Data.Results, nil
}

type ParasCollectionStatsResp struct {
	ID          string           `json:"_id"`
	AccountID   string           `json:"accountId"`
	Volume      numeric.BigInt   `json:"volume"`
	VolumeUsd   numeric.BigFloat `json:"volume_usd"`
	AvgPrice    numeric.BigInt   `json:"avg_price"`
	AvgPriceUsd numeric.BigFloat `json:"avg_price_usd"`
	FloorPrice  numeric.BigInt   `json:"floor_price"`
}

func (c *Client) GetParasCollectionStats(contractID string) (*ParasCollectionStatsResp, error) {
	var result struct {
		Data struct {
			Results *ParasCollectionStatsResp `json:"results"`
		} `json:"data"`
	}
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("https://api-v2-mainnet.paras.id/collection-stats?collection_id=%s", contractID))
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Data.Results, nil
}

type NftbankSaleResp struct {
	BlockTimestamp  string           `json:"block_timestamp"`
	BuyerAddress    string           `json:"buyer_address"`
	SellerAddress   string           `json:"seller_address"`
	TransactionHash string           `json:"transaction_hash"`
	SoldPriceEth    numeric.BigFloat `json:"sold_price_eth"`
}

func (c *Client) GetNftbankSaleHistories(contractID string, tokenID string, chainID string) ([]*NftbankSaleResp, error) {
	var result struct {
		Data []*NftbankSaleResp `json:"data"`
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.nftbank.ai/estimates-v2/asset-events/%s/%s?chain_id=%s", contractID, tokenID, chainID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", c.NftbankKey)
	client := &http.Client{}
	resp, err := client.Do(req)
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

type NftbankFloorPriceResp struct {
	EstimatedAt string `json:"estimated_at"`
	FloorPrice  []*struct {
		CurrencySymbol string           `json:"currency_symbol"`
		FloorPrice     numeric.BigFloat `json:"floor_price"`
	} `json:"floor_price"`
}

func (c *Client) GetNftbankFloorPrice(contractID string, chainID string) ([]*NftbankFloorPriceResp, error) {
	var result struct {
		Data []*NftbankFloorPriceResp `json:"data"`
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.nftbank.ai/estimates-v2/floor_price/%s?chain_id=%s", contractID, chainID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", c.NftbankKey)
	client := &http.Client{}
	resp, err := client.Do(req)
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) GetCoingeckoPrice(symbol string) (float64, error) {
	switch symbol {
	case "ETH":
		{
			symbol = "ETHEREUM"
		}
	}
	result := map[string]interface{}{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=USD", symbol), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode >= 300 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("http response bad status %d %s", resp.StatusCode, err.Error())
		}
		return 0, fmt.Errorf("http response bad status %d %s", resp.StatusCode, string(bodyBytes))
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result[strings.ToLower(symbol)].(map[string]interface{})["usd"].(float64), nil
}
