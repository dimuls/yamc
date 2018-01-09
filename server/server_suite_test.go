package server

import (
	"io/ioutil"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	gin.SetMode("test")
	gin.DefaultWriter = ioutil.Discard
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}
