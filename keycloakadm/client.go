// package keycloakadm contains a client and relevant data structs for interacting
// with the Keycloak Admin REST API
//
// For mapping, see https://www.keycloakadmorg/docs-api/4.0/rest-api/index.html
package keycloakadm

import (
	"net/http"
	"net/url"

	"context"
	"fmt"

	"gopkg.in/resty.v1"
)

const userAgent = "go/keycloak-admin"

// Client is the API client for talking to keycloak admin
type Client struct {
	Server     url.URL
	restClient *resty.Client
	Realm      string
}

// NewClient creates a new client instance set to talk to the keycloak service
// as well as the various services for working with specific resources
func NewClient(u url.URL, c *http.Client, realm string) *Client {
	u.Path = "/auth/admin"
	restClient := resty.NewWithClient(c)

	client := &Client{
		Server:     u,
		restClient: restClient,
		Realm:      realm,
	}

	return client
}

// Debug enables debugging for requests
func (c *Client) Debug() {
	c.restClient.SetDebug(true)
}

// newRequest creates a new request
func (c *Client) newRequest(ctx context.Context) *resty.Request {

	if c.restClient == nil {
		c.restClient = resty.NewWithClient(http.DefaultClient)
	}

	return c.restClient.
		// Set base url per request
		SetHostURL(c.Server.String()).
		// Set redirect policy based on host name
		SetRedirectPolicy(resty.DomainCheckRedirectPolicy(c.Server.Hostname())).
		// Setup error handling for non <= 399 codes
		OnAfterResponse(handleResponse).
		R().
		SetContext(ctx).
		SetHeader("UserAgent", userAgent)
}

// handleResponse handles 400+ http error codes
func handleResponse(i *resty.Client, response *resty.Response) error {
	if response.StatusCode() < 400 {
		return nil
	}

	return Error{
		Message: fmt.Sprintf("%s %s: %s", response.Request.Method, response.Request.URL, response.Status()),
		Code:    response.StatusCode(),
	}
}
