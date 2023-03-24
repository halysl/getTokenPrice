# Get Token Price

Get token price from okx.com by [API](https://www.okx.com/docs-v5/zh/#overview).

Supply to tg-bot.

## Usage

```shell
$ get-token-price --help
$ get-token-price --token-pairs BTC-USD(,FIL-USD) --format txt
$ get-token-price --token-pairs BTC-USD(,FIL-USD) --format json 
// $ get-token-price --file ./tokens.json --format txt
// $ get-token-price --file ./tokens.json --format json
```

## tokens.json

```json
[
  {
    "token_pair": "BTC-USDT",
    "format": "txt"
  },
  {
    "token_pair": "ETH-USDT",
    "format": "json"
  }
]
```
g