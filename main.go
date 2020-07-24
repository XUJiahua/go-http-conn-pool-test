package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

var m = make(map[string]int)
var ch = make(chan string, 10)

func count() {
	for s := range ch {
		m[s]++
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	logrus.Info(r.RemoteAddr)
	ch <- r.RemoteAddr
	// time.Sleep(time.Second)
	w.Write([]byte("helloworld"))
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func graceClose() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(ch)
		time.Sleep(time.Second)
		spew.Dump(m)
		os.Exit(0)
	}()
}

func main() {
	graceClose()
	go count()
	port := flag.Int("port", 8087, "")
	flag.Parse()

	logrus.Println("Listen port:", *port)

	http.HandleFunc("/", home)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}
}
