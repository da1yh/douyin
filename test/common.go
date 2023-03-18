package test

import (
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

var serverAddr = "http://localhost:8080"

func newExpect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client:   http.DefaultClient,
		BaseURL:  serverAddr,
		Reporter: httpexpect.Reporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
