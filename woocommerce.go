package woocommerce

import (
  "bytes"
  "encoding/base64"
  "encoding/json"
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "net/url"
  "strconv"
  "strings"
  "time"

  "github.com/JoshuaDoes/crunchio"
  "github.com/google/go-querystring/query"
)

const (
  defaultRestEndpointVersion   = "v3"
  defaultHeaderName            = "Authorization"
  acceptedContentType          = "application/json"
  userAgent                    = "go-woocommerce-api/1.1"
  clientRequestRetryAttempts   = 2
  clientRequestRetryHoldMillis = 1000
)

var errorDoAllAttemptsExhausted = errors.New("all request attempts were exhausted")
var errorDoAttemptNilRequest = errors.New("request could not be constructed")

type ClientConfig struct {
  HttpClient          *http.Client
  RestEndpointURL     string
  RestEndpointVersion string
}

type auth struct {
  HeaderName string
  ApiKey     string
}

type Client struct {
  config  *ClientConfig
  client  *http.Client
  auth    *auth
  baseURL *url.URL
  rl      *RateLimit

  Categories *CategoriesService
  Coupons    *CouponsService
  Customers  *CustomersService
  OrderNotes *OrderNotesService
  Orders     *OrdersService
  Products   *ProductsService
  Refunds    *RefundsService
  Tags       *TagsService
  Webhooks   *WebhookService
}

type service struct {
  client *Client
}

type errorResponse struct {
  Response *http.Response

  Code    string    `json:"code"`
  Message string    `json:"message"`
  Data    ErrorData `json:"data"`
}

type ErrorData struct {
  Status int `json:"status"`
}

func (response *errorResponse) Error() string {
  return fmt.Sprintf("%v %v: %d %v",
    response.Response.Request.Method, response.Response.Request.URL,
    response.Response.StatusCode, response.Message)
}

func New(shopURL string) (*Client, error) {
  if shopURL == "" {
    return nil, errors.New("store url is required")
  }

  config := ClientConfig{
    HttpClient: http.DefaultClient,
  }

  config.HttpClient = http.DefaultClient
  config.RestEndpointURL = shopURL
  config.RestEndpointVersion = defaultRestEndpointVersion

  // Create client
  baseURL, err := url.Parse(config.RestEndpointURL + "/wp-json/wc/" + defaultRestEndpointVersion)

  if err != nil {
    return nil, err
  }

  client := &Client{config: &config, client: config.HttpClient, auth: &auth{}, baseURL: baseURL, rl: &RateLimit{}}

  // Map services
  client.Categories = &CategoriesService{client: client}
  client.Coupons = &CouponsService{client: client}
  client.Customers = &CustomersService{client: client}
  client.OrderNotes = &OrderNotesService{client: client}
  client.Orders = &OrdersService{client: client}
  client.Products = &ProductsService{client: client}
  client.Refunds = &RefundsService{client: client}
  client.Tags = &TagsService{client: client}
  client.Webhooks = &WebhookService{client: client}

  return client, nil
}

// Authenticate saves authenitcation parameters for user
func (client *Client) Authenticate(consumer_key string, consumer_secret string) {
  client.auth.HeaderName = defaultHeaderName

  client.auth.ApiKey = "Basic " + base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{consumer_key, consumer_secret}, ":")))
}

// NewRequest creates an API request
func (client *Client) NewRequest(method, urlStr string, opts interface{}, body interface{}) (*http.Request, error) {
  // Append Query Params to URL
  if opts != nil {
    queryParams, err := query.Values(opts)

    if err != nil {
      return nil, err
    }

    rawQuery := queryParams.Encode()

    if rawQuery != "" {
      urlStr += "?" + rawQuery
    }
  }

  rel, err := url.Parse(client.config.RestEndpointVersion + urlStr)
  if err != nil {
    return nil, err
  }

  url := client.baseURL.ResolveReference(rel)

  var buf io.ReadWriter
  if body != nil {
    buf = new(bytes.Buffer)

    err := json.NewEncoder(buf).Encode(body)
    if err != nil {
      return nil, err
    }
  }

  req, err := http.NewRequest(method, url.String(), buf)
  if err != nil {
    return nil, err
  }

  req.Header.Add(client.auth.HeaderName, client.auth.ApiKey)
  req.Header.Add("Accept", acceptedContentType)
  req.Header.Add("Content-type", acceptedContentType)
  req.Header.Add("User-Agent", userAgent)

  return req, nil
}

type RateLimit struct {
  Limit     int       //Total amount of requests per rate
  Remaining int       //Remaining requests for this rate window
  Reset     int       //Seconds until the rate limit resets
  After     int       //Seconds until the exceeded rate limit is unblocked
  LastReq   time.Time //The last recorded request time for calculating the rate limit
}

