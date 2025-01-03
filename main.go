package pm_rules

import (
	"github.com/horcu/pm-models/enums"
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

	ab := gamer.Abilities[ability.Bin]

	if ab.CyclesUsedIndex == nil {
		ab.CyclesUsedIndex = []int{}
	}

	ab.CyclesUsedIndex = append(ab.CyclesUsedIndex, g.NightCycles)
	ab.TimesUsed++

	//save the gamer to the store
	return re.store.UpdateGamerAbilities(g.Bin, gamer.Bin, ab.Bin, ab)

}

func (re *RulesEngine) ApplyAbility(g *m.Game, sourceChar *m.GameCharacter, targetGamer *m.Gamer, ability string) (bool, string) {

	if !targetGamer.IsAlive {
		return false, "target is dead"
	}

	switch ability {
	case enums.Kill.String():
		re.store.ApplyAbility("kill", g.Bin, targetGamer.Bin)
		return true, targetGamer.Name + " was marked for death by the " + sourceChar.Name
	case enums.Hide.String():
		re.store.ApplyAbility("hide", g.Bin, targetGamer.Bin)
		return true, targetGamer.Name + " hidden from danger by the " + sourceChar.Name
	case enums.Trick.String():
		re.store.ApplyAbility("trick", g.Bin, targetGamer.Bin)
		return true, targetGamer.Name + " was tricked into their decision by the " + sourceChar.Name
	case enums.Mimic.String():
		re.store.ApplyAbility("mimic", g.Bin, targetGamer.Bin)
		return true, targetGamer.Name + " had their abilities copied by the " + sourceChar.Name
	case enums.Heal.String():
		if re.CanBeHealed(g, targetGamer) {
			re.store.ApplyAbility("heal", g.Bin, targetGamer.Bin)
			return true, targetGamer.Name + " was healed by the " + sourceChar.Name
		}
	case enums.Poison.String():
		re.store.ApplyAbility("poison", g.Bin, targetGamer.Bin)
		return true, targetGamer.Name + " was poisoned by the " + sourceChar.Name
	case enums.Block.String():
		re.store.ApplyAbility("block", g.Bin, targetGamer.Bin)
		return true, targetGamer.Name + " was blocked by the " + sourceChar.Name
	case enums.Retaliate.String():
		re.store.ApplyAbility("retaliate", g.Bin, targetGamer.Bin)
		return true, targetGamer.Name + " killed for targetting the " + sourceChar.Name
	case enums.Investigate.String():
		if targetGamer.IsAlive {
			re.store.ApplyAbility("investigate", g.Bin, targetGamer.Bin)
			return true, targetGamer.Name + " was investigated by the " + sourceChar.Name
		}
	//case enums..String():
	//	if targetGamer.IsAlive {
	//		re.store.ApplyAbility("meet", g.Bin, targetGamer.Bin)
	//		return true, targetGamer.Name + " met with the " + sourceChar.Name
	//	}
	case enums.Mark.String():
		if targetGamer.IsAlive && !re.GamerWasHidden(g, targetGamer) {
			re.store.ApplyAbility("mark", g.Bin, targetGamer.Bin)
			return true, targetGamer.Name + " was marked by the " + sourceChar.Name
		}
	case enums.Direct.String():
		// get the character to ensure that they are of type villains
		char, err := re.store.GetCharacterByBin(targetGamer.CharacterId)
		if err != nil {
			return false, "error getting character"
		}
		if char.SideId == 0 {
			re.store.ApplyAbility("direct", g.Bin, targetGamer.Bin)
			return true, targetGamer.Name + " was directed in their decision by the " + sourceChar.Name
		}

	default:
		return false, "ability not usable"
	}
	return false, "invalid ability"
}

func (re *RulesEngine) CanBeKilled(g *m.Game, gamer *m.Gamer) bool {

	return gamer.IsAlive
}

func (re *RulesEngine) CanBeTricked(g *m.Game, gamer *m.Gamer) bool {
	return gamer.IsAlive
}

func (re *RulesEngine) CanBeMimicked(g *m.Game, gamer *m.Gamer) bool {
	return gamer.IsAlive
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

func (re *RulesEngine) GamerNeedsHealing(g *m.Game, gamer *m.Gamer) bool {

	if gamer.Fate != nil {
		return gamer.Fate.AbilityBin == "afef47b0-a025-4feb-8961-8b8396018375"
	}
	return false
}

func (re *RulesEngine) GamerWasHidden(g *m.Game, gamer *m.Gamer) bool {
	if gamer.Fate != nil {
		return gamer.Fate.AbilityBin == "724b1b7f-f42e-4b4f-8393-b2ae38cb0d70" || gamer.Fate.AbilityBin == "f5b7c9e1-0a64-4e9c-9291-7e3b36df891d"
	}
	return false
}
