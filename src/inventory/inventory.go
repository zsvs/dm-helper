package inventory

import (
	"fmt"
	"log"

	"dnd-helper/src/abilities"
	"dnd-helper/src/condition"
)

const (
	// Item defaults
	DefaultItemName        = ""
	DefaultItemQuantity    = 0.0
	DefaultItemDescription = ""
	DefaultItemCondition   = condition.Condition("N/A")

	// Ability settings for items
	MinItemAbilityValue = 1
	MaxItemAbilityValue = 4
)

// Item represents a single item in the inventory
type Item struct {
	Name        string
	quantity    int
	abilities   *abilities.Abilities
	condition   condition.Condition
	description string
}

func (i *Item) SetName(name string) {
	i.Name = name
}

func (i *Item) GetName() string {
	return i.Name
}

func (i *Item) SetQuantity(quantity int) {
	i.quantity = quantity
}

func (i *Item) GetQuantity() int {
	return i.quantity
}

func (i *Item) SetAbilities(abs *abilities.Abilities) {
	i.abilities = abs
}

func (i *Item) GetAbilities() *abilities.Abilities {
	return i.abilities
}

func (i *Item) SetCondition(cond condition.Condition) {
	i.condition = cond
}

func (i *Item) GetCondition() condition.Condition {
	return i.condition
}

func (i *Item) SetDescription(description string) {
	i.description = description
}

func (i *Item) GetDescription() string {
	return i.description
}

// Inventory represents a collection of items
type Inventory struct {
	Items []Item
}

// NewItem creates a new item with validation
func NewItem(name string, quantity int, abilities *abilities.Abilities, condition condition.Condition, description string) (Item, error) {
	if quantity <= 0 {
		return Item{}, fmt.Errorf("item quantity cannot be negative or zero")
	}

	// Validate abilities if provided
	if abilities != nil {
		abs := abilities.GetAllAbilities()
		itemAbilities := []struct {
			name  string
			value int
		}{
			{"strength", abs["strength"]},
			{"luck", abs["luck"]},
			{"charisma", abs["charisma"]},
			{"agility", abs["agility"]},
			{"perception", abs["perception"]},
			{"intelligence", abs["intelligence"]},
		}

		for _, ability := range itemAbilities {
			if ability.value != 0 && (ability.value < MinItemAbilityValue || ability.value > MaxItemAbilityValue) {
				return Item{}, fmt.Errorf("item ability %s value %d must be 0 or in range [%d, %d]",
					ability.name, ability.value, MinItemAbilityValue, MaxItemAbilityValue)
			}
		}
	}

	return Item{
		Name:        name,
		quantity:    quantity,
		abilities:   abilities,
		condition:   condition,
		description: description,
	}, nil
}

// NewInventory creates a new empty inventory
func NewInventory() *Inventory {
	return &Inventory{
		Items: []Item{},
	}
}

// AddItem adds an item to the inventory
func (inv *Inventory) AddItem(item Item) {
	// Check if item with same name already exists
	for i := range inv.Items {
		if inv.Items[i].Name == item.Name && inv.Items[i].condition == item.condition {
			// Stack items by adding quantities
			inv.Items[i].quantity += item.quantity
			log.Printf("Added %d of %s to existing stack. New quantity: %d", item.quantity, item.Name, inv.Items[i].quantity)
			return
		}
	}
	// Add as new item
	inv.Items = append(inv.Items, item)
	log.Printf("Added new item: %s (quantity: %d)", item.Name, item.quantity)
}

// RemoveItem removes a specific quantity of an item from inventory
func (inv *Inventory) RemoveItem(name string, quantity int) error {
	for i := range inv.Items {
		if inv.Items[i].Name == name {
			if inv.Items[i].quantity < quantity {
				return fmt.Errorf("insufficient quantity: have %d, need %d", inv.Items[i].quantity, quantity)
			}
			inv.Items[i].quantity -= quantity
			if inv.Items[i].quantity == 0 {
				// Remove item from inventory if quantity reaches 0
				inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
				log.Printf("Removed %s from inventory (depleted)", name)
			} else {
				log.Printf("Removed %d of %s. Remaining: %d", quantity, name, inv.Items[i].quantity)
			}
			return nil
		}
	}
	return fmt.Errorf("item %s not found in inventory", name)
}

// GetItem returns a pointer to an item by name, or nil if not found
func (inv *Inventory) GetItem(name string) *Item {
	for i := range inv.Items {
		if inv.Items[i].Name == name {
			return &inv.Items[i]
		}
	}
	return nil
}

// GetAllItems returns all items in the inventory
func (inv *Inventory) GetAllItems() []Item {
	return inv.Items
}

// HasItem checks if an item exists in the inventory with sufficient quantity
func (inv *Inventory) HasItem(name string, quantity int) bool {
	for _, item := range inv.Items {
		if item.Name == name && item.quantity >= quantity {
			return true
		}
	}
	return false
}

// ChangeItem modifies fields of an item identified by name
func (inv *Inventory) ChangeItem(name string, fields []string, newVal any) *Item {
	if len(inv.Items) == 0 {
		return nil
	}

	item := inv.GetItem(name)
	if item == nil {
		return nil
	}

	for _, field := range fields {
		switch field {
		case "name":
			if v, ok := newVal.(string); ok {
				item.SetName(v)
			}
		case "quantity":
			if v, ok := newVal.(int); ok {
				item.SetQuantity(v)
			}
		case "condition":
			if v, ok := newVal.(condition.Condition); ok {
				item.SetCondition(v)
			}
		case "description":
			if v, ok := newVal.(string); ok {
				item.SetDescription(v)
			}
		case "abilities":
			if v, ok := newVal.(*abilities.Abilities); ok {
				item.SetAbilities(v)
			}
		default:
			log.Printf("Unknown field: %s", field)
			return nil
		}
	}
	return item
}

// GetTotalWeight returns the total quantity of all items (if representing weight)
func (inv *Inventory) GetTotalWeight() int {
	total := 0
	for _, item := range inv.Items {
		total += item.quantity
	}
	return total
}

// Clear removes all items from the inventory
func (inv *Inventory) Clear() {
	inv.Items = []Item{}
	log.Printf("Inventory cleared")
}

func (inv *Inventory) String() string {
	log.Printf("Inventory contains %d items", len(inv.Items))
	result := "Inventory:\n"
	for _, item := range inv.Items {
		result += fmt.Sprintf("Name: %s, Quantity: %d, Condition: %s, Description: %s\n", item.Name, item.quantity, item.condition.String(), item.description)
	}
	result += fmt.Sprintf("Total weight: %d", inv.GetTotalWeight())
	return result
}
