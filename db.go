package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
)

func storeManagers() error {
	file, err := os.Open("managers.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, record := range records {
		email := record[0]
		name := record[1]

		var count int
		checkStatement := `SELECT COUNT(email) FROM public.managers WHERE email = $1`
		err = db.QueryRow(checkStatement, email).Scan(&count)
		if err != nil {
			return err
		}

		if count > 0 {
			fmt.Printf("Manager with email %s already exists in the database\n", email)
			continue
		}

		_, err = db.Exec(`INSERT INTO public.managers (email, name) VALUES ($1, $2)`, email, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func getManagerNameByEmail(email string) (string, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var name string
	err = db.QueryRow(`SELECT name FROM managers WHERE email = $1`, email).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			// No manager with the given email found
			return "", nil
		} else {
			// A database error occurred
			return "", err
		}
	}

	return name, nil
}
