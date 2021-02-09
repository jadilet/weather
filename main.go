package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jadilet/weather/handler"
)

func main() {
	l := log.New(os.Stdout, "weather", log.LstdFlags)
	tmpl := template.Must(template.ParseFiles("template/index.html"))

	sm := mux.NewRouter()
	weathers := handler.NewWeathers(l, tmpl)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", weathers.GetIndexPage)
	getRouter.Handle("/static/{rest}", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))

	server := http.Server{
		Addr:         ":8080",
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		Handler:      sm,
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Received terminate, gracefully shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	defer cancel()
	server.Shutdown(tc)
}
