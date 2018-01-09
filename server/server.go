package server

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/someanon/yamc/store"
	"gopkg.in/yaml.v2"
)

// NewRouter creates gin router with server with binded routes and handlers
func NewRouter(a gin.Accounts, st store.Store) *gin.Engine {
	s := &server{store: st}

	r := gin.Default()

	ar := r.Group("/", gin.BasicAuth(a))

	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	ar.GET("/key", s.getKey)
	ar.PUT("/key", s.putKey)
	ar.DELETE("/key", s.delete)

	ar.GET("/list", s.getList)
	ar.PUT("/list", s.putList)
	ar.DELETE("/list", s.delete)

	ar.GET("/dict", s.getDict)
	ar.PUT("/dict", s.putDict)
	ar.DELETE("/dict", s.delete)

	ar.GET("/keys", s.getKeys)

	return r
}

// server is memory cache server
type server struct {
	store store.Store
}

// getKey handles GET /key request. This request corresponds to store's Get method. Required params: key.
// Returns value corresponded to key
func (s *server) getKey(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	value, err := s.store.Get(key)
	if err != nil {
		switch err {
		case store.ErrKeyNotExists, store.ErrNotKeyItem:
			c.AbortWithStatus(http.StatusNotFound)
		default:
			c.AbortWithError(http.StatusInternalServerError, errStoreError.causedBy(err))
		}
		return
	}
	c.String(http.StatusOK, value)
}

// putKey handles PUT /key request. This request corresponds to store's Set method. Required params: key, ttl
func (s *server) putKey(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	ttlStr, exists := c.GetQuery("ttl")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errTTLRequired)
		return
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidTTL.causedBy(err))
		return
	}
	valueBts, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errFailToReadAllBody.causedBy(err))
		return
	}
	s.store.Set(key, string(valueBts), ttl)
	c.Status(http.StatusOK)
}

// getList handles GET /list request. This request corresponds to store's ListGet method. Required params: key, index.
// Returns value corresponded to key and index
func (s *server) getList(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	indexStr, exists := c.GetQuery("index")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errIndexRequired)
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidIndex.causedBy(err))
		return
	}
	value, err := s.store.ListGet(key, index)
	if err != nil {
		switch err {
		case store.ErrKeyNotExists, store.ErrNotListItem, store.ErrListIndexNotExists:
			c.AbortWithStatus(http.StatusNotFound)
		case store.ErrInvalidListIndex:
			c.AbortWithError(http.StatusBadRequest, errStoreError.causedBy(err))
		default:
			c.AbortWithError(http.StatusInternalServerError, errStoreError.causedBy(err))
		}
	}
	c.String(http.StatusOK, value)
}

// putKey handles PUT /list request. This request corresponds to store's ListSet method. Required params: key, ttl
// and YAML formatted list in body
func (s *server) putList(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	ttlStr, exists := c.GetQuery("ttl")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errTTLRequired)
		return
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidTTL.causedBy(err))
		return
	}
	listYAML, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errFailToReadAllBody.causedBy(err))
		return
	}
	var list []string
	if err := yaml.Unmarshal(listYAML, &list); err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidListYAML.causedBy(err))
		return
	}
	s.store.ListSet(key, list, ttl)
	c.Status(http.StatusOK)
}

// getDict handles GET /dict request. This request corresponds to store's DictGet method. Required params: key, dkey.
// Returns value corresponded to key and dkey
func (s *server) getDict(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	dkey, exists := c.GetQuery("dkey")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errDKeyRequired)
		return
	}
	value, err := s.store.DictGet(key, dkey)
	if err != nil {
		switch err {
		case store.ErrKeyNotExists, store.ErrNotDictItem, store.ErrDictKeyNotExists:
			c.AbortWithStatus(http.StatusNotFound)
		default:
			c.AbortWithError(http.StatusInternalServerError, errStoreError.causedBy(err))
		}
	}
	c.String(http.StatusOK, value)
}

// putDict handles PUT /dict request. This request corresponds to store's DictSet method. Required params: key, ttl
// and YAML formatted dict in body
func (s *server) putDict(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	ttlStr, exists := c.GetQuery("ttl")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errTTLRequired)
		return
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidTTL.causedBy(err))
		return
	}
	dictYAML, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errFailToReadAllBody.causedBy(err))
		return
	}
	var dict map[string]string
	if err := yaml.Unmarshal(dictYAML, &dict); err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidDictYAML.causedBy(err))
		return
	}
	s.store.DictSet(key, dict, ttl)
	c.Status(http.StatusOK)
}

// delete handles DELETE /key, DELETE /list, DELETE /dict requests. This request corresponds store's Remove method.
// Required params: key
func (s *server) delete(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	s.store.Remove(key)
	c.Status(http.StatusOK)
}

// getKeys handles GET /keys request. This request corresponds to store's Keys method.
// Returns YAML formatted body with keys list
func (s *server) getKeys(c *gin.Context) {
	keysBytes, err := yaml.Marshal(s.store.Keys())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, "%s", keysBytes)
}
