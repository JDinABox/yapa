package main

import (
	"fmt"
	"net/http"

	"github.com/JDinABox/yapa/internal/app/config"
	"github.com/JDinABox/yapa/internal/app/router"
	"github.com/JDinABox/yapa/internal/ilog"
)

func main() {
	config := config.New(config.FromEnv(), config.FromFlags())
	r := router.New(router.WithConfig(config))
	mux := r.Init()
	fmt.Printf("Listening on %s\n", config.Address)
	if err := http.ListenAndServe(config.Address, mux); err != nil {
		ilog.Exit(err)
		return
	}
}
