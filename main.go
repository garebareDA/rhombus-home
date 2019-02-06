package main

import (
	"log"
	"net/http"
	"regexp"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
)

func main(){
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket){
		log.Println("on connection")
		so.Join("chat")

		so.On("message", func(msg string){
			log.Print(msg)

			postMessage := regexp.MustCompile(msg)

			if postMessage.MatchString("って何") {
				so.Emit("message", "hoge")

			}else if postMessage.MatchString("こんにちは") {
				so.Emit("message", "こんにちは")

			}else {
				so.Emit("message", "すみませんその機能はありません")
			}

			log.Println(msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "home",
	})
})

	r.GET("/socket.io/", gin.WrapH(server))
	r.Static("/public/", "./public/")

	r.Run()
}