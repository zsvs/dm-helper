package main

import (
	// abts "dnd-helper/src/abilities"
	char "dnd-helper/src/character"
	// cond "dnd-helper/src/condition"
	// inv "dnd-helper/src/inventory"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func mockSendDbRequest(data any) error {
	// Simulate sending data to a database
	log.Printf("Mock sending data to DB: %v", data)
	return nil
}

func withRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("%s %s in %s", r.Method, r.URL.Path, time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("panic: %v\n%s", x, debug.Stack())
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	var characters []char.Character
	mux := http.NewServeMux()
	handler := withRecovery(withRequestLogging(mux))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	mux.HandleFunc("/create-character", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var charReq []char.Character

		// Parse JSON request body directly into Character structs
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&charReq); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var responseData []map[string]interface{}

		// Process each character
		for _, character := range charReq {
			// Validate character
			if err := character.ValidateCharacter(); err != nil {
				http.Error(w, fmt.Sprintf("Invalid character: %v", err), http.StatusBadRequest)
				return
			}

			// Add to characters list
			characters = append(characters, character)

			// Prepare response for this character
			responseData = append(responseData, map[string]interface{}{
				"name":       character.Name,
				"race":       character.Race,
				"class":      character.Class,
				"abilities":  character.Abilities,
				"manaPoints": character.ManaPoints,
				"condition":  character.Condition,
				"inventory":  character.Inventory,
			})

			// Mock sending character data to a database
			charObj, err := json.MarshalIndent(character, "", "  ")
			if err != nil {
				log.Printf("Error marshaling character data: %v", err)
				continue
			}
			mockSendDbRequest(string(charObj))
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":    fmt.Sprintf("Successfully created %d character(s)", len(charReq)),
			"count":      len(charReq),
			"characters": responseData,
		})
	})

	mux.HandleFunc("/get-chars", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		log.Printf("Returning %d characters", len(characters))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"count":      len(characters),
			"characters": characters,
		})
	})
	log.Println("Starting server")
	log.Println("Listen on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
