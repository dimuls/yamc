package client

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	gourl "net/url"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// query paran keys
	keyKey   = "key"
	ttlKey   = "ttl"
	dkeyKey  = "dkey"
	indexKey = "index"

	// api paths
	keyPath  = "/key"
	listPath = "/list"
	dictPath = "/dict"
	keysPath = "/keys"
)

var (
	// server side errors
	ErrNotFound              = errors.New("key not found")
	ErrInvalidParams         = errors.New("invalid params")
	ErrInternalServerError   = errors.New("internal server error")
	ErrUnknownResponseStatus = errors.New("unknown response status")
	ErrInvalidServerResponse = errors.New("invalid server response")
)

// Client is a memory cache server client
type Client struct {
	method string
	url    *gourl.URL
	query  gourl.Values
}

// NewClient constructs memory cache server client
func NewClient(url string) (Client, error) {
	c := Client{}
	var err error
	if c.url, err = gourl.Parse(url); err != nil {
		return c, errors.New("failed to parse url: " + err.Error())
	}
	if c.query, err = gourl.ParseQuery(c.url.RawQuery); err != nil {
		return c, errors.New("failed to parse query params: " + err.Error())
	}
	return c, nil
}

// doReq performs request according Client data and given body
func (c Client) doReq(body []byte) (string, error) {
	req, err := http.NewRequest(c.method, c.url.String(), bytes.NewReader(body))
	if err != nil {
		return "", errors.New("failed to create http doReq: " + err.Error())
	}
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", errors.New("failed to do http doReq: " + err.Error())
	}
	switch res.StatusCode {
	case http.StatusOK:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", errors.New("failed to read response body: " + err.Error())
		}
		if res.Body.Close(); err != nil {
			return "", errors.New("failed to close response body: " + err.Error())
		}
		return string(body), nil
	case http.StatusNotFound:
		return "", ErrNotFound
	case http.StatusBadRequest:
		return "", ErrInvalidParams
	case http.StatusInternalServerError:
		return "", ErrInternalServerError
	}
	return "", ErrUnknownResponseStatus
}

// Get gets value by key
func (c Client) Get(key string) (string, error) {
	c.method = http.MethodGet
	c.url.Path = keyPath
	c.query.Set(keyKey, key)
	c.url.RawQuery = c.query.Encode()
	return c.doReq(nil)
}

// Set sets key to value with time to live ttl
func (c Client) Set(key string, value string, ttl time.Duration) error {
	c.method = http.MethodPut
	c.url.Path = keyPath
	c.query.Set(keyKey, key)
	c.query.Set(ttlKey, ttl.String())
	c.url.RawQuery = c.query.Encode()
	_, err := c.doReq([]byte(value))
	return err
}

// ListGet gets value by key and index
func (c Client) ListGet(key string, index uint) (string, error) {
	c.method = http.MethodGet
	c.url.Path = listPath
	c.query.Set(keyKey, key)
	c.query.Set(indexKey, fmt.Sprintf("%d", index))
	c.url.RawQuery = c.query.Encode()
	return c.doReq(nil)
}

// ListSet sets string list to the key
func (c Client) ListSet(key string, list []string, ttl time.Duration) error {
	c.method = http.MethodPut
	c.url.Path = listPath
	c.query.Set(keyKey, key)
	c.query.Set(ttlKey, ttl.String())
	c.url.RawQuery = c.query.Encode()
	listYAML, err := yaml.Marshal(list)
	if err != nil {
		return errors.New("failed to marshal list: " + err.Error())
	}
	_, err = c.doReq(listYAML)
	return err
}

// DictGet returns value by key and dkey
func (c Client) DictGet(key string, dkey string) (string, error) {
	c.method = http.MethodGet
	c.url.Path = dictPath
	c.query.Set(keyKey, key)
	c.query.Set(dkeyKey, dkey)
	c.url.RawQuery = c.query.Encode()
	return c.doReq(nil)
}

// DictSet sets string dict to the key
func (c Client) DictSet(key string, dict map[string]string, ttl time.Duration) error {
	c.method = http.MethodPut
	c.url.Path = dictPath
	c.query.Set(keyKey, key)
	c.query.Set(ttlKey, ttl.String())
	c.url.RawQuery = c.query.Encode()
	dictYAML, err := yaml.Marshal(dict)
	if err != nil {
		return errors.New("failed to marshal dict: " + err.Error())
	}
	_, err = c.doReq(dictYAML)
	return err
}

// Remove removes value by key
func (c Client) Remove(key string) error {
	c.method = http.MethodDelete
	c.url.Path = keyPath
	c.query.Set(keyKey, key)
	c.url.RawQuery = c.query.Encode()
	_, err := c.doReq(nil)
	return err
}

// Keys returns all keys list
func (c Client) Keys() ([]string, error) {
	c.method = http.MethodGet
	c.url.Path = keysPath
	c.url.RawQuery = c.query.Encode()
	body, err := c.doReq(nil)
	if err != nil {
		return nil, err
	}
	var keys []string
	if err := yaml.Unmarshal([]byte(body), &keys); err != nil {
		return nil, ErrInvalidServerResponse
	}
	return keys, nil
}
