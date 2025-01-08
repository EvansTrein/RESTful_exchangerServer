package main

import (
	"fmt"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
)

var conf *config.Config

func init() {
	conf = config.MustLoad()
}

func main() {
	fmt.Println(conf)
}