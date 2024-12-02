package pm_rules

import (
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

func (re *RulesEngine) UpdateAbilityUsage(g *m.Game, gamer *m.Gamer, ability *m.Ability) bool {

	ab := gamer.Abilities[ability.Name]

	if ab.CyclesUsedIndex == nil {
		ab.CyclesUsedIndex = []int{}
	}

	ab.CyclesUsedIndex = append(ab.CyclesUsedIndex, g.NightCycles)
	ab.TimesUsed++

	// create an update gamer map
	gamer.Abilities[ability.Name] = ab

	mp := map[string]interface{}{
		"abilities": gamer.Abilities,
	}

	//save the gamer to the store
	return re.store.UpdateGamer(g.Bin, mp)

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

	return gamer.IsAlive && !re.FateSealed(g, gamer)
}

func (re *RulesEngine) CanBeTricked(g *m.Game, gamer *m.Gamer) bool {
	return gamer.IsAlive && !re.FateSealed(g, gamer)
}

func (re *RulesEngine) CanBeMimicked(g *m.Game, gamer *m.Gamer) bool {
	return gamer.IsAlive && !re.FateSealed(g, gamer)
}

func (re *RulesEngine) CanBeHealed(g *m.Game, gamer *m.Gamer) bool {
	return gamer.IsAlive && re.GamerNeedsHealing(g, gamer)
}

func (re *RulesEngine) CanBePoisoned(g *m.Game, gamer *m.Gamer) bool {
	return gamer.IsAlive && !re.GamerWasHidden(g, gamer)
}

func (re *RulesEngine) CanBeBlocked(g *m.Game, gamer *m.Gamer) bool {
	return gamer.IsAlive
}

func (re *RulesEngine) FateSealed(game *m.Game, gamer *m.Gamer) bool {

	if gamer.Fate == nil {
		return true
	}

	// fate cant be sealed if healed or hidden
	var nonSealable []string

	nonSealable = append(nonSealable, "74668334-c229-4a8a-86d6-ff1943f34754") // heal
	nonSealable = append(nonSealable, "724b1b7f-f42e-4b4f-8393-b2ae38cb0d70") // hide once
	nonSealable = append(nonSealable, "f5b7c9e1-0a64-4e9c-9291-7e3b36df891d") // hide every other

	for _, n := range nonSealable {
		if gamer.Fate.AbilityBin == n {
			return false
		}
	}
	return true

}

func (re *RulesEngine) GamerNeedsHealing(g *m.Game, gamer *m.Gamer) bool {
	if gamer.Fate.AbilityBin == "afef47b0-a025-4feb-8961-8b8396018375" { // poisoned
		return true
	}

	return false
}

func (re *RulesEngine) GamerWasHidden(g *m.Game, gamer *m.Gamer) bool {
	return gamer.Fate != nil &&
		(gamer.Fate.AbilityBin == "724b1b7f-f42e-4b4f-8393-b2ae38cb0d70" || gamer.Fate.AbilityBin == "f5b7c9e1-0a64-4e9c-9291-7e3b36df891d")
}
