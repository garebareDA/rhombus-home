package main

import (
	"log"
	"net/http"
	"regexp"
	"time"
	"strings"
	"strconv"
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
				hour, _ := getTime()

				if hour > 12 && 18 > hour {
					so.Emit("message", "こんにちは")
				}else if 18 < hour {
					so.Emit("message", "こんにちは, もう暗いですね")
				}else if hour > 0 {
					so.Emit("message", "こんにちはまだ昼じゃないですよ")
				}

			} else if postMessage.MatchString("こんばんは") {
				hour, _ := getTime()

				if hour > 12 && 18 > hour {
					so.Emit("message", "こんばんは、まだ早いですね")
				}else if 18 < hour {
					so.Emit("message", "こんばんは")
				}else if hour > 0  {
					so.Emit("message", "こんばんは、まだ明るいですね")
				}

			} else if postMessage.MatchString("おはよう") {
				hour, _:= getTime()

				if hour > 12 && 18 > hour {
					so.Emit("message", "おはようございます、朝ですか？")
				}else if 18 < hour {
					so.Emit("message", "おはようございます、遅いですね")
				}else if hour > 0 {
					so.Emit("message", "おはようございます")
				}

			}else if postMessage.MatchString("今何時"){
				hour, minute := getTime()
				so.Emit("message", "現在" + strconv.Itoa(hour) + "時" + strconv.Itoa(minute) + "分です")

			}else if postMessage.MatchString("今日は何日"){
				weeks := [...]string{"日", "月", "火", "水", "木", "金", "土"}

				now := time.Now().UTC()
				jst := time.FixedZone("Asia/Tokyo", 9*60*60)
				times := now.In(jst)

				week := weeks[times.Weekday()]
				var month = int (times.Month())
				day := times.Day()

				so.Emit("message","今日は" + strconv.Itoa(month) + "月" + strconv.Itoa(day) + "日" + week + "曜日です")
			}else{
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

func getTime() (int, int) {
	now := time.Now().UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	times := now.In(jst)
	hour := times.Hour()
	minute := times.Minute()

	return hour, minute
}