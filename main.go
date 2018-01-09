package main

import (
	"io/ioutil"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/someanon/yamc/server"
	"github.com/someanon/yamc/store"
	"gopkg.in/yaml.v2"
)

func main() {

	var args struct {
		AccountsPath   string        `arg:"--accounts-path" help:"accounts file path"`
		CleaningPeriod time.Duration `arg:"--cleaning-period" help:"store cleaning period, must be >= 100ms"`
		DumpingPeriod  time.Duration `arg:"--dumping-period" help:"store dumping period, must be >= 60s"`
		DumpPath       string        `arg:"--dump-path" help:"store dump file path"`
	}

	args.AccountsPath = "./accounts"
	args.CleaningPeriod = 60 * time.Second
	args.DumpingPeriod = 60 * time.Second
	args.DumpPath = "./dump"

	arg.MustParse(&args)

	accountsYAML, err := ioutil.ReadFile(args.AccountsPath)
	if err != nil {
		panic("failed to read accounts file: " + err.Error())
	}

	var a gin.Accounts
	if err := yaml.Unmarshal(accountsYAML, &a); err != nil {
		panic("failed to parse accounts YAML: " + err.Error())
	}

	p := store.Params{
		CleaningPeriod: args.CleaningPeriod,
		DumpingPeriod:  args.DumpingPeriod,
	}

	s, err := store.NewStore(p, store.SystemClock{}, store.FileDumper(args.DumpPath))
	if err != nil {
		panic("unexpected store.NewStore() error: " + err.Error())
	}

	if err := s.StartCleaning(); err != nil {
		panic("unexpected store.Store.StartCleaning() error: " + err.Error())
	}

	if err := s.StartDumping(); err != nil {
		panic("unexpected store.Store.StartDumping() error: " + err.Error())
	}

	gin.SetMode("release")

	r := server.NewRouter(a, s)

	if err := r.Run(); err != nil {
		panic("failed to gin.Engine.Run(): " + err.Error())
	}
}
