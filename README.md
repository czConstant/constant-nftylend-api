# API Reference

## API for collection

#### Get collection list

```http
GET /api/collections/list
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `page` | `number` | page number |
| `limt` | `number` | page size |

```json
{
    "result": [
        CollectionResponse,
        CollectionResponse
    ],
    "error": null,
    "count": 8
}
```

#### Get collection detail

```http
GET /api/collections/detail/${seo_url}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `seo_url`      | `string` | **Required**. seo_url of collection to fetch |

```json
{
    "result": CollectionResponse,
    "error": null
}
```

## API for asset

#### Get asset detail

```http
GET /api/assets/detail/${seo_url}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `seo_url`      | `string` | **Required**. seo_url of asset to fetch |

```json
{
    "result": AssetResponse,
    "error": null
}
```

## API for loan

#### Get all listed loans

```http
GET /api/loans/listing
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `network` | `string` | network (NEAR, MATIC, ...)|
| `page` | `number` | page number |
| `limt` | `number` | page size |

```json
{
    "result": [
        LoanResponse,
        LoanResponse,
    ],
    "error": null,
    "count": 5
}
```

#### Get list loan

```http
GET /api/loans/list
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `network`      | `string` | network (NEAR, MATIC, ...) |
| `owner`      | `string` | borrower address |
| `lender`      | `string` | lender address |
| `status`      | `string` | loan status (new, created, done, liquidated) |
| `page` | `number` | page number |
| `limt` | `number` | page size |

```json
{
    "result": [
        LoanResponse,
        LoanResponse,
    ],
    "error": null,
    "count": 1
}
```

#### Get list offer

```http
GET /api/loans/offers
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `network`      | `string` | network (NEAR, MATIC, ...) |
| `owner`      | `string` | borrower address |
| `lender`      | `string` | lender address |
| `status`      | `string` | loan status (new, created, done, liquidated) |
| `page` | `number` | page number |
| `limt` | `number` | page size |

```json
{
    "result": [
        LoanOfferResponse,
        LoanOfferResponse,
    ],
    "error": null,
    "count": 4
}
```

#### Get borrower stats

```http
GET /api/loans/borrower-stats/${address}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `address`      | `string` | borrower address |

```json
{
    "result": {
        "total_loans": 3,
        "total_done_loans": 3,
        "total_volume": 3.0000000000
    },
    "error": null
}
```

#### Get loan transactions

```http
GET /api/loans/transactions
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `asset_id`      | `number` | id of asset to fetch |
| `page` | `number` | page number |
| `limt` | `number` | page size |

```json
{
    "result": [
        {
            "id": 7,
            "created_at": "2022-05-24T03:23:19Z",
            "updated_at": "2022-05-24T03:23:19Z",
            "network": "NEAR",
            "loan_id": 7,
            "loan": LoanResponse,
            "type": "listed",
            "borrower": "trihuynh.near",
            "lender": "",
            "started_at": "2022-05-24T03:23:15Z",
            "duration": 864000,
            "expired_at": "2022-06-03T03:23:15Z",
            "principal_amount": 1.0000000000,
            "interest_rate": 0.02,
            "tx_hash": ""
        }
    ],
    "error": null,
    "count": 3
}
```

# API Response

#### CurrencyResponse

```json
{
    "id": 7,
    "created_at": "2022-03-03T14:16:53Z",
    "updated_at": "2022-03-03T14:16:55Z",
    "network": "NEAR",
    "contract_address": "dac17f958d2ee523a2206206994597c13d831ec7.factory.bridge.near",
    "decimals": 6,
    "symbol": "USDT",
    "name": "Tether",
    "icon_url": "https://s2.coinmarketcap.com/static/img/coins/64x64/825.png",
    "admin_fee_address": "",
    "price": 1
}
```

| Key | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `number` | collection id |
| `created_at` | `date` | created time |
| `updated_at` | `date` | update time |
| `network` | `string` | network of currency |
| `contract_address` | `string` | contract address of currency |
| `decimals` | `number` | decimals of currency |
| `symbol` | `string` | symbol of currency |
| `name` | `string` | name of currency |
| `icon_url` | `string` | image url of currency |
| `price` | `number` | price of currency |

#### CollectionResponse

```json
{
    "id": 26,
    "created_at": "2022-05-24T16:42:24Z",
    "updated_at": "2022-05-24T16:42:24Z",
    "network": "NEAR",
    "seo_url": "near-x-paras-near-metamorphoses-by-ludanear",
    "name": "MetaMorphoses",
    "description": "The new Metaverse is approaching! \n\nWe invite you to take your unique chance and steal an egg! Exactly in a week a new story will begin and the egg will disappear.",
    "verified": false,
    "listing_asset": null,
    "listing_total": 0,
    "total_volume": 0.0000000000,
    "total_listed": 0,
    "avg24h_amount": 0.0000000000
}
```

