package dyn

import "net/http"

// Session is used to do session operations on Dyn
type Session struct {
	c *Client
}

// Session returns a client from session API calls
func (c *Client) Session() *Session {
	return &Session{c: c}
}

// Login performs an POST on /REST/Session
func (c *Session) Login(username, password string) error {
	req := NewLoginRequest{
		CustomerName: c.c.CustomerName,
		Username:     username,
		Password:     password,
	}

	var resp SessionLoginResponse

	if err := c.c.Do(http.MethodPost, ObjSession, "", req, &resp); err != nil {
		return err
	}

	c.c.Token = resp.Data.Token

	return nil
}

// Logout performs an DELETE on /REST/Session
func (c *Session) Logout() error {
	return c.c.Do(http.MethodDelete, ObjSession, "", nil, nil)
}

// Renew performs an Put on /REST/Session
func (c *Session) Renew() error {
	return c.c.Do(http.MethodPut, ObjSession, "", nil, nil)
}

// Status performs an GET on /REST/Session
func (c *Session) Status() error {
	return c.c.Do(http.MethodGet, ObjSession, "", nil, nil)
}
