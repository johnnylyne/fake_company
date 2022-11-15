package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/johnnylyne/fake_company/database"
	"log"
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
	json.NewEncoder(w).Encode(responseObj)
}

type Branch struct {
	ID     int `json:"id"`
	Address string `json:"address"`
}

type Department struct {
	Name     string `json:"name"`
	Branch string `json:"branch"`
}

type Employee struct {
	Forename     string `json:"forename"`
	Surname     string `json:"surname"`
	Department string `json:"address"`
}

func retrieveBranchesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(database.DB)
	var branches []Branch

	rows, rowErr := database.DB.Query("SELECT * FROM main.branch")
	if rowErr != nil {
		log.Fatalf("could not query data: %v", rowErr)
	}

	for rows.Next() {
		branch := Branch{}
		
		if err := rows.Scan(&branch.ID, &branch.Address); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		fmt.Printf("branch: %v", branch)
		// append the current instance to the slice of birds
		branches = append(branches, branch)
	}
	json.NewEncoder(w).Encode(branches)
}
