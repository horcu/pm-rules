package main

import (
	"fmt"
	m "github.com/horcu/mafia-models"
)

type RulesEngine struct {
	game *m.Game
}

func New(game *m.Game) (*RulesEngine, error) {
	return &RulesEngine{game: game}, nil
}

func (re *RulesEngine) GetAllowedAbilities(playerBin string, currentStep m.Step) ([]*m.Ability, error) {
	// 1. Find the current step

	for _, s := range re.game.Steps {
		if s.Bin == currentStep.Bin {
			currentStep = *s
			break
		}
	}
	if currentStep.Bin == "" {
		return nil, fmt.Errorf("invalid step: %s", currentStep.Bin)
	}

	// 2. Find the player's abilities
	var gamer *m.Gamer
	for _, c := range re.game.Gamers {
		if c.Bin == playerBin {
			gamer = c
			break
		}
	}
	if gamer == nil {
		return nil, fmt.Errorf("invalid player: %s", playerBin)
	}

	// 3. Determine allowed abilities
	var allowedAbilities []*m.Ability
	for _, ability := range gamer.Abilities {
		// Check if the ability is allowed in the current step
		isAllowedInStep := false
		for _, allowed := range currentStep.Allowed {
			if allowed == ability.Name {
				isAllowedInStep = true
				break
			}
		}

		if isAllowedInStep {
			// Check ability frequency
			switch ability.Frequency {
			case "every":
				allowedAbilities = append(allowedAbilities, ability)
			case "once":
				if ability.TimesUsed > -1 && ability.TimesUsed < 1 {
					allowedAbilities = append(allowedAbilities, ability)
				}
			case "twice":
				if ability.TimesUsed > -1 && ability.TimesUsed < 2 {
					allowedAbilities = append(allowedAbilities, ability)
				}
			case "every_other":
				if ability.TimesUsed == 0 {
					allowedAbilities = append(allowedAbilities, ability)
					ability.TimesUsed++
					ability.CycleUsedIndex = re.game.Cycles
				} else {
					if ability.CycleUsedIndex%2 == 0 {
						if re.game.Cycles%2 == 0 {
							allowedAbilities = append(allowedAbilities, ability)
						}
					} else {
						if 1%re.game.Cycles == 0 {
							allowedAbilities = append(allowedAbilities, ability)
						}
					}
				}
			default:
				return nil, fmt.Errorf("invalid ability frequency: %s", ability.Frequency)
			}
		}
	}

	return allowedAbilities, nil
}
