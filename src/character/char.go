package character

import (
	"dnd-helper/src/abilities"
	"dnd-helper/src/condition"
	"dnd-helper/src/inventory"
	"fmt"
	"log"
)

type Character struct {
	Race       string              `json:"race"`
	Name       string              `json:"name"`
	Class      string              `json:"class"`
	Abilities  abilities.Abilities `json:"abilities"`
	Inventory  inventory.Inventory `json:"inventory"`
	Condition  condition.Condition `json:"condition"`
	ManaPoints int                 `json:"manaPoints"`
}

func NewCharacter(race string, name string, class string, abs abilities.Abilities, inv inventory.Inventory, cond condition.Condition) *Character {
	log.Printf("Creating new character %s %s with class %s, \nabilities %v, \ninventory %v \nand in %v condition", race, name, class, abs.String(), inv.String(), cond)
	return &Character{
		Race:       race,
		Name:       name,
		Class:      class,
		Abilities:  abs,
		Inventory:  inv,
		Condition:  cond,
		ManaPoints: abs.GetIntelligence() * 50,
	}
}

func NewDefaultCharacter(race string, name string, class string) *Character {
	defaultAbilities := abilities.NewDefaultAbilities()
	defaultInventory := inventory.NewInventory()
	defaultCondition := condition.NewCondition("Healthy")
	return &Character{
		Race:       race,
		Name:       name,
		Class:      class,
		Abilities:  defaultAbilities,
		Inventory:  *defaultInventory,
		Condition:  defaultCondition,
		ManaPoints: defaultAbilities.GetIntelligence() * 50,
	}
}

func (c *Character) GetName() string {
	return c.Name
}

func (c *Character) GetRace() string {
	return c.Race
}

func (c *Character) GetClass() string {
	return c.Class
}

func (c *Character) GetAbilities() abilities.Abilities {
	return c.Abilities
}

func (c *Character) GetInventory() inventory.Inventory {
	return c.Inventory
}

func (c *Character) GetCondition() condition.Condition {
	return c.Condition
}

func (c *Character) GetManaPoints() int {
	return c.ManaPoints
}

func (c *Character) SetName(newName string) {
	if newName != "" {
		c.Name = newName
		log.Printf("Name changed to: %s", newName)
	} else {
		log.Println("Name not changed, new name is empty")
	}
}

func (c *Character) SetClass(newClass string) {
	if newClass != "" {
		c.Class = newClass
		log.Printf("Class changed to: %s", newClass)
	} else {
		log.Println("Class not changed, new class is empty")
	}

}

func (c *Character) SetCondition(newCondition condition.Condition) {
	if newCondition.String() != "" {
		c.Condition = newCondition
		log.Printf("Condition changed to: %s", newCondition.String())
	} else {
		log.Println("Condition not changed, new condition is empty")
	}
}

func (c *Character) SetInventory(newItem inventory.Item) {

	c.Inventory.AddItem(newItem)
}

func (c *Character) ValidateCharacter() error {
	log.Printf("Validating character: %s", c.Name)
	if c.Name == "" || c.Race == "" || c.Class == "" {
		errMsg := "Character validation failed: name, race, or class cannot be empty"
		log.Println(errMsg)
		return fmt.Errorf(errMsg, nil)
	}
	if err := c.Abilities.ValidateAbilities(); err != nil {
		log.Printf("Character validation failed: %v", err)
		return err
	}
	return nil
}