// Sleep: TODO: Sleep properly
func (rl *RateLimit) Sleep() {
  /*if rl.LastReq.IsZero() || rl.Remaining > 0 {
      rl.LastReq = time.Now()
      return
    }

    now := time.Now().Unix()
    since := rl.LastReq.Unix()
    remaining := int(now - since)
    seconds := time.Second * time.Duration(remaining-rl.Reset)

    if seconds > 0 {
      fmt.Println("[WOOCOMMERCE] Sleeping for", seconds)
      time.Sleep(seconds)
    }*/
}

func (rl *RateLimit) String() string {
  return fmt.Sprintf("RateLimit{Limit:%d Remaining:%d Reset:%d After:%d LastReq:%d}",
    rl.Limit, rl.Remaining, rl.Reset, rl.After, rl.LastReq.Unix())
}

func (client *Client) RateLimitUpdate(resp *http.Response) {
  if client.rl == nil {
    client.rl = new(RateLimit)
  }
  client.rl.LastReq = time.Now()

  if resp != nil {

    for key, val := range resp.Header {
      //fmt.Printf("[WOOCOMMERCE] Header: %s = %v\n", key, val)

      switch strings.ToLower(key) {
      case "x-ratelimit-limit":
        if len(val) <= 0 {
          continue
        }
        limit, err := strconv.Atoi(val[0])
        if err != nil {
          fmt.Println("[WOOCOMMERCE] Failed to convert limit:", val[0])
          continue
        }
        client.rl.Limit = limit
      case "x-ratelimit-remaining":
        if len(val) <= 0 {
          continue
        }
        remaining, err := strconv.Atoi(val[0])
        if err != nil {
          fmt.Println("[WOOCOMMERCE] Failed to convert remaining:", val[0])
          continue
        }
        client.rl.Remaining = remaining
      case "x-ratelimit-reset":
        if len(val) <= 0 {
          continue
        }
        reset, err := strconv.Atoi(val[0])
        if err != nil {
          fmt.Println("[WOOCOMMERCE] Failed to convert reset:", val[0])
          continue
        }
        client.rl.Reset = reset
      case "x-ratelimit-retry-after":
        if len(val) <= 0 {
          continue
        }
        after, err := strconv.Atoi(val[0])
        if err != nil {
          fmt.Println("[WOOCOMMERCE] Failed to convert after:", val[0])
          continue
        }
        client.rl.After = after
      }
    }

    //fmt.Println("[WOOCOMMERCE] Rate limit stats:", client.rl)
  }
}

// Do sends an API request
func (client *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
  var lastErr error

  attempts := 0

  for attempts < clientRequestRetryAttempts {
    attempts++

    // Sleep if the rate limit says we should
    client.rl.Sleep()

    // Dispatch request attempt
    resp, shouldRetry, err := client.doAttempt(req, v)

    // Store the newest rate limit stats
    client.RateLimitUpdate(resp)

    // Return response straight away? (we are done)
    if !shouldRetry {
      return resp, err
    }

    // Should retry: store last error (we are not done)
    lastErr = err
  }

  // Set default error? (all attempts failed, but no error is set)
  if lastErr == nil {
    lastErr = errorDoAllAttemptsExhausted
  }

  // All attempts failed, return last attempt error
  return nil, lastErr
}

func (client *Client) doAttempt(req *http.Request, v interface{}) (*http.Response, bool, error) {
  if req == nil {
    return nil, false, errorDoAttemptNilRequest
  }

  resp, err := client.client.Do(req)

  if checkRequestRetry(resp, err) {
    return resp, true, err
  }

  defer resp.Body.Close()

  err = checkResponse(resp)
  if err != nil {
    return resp, false, err
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return resp, false, err
  }
  bodyBuf := crunchio.NewBuffer(body)

  /*fmt.Printf("--- JSON DUMP ---\n\n%s\n\n--- ---- ---- ---\n", bodyBuf.String())
    if _, err := bodyBuf.Seek(0, 0); err != nil {
      return resp, false, err
    }*/

  if v != nil {
    if w, ok := v.(io.Writer); ok {
      io.Copy(w, bodyBuf)
    } else {
      err = json.NewDecoder(bodyBuf).Decode(v)
      if err == io.EOF {
        err = nil
      }
    }
  }

  return resp, false, err
}

// checkRequestRetry checks if should retry request
func checkRequestRetry(response *http.Response, err error) bool {
  // Low-level error, or response status is a server error? (HTTP 5xx)
  if err != nil || response.StatusCode >= 500 {
    return true
  }

  // No low-level error (should not retry)
  return false
}

// checkResponse checks response for errors
func checkResponse(response *http.Response) error {
  // No error in response? (HTTP 2xx)
  if code := response.StatusCode; 200 <= code && code <= 299 {
    return nil
  }

  // Map response error data (eg. HTTP 4xx)
  errorResponse := &errorResponse{Response: response}

  data, err := ioutil.ReadAll(response.Body)
  if err == nil && data != nil {
    json.Unmarshal(data, errorResponse)
  }

  return errorResponse
}
