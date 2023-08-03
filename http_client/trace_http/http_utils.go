package trace_http

import (
	"context"
	"encoding/json"
	"fmt"
	neturl "net/url"
	"prometheus-test/config"
	"prometheus-test/metrics"
	"time"

	"prometheus-test/lib/logger"
	"prometheus-test/lib/util"

	"github.com/go-resty/resty/v2"
)

var (
	EmptyByteArr []byte
	httpClient   *resty.Client
)

func Init() {
	httpClient = resty.New()
	httpClient.SetRetryCount(2)
	httpClient.SetTimeout(10 * time.Second)
	httpClient.SetDebug(config.Cfg.Log.Level == "debug")
	httpClient.SetLogger(logger.GetBizLogger())
	httpClient.EnableTrace()
}

type Client struct {
	client      *resty.Client
	RawResponse *resty.Response
	request     *resty.Request
	Result      []byte
	Err         error
	LogResult   bool
	NeedTrace   bool
}

type HttpResult struct {
	Status int
	Result []byte
	Err    error
}

type ShouldTrace func() bool

func FetchDefaultTraceClient() *Client {
	client := &Client{}
	client.initial().
		TraceData(true).
		WithTimeout(10 * time.Second).
		WithRetryCount(2).
		WithDebug(config.Cfg.Log.Level == "debug").
		WithLogger(logger.GetBizLogger())
	return client
}

func (h *Client) TraceData(needTrace bool) *Client {
	h.NeedTrace = needTrace
	return h
}

func (h *Client) ExportResult() *HttpResult {
	return &HttpResult{Status: h.RawResponse.StatusCode(), Result: h.Result, Err: h.Err}
}

func (h *Client) initial() *Client {
	h.client = httpClient
	h.request = httpClient.R()
	return h
}

func (h *Client) DeserializeResult(r any) {
	err := json.Unmarshal(h.Result, r)
	if err != nil {
		logger.GetBizLogger().Errorf("DeserializeResult failed,err=%v", err)
	}
}

func (h *Client) prometheusMetrics(url string, start time.Time, err *error) {
	path := url
	uri, errPro := neturl.Parse(url)
	if errPro == nil {
		path = uri.Path
	}
	metrics.UpdateDependence("all", path, time.Since(start).Milliseconds(), *err)
	metrics.UpdateDependenceQPS("all", path, h.RawResponse.StatusCode(), 1)
}

func (h *Client) GET(ctx context.Context, url string) *Client {
	h.WithQueryParam("request_id", util.GetRequestId(ctx))

	elapsed := time.Now()
	var err error
	defer h.prometheusMetrics(url, elapsed, &err)
	resp, err := h.request.Get(url)
	h.HandleResponse(ctx, err, url, resp)
	if h.LogResult {
		ExportResultLog(ctx, url, []byte{}, resp, err, "GET")
	}
	return h
}

func (h *Client) POST(ctx context.Context, url string, body []byte) *Client {
	h.WithQueryParam("request_id", util.GetRequestId(ctx))
	elapsed := time.Now()
	var err error
	defer h.prometheusMetrics(url, elapsed, &err)

	resp, err := h.request.SetBody(body).Post(url)
	h.HandleResponse(ctx, err, url, resp)
	if h.LogResult {
		ExportResultLog(ctx, url, body, resp, err, "POST")
	}
	return h
}

func (h *Client) PUT(ctx context.Context, url string, body []byte) *Client {
	h.WithQueryParam("request_id", util.GetRequestId(ctx))
	elapsed := time.Now()
	var err error
	defer h.prometheusMetrics(url, elapsed, &err)
	resp, err := h.request.SetBody(body).Put(url)
	h.HandleResponse(ctx, err, url, resp)
	if h.LogResult {
		ExportResultLog(ctx, url, body, resp, err, "PUT")
	}
	return h
}

func ExportResultLog(ctx context.Context, url string, body []byte,
	resp *resty.Response, err error, method string) {
	logger.Infof(ctx, "resty client access: %s, with param: %s, method: %s, status: %s, result: %s, err: %s",
		url, string(body), method, resp.Status(), util.ToJsonStr(resp.Body()), fmt.Sprintf("%s", err))
}

func (h *Client) Log(logResult bool) *Client {
	h.LogResult = logResult
	return h
}

func (h *Client) HandleResponse(ctx context.Context,
	err error, url string, resp *resty.Response) *Client {
	if err != nil {
		print(fmt.Sprintf("Fail to access: %s cause: %s", url, err))
		h.Result = EmptyByteArr
		h.Err = err
	}
	if resp == nil {
		print(fmt.Sprintf("Empty response while access: %s", url))
		h.Result = EmptyByteArr
		h.Err = err
	}
	h.RawResponse = resp
	h.Result = resp.Body()
	h.Err = nil
	return h
}

func (h *Client) WithHeaders(headers map[string]string) *Client {
	h.request.SetHeaders(headers)
	return h
}
func (h *Client) WithTimeout(timeout time.Duration) *Client {
	h.client.SetTimeout(timeout)
	return h
}

func (h *Client) WithQueryParam(username, password string) *Client {
	h.request.SetQueryParam(username, password)
	return h
}

func (h *Client) WithDebug(t bool) *Client {
	h.client.SetDebug(t)
	return h
}

func (h *Client) WithLogger(t resty.Logger) *Client {
	h.client.SetLogger(t)
	return h
}

func (h *Client) WithRetryCount(t int) *Client {
	h.client.SetRetryCount(t)
	return h
}
