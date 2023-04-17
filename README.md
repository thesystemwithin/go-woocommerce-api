# go-woocommerce-api
A Woocommerce API Golang Wrapper for the [Woocommerce Rest API (v3)](https://woocommerce.github.io/woocommerce-rest-api-docs/)

# Install

```console
$ go get github.com/dinistavares/go-woocommerce-api
```

# Usage

Create a new API client and authenticate with your REST API keys.

```go
import (
  "github.com/dinistavares/go-woocommerce-api"
)

func main(){
  shopURL := "https://example.com"
  key := "ck_xxxxxxx"
  secret := "cs_xxxxxxx"

  client, err := woocommerce.New(shopURL)

  if err != nil {
    // handle error

    return
  }

  client.Authenticate(key, secret)
}

```

List Orders by customerID and page number.

```go
func (client *woocommerce.Client) listOrders() {
  customerID := 3
  pageNumber := 1

  opts := woocommerce.ListOrdersParams{
    Customer: customerID,
    Page: pageNumber,
  }

  orders, resp, err := client.Orders.List(opts)

  if err != nil {
    // Handle errors

    return
  }

  // Pagination headers.
  totalPages := resp.Header.Get("X-Wp-Totalpages")
  totalItems := resp.Header.Get("X-Wp-Total")

  // ....

}

```

