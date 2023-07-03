# Exchange Rate API Documentation

This documentation provides details about the Exchange Rate API endpoints and their usage.

## Get Exchange Rate

### Endpoint

`GET /rates/{cryptocurrency}/{fiat}`

Retrieve the current exchange rate for a specific cryptocurrency against a fiat currency.

#### Parameters

- `{cryptocurrency}`: The name of the cryptocurrency (e.g., bitcoin, ethereum).
- `{fiat}`: The fiat currency to get the exchange rate against (e.g., usd, eur).

#### Response

- **Status Code**: 200 OK
- **Body**: JSON object containing the exchange rate data.

#### Example

**Request**

`GET /rates/bitcoin/usd`

**Response**

```json
{
  "cryptocurrency": "bitcoin",
  "fiat": "usd",
  "rate": 30583.0
}
```

# Exchange Rate API Documentation

This documentation provides details about the Exchange Rate API endpoints and their usage.

## Get Exchange Rate for Cryptocurrency

### Endpoint

`GET /rates/{cryptocurrency}`

Retrieve the current exchange rates for a specific cryptocurrency against all available fiat currencies.

#### Parameters

- `{cryptocurrency}`: The name of the cryptocurrency (e.g., bitcoin, ethereum).

#### Response

- **Status Code**: 200 OK
- **Body**: JSON object containing the exchange rates for the specified cryptocurrency.

#### Example

**Request**

`GET /rates/bitcoin`

**Response**

```json
{
  "cryptocurrency": "bitcoin",
  "fiat_rates": {
    "usd": 30583.0,
    "eur": 28027.0,
    "gbp": 24069.0
  }
}

