package dyn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// ObjSession Session
	ObjSession = "Session"

	// JobStatusIncomplete incomplete
	JobStatusIncomplete = "incomplete"

	// JobStatusSuccess sucess
	JobStatusSuccess = "success"

	// JobStatusFailure failure
	JobStatusFailure = "failure"
)

// Client struct
type Client struct {
	Token        string
	CustomerName string
	Transport    http.RoundTripper
	URL          string
}

// NewClient creates a new Dyn Client
func NewClient(customerName string) *Client {
	return &Client{
		CustomerName: customerName,
		Transport:    &http.Transport{Proxy: http.ProxyFromEnvironment},
		URL:          "https://api.dynect.net/REST",
	}
}

func (c *Client) isLoggedIn() bool {
	return c.Token != ""
}

// Login creates a new session
// TODO make short for Session().Login()
func (c *Client) Login(username, password string) error {
	return c.Session().Login(username, password)
}

// Logout is a shorthand method to Session().Logout()
func (c *Client) Logout() error {
	return c.Session().Logout()
}

// Do is a raw function to call the dyn API, this should not be used by
// external packages but is exported for convience
func (c *Client) Do(method, object, indentifier string, request, response interface{}) error {
	if c.Token == "" && object != ObjSession && method != http.MethodPost {
		return errors.New("Refusing request without Login")
	}

	var payload []byte

	if request != nil {
		var err error
		payload, err = json.Marshal(request)
		if err != nil {
			return err
		}
	} else {
		payload = []byte("")
	}

	url := strings.Join([]string{c.URL, object, indentifier}, "/")
	log.Debugf("requesting url: %s", url)

	req, err := c.newRequest(method, url, payload)
	if err != nil {
		return err
	}

	resp, err := c.Transport.RoundTrip(req)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if resp.ContentLength == 0 {
			return nil
		}

		body, err := ioutil.ReadAll(resp.Body)
		log.Debugf("%s", body)

		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}

		return nil
	case http.StatusTemporaryRedirect:
		location := resp.Header.Get("Location")

		req, err := c.newRequest("GET", location, nil)
		if err != nil {
			return err
		}

		var jobResponse JobResponse

		for {
			select {
			case <-time.After(time.Second):
				resp, err := c.Transport.RoundTrip(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				if err := json.Unmarshal(body, &jobResponse); err != nil {
					return err
				}

				switch jobResponse.Status {
				case JobStatusIncomplete:
					continue
				case JobStatusSuccess:
					if err := json.Unmarshal(body, &response); err != nil {
						return err
					}
					return nil
				case JobStatusFailure:
					return fmt.Errorf("request failed: %v", jobResponse.Messages)

				}
			}
		}
	case http.StatusTooManyRequests:
		return errors.New("Too Many Requests")
	}

	reason, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return fmt.Errorf("server responded with %s: %s", resp.Status, reason)
}

func (c *Client) newRequest(method, url string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Auth-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
