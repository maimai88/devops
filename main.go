// http.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	server := &http.Server{
		Addr:         ":8888",
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/bye", sayBye)

	server.Handler = mux

	go func() {
		<-quit
		if err := server.Close(); err != nil {
			log.Fatal("Close Server", err)
		}
	}()

	log.Println("Starting server ...")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Println("Http Server Closed!")
		} else {
			log.Fatal("Http Unexcept!")
		}
	}

	log.Println("Http Run end!")
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, this is http server! " + r.URL.String()))
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bye bye, this is http server! " + r.URL.String()))

	// 获取url参数 ，可以获取同名参数数组
	vars := r.URL.Query()
	a, ok := vars["a"]
	if !ok {
		fmt.Printf("param a does not exist\n")
	} else {
		fmt.Printf("param a value is [%s]\n", a)
	}

	// 获取url参数
	b := vars.Get("a")
	fmt.Printf("param a value is [%s]\n", b)
}
