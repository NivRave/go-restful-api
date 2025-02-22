package main

import (
	"fmt"
	"net/http"
	"os"
	"restful_api/handlers"
	"restful_api/routers"
)

func main() {
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/users", routers.UsersRouter)
	http.HandleFunc("/users/", routers.UsersRouter)
	err := http.ListenAndServe("localhost:11111", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