| Key | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `number` | collection id |
| `created_at` | `date` | created time |
| `updated_at` | `date` | update time |
| `network` | `string` | network of collection |
| `seo_url` | `string` | seo url of collection |
| `name` | `string` | name of collection |
| `description` | `string` | description of collection |
| `total_volume` | `number` | volume loan of collection |
| `total_listed` | `number` | number listing loan of collection |
| `avg24h_amount` | `number` | average loan volume in 24 hour |

#### AssetResponse

```json
{
    "id": 17,
    "created_at": "2022-05-24T10:50:35Z",
    "updated_at": "2022-05-24T16:39:52Z",
    "network": "NEAR",
    "collection_id": 25,
    "collection": CollectionResponse,
    "seo_url": "near-tokodao-near-1525",
    "contract_address": "tokodao.near",
    "token_url": "https://gateway.pinata.cloud/ipfs/QmehZFCwtyubKgPBRpiJ4BHURMkgWFuU2UUg4nw66bqvpb/1525.png",
    "token_id": "1525",
    "name": "Tokonami #1525",
    "description": "",
    "seller_fee_rate": 0,
    "attributes": [
        {
            "trait_type": "Background",
            "value": "Crimson Desert Lands"
        },
        {
            "trait_type": "Wing",
            "value": "None"
        },
        {
            "trait_type": "Rightarm",
            "value": "General Healing Shoulder Kit"
        },
        {
            "trait_type": "Chest",
            "value": "Lightweight blasting Chest Plate"
        },
        {
            "trait_type": "Head",
            "value": "Turbo Assisted head"
        },
        {
            "trait_type": "Helmet",
            "value": "Pain Coated Horns Helmet"
        },
        {
            "trait_type": "Visor",
            "value": "Rapid Gliding Sunvisor"
        },
        {
            "trait_type": "Medal",
            "value": "Ground Unit"
        },
        {
            "trait_type": "Leftarm",
            "value": "General Healing Shoulder Kit"
        },
        {
            "trait_type": "Weapon",
            "value": "Blasting Warhammer"
        }
    ],
    "new_loan": LoanRespnse,
    "stats": null
}
```

| Key | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `number` | asset id |
| `created_at` | `date` | created time |
| `updated_at` | `date` | update time |
| `network` | `string` | network of asset |
| `collection_id` | `number` | collection id |
| `collection` | `CollectionResponse` | collection detail |
| `seo_url` | `string` | seo url of asset |
| `contract_address` | `string` | contract address of asset |
| `token_id` | `string` | token id of asset |
| `token_url` | `string` | image url of asset |
| `name` | `string` | name of asset |
| `description` | `string` | description of asset |
| `seller_fee_rate` | `number` | loyalty rate of asset |
| `attributes` | `object` | attributes of asset |
| `new_loan` | `LoanRespnse` | listing loan of asset |

#### LoanResponse

```json
{
    "id": 19,
    "created_at": "2022-05-24T10:50:55Z",
    "updated_at": "2022-05-24T10:50:55Z",
    "network": "NEAR",
    "owner": "steven4293.near",
    "lender": "",
    "asset_id": 17,
    "asset": AssetResponse,
    "currency_id": 4,
    "currency": CurrencyResponse,
    "started_at": "2022-05-24T10:50:54Z",
    "duration": 864000,
    "expired_at": "2022-06-03T10:50:54Z",
    "finished_at": null,
    "principal_amount": 5.0000000000,
    "interest_rate": 0.12,
    "interest_amount": 0.0000000000,
    "valid_at": "2022-05-25T10:50:32Z",
    "config": 111,
    "fee_rate": 0,
    "fee_amount": 0.0000000000,
    "nonce_hex": "1653389454",
    "signature": "",
    "status": "new",
    "data_loan_address": "2",
    "data_asset_address": "",
    "offers": []LoanOfferResponse,
    "approved_offer": LoanOfferResponse,
    "offer_started_at": null,
    "offer_duration": 0,
    "offer_expired_at": null,
    "offer_principal_amount": 0.0000000000,
    "offer_interest_rate": 0,
    "init_tx_hash": "",
    "cancel_tx_hash": "",
    "pay_tx_hash": "",
    "liquidate_tx_hash": ""
}
```

