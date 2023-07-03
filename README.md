# firebond_assignment

# Exchange Rate API Documentation

This documentation provides details about the Exchange Rate API endpoints and their usage.

## Get Exchange Rate

### Endpoint

`GET /exchange-rate/{cryptocurrency}/{fiat}`

Retrieve the current exchange rate for a specific cryptocurrency against a fiat currency.

#### Parameters

- `{cryptocurrency}`: The name of the cryptocurrency (e.g., bitcoin, ethereum).
- `{fiat}`: The fiat currency to get the exchange rate against (e.g., usd, eur).

#### Response

- **Status Code**: 200 OK
- **Body**: JSON object containing the exchange rate data.

#### Example

**Request**

`GET /exchange-rate/bitcoin/usd`

**Response**

```json
{
  "cryptocurrency": "bitcoin",
  "fiat": "usd",
  "rate": 30583.0
}
