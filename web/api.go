package web

import (
	"echo-demo/model"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var (
	timeLayout = time.RFC3339
	dateLayout = "2006-01-02"
)

type API struct {
	Repo   model.Repository
	AppURL *url.URL
}

// NewAPI API実装を初期化して生成します
func NewAPI(repo model.Repository, appURL *url.URL) (*API, error) {
	api := &API{
		Repo:   repo,
		AppURL: appURL,
	}
	return api, nil
}
