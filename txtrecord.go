package dyn

import (
	"fmt"
	"net/http"
)

// RecordFunc struct
type RecordFunc struct {
	c      *Client
	zone   *Zone
	fqdn   string
	object string
}

type AllRecordsResponse struct {
	ResponseBlock
	Data []string `json:"data"`
}

// Get a record
func (c *RecordFunc) Get(recordID string) (*RecordResponse, error) {
	result := &RecordResponse{}

	if err := c.c.Do(http.MethodGet, c.object, fmt.Sprintf("%s/%s/%s", c.zone.name, c.fqdn, recordID), nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAll result a List of all records
func (c *RecordFunc) GetAll() (*AllRecordsResponse, error) {
	result := &AllRecordsResponse{}

	if err := c.c.Do(http.MethodGet, c.object, fmt.Sprintf("%s/%s", c.zone.name, c.fqdn), nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Create a new record
func (c *RecordFunc) Create(record *Record) (*RecordResponse, error) {
	result := &RecordResponse{}

	if err := c.c.Do(http.MethodPost, c.object, fmt.Sprintf("%s/%s", c.zone.name, c.fqdn), record, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete a TXT Record
func (c *RecordFunc) Delete(recordID string) error {
	return c.c.Do(http.MethodDelete, c.object, fmt.Sprintf("%s/%s/%s", c.zone.name, c.fqdn, recordID), nil, nil)
}

//DeleteAll deletes all txt records under a zone
func (c *RecordFunc) DeleteAll() error {
	return c.c.Do(http.MethodGet, c.object, fmt.Sprintf("%s/%s", c.zone.name, c.fqdn), nil, nil)
}

// Update a TXT record
func (c *RecordFunc) Update(recordID string, record *Record) (*RecordResponse, error) {
	result := &RecordResponse{}

	if err := c.c.Do(http.MethodPut, c.object, fmt.Sprintf("%s/%s/%s", c.zone.name, c.fqdn, recordID), record, &result); err != nil {
		return nil, err
	}

	return result, nil
}