| Key | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `number` | loan id |
| `created_at` | `date` | created time |
| `updated_at` | `date` | update time |
| `network` | `string` | network of loan |
| `owner` | `string` | borrower address |
| `lender` | `string` | lender address |
| `asset_id` | `number` | asset id |
| `asset` | `AssetResponse` | asset detail |
| `currency_id` | `number` | currency id |
| `currency` | `CurrencyResponse` | currency detail |
| `started_at` | `date` | loan started time |
| `duration` | `number` | loan duration |
| `expired_at` | `date` | loan expired time |
| `finished_at` | `date` | loan finished time |
| `principal_amount` | `number` | loan principal amount |
| `interest_rate` | `number` | loan interest rate |
| `interest_amount` | `number` | loan interest amount |
| `valid_at` | `date` | loan valid time |
| `config` | `number` | loan config |
| `fee_rate` | `string` | loan platform fee rate |
| `fee_amount` | `string` | loan platform fee amount |
| `nonce_hex` | `string` | random hex client |
| `signature` | `string` | borrower signature |
| `status` | `string` | loan status |
| `data_loan_address` | `string` | loan data onchain address |
| `data_asset_address` | `string` | loan data onchain address |
| `offers` | `string` | list of loan offer |
| `approved_offer` | `string` | actived offer |
| `offer_started_at` | `string` | started time of offer |
| `offer_duration` | `string` | duration of offer |
| `offer_expired_at` | `string` | expired time of offer |
| `offer_principal_amount` | `string` | principal amount of offer |
| `offer_interest_rate` | `string` | interest rate of offer |
| `init_tx_hash` | `string` | onchain hash when new loan |
| `cancel_tx_hash` | `string` | onchain hash when cancel loan |
| `pay_tx_hash` | `string` | onchain hash when repaid loan |
| `liquidate_tx_hash` | `string` | onchain hash when liquidate loan |

#### LoanOfferResponse

```json
{
    "id": 1,
    "created_at": "2022-05-23T09:03:03Z",
    "updated_at": "2022-05-23T09:51:48Z",
    "loan_id": 1,
    "loan": LoanResponse,
    "lender": "lukhuynh.near",
    "principal_amount": 1.0000000000,
    "interest_rate": 0.1,
    "valid_at": "2022-05-23T09:02:55Z",
    "duration": 2592000,
    "nonce_hex": "1",
    "signature": "",
    "status": "done",
    "data_offer_address": "",
    "data_currency_address": "",
    "make_tx_hash": "",
    "accept_tx_hash": "",
    "cancel_tx_hash": "",
    "close_tx_hash": ""
}
```

| Key | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `number` | offer id |
| `created_at` | `date` | created time |
| `updated_at` | `date` | update time |
| `network` | `string` | network of loan offer |
| `loan_id` | `number` | loan id |
| `loan` | `LoanResponse` | loan detail |
| `lender` | `string` | lender of loan offer |
| `duration` | `number` | offer duration |
| `expired_at` | `date` | offer expired time |
| `finished_at` | `date` | offer finished time |
| `principal_amount` | `number` | offer principal amount |
| `interest_rate` | `number` | offer interest rate |
| `interest_amount` | `number` | offer interest amount |
| `nonce_hex` | `string` | random hex client |
| `signature` | `string` | lender signature |
| `status` | `string` | offer status |

#### LoanTransactionResponse

```json
{
    "id": 7,
    "created_at": "2022-05-24T03:23:19Z",
    "updated_at": "2022-05-24T03:23:19Z",
    "network": "NEAR",
    "loan_id": 7,
    "loan": LoanResponse,
    "type": "listed",
    "borrower": "trihuynh.near",
    "lender": "",
    "started_at": "2022-05-24T03:23:15Z",
    "duration": 864000,
    "expired_at": "2022-06-03T03:23:15Z",
    "principal_amount": 1.0000000000,
    "interest_rate": 0.02,
    "tx_hash": ""
}
```

| Key | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `number` | loan id |
| `created_at` | `date` | created time |
| `updated_at` | `date` | update time |
| `network` | `string` | network of loan |
| `loan_id` | `number` | loan id |
| `loan` | `LoanResponse` | loan detail |
| `type` | `string` | transaction type (listed, cancelled, offered, repaid, liquidated) |
| `borrower` | `string` | borrower address |
| `lender` | `string` | lender address |
| `duration` | `number` | loan duration |
| `expired_at` | `date` | loan expired time |
| `principal_amount` | `number` | loan principal amount |
| `interest_rate` | `number` | loan interest rate |
| `tx_hash` | `string` | transaction hash |