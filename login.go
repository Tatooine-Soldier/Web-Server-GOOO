package main

import (
	"fmt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println("username:", r.Form["username"])
	fmt.Println("password:", r.Form["password"])
}
