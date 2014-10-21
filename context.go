package maps

import (
	"net/http"

	"code.google.com/p/go.net/context"
)

type contextKey int

// NewContext returns a new context that uses the provided http.Client and API key.
func NewContext(key string, c *http.Client) context.Context {
	return WithContext(context.Background(), key, c)
}

// WithContext returns a new context in a similar way NewContext does, but initiates the new context with the specified parent.
func WithContext(parent context.Context, key string, c *http.Client) context.Context {
	vals := map[string]interface{}{}
	vals["key"] = key
	vals["httpClient"] = c
	return context.WithValue(parent, contextKey(0), vals)
}

// NewWorkContext returns a new context that uses the provided http.Client and Google Maps API for Work client ID and private key.
func NewWorkContext(clientID, privateKey string, c *http.Client) context.Context {
	return WithWorkCredentials(context.Background(), clientID, privateKey, c)
}

// WithWorkCredentials returns a new context in a similar way NewWorkContext does, but initiates the new context with the specified parent.
func WithWorkCredentials(parent context.Context, clientID, privateKey string, c *http.Client) context.Context {
	vals := map[string]interface{}{}
	vals["workClientID"] = clientID
	vals["workPrivKey"] = privateKey
	vals["httpClient"] = c
	return context.WithValue(parent, contextKey(0), vals)
}

func key(ctx context.Context) string {
	k, found := ctx.Value(contextKey(0)).(map[string]interface{})["key"]
	if found {
		return k.(string)
	}
	return ""
}

func workCreds(ctx context.Context) (string, string) {
	var clientID string
	cid, found := ctx.Value(contextKey(0)).(map[string]interface{})["workClientID"]
	if found {
		clientID = cid.(string)
	}

	var privateKey string
	pkey, found := ctx.Value(contextKey(0)).(map[string]interface{})["workPrivKey"]
	if found {
		privateKey = pkey.(string)
	}
	return clientID, privateKey
}

func httpClient(ctx context.Context) *http.Client {
	cl, found := ctx.Value(contextKey(0)).(map[string]interface{})["httpClient"]
	if found && cl != nil {
		return cl.(*http.Client)
	}
	return &http.Client{Transport: &backoff{Transport: http.DefaultTransport}}
}
