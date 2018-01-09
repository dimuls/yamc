package client_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/someanon/yamc/client"
)

var _ = Describe("Client", func() {
	var (
		s  *testServer
		ts *httptest.Server
		c  Client
	)
	BeforeEach(func() {
		s = &testServer{}
		ts = httptest.NewServer(s)
		var err error
		c, err = NewClient(ts.URL, "tlogin", "tpassword")
		Expect(err).ToNot(HaveOccurred())
	})
	Describe("Get", func() {
		Specify("unknown response status", func() {
			s.status = http.StatusCreated
			_, err := c.Get("a")
			Expect(err).To(MatchError(ErrUnknownResponseStatus))
			s.expReq(http.MethodGet, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("internal server error", func() {
			s.status = http.StatusInternalServerError
			_, err := c.Get("a")
			Expect(err).To(MatchError(ErrInternalServerError))
			s.expReq(http.MethodGet, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("invalid params error", func() {
			s.status = http.StatusBadRequest
			_, err := c.Get("a")
			Expect(err).To(MatchError(ErrInvalidParams))
			s.expReq(http.MethodGet, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("not found error", func() {
			s.status = http.StatusNotFound
			_, err := c.Get("a")
			Expect(err).To(MatchError(ErrNotFound))
			s.expReq(http.MethodGet, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("unauthorized error", func() {
			s.status = http.StatusUnauthorized
			_, err := c.Get("a")
			Expect(err).To(MatchError(ErrUnauthorized))
			s.expReq(http.MethodGet, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("succeed", func() {
			s.status = http.StatusOK
			s.body = "v"
			v, err := c.Get("a")
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal("v"))
			s.expReq(http.MethodGet, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
	})
	Describe("Set", func() {
		Specify("unknown response status", func() {
			s.status = http.StatusCreated
			Expect(c.Set("a", "v", 10*time.Second)).To(MatchError(ErrUnknownResponseStatus))
			s.expReq(http.MethodPut, "/key", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "v")
			s.expNoReq()
		})
		Specify("internal server error", func() {
			s.status = http.StatusInternalServerError
			Expect(c.Set("a", "v", 10*time.Second)).To(MatchError(ErrInternalServerError))
			s.expReq(http.MethodPut, "/key", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "v")
			s.expNoReq()
		})
		Specify("invalid params error", func() {
			s.status = http.StatusBadRequest
			Expect(c.Set("a", "v", 10*time.Second)).To(MatchError(ErrInvalidParams))
			s.expReq(http.MethodPut, "/key", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "v")
			s.expNoReq()
		})
		Specify("not found error", func() {
			s.status = http.StatusNotFound
			Expect(c.Set("a", "v", 10*time.Second)).To(MatchError(ErrNotFound))
			s.expReq(http.MethodPut, "/key", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "v")
			s.expNoReq()
		})
		Specify("unauthorized error", func() {
			s.status = http.StatusUnauthorized
			Expect(c.Set("a", "v", 10*time.Second)).To(MatchError(ErrUnauthorized))
			s.expReq(http.MethodPut, "/key", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "v")
			s.expNoReq()
		})
		Specify("succeed", func() {
			s.status = http.StatusOK
			Expect(c.Set("a", "v", 10*time.Second)).ToNot(HaveOccurred())
			s.expReq(http.MethodPut, "/key", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "v")
			s.expNoReq()
		})
	})
	Describe("ListGet", func() {
		Specify("unknown response status", func() {
			s.status = http.StatusCreated
			_, err := c.ListGet("a", 0)
			Expect(err).To(MatchError(ErrUnknownResponseStatus))
			s.expReq(http.MethodGet, "/list", "tlogin", "tpassword", []string{"key=a", "index=0"}, "")
			s.expNoReq()
		})
		Specify("internal server error", func() {
			s.status = http.StatusInternalServerError
			_, err := c.ListGet("a", 0)
			Expect(err).To(MatchError(ErrInternalServerError))
			s.expReq(http.MethodGet, "/list", "tlogin", "tpassword", []string{"key=a", "index=0"}, "")
			s.expNoReq()
		})
		Specify("invalid params error", func() {
			s.status = http.StatusBadRequest
			_, err := c.ListGet("a", 0)
			Expect(err).To(MatchError(ErrInvalidParams))
			s.expReq(http.MethodGet, "/list", "tlogin", "tpassword", []string{"key=a", "index=0"}, "")
			s.expNoReq()
		})
		Specify("not found error", func() {
			s.status = http.StatusNotFound
			_, err := c.ListGet("a", 0)
			Expect(err).To(MatchError(ErrNotFound))
			s.expReq(http.MethodGet, "/list", "tlogin", "tpassword", []string{"key=a", "index=0"}, "")
			s.expNoReq()
		})
		Specify("unauthorized error", func() {
			s.status = http.StatusUnauthorized
			_, err := c.ListGet("a", 0)
			Expect(err).To(MatchError(ErrUnauthorized))
			s.expReq(http.MethodGet, "/list", "tlogin", "tpassword", []string{"key=a", "index=0"}, "")
			s.expNoReq()
		})
		Specify("succeed", func() {
			s.status = http.StatusOK
			s.body = "v"
			v, err := c.ListGet("a", 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal("v"))
			s.expReq(http.MethodGet, "/list", "tlogin", "tpassword", []string{"key=a", "index=0"}, "")
			s.expNoReq()
		})
	})
	Describe("ListSet", func() {
		Specify("unknown response status", func() {
			s.status = http.StatusCreated
			Expect(c.ListSet("a", []string{"a", "b", "c"}, 10*time.Second)).To(MatchError(ErrUnknownResponseStatus))
			s.expReq(http.MethodPut, "/list", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "- a\n- b\n- c\n")
			s.expNoReq()
		})
		Specify("internal server error", func() {
			s.status = http.StatusInternalServerError
			Expect(c.ListSet("a", []string{"a", "b", "c"}, 10*time.Second)).To(MatchError(ErrInternalServerError))
			s.expReq(http.MethodPut, "/list", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "- a\n- b\n- c\n")
			s.expNoReq()
		})
		Specify("invalid params error", func() {
			s.status = http.StatusBadRequest
			Expect(c.ListSet("a", []string{"a", "b", "c"}, 10*time.Second)).To(MatchError(ErrInvalidParams))
			s.expReq(http.MethodPut, "/list", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "- a\n- b\n- c\n")
			s.expNoReq()
		})
		Specify("not found error", func() {
			s.status = http.StatusNotFound
			Expect(c.ListSet("a", []string{"a", "b", "c"}, 10*time.Second)).To(MatchError(ErrNotFound))
			s.expReq(http.MethodPut, "/list", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "- a\n- b\n- c\n")
			s.expNoReq()
		})
		Specify("unauthorized error", func() {
			s.status = http.StatusUnauthorized
			Expect(c.ListSet("a", []string{"a", "b", "c"}, 10*time.Second)).To(MatchError(ErrUnauthorized))
			s.expReq(http.MethodPut, "/list", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "- a\n- b\n- c\n")
			s.expNoReq()
		})
		Specify("succeed", func() {
			s.status = http.StatusOK
			Expect(c.ListSet("a", []string{"a", "b", "c"}, 10*time.Second)).ToNot(HaveOccurred())
			s.expReq(http.MethodPut, "/list", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "- a\n- b\n- c\n")
			s.expNoReq()
		})
	})
	Describe("DictGet", func() {
		Specify("unknown response status", func() {
			s.status = http.StatusCreated
			_, err := c.DictGet("a", "b")
			Expect(err).To(MatchError(ErrUnknownResponseStatus))
			s.expReq(http.MethodGet, "/dict", "tlogin", "tpassword", []string{"key=a", "dkey=b"}, "")
			s.expNoReq()
		})
		Specify("internal server error", func() {
			s.status = http.StatusInternalServerError
			_, err := c.DictGet("a", "b")
			Expect(err).To(MatchError(ErrInternalServerError))
			s.expReq(http.MethodGet, "/dict", "tlogin", "tpassword", []string{"key=a", "dkey=b"}, "")
			s.expNoReq()
		})
		Specify("invalid params error", func() {
			s.status = http.StatusBadRequest
			_, err := c.DictGet("a", "b")
			Expect(err).To(MatchError(ErrInvalidParams))
			s.expReq(http.MethodGet, "/dict", "tlogin", "tpassword", []string{"key=a", "dkey=b"}, "")
			s.expNoReq()
		})
		Specify("not found error", func() {
			s.status = http.StatusNotFound
			_, err := c.DictGet("a", "b")
			Expect(err).To(MatchError(ErrNotFound))
			s.expReq(http.MethodGet, "/dict", "tlogin", "tpassword", []string{"key=a", "dkey=b"}, "")
			s.expNoReq()
		})
		Specify("unauthorized error", func() {
			s.status = http.StatusUnauthorized
			_, err := c.DictGet("a", "b")
			Expect(err).To(MatchError(ErrUnauthorized))
			s.expReq(http.MethodGet, "/dict", "tlogin", "tpassword", []string{"key=a", "dkey=b"}, "")
			s.expNoReq()
		})
		Specify("succeed", func() {
			s.status = http.StatusOK
			s.body = "v"
			v, err := c.DictGet("a", "b")
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal("v"))
			s.expReq(http.MethodGet, "/dict", "tlogin", "tpassword", []string{"key=a", "dkey=b"}, "")
			s.expNoReq()
		})
	})
	Describe("DictSet", func() {
		Specify("unknown response status", func() {
			s.status = http.StatusCreated
			Expect(c.DictSet("a", map[string]string{"a": "a", "b": "b"}, 10*time.Second)).To(MatchError(ErrUnknownResponseStatus))
			s.expReq(http.MethodPut, "/dict", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "a: a\nb: b\n")
			s.expNoReq()
		})
		Specify("internal server error", func() {
			s.status = http.StatusInternalServerError
			Expect(c.DictSet("a", map[string]string{"a": "a", "b": "b"}, 10*time.Second)).To(MatchError(ErrInternalServerError))
			s.expReq(http.MethodPut, "/dict", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "a: a\nb: b\n")
			s.expNoReq()
		})
		Specify("invalid params error", func() {
			s.status = http.StatusBadRequest
			Expect(c.DictSet("a", map[string]string{"a": "a", "b": "b"}, 10*time.Second)).To(MatchError(ErrInvalidParams))
			s.expReq(http.MethodPut, "/dict", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "a: a\nb: b\n")
			s.expNoReq()
		})
		Specify("not found error", func() {
			s.status = http.StatusNotFound
			Expect(c.DictSet("a", map[string]string{"a": "a", "b": "b"}, 10*time.Second)).To(MatchError(ErrNotFound))
			s.expReq(http.MethodPut, "/dict", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "a: a\nb: b\n")
			s.expNoReq()
		})
		Specify("unauthorized error", func() {
			s.status = http.StatusUnauthorized
			Expect(c.DictSet("a", map[string]string{"a": "a", "b": "b"}, 10*time.Second)).To(MatchError(ErrUnauthorized))
			s.expReq(http.MethodPut, "/dict", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "a: a\nb: b\n")
			s.expNoReq()
		})
		Specify("succeed", func() {
			s.status = http.StatusOK
			Expect(c.DictSet("a", map[string]string{"a": "a", "b": "b"}, 10*time.Second)).ToNot(HaveOccurred())
			s.expReq(http.MethodPut, "/dict", "tlogin", "tpassword", []string{"key=a", "ttl=10s"}, "a: a\nb: b\n")
			s.expNoReq()
		})
	})
	Describe("Remove", func() {
		Specify("unknown response status", func() {
			s.status = http.StatusCreated
			Expect(c.Remove("a")).To(MatchError(ErrUnknownResponseStatus))
			s.expReq(http.MethodDelete, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("internal server error", func() {
			s.status = http.StatusInternalServerError
			Expect(c.Remove("a")).To(MatchError(ErrInternalServerError))
			s.expReq(http.MethodDelete, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("invalid params error", func() {
			s.status = http.StatusBadRequest
			Expect(c.Remove("a")).To(MatchError(ErrInvalidParams))
			s.expReq(http.MethodDelete, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("not found error", func() {
			s.status = http.StatusNotFound
			Expect(c.Remove("a")).To(MatchError(ErrNotFound))
			s.expReq(http.MethodDelete, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("unauthorized error", func() {
			s.status = http.StatusUnauthorized
			Expect(c.Remove("a")).To(MatchError(ErrUnauthorized))
			s.expReq(http.MethodDelete, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
		Specify("succeed", func() {
			s.status = http.StatusOK
			Expect(c.Remove("a")).ToNot(HaveOccurred())
			s.expReq(http.MethodDelete, "/key", "tlogin", "tpassword", []string{"key=a"}, "")
			s.expNoReq()
		})
	})
	Describe("Keys", func() {
		Specify("unknown response status", func() {
			s.status = http.StatusCreated
			_, err := c.Keys()
			Expect(err).To(MatchError(ErrUnknownResponseStatus))
			s.expReq(http.MethodGet, "/keys", "tlogin", "tpassword", nil, "")
			s.expNoReq()
		})
		Specify("internal server error", func() {
			s.status = http.StatusInternalServerError
			s.body = ""
			_, err := c.Keys()
			Expect(err).To(MatchError(ErrInternalServerError))
			s.expReq(http.MethodGet, "/keys", "tlogin", "tpassword", nil, "")
			s.expNoReq()
		})
		Specify("invalid params error", func() {
			s.status = http.StatusBadRequest
			_, err := c.Keys()
			Expect(err).To(MatchError(ErrInvalidParams))
			s.expReq(http.MethodGet, "/keys", "tlogin", "tpassword", nil, "")
			s.expNoReq()
		})
		Specify("not found error", func() {
			s.status = http.StatusNotFound
			_, err := c.Keys()
			Expect(err).To(MatchError(ErrNotFound))
			s.expReq(http.MethodGet, "/keys", "tlogin", "tpassword", nil, "")
			s.expNoReq()
		})
		Specify("response YAML parse error", func() {
			s.status = http.StatusOK
			s.body = "asd"
			_, err := c.Keys()
			Expect(err).To(MatchError(ErrInvalidServerResponse))
			s.expReq(http.MethodGet, "/keys", "tlogin", "tpassword", nil, "")
			s.expNoReq()
		})
		Specify("unauthorized error", func() {
			s.status = http.StatusUnauthorized
			_, err := c.Keys()
			Expect(err).To(MatchError(ErrUnauthorized))
			s.expReq(http.MethodGet, "/keys", "tlogin", "tpassword", nil, "")
			s.expNoReq()
		})
		Specify("succeed", func() {
			s.status = http.StatusOK
			s.body = "- a\n- b\n"
			v, err := c.Keys()
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal([]string{"a", "b"}))
			s.expReq(http.MethodGet, "/keys", "tlogin", "tpassword", nil, "")
			s.expNoReq()
		})
	})
})

type request struct {
	method   string
	path     string
	login    string
	password string
	query    []string
	body     string
}

type testServer struct {
	requests []request
	status   int
	body     string
}

func (s *testServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic("failed to read request body: " + err.Error())
	}
	var q []string
	if req.URL.RawQuery != "" {
		q = strings.Split(req.URL.RawQuery, "&")
		sort.Strings(q)
	}
	login, password, _ := req.BasicAuth()
	s.requests = append(s.requests, request{
		method:   req.Method,
		path:     req.URL.Path,
		login:    login,
		password: password,
		query:    q,
		body:     string(b),
	})
	res.WriteHeader(s.status)
	res.Write([]byte(s.body))
}

func (s *testServer) expReq(method string, path string, login string, password string, query []string, body string) {
	ExpectWithOffset(1, s.requests).ToNot(BeEmpty(), "requests")
	r := s.requests[0]
	s.requests = s.requests[1:]
	ExpectWithOffset(1, r.method).To(BeIdenticalTo(method), "method")
	ExpectWithOffset(1, r.path).To(BeIdenticalTo(path), "path")
	ExpectWithOffset(1, r.login).To(BeIdenticalTo(login), "login")
	ExpectWithOffset(1, r.password).To(BeIdenticalTo(password), "login")
	sort.Strings(query)
	ExpectWithOffset(1, r.query).To(Equal(query), "query params")
	ExpectWithOffset(1, r.body).To(BeIdenticalTo(body), "body")
}

func (s *testServer) expNoReq() {
	ExpectWithOffset(1, s.requests).To(BeEmpty())
}
