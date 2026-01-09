package abilities

import (
	"fmt"
	"log"
)

/*
ability: &abilities
  minAbilityValue: 1
  maxAbilityValue: 10
  rules: ["ability MUST BE IN RANGE minAbilityValue <= ability <= maxAbilityValue", "if item adds value more than 10 to ability then set it to 10 and show warning('Нельзя повысить способность более 10')", "All ability values are derived from a fixed base value and a shared ability point budget. The total sum of all positive and negative modifications must exactly match the available ability point"]
  defaults:
    valueOfEachAbility: 5
    abilityPoints: 5
  strength: integer
    description:"Define how much character can hit, how much can carry or lift"
  luck: integer
    description:"Define the percent of positive random encounter(help character find profit of any situation)"
  charisma: integer
    description:"Define how efficently character can speak to NPC(persuade etc.)"
  agility: integer
    description:"Define the ability to avoid enemy hits, perform hard atacks etc."
  perception: integer
    description:"Define the ability to feel the enemy nearby or to find paths and solutions"
  intelligence: integer
    description:"Define how much character can understand and act as main mana ability"

*/

const (
	MinAbilityValue     = 1
	MaxAbilityValue     = 10
	DefaultAbilityValue = 5
	AbilityPointBudget  = 5
)

type Abilities struct {
	pointsPool   int //counter for ability points spent by character creator UI
	strength     int
	luck         int
	charisma     int
	agility      int
	perception   int
	intelligence int
}

// NewDefaultAbilities creates an Abilities instance with all default values
func NewDefaultAbilities() Abilities {
	return Abilities{
		pointsPool:   AbilityPointBudget,
		strength:     DefaultAbilityValue,
		luck:         DefaultAbilityValue,
		charisma:     DefaultAbilityValue,
		agility:      DefaultAbilityValue,
		perception:   DefaultAbilityValue,
		intelligence: DefaultAbilityValue,
	}
}

// NewAbilities creates an Abilities instance with validation
func NewAbilities(strength int, luck int, charisma int, agility int, perception int, intelligence int) (Abilities, error) {
	// Validate each ability is in range
	abilities := []struct {
		name  string
		value int
	}{
		{"strength", strength},
		{"luck", luck},
		{"charisma", charisma},
		{"agility", agility},
		{"perception", perception},
		{"intelligence", intelligence},
	}

	for _, ability := range abilities {
		if ability.value < MinAbilityValue || ability.value > MaxAbilityValue {
			return Abilities{}, fmt.Errorf("ability %s value %d must be in range [%d, %d]",
				ability.name, ability.value, MinAbilityValue, MaxAbilityValue)
		}
	}

	// Calculate total sum of abilities
	totalAbilitySum := strength + luck + charisma + agility + perception + intelligence
	expectedSum := (6 * DefaultAbilityValue) + AbilityPointBudget

	if totalAbilitySum != expectedSum {
		return Abilities{}, fmt.Errorf("total ability points (%d) must equal %d (6×%d base + %d bonus points)",
			totalAbilitySum, expectedSum, DefaultAbilityValue, AbilityPointBudget)
	}

	// Calculate remaining points in pool
	pointsSpent := (strength - DefaultAbilityValue) + (luck - DefaultAbilityValue) +
		(charisma - DefaultAbilityValue) + (agility - DefaultAbilityValue) +
		(perception - DefaultAbilityValue) + (intelligence - DefaultAbilityValue)
	remainingPoints := AbilityPointBudget - pointsSpent

	return Abilities{
		pointsPool:   remainingPoints,
		strength:     strength,
		luck:         luck,
		charisma:     charisma,
		agility:      agility,
		perception:   perception,
		intelligence: intelligence,
	}, nil
}

// AddToAbility adds value to a specific ability using pointsPool for tracking
func (a *Abilities) AddToAbility(abilityName string, value int) error {
	getCurrentValue := func() int {
		switch abilityName {
		case "strength":
			return a.strength
		case "luck":
			return a.luck
		case "charisma":
			return a.charisma
		case "agility":
			return a.agility
		case "perception":
			return a.perception
		case "intelligence":
			return a.intelligence
		default:
			return 0
		}
	}

	currentValue := getCurrentValue()
	newValue := currentValue + value

	// Validate range
	if newValue < MinAbilityValue {
		return fmt.Errorf("cannot decrease %s below minimum (%d)", abilityName, MinAbilityValue)
	}
	if newValue > MaxAbilityValue {
		log.Printf("cannot increase ability more than 10")
		return fmt.Errorf("cannot increase %s above maximum (%d)", abilityName, MaxAbilityValue)
	}

	// Calculate point cost (relative to default value of 5)
	currentCost := currentValue - DefaultAbilityValue
	newCost := newValue - DefaultAbilityValue
	pointDelta := newCost - currentCost

	// Check if we have enough points in pool
	if pointDelta > 0 && a.pointsPool < pointDelta {
		return fmt.Errorf("insufficient points in pool: need %d, have %d", pointDelta, a.pointsPool)
	}

	// Update the ability and pointsPool
	switch abilityName {
	case "strength":
		a.strength = newValue
	case "luck":
		a.luck = newValue
	case "charisma":
		a.charisma = newValue
	case "agility":
		a.agility = newValue
	case "perception":
		a.perception = newValue
	case "intelligence":
		a.intelligence = newValue
	default:
		return fmt.Errorf("unknown ability: %s", abilityName)
	}

	// Update points pool (if value decreased, points return to pool)
	a.pointsPool -= pointDelta
	log.Printf("Updated %s: %d -> %d (points pool: %d)", abilityName, currentValue, newValue, a.pointsPool)

	return nil
}

