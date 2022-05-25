
## API Reference

#### Get all listed loans

```http
  GET /api/loans/listing
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `network` | `string` | network (NEAR, MATIC, ...)|
| `page` | `number` | page number |
| `limt` | `number` | page size |

```javascript
{
    "result": [
        {
            "id": 19,
            "created_at": "2022-05-24T10:50:55Z",
            "updated_at": "2022-05-24T10:50:55Z",
            "network": "NEAR",
            "owner": "steven4293.near",
            "lender": "",
            "asset_id": 17,
            "asset": {
                "id": 17,
                "created_at": "2022-05-24T10:50:35Z",
                "updated_at": "2022-05-24T16:39:52Z",
                "network": "NEAR",
                "collection_id": 25,
                "collection": {
                    "id": 25,
                    "created_at": "2022-05-24T10:50:35Z",
                    "updated_at": "2022-05-24T10:50:35Z",
                    "network": "NEAR",
                    "seo_url": "near-tokodao-near",
                    "name": "Tokonami",
                    "description": "2331 TOKONAMI Ready for the Revolution",
                    "verified": true,
                    "listing_asset": null,
                    "listing_total": 0,
                    "total_volume": 0.0000000000,
                    "total_listed": 0,
                    "avg24h_amount": 0.0000000000,
                    "origin_network": "",
                    "origin_contract_address": "",
                    "rand_asset": null
                },
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
                "origin_network": "",
                "origin_contract_address": "",
                "origin_token_id": "",
                "new_loan": null,
                "stats": null
            },
            "description": "",
            "currency_id": 4,
            "currency": {
                "id": 4,
                "created_at": "2022-03-03T14:16:53Z",
                "updated_at": "2022-03-03T14:16:55Z",
                "network": "NEAR",
                "contract_address": "usn",
                "decimals": 18,
                "symbol": "USN",
                "name": "USN",
                "icon_url": "https://s2.coinmarketcap.com/static/img/coins/64x64/19682.png",
                "admin_fee_address": "",
                "price": 1
            },
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
            "image_url": "",
            "signature": "",
            "status": "new",
            "data_loan_address": "2",
            "data_asset_address": "",
            "offers": [],
            "approved_offer": null,
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
    ],
    "error": null,
    "count": 5
}
```

#### Get collection detail

```http
  GET /api/collections/detail/${seo_url}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `seo_url`      | `string` | **Required**. seo_url of collection to fetch |

```
{
    "result": {
        "id": 25,
        "created_at": "2022-05-24T10:50:35Z",
        "updated_at": "2022-05-24T10:50:35Z",
        "network": "NEAR",
        "seo_url": "near-tokodao-near",
        "name": "Tokonami",
        "description": "2331 TOKONAMI Ready for the Revolution",
        "verified": true,
        "listing_asset": null,
        "listing_total": 0,
        "total_volume": 0.0000000000,
        "total_listed": 0,
        "avg24h_amount": 0.0000000000,
        "origin_network": "",
        "origin_contract_address": "",
        "rand_asset": {
            "id": 17,
            "created_at": "2022-05-24T10:50:35Z",
            "updated_at": "2022-05-24T16:39:52Z",
            "network": "NEAR",
            "collection_id": 25,
            "collection": null,
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
            "origin_network": "",
            "origin_contract_address": "",
            "origin_token_id": "",
            "new_loan": null,
            "stats": null
        }
    },
    "error": null
}
```

#### Get asset detail

```http
  GET /api/assets/detail/${seo_url}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `seo_url`      | `string` | **Required**. seo_url of asset to fetch |

```
{
    "result": {
        "id": 17,
        "created_at": "2022-05-24T10:50:35Z",
        "updated_at": "2022-05-24T16:39:52Z",
        "network": "NEAR",
        "collection_id": 25,
        "collection": {
            "id": 25,
            "created_at": "2022-05-24T10:50:35Z",
            "updated_at": "2022-05-24T10:50:35Z",
            "network": "NEAR",
            "seo_url": "near-tokodao-near",
            "name": "Tokonami",
            "description": "2331 TOKONAMI Ready for the Revolution",
            "verified": true,
            "listing_asset": null,
            "listing_total": 0,
            "total_volume": 0.0000000000,
            "total_listed": 0,
            "avg24h_amount": 0.0000000000,
            "origin_network": "",
            "origin_contract_address": "",
            "rand_asset": null
        },
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
        "origin_network": "",
        "origin_contract_address": "",
        "origin_token_id": "",
        "new_loan": {
            "id": 19,
            "created_at": "2022-05-24T10:50:55Z",
            "updated_at": "2022-05-24T10:50:55Z",
            "network": "NEAR",
            "owner": "steven4293.near",
            "lender": "",
            "asset_id": 17,
            "asset": null,
            "description": "",
            "currency_id": 4,
            "currency": {
                "id": 4,
                "created_at": "2022-03-03T14:16:53Z",
                "updated_at": "2022-03-03T14:16:55Z",
                "network": "NEAR",
                "contract_address": "usn",
                "decimals": 18,
                "symbol": "USN",
                "name": "USN",
                "icon_url": "https://s2.coinmarketcap.com/static/img/coins/64x64/19682.png",
                "admin_fee_address": "",
                "price": 1
            },
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
            "image_url": "",
            "signature": "",
            "status": "new",
            "data_loan_address": "2",
            "data_asset_address": "",
            "offers": [],
            "approved_offer": null,
            "offer_started_at": null,
            "offer_duration": 0,
            "offer_expired_at": null,
            "offer_principal_amount": 0.0000000000,
            "offer_interest_rate": 0,
            "init_tx_hash": "",
            "cancel_tx_hash": "",
            "pay_tx_hash": "",
            "liquidate_tx_hash": ""
        },
        "stats": {
            "id": 0,
            "floor_price": 6.4000000000,
            "avg_price": 0.0000000000,
            "currency": {
                "id": 5,
                "created_at": "2022-03-03T14:16:53Z",
                "updated_at": "2022-05-25T02:00:05Z",
                "network": "NEAR",
                "contract_address": "near",
                "decimals": 24,
                "symbol": "NEAR",
                "name": "NEAR",
                "icon_url": "https://s2.coinmarketcap.com/static/img/coins/64x64/3408.png",
                "admin_fee_address": "",
                "price": 5.9
            }
        }
    },
    "error": null
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

```
{
    "result": [
        {
            "id": 19,
            "created_at": "2022-05-24T10:50:55Z",
            "updated_at": "2022-05-24T10:50:55Z",
            "network": "NEAR",
            "owner": "steven4293.near",
            "lender": "",
            "asset_id": 17,
            "asset": {
                "id": 17,
                "created_at": "2022-05-24T10:50:35Z",
                "updated_at": "2022-05-24T16:39:52Z",
                "network": "NEAR",
                "collection_id": 25,
                "collection": {
                    "id": 25,
                    "created_at": "2022-05-24T10:50:35Z",
                    "updated_at": "2022-05-24T10:50:35Z",
                    "network": "NEAR",
                    "seo_url": "near-tokodao-near",
                    "name": "Tokonami",
                    "description": "2331 TOKONAMI Ready for the Revolution",
                    "verified": true,
                    "listing_asset": null,
                    "listing_total": 0,
                    "total_volume": 0.0000000000,
                    "total_listed": 0,
                    "avg24h_amount": 0.0000000000,
                    "origin_network": "",
                    "origin_contract_address": "",
                    "rand_asset": null
                },
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
                "origin_network": "",
                "origin_contract_address": "",
                "origin_token_id": "",
                "new_loan": null,
                "stats": null
            },
            "description": "",
            "currency_id": 4,
            "currency": {
                "id": 4,
                "created_at": "2022-03-03T14:16:53Z",
                "updated_at": "2022-03-03T14:16:55Z",
                "network": "NEAR",
                "contract_address": "usn",
                "decimals": 18,
                "symbol": "USN",
                "name": "USN",
                "icon_url": "https://s2.coinmarketcap.com/static/img/coins/64x64/19682.png",
                "admin_fee_address": "",
                "price": 1
            },
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
            "image_url": "",
            "signature": "",
            "status": "new",
            "data_loan_address": "2",
            "data_asset_address": "",
            "offers": [],
            "approved_offer": null,
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

```
{
    "result": [
        {
            "id": 4,
            "created_at": "2022-05-24T07:14:15Z",
            "updated_at": "2022-05-24T07:14:15Z",
            "loan_id": 13,
            "loan": {
                "id": 13,
                "created_at": "2022-05-24T07:14:15Z",
                "updated_at": "2022-05-24T07:14:15Z",
                "network": "NEAR",
                "owner": "trihuynh.near",
                "lender": "lukhuynh.near",
                "asset_id": 11,
                "asset": {
                    "id": 11,
                    "created_at": "2022-05-24T07:14:15Z",
                    "updated_at": "2022-05-24T07:14:15Z",
                    "network": "NEAR",
                    "collection_id": 18,
                    "collection": {
                        "id": 18,
                        "created_at": "2022-05-24T07:14:15Z",
                        "updated_at": "2022-05-24T07:14:15Z",
                        "network": "NEAR",
                        "seo_url": "near-x-paras-near-furr-fighters-by-vouuunear",
                        "name": "FURR FIGHTERS",
                        "description": "FURR FIGHTERS",
                        "verified": false,
                        "listing_asset": null,
                        "listing_total": 0,
                        "total_volume": 0.0000000000,
                        "total_listed": 0,
                        "avg24h_amount": 0.0000000000,
                        "origin_network": "",
                        "origin_contract_address": "",
                        "rand_asset": null
                    },
                    "seo_url": "near-x-paras-near-309189-34",
                    "contract_address": "x.paras.near",
                    "token_url": "https://ipfs.fleek.co/ipfs/bafybeihjpku3z3x4txnqt6vj7yr5gb5tkezirzdbz6uhhahcudmorthnga",
                    "token_id": "309189:34",
                    "name": "FREE PROMOTION CARD#3 #34",
                    "description": "",
                    "seller_fee_rate": 0,
                    "attributes": null,
                    "origin_network": "",
                    "origin_contract_address": "",
                    "origin_token_id": "",
                    "new_loan": null,
                    "stats": null
                },
                "description": "",
                "currency_id": 4,
                "currency": {
                    "id": 4,
                    "created_at": "2022-03-03T14:16:53Z",
                    "updated_at": "2022-03-03T14:16:55Z",
                    "network": "NEAR",
                    "contract_address": "usn",
                    "decimals": 18,
                    "symbol": "USN",
                    "name": "USN",
                    "icon_url": "https://s2.coinmarketcap.com/static/img/coins/64x64/19682.png",
                    "admin_fee_address": "",
                    "price": 1
                },
                "started_at": "2022-05-23T09:56:47Z",
                "duration": 2592000,
                "expired_at": "2022-06-22T09:56:47Z",
                "finished_at": null,
                "principal_amount": 1.0000000000,
                "interest_rate": 0.05,
                "interest_amount": 0.0000000000,
                "valid_at": "2022-05-25T09:56:32Z",
                "config": 111,
                "fee_rate": 0,
                "fee_amount": 0.0000000000,
                "nonce_hex": "1653299807",
                "image_url": "",
                "signature": "",
                "status": "done",
                "data_loan_address": "1",
                "data_asset_address": "",
                "offers": [],
                "approved_offer": null,
                "offer_started_at": "2022-05-24T03:02:07Z",
                "offer_duration": 2592000,
                "offer_expired_at": "2022-06-23T03:02:07Z",
                "offer_principal_amount": 1.0000000000,
                "offer_interest_rate": 0.05,
                "init_tx_hash": "",
                "cancel_tx_hash": "",
                "pay_tx_hash": "",
                "liquidate_tx_hash": ""
            },
            "lender": "lukhuynh.near",
            "principal_amount": 1.0000000000,
            "interest_rate": 0.05,
            "valid_at": "2022-05-26T03:00:50Z",
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
    ],
    "error": null,
    "count": 4
}
```