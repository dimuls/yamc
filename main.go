package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/someanon/yams/server"
	"github.com/someanon/yams/store"
)

func main() {
	s, err := store.NewStore(store.Params{CleaningPeriod: 10 * time.Second}, store.SystemClock{})
	if err != nil {
		panic("unexpected store.NewStore() error: " + err.Error())
	}

	if err := s.StartCleaner(); err != nil {
		panic("unexpected store.Store.StartCleaner() error: " + err.Error())
	}

	gin.SetMode("release")

	r := server.NewRouter(s)
	r.Use(gin.Recovery())

	if err := r.Run(); err != nil {
		panic("failed to gin.Engine.Run(): " + err.Error())
	}
}
