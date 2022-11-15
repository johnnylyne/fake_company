package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/johnnylyne/fake_company/database"
	"log"
)

type Branch struct {
	Name string `json:"name"`
	Address string `json:"address"`
}

type Department struct {
	Name     string `json:"name"`
	BranchName string `json:"branch"`
}

type Employee struct {
	Forename     string `json:"forename"`
	Surname     string `json:"surname"`
	Department string `json:"department"`
}

func retrieveBranchesHandler(w http.ResponseWriter, r *http.Request) {
	var branches []Branch

	rows, rowErr := database.DB.Query("SELECT br.name, br.address FROM main.branch br")
	if rowErr != nil {
		log.Fatalf("could not query data: %v", rowErr)
	}

	for rows.Next() {
		branch := Branch{}
		
		if err := rows.Scan(&branch.Name, &branch.Address); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		branches = append(branches, branch)
	}
	json.NewEncoder(w).Encode(branches)
}

func retrieveDepartmentsHandler(w http.ResponseWriter, r *http.Request) {
	var departments []Department

	parameters, success := r.URL.Query()["branch_id"]
	var branchId int

	if success && len(parameters) > 0 {
		branchParameter := parameters[0]
		var err error
		branchId, err = strconv.Atoi(branchParameter)
		if err != nil {
			log.Fatalf("Error parsing provided branchId")
		}
		rows, rowErr := database.DB.Query("SELECT dep.name, br.name FROM main.branch br JOIN main.department dep ON br.id = dep.branch_id WHERE br.id = $1", branchId)
		if rowErr != nil {
			log.Fatalf("could not query data: %v", rowErr)
		}

		defer rows.Close()
	
		for rows.Next() {
			department := Department{}
			
			if err := rows.Scan(&department.Name, &department.BranchName); err != nil {
				log.Fatalf("could not scan row: %v", err)
			}
			departments = append(departments, department)
		}
	} else {
		rows, rowErr := database.DB.Query("SELECT dep.name, br.address FROM main.branch br JOIN main.department dep ON br.id = dep.branch_id")
		if rowErr != nil {
			log.Fatalf("could not query data: %v", rowErr)
		}

		for rows.Next() {
			department := Department{}
		
			if err := rows.Scan(&department.Name, &department.BranchName); err != nil {
				log.Fatalf("could not scan row: %v", err)
			}
			departments = append(departments, department)
		}
	}
	json.NewEncoder(w).Encode(departments)
}

func retrieveEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	var employees []Employee

	rows, rowErr := database.DB.Query("SELECT emp.forename, emp.surname, dep.name FROM main.department dep JOIN main.employee emp ON dep.id = emp.department_id")
	if rowErr != nil {
		log.Fatalf("could not query data: %v", rowErr)
	}

	for rows.Next() {
		employee := Employee{}
		
		if err := rows.Scan(&employee.Forename, &employee.Surname, &employee.Department); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		employees = append(employees, employee)
	}
	json.NewEncoder(w).Encode(employees)
}


func createBranchHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Id int `json:"id"`
	}

	response := Response{}

	err := r.ParseForm()

	// In case of any error, we respond with an error to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	branch := Branch{}

	branch.Name = r.Form.Get("name")
	branch.Address = r.Form.Get("address")
	
	tx, transactionErr := database.DB.Begin()
	if transactionErr != nil {
		log.Fatalf("Error starting database transaction")
	}


	query := "INSERT INTO main.branch(name, address) VALUES ($1, $2) RETURNING ID"
	stmt, statementErr := tx.Prepare(query)
	if statementErr != nil {
		log.Fatalf("could not formulate statement: %v", query)
	}

	insertErr := stmt.QueryRow(
		&branch.Name,
		&branch.Address,
	).Scan(&response.Id)
	
	if insertErr != nil {
		log.Fatalf("could not insert data: %v", branch)
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		log.Fatalf("Error commiting transaction")
	}
	log.Printf("Successfully committed %v", branch)

	http.Redirect(w, r, "/assets/", http.StatusFound)
}
