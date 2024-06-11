package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ClassData struct {
	Year       int    `json:"YEAR"`
	Semester   string `json:"SEMESTER"`
	Instructor string `json:"INSTRUCTOR"`
	Subject    string `json:"SUBJECT"`
	Number     string `json:"NUMBER"`
	CourseID   string `json:"COURSE ID"`
	A          int    `json:"A"`
	B          int    `json:"B"`
	C          int    `json:"C"`
	D          int    `json:"D"`
	F          int    `json:"F"`
	W          int    `json:"W"`
}

type InstructorData map[string]map[string]ClassData

func main() {
	http.HandleFunc("/get-data", func(w http.ResponseWriter, r *http.Request) {
		subject := r.URL.Query().Get("subject")
		number := r.URL.Query().Get("number")
		instructor := r.URL.Query().Get("instructor")

		if subject == "" || instructor == "" {
			http.Error(w, "Subject and instructor are required", http.StatusBadRequest)
			return
		}

		filePath := fmt.Sprintf("pgd/pgd_%s.json", subject)
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		var instructorData InstructorData
		if err := json.NewDecoder(file).Decode(&instructorData); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
			return
		}

		var results []ClassData
		for prof, classes := range instructorData {
			if prof == instructor {
				for _, class := range classes {
					if class.Subject == subject && (number == "" || class.Number == number) {
						results = append(results, class)
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	http.ListenAndServe(":8080", nil)
}
