package discovery_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	discovery "github.com/zeiss/fiber-discovery"

	"github.com/stretchr/testify/require"
)

type response struct {
	OidcIssuer string `json:"oidc_issuer"`
	APIURL     string `json:"api_url"`
}

func TestClient_Discover(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		oidcIssuer string
		apiURL     string
		statusCode int
	}{
		{
			name:       "successful discovery",
			response:   `{"oidc_issuer": "https://example.com", "api_url": "https://example.com/api/v1"}`,
			oidcIssuer: "https://example.com",
			apiURL:     "https://example.com/api/v1",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Add("content-type", "application/json; charset=utf-8")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response)) //nolint:errcheck
			}))
			defer server.Close()

			client := discovery.NewClient(discovery.WithClient(server.Client()))

			var res response
			err := client.Discover(context.Background(), server.URL, &res)
			require.NoError(t, err)
			require.Equal(t, tt.apiURL, res.APIURL)
			require.Equal(t, tt.oidcIssuer, res.OidcIssuer)
		})
	}
}
