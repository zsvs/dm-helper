package character

import (
	"dnd-helper/src/abilities"
	"dnd-helper/src/condition"
	"dnd-helper/src/inventory"
	"fmt"
	"log"
)

type Character struct {
	race       string
	name       string
	class      string
	abilities  abilities.Abilities
	inventory  inventory.Inventory
	condition  condition.Condition
	manaPoints int
}

func NewCharacter(race string, name string, class string, abs abilities.Abilities, inv inventory.Inventory, cond condition.Condition) *Character {
	log.Printf("Creating new character %s %s with class %s, \nabilities %v, \ninventory %v \nand in %v condition", race, name, class, abs.String(), inv.String(), cond)
	return &Character{
		race:       race,
		name:       name,
		class:      class,
		abilities:  abs,
		inventory:  inv,
		condition:  cond,
		manaPoints: abs.GetIntelligence() * 50,
	}
}

func NewDefaultCharacter(race string, name string, class string) *Character {
	defaultAbilities := abilities.NewDefaultAbilities()
	defaultInventory := inventory.NewInventory()
	defaultCondition := condition.NewCondition("Healthy")
	return &Character{
		race:       race,
		name:       name,
		class:      class,
		abilities:  defaultAbilities,
		inventory:  *defaultInventory,
		condition:  defaultCondition,
		manaPoints: defaultAbilities.GetIntelligence() * 50,
	}
}

func (c *Character) GetName() string {
	return c.name
}

func (c *Character) GetRace() string {
	return c.race
}

func (c *Character) GetClass() string {
	return c.class
}

func (c *Character) GetAbilities() abilities.Abilities {
	return c.abilities
}

func (c *Character) GetInventory() inventory.Inventory {
	return c.inventory
}

func (c *Character) GetCondition() condition.Condition {
	return c.condition
}

func (c *Character) GetManaPoints() int {
	return c.manaPoints
}

func (c *Character) SetName(newName string) {
	if newName != "" {
		c.name = newName
		log.Printf("Name changed to: %s", newName)
	} else {
		log.Println("Name not changed, new name is empty")
	}
}

func (c *Character) SetClass(newClass string) {
	if newClass != "" {
		c.class = newClass
		log.Printf("Class changed to: %s", newClass)
	} else {
		log.Println("Class not changed, new class is empty")
	}

}

func (c *Character) SetCondition(newCondition condition.Condition) {
	if newCondition.String() != "" {
		c.condition = newCondition
		log.Printf("Condition changed to: %s", newCondition.String())
	} else {
		log.Println("Condition not changed, new condition is empty")
	}
}

func (c *Character) SetInventory(newItem inventory.Item) {

	c.inventory.AddItem(newItem)
}

func (c *Character) ValidateCharacter() error {
	log.Printf("Validating character: %s", c.name)
	if c.name == "" || c.race == "" || c.class == "" {
		errMsg := "Character validation failed: name, race, or class cannot be empty"
		log.Println(errMsg)
		return fmt.Errorf(errMsg, nil)
	}
	if err := c.abilities.ValidateAbilities(); err != nil {
		log.Printf("Character validation failed: %v", err)
		return err
	}
	return nil
}
