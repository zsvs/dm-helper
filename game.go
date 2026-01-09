package main

import (
	abts "dnd-helper/src/abilities"
	char "dnd-helper/src/character"
	cond "dnd-helper/src/condition"
	inv "dnd-helper/src/inventory"
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

		// Define request structure matching character structure
		type ItemDTO struct {
			Name        string `json:"name"`
			Quantity    int    `json:"quantity"`
			Condition   string `json:"condition"`
			Description string `json:"description"`
			Abilities   *struct {
				Strength     int `json:"strength"`
				Luck         int `json:"luck"`
				Charisma     int `json:"charisma"`
				Agility      int `json:"agility"`
				Perception   int `json:"perception"`
				Intelligence int `json:"intelligence"`
			} `json:"abilities,omitempty"`
		}

		type CreateCharacterRequest struct {
			Race      string `json:"race"`
			Name      string `json:"name"`
			Class     string `json:"class"`
			Inventory struct {
				Items []ItemDTO `json:"items"`
			} `json:"inventory"`
			Abilities struct {
				Strength     int `json:"strength"`
				Luck         int `json:"luck"`
				Charisma     int `json:"charisma"`
				Agility      int `json:"agility"`
				Perception   int `json:"perception"`
				Intelligence int `json:"intelligence"`
			} `json:"abilities"`
			Condition string `json:"condition"`
		}

		var charReq []CreateCharacterRequest

		// Parse JSON request body
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&charReq); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Create each character from request data
		for _, req := range charReq {
			abilities, err := abts.NewAbilities(
				req.Abilities.Strength,
				req.Abilities.Luck,
				req.Abilities.Charisma,
				req.Abilities.Agility,
				req.Abilities.Perception,
				req.Abilities.Intelligence,
			)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid abilities: %v", err), http.StatusBadRequest)
				return
			}

			// Create inventory and add items
			inventory := inv.NewInventory()
			for _, itemDTO := range req.Inventory.Items {
				var itemAbilities *abts.Abilities
				if itemDTO.Abilities != nil {
					itemAbs, err := abts.NewAbilities(
						itemDTO.Abilities.Strength,
						itemDTO.Abilities.Luck,
						itemDTO.Abilities.Charisma,
						itemDTO.Abilities.Agility,
						itemDTO.Abilities.Perception,
						itemDTO.Abilities.Intelligence,
					)
					if err != nil {
						http.Error(w, fmt.Sprintf("Invalid item abilities: %v", err), http.StatusBadRequest)
						return
					}
					itemAbilities = &itemAbs
				}

				item, err := inv.NewItem(
					itemDTO.Name,
					itemDTO.Quantity,
					itemAbilities,
					cond.NewCondition(itemDTO.Condition),
					itemDTO.Description,
				)
				if err != nil {
					http.Error(w, fmt.Sprintf("Invalid item: %v", err), http.StatusBadRequest)
					return
				}
				inventory.AddItem(item)
			}

			// Create condition and character
			condition := cond.NewCondition(req.Condition)
			character := char.NewCharacter(req.Race, req.Name, req.Class, abilities, *inventory, condition)
			characters = append(characters, *character)
			// Get abilities and inventory
			charAbilities := character.GetAbilities()
			charInventory := character.GetInventory()

			// Prepare inventory response (only public fields will serialize)
			inventoryItems := []map[string]interface{}{}
			for _, item := range charInventory.GetAllItems() {
				itemData := map[string]interface{}{
					"name": item.Name,
				}
				inventoryItems = append(inventoryItems, itemData)
			}

			// Prepare response data once
			responseData := map[string]interface{}{
				"message": "Character created successfully",
				"character": map[string]interface{}{
					"name":       character.GetName(),
					"race":       character.GetRace(),
					"class":      character.GetClass(),
					"abilities":  charAbilities.GetAllAbilities(),
					"manaPoints": character.GetManaPoints(),
					"condition":  character.GetCondition().String(),
					"inventory": map[string]interface{}{
						"items": inventoryItems,
					},
				},
			}

			// Return success response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(responseData)
			// Mock sending character data to a database
			charObj, err := json.MarshalIndent(responseData, "", "  ")
			if err != nil {
				log.Printf("Error marshaling character data: %v", err)
				return
			}
			mockSendDbRequest(string(charObj))
		}
	})

	mux.HandleFunc("/get-chars", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var responseData []map[string]interface{}

		for _, character := range characters {
			// Get character data
			charAbilities := character.GetAbilities()
			charInventory := character.GetInventory()

			// Prepare inventory items
			for _, item := range charInventory.GetAllItems() {
				// inventoryItems = append(inventoryItems, itemData)
				// Add character data to response
				responseData = append(responseData, map[string]interface{}{
					"name":       character.GetName(),
					"race":       character.GetRace(),
					"class":      character.GetClass(),
					"abilities":  charAbilities.GetAllAbilities(),
					"manaPoints": character.GetManaPoints(),
					"condition":  character.GetCondition().String(),
					"inventory": map[string]interface{}{
						"items": map[string]interface{}{
							"name":        item.Name,
							"qantity":     item.GetQuantity(),
							"condition":   item.GetCondition().String(),
							"description": item.GetDescription(),
						},
					},
				})
			}

		}

		log.Printf("Returning %d characters", len(characters))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"count":      len(characters),
			"characters": responseData,
		})
	})
	log.Println("Starting server")
	log.Println("Listen on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
