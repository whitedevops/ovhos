package ovhos

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"path"
	"strings"
)

// Errors
var (
	ErrRequest = errors.New("ovhos: request failed")
)

// Client is an OVH Object Storage client.
// All fields are required for a successful connection.
type Client struct {
	Region       string // Region must be "BHS1", "GRA1" or "SBG1" (according to the container region).
	Container    string // Container is the name of the targetted container.
	TenantID     string // TenantID is the "AUTH_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" section of the container URL, but without the "AUTH_" part.
	Username     string // Username is an OpenStack username.
	Password     string // Password is the OpenStack password for user.
	currentToken token  // currentToken contains the token used for requests.
}

// request makes a new client request to the OVH Object Storage.
func (c *Client) request(method, object string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, c.URL(object), body)
	if err != nil {
		return nil, err
	}

	t, err := c.token()
	if err != nil {
		return nil, err
	}
	r.Header.Set("X-Auth-Token", t)

	return (&http.Client{}).Do(r)
}

// URL returns the full object address.
func (c *Client) URL(object string) string {
	return "https://" + path.Join("storage."+strings.ToLower(c.Region)+".cloud.ovh.net/v1/AUTH_"+c.TenantID+"/"+c.Container, object)
}

// get returns the response of a GET request.
//
// CURL equivalent:
// 	curl https://storage.$REGION.cloud.ovh.net/v1/AUTH_$TENANTID/$CONTAINER -X GET -H "X-Auth-Token: $TOKEN"
func (c *Client) get() (r *http.Response, err error) {
	r, err = c.request("GET", "", nil)
	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusNoContent {
		err = ErrRequest
	}
	return
}

// Ping checks if the connection is OK for the client credentials.
func (c *Client) Ping() (err error) {
	_, err = c.get()
	return
}

// List returns a slice of all objects in the container.
//
// CURL equivalent:
// 	curl https://storage.$REGION.cloud.ovh.net/v1/AUTH_$TENANTID/$CONTAINER -X GET -H "X-Auth-Token: $TOKEN"
func (c *Client) List() ([]string, error) {
	r, err := c.get()
	if err != nil {
		return nil, err
	}

	var l []string
	s := bufio.NewScanner(r.Body)
	for s.Scan() {
		l = append(l, s.Text())
	}

	return l, s.Err()
}

// Exists checks if the object exists in the container.
func (c *Client) Exists(object string) (bool, error) {
	r, err := c.request("HEAD", object, nil)
	if err != nil {
		return false, err
	}
	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusNotFound {
		return false, ErrRequest
	}

	return r.StatusCode == http.StatusOK, nil
}

// Upload puts a new object in the container.
//
// CURL equivalent:
// 	curl https://storage.$REGION.cloud.ovh.net/v1/AUTH_$TENANTID/$CONTAINER/$OBJECT -X PUT -H "X-Auth-Token: $TOKEN"
func (c *Client) Upload(object string, body io.Reader) error {
	r, err := c.request("PUT", object, body)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusCreated {
		return ErrRequest
	}
	return nil
}

// Delete removes an object from the container.
//
// CURL equivalent:
// 	curl https://storage.$REGION.cloud.ovh.net/v1/AUTH_$TENANTID/$CONTAINER/$OBJECT -X DELETE -H "X-Auth-Token: $TOKEN"
func (c *Client) Delete(object string) error {
	r, err := c.request("DELETE", object, nil)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusNoContent && r.StatusCode != http.StatusNotFound {
		return ErrRequest
	}
	return nil
}
