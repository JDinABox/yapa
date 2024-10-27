package main

import (
	"fmt"

	"github.com/JDinABox/yapa/internal/app/config"
)

func main() {
	config := config.New(config.FromEnv(), config.FromFlags())
	fmt.Println(config)
}
