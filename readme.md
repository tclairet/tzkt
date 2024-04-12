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
