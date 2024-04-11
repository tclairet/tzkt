# Exercise: Tezos Delegation Service

## Run

To start the service run

```bash
make run
```

It will start the service on port `3333` (can be configured with env var `TZKT_PORT`)

## Query

```bash
curl --request GET "http://127.0.0.1:3333/xtz/delegations"
```

You can filter on the year with the query param `year=YYYY`

```bash
curl --request GET "http://127.0.0.1:3333/xtz/delegations?year=2024"
```

The response format is as follows. It sorted from the newest to the oldest.

```json
{
  "data":[
    {
      "amount":"100180591",
      "delegator":"tz1LzEqJ1sojuyp7nnBtyVWaPhAsbQfMPscb",
      "block":"5415937",
      "timestamp":"2024-04-11T14:57:49Z"
    }
  ]
}
```


In this exercise, you will build a Golang service that gathers new [delegations](https://opentezos.com/node-baking/baking/delegating/) made on the Tezos protocol and exposes them through a public API.

## Requirements:

- The service will poll the new delegations from this Tzkt API endpoint: https://api.tzkt.io/#operation/Operations_GetDelegations
- The data aggregation service must store the delegation data in a store of your choice.
- The API must read data from that store.
- For each delegation, save the following information: sender's address, timestamp, amount, and block.
- Expose the collected data through a public API at the endpoint `/xtz/delegations`.
    - The expected response format is:
```json
{
   "data":[
      {
         "timestamp":"2022-05-05T06:29:14Z",
         "amount":"125896",
         "delegator":"tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
         "block":"2338084"
      },
      {
         "timestamp":"2021-05-07T14:48:07Z",
         "amount":"9856354",
         "delegator":"KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
         "block":"1461334"
      }
   ]
}
```
