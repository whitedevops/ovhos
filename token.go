package ovhos

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Errors
var (
	ErrToken = errors.New("ovhos: new token request failed")
)

// token is an OVH request token.
type token struct {
	ID     string
	Expiry time.Time
}

// token returns a working token.
// A new one is requested when the current is expired.
//
// CURL equivalent:
// 	curl https://auth.cloud.ovh.net/v2.0/tokens -X POST -H "content-type: application/json" -d '{"auth": {"passwordCredentials": {"username": "$USERNAME", "password": "$PASSWORD"}, "tenantId": "$TENANTID"}}'
func (c *Client) token() (string, error) {
	// Take the current token if not ready to expire.
	if c.currentToken.ID != "" && c.currentToken.Expiry.Before(time.Now().Truncate(5*time.Minute)) {
		return c.currentToken.ID, nil
	}

	// Prepare the request.
	req, err := http.NewRequest("POST", "https://auth.cloud.ovh.net/v2.0/tokens", strings.NewReader(fmt.Sprintf(`{"auth": {"passwordCredentials": {"username": "%s", "password": "%s"}, "tenantId": "%s"}}`, c.Username, c.Password, c.TenantID)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the request.
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", ErrToken
	}

	// Unmarshall the response.
	js := &struct {
		Access struct {
			Token struct {
				ID      string `json:"id"`
				Expires string `json:"expires"`
			} `json:"token"`
		} `json:"access"`
	}{}
	if err = json.NewDecoder(res.Body).Decode(js); err != nil {
		return "", err
	}

	// Set and return the new token.
	c.currentToken.ID = js.Access.Token.ID
	if c.currentToken.Expiry, err = time.Parse(time.RFC3339, js.Access.Token.Expires); err != nil {
		return "", err
	}

	return c.currentToken.ID, nil
}
