package main

import (
	"flag"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/someanon/yamc/server"
	"github.com/someanon/yamc/store"
	"gopkg.in/yaml.v2"
)

func main() {
	accountsPath := flag.String("accounts", "./accounts", "accounts file path")
	cleaningPeriod := flag.Duration("cleaning-period", 60*time.Second, "store cleaning period")

	flag.Parse()

	accountsYAML, err := ioutil.ReadFile(*accountsPath)
	if err != nil {
		panic("failed to read accounts file: " + err.Error())
	}

	var a gin.Accounts
	if err := yaml.Unmarshal(accountsYAML, &a); err != nil {
		panic("failed to parse accounts YAML: " + err.Error())
	}

	p := store.Params{
		CleaningPeriod: *cleaningPeriod,
	}

	s, err := store.NewStore(p, store.SystemClock{})
	if err != nil {
		panic("unexpected store.NewStore() error: " + err.Error())
	}

	if err := s.StartCleaner(); err != nil {
		panic("unexpected store.Store.StartCleaner() error: " + err.Error())
	}

	gin.SetMode("release")

	r := server.NewRouter(a, s)
	r.Use(gin.Recovery())

	if err := r.Run(); err != nil {
		panic("failed to gin.Engine.Run(): " + err.Error())
	}
}
