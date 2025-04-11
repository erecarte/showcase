package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/erecarte/showcase/internal/numeral/api/payment_orders"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var users = map[string]string{
	"RECARTE": "xxxx",
}

type HttpServer struct {
	server     http.Server
	port       int64
	paymentApi *payment_orders.Api
}

func NewHttpServer(port int64, service *payment_orders.Service) *HttpServer {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}

		pwd, ok := users[username]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}
		if password != pwd {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}
		c.Next()
	})
	v1Route := r.Group("/v1")
	_, err := payment_orders.NewApi(v1Route, service)
	if err != nil {
		log.Fatalln(err)
	}

	return &HttpServer{
		port: port,
		server: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
	}
}

func (s *HttpServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

}
func (s *HttpServer) Stop() {
	if err := s.server.Shutdown(context.TODO()); err != nil { // Use here context with a required timeout
		log.Printf("server shutdown returned an err: %v\n", err)
	}
}
