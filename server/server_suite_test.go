package server

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServer(t *testing.T) {
	gin.SetMode("test")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}
