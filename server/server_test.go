package server

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/someanon/yams/store"
)

var _ = Describe("Router", func() {
	var (
		s      *testStore
		r      *gin.Engine
		method string
		path   string
		res    *httptest.ResponseRecorder
	)
	req := func(params ...string) *http.Request {
		query := path
		if len(params) > 0 {
			query += "?" + strings.Join(params, "&")
		}
		return httptest.NewRequest(method, query, nil)
	}
	body := func(body string) io.ReadCloser {
		return ioutil.NopCloser(strings.NewReader(body))
	}
	BeforeEach(func() {
		s = &testStore{}
		r = NewRouter(s)
		res = httptest.NewRecorder()
	})
	Describe("getKey", func() {
		BeforeEach(func() {
			method = http.MethodGet
			path = "/key"
		})
		Specify("no key query param error", func() {
			r.ServeHTTP(res, req())
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("store key not exists error", func() {
			s.error = store.ErrKeyNotExists
			r.ServeHTTP(res, req("key=a"))
			s.expectGet("a")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusNotFound))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("store not key item error", func() {
			s.error = store.ErrNotKeyItem
			r.ServeHTTP(res, req("key=a"))
			s.expectGet("a")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusNotFound))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("other store error", func() {
			s.error = errors.New("error")
			r.ServeHTTP(res, req("key=a"))
			s.expectGet("a")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success", func() {
			s.value = "v"
			r.ServeHTTP(res, req("key=a"))
			s.expectGet("a")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal("v"))
		})
		Specify("success with empty key", func() {
			s.value = "v"
			r.ServeHTTP(res, req("key="))
			s.expectGet("")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal("v"))
		})
	})
	Describe("putKey", func() {
		BeforeEach(func() {
			method = http.MethodPut
			path = "/key"
		})
		Specify("no key query param error", func() {
			r.ServeHTTP(res, req())
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("no ttl query param error", func() {
			r.ServeHTTP(res, req("key=a"))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("empty ttl error", func() {
			r.ServeHTTP(res, req("key=a", "ttl="))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("ttl parse error", func() {
			r.ServeHTTP(res, req("key=a", "ttl=asd"))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("body read error", func() {
			rq := req("key=a", "ttl=10s")
			rq.Body = bodyReadErr("read error")
			r.ServeHTTP(res, rq)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success with empty key", func() {
			rq := req("key=", "ttl=10s")
			rq.Body = body("v")
			r.ServeHTTP(res, rq)
			s.expectSet("", "v", 10*time.Second)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success with empty body", func() {
			r.ServeHTTP(res, req("key=a", "ttl=10s"))
			s.expectSet("a", "", 10*time.Second)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success with body", func() {
			rq := req("key=a", "ttl=10s")
			rq.Body = body("v")
			r.ServeHTTP(res, rq)
			s.expectSet("a", "v", 10*time.Second)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.Len()).To(BeZero())
		})
	})
	Describe("getList", func() {
		BeforeEach(func() {
			method = http.MethodGet
			path = "/list"
		})
		Specify("no key query param error", func() {
			r.ServeHTTP(res, req())
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("no index query param error", func() {
			r.ServeHTTP(res, req("key=a"))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("index parse error", func() {
			r.ServeHTTP(res, req("key=a", "index=asd"))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("store key not exists error", func() {
			s.error = store.ErrKeyNotExists
			r.ServeHTTP(res, req("key=a", "index=0"))
			s.expectListGet("a", 0)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusNotFound))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("store not list item error", func() {
			s.error = store.ErrNotListItem
			r.ServeHTTP(res, req("key=a", "index=0"))
			s.expectListGet("a", 0)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusNotFound))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("store list index not exists error", func() {
			s.error = store.ErrListIndexNotExists
			r.ServeHTTP(res, req("key=a", "index=1"))
			s.expectListGet("a", 1)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusNotFound))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("store invalid list index error", func() {
			s.error = store.ErrInvalidListIndex
			r.ServeHTTP(res, req("key=a", "index=-1"))
			s.expectListGet("a", -1)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("other store error", func() {
			s.error = errors.New("error")
			r.ServeHTTP(res, req("key=a", "index=5"))
			s.expectListGet("a", 5)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success", func() {
			s.value = "v"
			r.ServeHTTP(res, req("key=a", "index=10"))
			s.expectListGet("a", 10)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal("v"))
		})
		Specify("success with empty key", func() {
			s.value = "v"
			r.ServeHTTP(res, req("key=", "index=10"))
			s.expectListGet("", 10)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal("v"))
		})
	})
	Describe("putList", func() {
		BeforeEach(func() {
			method = http.MethodPut
			path = "/list"
		})
		Specify("no key query param error", func() {
			r.ServeHTTP(res, req())
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("no ttl query param error", func() {
			r.ServeHTTP(res, req("key=a"))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("empty ttl error", func() {
			r.ServeHTTP(res, req("key=a", "ttl="))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("ttl parse error", func() {
			r.ServeHTTP(res, req("key=a", "ttl=asd"))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("body read error", func() {
			rq := req("key=a", "ttl=10s")
			rq.Body = bodyReadErr("read error")
			r.ServeHTTP(res, rq)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("body YAML list parse error", func() {
			rq := req("key=", "ttl=10s")
			rq.Body = body(`"`)
			r.ServeHTTP(res, rq)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("body YAML list parse error", func() {
			rq := req("key=", "ttl=10s")
			rq.Body = body(`asd`)
			r.ServeHTTP(res, rq)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("body YAML list parse error", func() {
			rq := req("key=", "ttl=10s")
			rq.Body = body(`a: a\nb: b`)
			r.ServeHTTP(res, rq)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success with empty key", func() {
			rq := req("key=", "ttl=10s")
			rq.Body = body("- a\n- b")
			r.ServeHTTP(res, rq)
			s.expectListSet("", []string{"a", "b"}, 10*time.Second)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success with empty body", func() {
			rq := req("key=a", "ttl=10s")
			r.ServeHTTP(res, rq)
			s.expectListSet("a", nil, 10*time.Second)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success with body", func() {
			rq := req("key=a", "ttl=10s")
			rq.Body = body("- a\n- b")
			r.ServeHTTP(res, rq)
			s.expectListSet("a", []string{"a", "b"}, 10*time.Second)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.Len()).To(BeZero())
		})
	})
	Describe("getDict", func() {
		BeforeEach(func() {
			method = http.MethodGet
			path = "/dict"
		})
		Specify("no key query param error", func() {
			r.ServeHTTP(res, req())
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("no dkey query param error", func() {
			r.ServeHTTP(res, req("key=a"))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("store key not exists error", func() {
			s.error = store.ErrKeyNotExists
			r.ServeHTTP(res, req("key=a", "dkey=b"))
			s.expectDictGet("a", "b")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusNotFound))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("store not dict item error", func() {
			s.error = store.ErrNotDictItem
			r.ServeHTTP(res, req("key=a", "dkey=b"))
			s.expectDictGet("a", "b")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusNotFound))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("store dict key not exists error", func() {
			s.error = store.ErrDictKeyNotExists
			r.ServeHTTP(res, req("key=a", "dkey=b"))
			s.expectDictGet("a", "b")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusNotFound))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("other store error", func() {
			s.error = errors.New("error")
			r.ServeHTTP(res, req("key=a", "dkey=b"))
			s.expectDictGet("a", "b")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success", func() {
			s.value = "v"
			r.ServeHTTP(res, req("key=a", "dkey=b"))
			s.expectDictGet("a", "b")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal("v"))
		})
		Specify("success with empty key", func() {
			s.value = "v"
			r.ServeHTTP(res, req("key=", "dkey=b"))
			s.expectDictGet("", "b")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal("v"))
		})
		Specify("success with empty dkey", func() {
			s.value = "v"
			r.ServeHTTP(res, req("key=a", "dkey="))
			s.expectDictGet("a", "")
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal("v"))
		})
	})
	Describe("putDict", func() {
		BeforeEach(func() {
			method = http.MethodPut
			path = "/dict"
		})
		Specify("no key query param error", func() {
			r.ServeHTTP(res, req())
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("no ttl query param error", func() {
			r.ServeHTTP(res, req("key=a"))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("empty ttl error", func() {
			r.ServeHTTP(res, req("key=a", "ttl="))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("ttl parse error", func() {
			r.ServeHTTP(res, req("key=a", "ttl=asd"))
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("body read error", func() {
			rq := req("key=a", "ttl=10s")
			rq.Body = bodyReadErr("read error")
			r.ServeHTTP(res, rq)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusInternalServerError))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("body YAML dict parse error", func() {
			rq := req("key=", "ttl=10s")
			rq.Body = body(`"`)
			r.ServeHTTP(res, rq)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("body YAML dict parse error", func() {
			rq := req("key=", "ttl=10s")
			rq.Body = body(`asd`)
			r.ServeHTTP(res, rq)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("body YAML dict parse error", func() {
			rq := req("key=", "ttl=10s")
			rq.Body = body(`- a\n- b`)
			r.ServeHTTP(res, rq)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success with empty key", func() {
			rq := req("key=", "ttl=10s")
			rq.Body = body("a: a\nb: b")
			r.ServeHTTP(res, rq)
			s.expectDictSet("", map[string]string{"a": "a", "b": "b"}, 10*time.Second)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success with empty body", func() {
			rq := req("key=a", "ttl=10s")
			r.ServeHTTP(res, rq)
			s.expectDictSet("a", nil, 10*time.Second)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.Len()).To(BeZero())
		})
		Specify("success with body", func() {
			rq := req("key=a", "ttl=10s")
			rq.Body = body("a: a\nb: b")
			r.ServeHTTP(res, rq)
			s.expectDictSet("a", map[string]string{"a": "a", "b": "b"}, 10*time.Second)
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.Len()).To(BeZero())
		})
	})
	Describe("delete", func() {
		BeforeEach(func() {
			method = http.MethodDelete
		})
		Describe("key", func() {
			BeforeEach(func() {
				path = "/key"
			})
			Specify("no key query param error", func() {
				r.ServeHTTP(res, req())
				s.expectNoCalls()
				Expect(res.Code).To(Equal(http.StatusBadRequest))
				Expect(res.Body.Len()).To(BeZero())
			})
			Specify("success", func() {
				r.ServeHTTP(res, req("key=a"))
				s.expectRemove("a")
				s.expectNoCalls()
				Expect(res.Code).To(Equal(http.StatusOK))
				Expect(res.Body.Len()).To(BeZero())
			})
		})
		Describe("list", func() {
			BeforeEach(func() {
				path = "/list"
			})
			Specify("no key query param error", func() {
				r.ServeHTTP(res, req())
				s.expectNoCalls()
				Expect(res.Code).To(Equal(http.StatusBadRequest))
				Expect(res.Body.Len()).To(BeZero())
			})
			Specify("success", func() {
				r.ServeHTTP(res, req("key=a"))
				s.expectRemove("a")
				s.expectNoCalls()
				Expect(res.Code).To(Equal(http.StatusOK))
				Expect(res.Body.Len()).To(BeZero())
			})
		})
		Describe("dict", func() {
			BeforeEach(func() {
				path = "/dict"
			})
			Specify("no key query param error", func() {
				r.ServeHTTP(res, req())
				s.expectNoCalls()
				Expect(res.Code).To(Equal(http.StatusBadRequest))
				Expect(res.Body.Len()).To(BeZero())
			})
			Specify("success", func() {
				r.ServeHTTP(res, req("key=a"))
				s.expectRemove("a")
				s.expectNoCalls()
				Expect(res.Code).To(Equal(http.StatusOK))
				Expect(res.Body.Len()).To(BeZero())
			})
		})
	})
	Describe("getKeys", func() {
		BeforeEach(func() {
			method = http.MethodGet
			path = "/keys"
		})
		Specify("success when no keys", func() {
			r.ServeHTTP(res, req())
			s.expectKeys()
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal("[]\n"))
		})
		Specify("success with keys", func() {
			s.keys = []string{"a", "b", "c"}
			r.ServeHTTP(res, req())
			s.expectKeys()
			s.expectNoCalls()
			Expect(res.Code).To(Equal(http.StatusOK))
			Expect(res.Body.String()).To(Equal("- a\n- b\n- c\n"))
		})
	})
})

type testStore struct {
	calls []call
	value string
	keys  []string
	error error
}

func (s *testStore) newCall(f interface{}, args ...interface{}) {
	s.calls = append(s.calls, call{
		method: funcToName(f),
		args:   args,
	})
}

func (s *testStore) popCall() call {
	if len(s.calls) == 0 {
		panic("no calls are left")
	}
	var c call
	c, s.calls = s.calls[0], s.calls[1:]
	return c
}

func (s *testStore) expectNoCalls() {
	ExpectWithOffset(1, s.calls).To(BeEmpty())
}

func (s *testStore) Get(key string) (string, error) {
	s.newCall(s.Get, key)
	return s.value, s.error
}

func (s *testStore) expectGet(key string) {
	ExpectWithOffset(1, s.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, s.popCall()).To(beCall(s.Get, key))
}

func (s *testStore) Set(key string, value string, ttl time.Duration) {
	s.newCall(s.Set, key, value, ttl)
}

func (s *testStore) expectSet(key string, value string, ttl time.Duration) {
	ExpectWithOffset(1, s.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, s.popCall()).To(beCall(s.Set, key, value, ttl))
}

func (s *testStore) ListGet(key string, index int) (string, error) {
	s.newCall(s.ListGet, key, index)
	return s.value, s.error
}

func (s *testStore) expectListGet(key string, index int) {
	ExpectWithOffset(1, s.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, s.popCall()).To(beCall(s.ListGet, key, index))
}

func (s *testStore) ListSet(key string, list []string, ttl time.Duration) {
	s.newCall(s.ListSet, key, list, ttl)
}

func (s *testStore) expectListSet(key string, list []string, ttl time.Duration) {
	ExpectWithOffset(1, s.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, s.popCall()).To(beCall(s.ListSet, key, list, ttl))
}

func (s *testStore) DictGet(key string, dkey string) (string, error) {
	s.newCall(s.DictGet, key, dkey)
	return s.value, s.error
}

func (s *testStore) expectDictGet(key string, dkey string) {
	ExpectWithOffset(1, s.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, s.popCall()).To(beCall(s.DictGet, key, dkey))
}

func (s *testStore) DictSet(key string, dict map[string]string, ttl time.Duration) {
	s.newCall(s.DictSet, key, dict, ttl)
}

func (s *testStore) expectDictSet(key string, dict map[string]string, ttl time.Duration) {
	ExpectWithOffset(1, s.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, s.popCall()).To(beCall(s.DictSet, key, dict, ttl))
}

func (s *testStore) Remove(key string) error {
	s.newCall(s.Remove, key)
	return s.error
}

func (s *testStore) expectRemove(key string) {
	ExpectWithOffset(1, s.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, s.popCall()).To(beCall(s.Remove, key))
}

func (s *testStore) Keys() []string {
	s.newCall(s.Keys)
	return s.keys
}

func (s *testStore) expectKeys() {
	ExpectWithOffset(1, s.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, s.popCall()).To(beCall(s.Keys))
}

func (s *testStore) StartCleaner() error {
	s.newCall(s.StartCleaner)
	return s.error
}

func (s *testStore) StopCleaner() error {
	s.newCall(s.StopCleaner)
	return s.error
}

type call struct {
	method string
	args   []interface{}
}

func beCall(method interface{}, args ...interface{}) types.GomegaMatcher {
	methodMatcher := WithTransform(func(c interface{}) string {
		return c.(call).method
	}, BeIdenticalTo(funcToName(method)))
	argsMatcher := BeEmpty()
	if len(args) > 0 {
		argsMatcher = Equal(args)
	}
	return And(methodMatcher, WithTransform(func(c interface{}) []interface{} {
		return c.(call).args
	}, argsMatcher))
}

func funcToName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	name = strings.TrimPrefix(name, "github.com/someanon/yams/")
	name = strings.TrimSuffix(name, "-fm")
	return name
}

type bodyReadErr string

func (e bodyReadErr) Read(_ []byte) (int, error) {
	return 0, errors.New(string(e))
}

func (e bodyReadErr) Close() error {
	return errors.New(string(e))
}
