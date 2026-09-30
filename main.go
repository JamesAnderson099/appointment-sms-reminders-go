package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	client, err := NewSMSClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/appointments/reminder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var a Appointment
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			http.Error(w, "invalid appointment", http.StatusBadRequest)
			return
		}
		id, err := SendReminder(client, a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message_id": id})
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
