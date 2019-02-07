package main

import (
	"log"
	"net/http"
	"regexp"
	"time"
	"strings"
	"github.com/PuerkitoBio/goquery"
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

			if regexp.MustCompile(`て何`).MatchString(msg) {
				req := strings.Replace(msg, `って何`, " ", 1)
				doc, err := goquery.NewDocument("https://ja.wikipedia.org/wiki/" + req)
				if err != nil {
					log.Fatal(err)
				}

				selection := doc.Find("p")
				text := selection.Text()

				so.Emit("message", text + ",,,wikiより引用")

			}else if postMessage.MatchString("こんにちは") {
				time.LoadLocation("Asia/Tokyo")
				hour := time.Now().Hour()

				if hour > 12 && 18 > hour {
					so.Emit("message", "こんにちは")
				}else if 18 > hour {
					so.Emit("message", "こんにちは, もう暗いですね")
				}else if hour > 0 {
					so.Emit("message", "こんにちはまだ昼じゃないですよ")
				}

			} else if postMessage.MatchString("こんばんは") {
				time.LoadLocation("Asia/Tokyo")
				hour := time.Now().Hour()

				if hour > 12 && 18 > hour {
					so.Emit("message", "こんばんは、まだ早いですね")
				}else if 18 > hour {
					so.Emit("message", "こんばんは")
				}else if hour > 0  {
					so.Emit("message", "こんばんは、まだ明るいですね")
				}

			} else if postMessage.MatchString("おはよう") {
				time.LoadLocation("Asia/Tokyo")
				hour := time.Now().Hour()

				if hour > 12 && 18 > hour {
					so.Emit("message", "おはようございます、朝ですか？")
				}else if 18 > hour {
					so.Emit("message", "おはようございます、遅いですね")
				}else if hour > 0 {
					so.Emit("message", "おはようございます")
				}

			}else {
				so.Emit("message", "すみませんその機能はありません")
			}
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