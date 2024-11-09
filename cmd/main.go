package main

import (
	"fmt"
	"url_shortener_svc/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
