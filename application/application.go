package application

import (
	"fmt"
	"net"
	"net/http"

	"github.com/tpp/msf/application/handler"
	"github.com/tpp/msf/application/router"
	"github.com/tpp/msf/shared/log"
)

type Config struct {
	Port int `yaml:"port"`
}

func Run(c *Config) {
	if c == nil {
		c = &Config{Port: 8080}
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	router.SetupHandler(handler.New())

	log.Info().Msg(fmt.Sprintf("Start http server at :%d", c.Port))
	log.Fatal().Err(http.Serve(listener, router.Router)).Msg("")
}
