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

func NewRouter(st store.Store) *gin.Engine {
	s := &server{store: st}

	r := gin.New()

	r.GET("/key", s.getKey)
	r.PUT("/key", s.putKey)
	r.DELETE("/key", s.delete)

	r.GET("/list", s.getList)
	r.PUT("/list", s.putList)
	r.DELETE("/list", s.delete)

	r.GET("/dict", s.getDict)
	r.PUT("/dict", s.putDict)
	r.DELETE("/dict", s.delete)

	r.GET("/keys", s.getKeys)

	return r
}

type server struct {
	store store.Store
}

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
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	c.String(http.StatusOK, value)
}

func (s *server) putKey(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	ttlStr, exists := c.GetQuery("ttl")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errTTLRrequired)
		return
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidTTL)
		return
	}
	valueBts, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	s.store.Set(key, string(valueBts), ttl)
	c.Status(http.StatusOK)
}

func (s *server) getList(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	indexStr, exists := c.GetQuery("index")
	if !exists {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	value, err := s.store.ListGet(key, index)
	if err != nil {
		switch err {
		case store.ErrKeyNotExists, store.ErrNotListItem, store.ErrListIndexNotExists:
			c.AbortWithStatus(http.StatusNotFound)
		case store.ErrInvalidListIndex:
			c.AbortWithError(http.StatusBadRequest, err)
		default:
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	c.String(http.StatusOK, value)
}

func (s *server) putList(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	ttlStr, exists := c.GetQuery("ttl")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errTTLRrequired)
		return
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidTTL)
		return
	}
	listYAML, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var list []string
	if err := yaml.Unmarshal(listYAML, &list); err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidListYAML)
		return
	}
	s.store.ListSet(key, list, ttl)
	c.Status(http.StatusOK)
}

func (s *server) getDict(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	dkey, exists := c.GetQuery("dkey")
	if !exists {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	value, err := s.store.DictGet(key, dkey)
	if err != nil {
		switch err {
		case store.ErrKeyNotExists, store.ErrNotDictItem, store.ErrDictKeyNotExists:
			c.AbortWithStatus(http.StatusNotFound)
		default:
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	c.String(http.StatusOK, value)
}

func (s *server) putDict(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	ttlStr, exists := c.GetQuery("ttl")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errTTLRrequired)
		return
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidTTL)
		return
	}
	dictYAML, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var dict map[string]string
	if err := yaml.Unmarshal(dictYAML, &dict); err != nil {
		c.AbortWithError(http.StatusBadRequest, errInvalidDictYAML)
		return
	}
	s.store.DictSet(key, dict, ttl)
	c.Status(http.StatusOK)
}

func (s *server) delete(c *gin.Context) {
	key, exists := c.GetQuery("key")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errKeyRequired)
		return
	}
	s.store.Remove(key)
	c.Status(http.StatusOK)
}

func (s *server) getKeys(c *gin.Context) {
	keysBytes, err := yaml.Marshal(s.store.Keys())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, "%s", keysBytes)
}
