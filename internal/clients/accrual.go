package clients

import (
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/LilLebowski/loyalty-system/internal/storage"
	"github.com/LilLebowski/loyalty-system/internal/utils"
)

type Client struct {
	client    *resty.Client
	serverURL string
}

func AccrualInit(client *resty.Client, serverURL string) *Client {
	return &Client{client: client, serverURL: serverURL}
}

func (acc *Client) CheckAccrual(number string) (*storage.Accrual, error) {
	accrual := storage.Accrual{}
	response, err := acc.client.R().
		SetResult(&accrual).
		SetRawPathParam("number", number).
		Get(acc.serverURL + "/api/orders/{number}")
	if response.StatusCode() == http.StatusTooManyRequests {
		return nil, utils.ErrTooManyRequests
	}
	if strings.Contains(response.Status(), http.StatusText(http.StatusNoContent)) {
		return nil, utils.ErrNoContent
	}
	if err != nil {
		return nil, err
	}
	return &accrual, nil
}