// SetAbility sets a specific ability value using pointsPool for tracking
func (a *Abilities) SetAbility(abilityName string, value int) error {
	if value < MinAbilityValue {
		return fmt.Errorf("cannot set %s below minimum (%d)", abilityName, MinAbilityValue)
	}
	if value > MaxAbilityValue {
		log.Printf("You can't set ability more than 10")
		return fmt.Errorf("cannot set %s above maximum (%d)", abilityName, MaxAbilityValue)
	}

	getCurrentValue := func() int {
		switch abilityName {
		case "strength":
			return a.strength
		case "luck":
			return a.luck
		case "charisma":
			return a.charisma
		case "agility":
			return a.agility
		case "perception":
			return a.perception
		case "intelligence":
			return a.intelligence
		default:
			return 0
		}
	}

	currentValue := getCurrentValue()

	// Calculate point cost change
	currentCost := currentValue - DefaultAbilityValue
	newCost := value - DefaultAbilityValue
	pointDelta := newCost - currentCost

	// Check if we have enough points
	if pointDelta > 0 && a.pointsPool < pointDelta {
		return fmt.Errorf("insufficient points in pool: need %d, have %d", pointDelta, a.pointsPool)
	}

	// Update the ability
	switch abilityName {
	case "strength":
		a.strength = value
	case "luck":
		a.luck = value
	case "charisma":
		a.charisma = value
	case "agility":
		a.agility = value
	case "perception":
		a.perception = value
	case "intelligence":
		a.intelligence = value
	default:
		return fmt.Errorf("unknown ability: %s", abilityName)
	}

	// Update points pool
	a.pointsPool -= pointDelta
	log.Printf("Set %s to %d (points pool: %d)", abilityName, value, a.pointsPool)

	return nil
}

// Getter methods for individual abilities
func (a *Abilities) GetStrength() int {
	return a.strength
}

func (a *Abilities) GetLuck() int {
	return a.luck
}

func (a *Abilities) GetCharisma() int {
	return a.charisma
}

func (a *Abilities) GetAgility() int {
	return a.agility
}

func (a *Abilities) GetPerception() int {
	return a.perception
}

func (a *Abilities) GetIntelligence() int {
	return a.intelligence
}

func (a *Abilities) GetAllAbilities() map[string]int {
	return map[string]int{
		"strength":     a.strength,
		"luck":         a.luck,
		"charisma":     a.charisma,
		"agility":      a.agility,
		"perception":   a.perception,
		"intelligence": a.intelligence,
	}
}

// String returns a string representation of all abilities
func (a *Abilities) String() string {
	log.Printf("Abilities: Strength=%d, Luck=%d, Charisma=%d, Agility=%d, Perception=%d, Intelligence=%d",
		a.strength, a.luck, a.charisma, a.agility, a.perception, a.intelligence)
	return fmt.Sprintf("Strength: %d, Luck: %d, Charisma: %d, Agility: %d, Perception: %d, Intelligence: %d",
		a.strength, a.luck, a.charisma, a.agility, a.perception, a.intelligence)
}

func (a *Abilities) GetPointsPool() int {
	return a.pointsPool
}

func (a *Abilities) ValidateAbilities() error {
	log.Println("Validating abilities")
	abilities := []struct {
		name  string
		value int
	}{
		{"strength", a.strength},
		{"luck", a.luck},
		{"charisma", a.charisma},
		{"agility", a.agility},
		{"perception", a.perception},
		{"intelligence", a.intelligence},
	}

	for _, ability := range abilities {
		if ability.value < MinAbilityValue || ability.value > MaxAbilityValue {
			errMsg := fmt.Sprintf("ability %s value %d must be in range [%d, %d]",
				ability.name, ability.value, MinAbilityValue, MaxAbilityValue)
			log.Println(errMsg)
			return fmt.Errorf(errMsg, nil)
		}
	}
	log.Println("All abilities are valid")
	return nil
}
