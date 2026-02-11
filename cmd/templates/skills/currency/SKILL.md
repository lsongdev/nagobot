---
name: currency
description: Currency exchange rates and conversion (free API, no key required).
tags: [currency, finance, cross-platform]
---
# Currency Exchange

Query real-time exchange rates and convert currencies using free APIs. No API key required.

## Get Exchange Rate

USD to other currencies:
```
exec: curl -s "https://api.exchangerate-api.com/v4/latest/USD" | jq '{base: .base, date: .date, CNY: .rates.CNY, EUR: .rates.EUR, GBP: .rates.GBP, JPY: .rates.JPY}'
```

From any base currency:
```
exec: curl -s "https://api.exchangerate-api.com/v4/latest/BASE_CURRENCY" | jq '.rates'
```

Specific pair:
```
exec: curl -s "https://api.exchangerate-api.com/v4/latest/USD" | jq '.rates.CNY'
```

## Convert Amount

Convert 100 USD to CNY:
```
exec: curl -s "https://api.exchangerate-api.com/v4/latest/USD" | jq '.rates.CNY * 100'
```

Convert with formatted output:
```
exec: RATE=$(curl -s "https://api.exchangerate-api.com/v4/latest/USD" | jq -r '.rates.CNY') && echo "100 USD = $(echo "$RATE * 100" | bc) CNY (Rate: $RATE)"
```

Generic conversion (replace AMOUNT, FROM, TO):
```
exec: RATE=$(curl -s "https://api.exchangerate-api.com/v4/latest/FROM_CURRENCY" | jq -r '.rates.TO_CURRENCY') && echo "AMOUNT FROM_CURRENCY = $(echo "$RATE * AMOUNT" | bc) TO_CURRENCY"
```

## List All Available Currencies

```
exec: curl -s "https://api.exchangerate-api.com/v4/latest/USD" | jq '.rates | keys[]'
```

## Alternative API: FloatRates

```
exec: curl -s "https://www.floatrates.com/daily/usd.json" | jq '{CNY: .cny.rate, EUR: .eur.rate, GBP: .gbp.rate, JPY: .jpy.rate}'
```

## Alternative API: Open Exchange Rates (free tier, key required)

```
exec: curl -s "https://openexchangerates.org/api/latest.json?app_id=YOUR_APP_ID" | jq '.rates | {CNY, EUR, GBP, JPY}'
```

## Compare Multiple Currencies

```
exec: curl -s "https://api.exchangerate-api.com/v4/latest/CNY" | jq '{base: .base, USD: .rates.USD, EUR: .rates.EUR, JPY: .rates.JPY, GBP: .rates.GBP, KRW: .rates.KRW, HKD: .rates.HKD}'
```

## Common Currency Codes

- `USD` — US Dollar
- `CNY` — Chinese Yuan
- `EUR` — Euro
- `GBP` — British Pound
- `JPY` — Japanese Yen
- `KRW` — Korean Won
- `HKD` — Hong Kong Dollar
- `TWD` — Taiwan Dollar
- `SGD` — Singapore Dollar
- `AUD` — Australian Dollar
- `CAD` — Canadian Dollar
- `CHF` — Swiss Franc

## Notes

- `exchangerate-api.com` free tier: 1500 requests/month, no key needed.
- `floatrates.com`: fully free, no key needed, updated daily.
- Rates are approximate and may have a slight delay vs live market rates.
- Uses `curl` and `jq` — works on macOS, Linux, Windows (WSL).
- For production use, consider paid APIs (exchangerate-api.com paid, currencylayer, fixer.io).
