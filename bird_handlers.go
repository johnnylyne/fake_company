package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Response struct {
	Counter string `json:"counter"`
}

var counter int

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit dummy handler")
	counter++
	fmt.Println(strconv.Itoa(counter))

	http.Redirect(w, r, "/assets/", http.StatusFound)
}

func returnCounterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Returning counter")
	responseObj := Response{Counter: strconv.Itoa(counter)}
	response, _ := json.Marshal(responseObj)
	w.Write(response)
}