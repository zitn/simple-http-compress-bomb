package main

import (
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"fmt"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

var blackHole = make([]byte, 8192)

func main() {
	r := gin.New()
	r.Any("/", handler)
	r.GET("/download", handler)
	err := r.Run(":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handler(c *gin.Context) {
	if c.Request.URL.Path == "/" {
		c.Header("content-type", "text/html; charset=UTF-8")
	}
	acceptEncodingsStr := c.GetHeader("accept-encoding")
	acceptEncodings := strings.Split(acceptEncodingsStr, ", ")
	if len(acceptEncodings) == 0 {
		return
	}
	switch acceptEncodings[0] {
	case "gzip":
		fmt.Println("gzip")
		c.Header("Content-Encoding", "gzip")
		writer, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			fmt.Println(err)
			return
		}
		for err == nil {
			_, err = writer.Write(blackHole)
		}
		fmt.Println(err)
		return
	case "deflate":
		fmt.Println("deflate")
		c.Header("Content-Encoding", "deflate")
		writer, err := flate.NewWriter(c.Writer, flate.BestSpeed)
		if err != nil {
			fmt.Println(err)
			return
		}
		for err == nil {
			_, err = writer.Write(blackHole)
		}
		fmt.Println(err)
		return
	case "br":
		fmt.Println("br")
		c.Header("Content-Encoding", "br")
		writer := brotli.NewWriterLevel(c.Writer, brotli.BestSpeed)
		var err error
		for err == nil {
			_, err = writer.Write(blackHole)
		}
		fmt.Println(err)
	case "compress":
		fmt.Println("compress")
		c.Header("Content-Encoding", "compress")
		writer := lzw.NewWriter(c.Writer, lzw.LSB, 8)
		var err error
		for err == nil {
			_, err = writer.Write(blackHole)
		}
		fmt.Println(err)
	default:
		return
	}
}
