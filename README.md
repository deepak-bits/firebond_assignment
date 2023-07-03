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
```



## Get All Exchange Rates

### Endpoint

`GET /rates`

Retrieve the current exchange rates for all supported cryptocurrency-fiat pairs.

#### Response

- **Status Code**: 200 OK
- **Body**: JSON object containing the exchange rates for all supported cryptocurrency-fiat pairs.

#### Example

**Request**

`GET /rates`

**Response**

```json
{
  "bitcoin": {
    "usd": 30583.0,
    "eur": 28027.0,
    "gbp": 24069.0
  },
  "ethereum": {
    "usd": 1919.16,
    "eur": 1758.76,
    "gbp": 1510.39
  },
  "litecoin": {
    "usd": 112.01,
    "eur": 102.65,
    "gbp": 88.15
  }
}
```



**GET /rates/history/{cryptocurrency}/{fiat}**


## Get Exchange Rate History

### Endpoint

`GET /rates/history/{cryptocurrency}/{fiat}`

Retrieve the exchange rate history between the specified cryptocurrency and fiat currency for the past 24 hours.

#### Parameters

- `{cryptocurrency}`: The name of the cryptocurrency (e.g., bitcoin, ethereum).
- `{fiat}`: The fiat currency symbol (e.g., usd, eur, gbp).

#### Response

- **Status Code**: 200 OK
- **Body**: JSON array containing the exchange rate history for the specified cryptocurrency and fiat currency.

#### Example

**Request**

`GET /rates/history/bitcoin/usd`

**Response**

```json
[
  {
    "timestamp": 1625256000,
    "rate": 30583.0
  },
  {
    "timestamp": 1625256300,
    "rate": 30583.5
  },
  {
    "timestamp": 1625256600,
    "rate": 30584.0
  },
  ...
]
```
