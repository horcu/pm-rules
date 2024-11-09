package main

import (
	"fmt"
	m "github.com/horcu/pm-models/types"
	st "github.com/horcu/pm-store"
)

type RulesEngine struct {
	game  *m.Game
	store *st.Store
}

func New(game *m.Game) (*RulesEngine, error) {
	return &RulesEngine{
		game:  game,
		store: st.NewStore(),
	}, nil
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
					ability.CycleUsedIndex = re.game.NightCycles
				} else {
					if ability.CycleUsedIndex%2 == 0 {
						if re.game.NightCycles%2 == 0 {
							allowedAbilities = append(allowedAbilities, ability)
						}
					} else {
						if re.game.NightCycles%1 == 0 {
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

func (re *RulesEngine) ApplyAbility(g *m.Game, targetGamer *m.Gamer, ability string) bool {

	switch ability {
	case "kill":
		if re.CanBeKilled(g, targetGamer) {
			re.store.ApplyAbility("kill", g.Bin, targetGamer.Bin)
			return true
		}
	case "trick":
		if re.CanBeTricked(g, targetGamer) {
			re.store.ApplyAbility("trick", g.Bin, targetGamer.Bin)
			return true
		}
	case "mimic":
		if re.CanBeMimicked(g, targetGamer) {
			re.store.ApplyAbility("mimic", g.Bin, targetGamer.Bin)
			return true
		}
	case "heal":
		if re.CanBeHealed(g, targetGamer) {
			re.store.ApplyAbility("heal", g.Bin, targetGamer.Bin)
			return true
		}
	case "poison":
		if re.CanBePoisoned(g, targetGamer) {
			re.store.ApplyAbility("poison", g.Bin, targetGamer.Bin)
			return true
		}
	case "block":
		if re.CanBeBlocked(g, targetGamer) {
			re.store.ApplyAbility("block", g.Bin, targetGamer.Bin)
			return true
		}
	default:
		return false
	}
	return false
}

func (re *RulesEngine) CanBeKilled(g *m.Game, gamer *m.Gamer) bool {
	return true
}

func (re *RulesEngine) CanBeTricked(g *m.Game, gamer *m.Gamer) bool {
	return true
}

func (re *RulesEngine) CanBeMimicked(g *m.Game, gamer *m.Gamer) bool {
	return true
}

func (re *RulesEngine) CanBeHealed(g *m.Game, gamer *m.Gamer) bool {
	return true
}

func (re *RulesEngine) CanBePoisoned(g *m.Game, gamer *m.Gamer) bool {
	return true
}

func (re *RulesEngine) CanBeBlocked(g *m.Game, gamer *m.Gamer) bool {
	return true
}
