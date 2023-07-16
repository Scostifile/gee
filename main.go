package main

import (
	"gee"
	"net/http"
)

/*
(1)
$ curl "http://localhost:9999"
Hello Timo

(2)
$ curl "http://localhost:9999/panic"
{"message":"Internal Server Error"}

(3)
$ curl "http://localhost:9999"
Hello Timo
*/

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Timo\n")
	})

	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"timo"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
