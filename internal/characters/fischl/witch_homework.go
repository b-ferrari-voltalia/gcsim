package fischl

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	e "github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

var (
	c6WBDuration = 600 // buff duration (10s)
	c6WBExpiry   = 0   // frame when the doubled C6 buff expires
)

// Activates the timer for the C6 buff
func (c *char) c6WitchBuffTimer() {
	c6WBExpiry = c.Core.F + 600 // 10s of doubled buff
}

// Registers all buffs from the Crimson Witch set
func (c *char) registerWitchBonus() {
	// ------------------------------------------------
	//            	  Overload → ATK%
	// ------------------------------------------------
	c.Core.Events.Subscribe(e.OnOverload, func(args ...any) bool {
		// ---------------------------
		// Fischl always receives the buff
		// ---------------------------
		m1 := character.StatMod{
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				vals := make([]float64, attributes.EndStatType)
				atk := 0.225

				// If doubled C6 buff is active
				if c.Core.F < c6WBExpiry {
					atk *= 2
				}

				vals[attributes.ATKP] = atk
				return vals, true
			},
		}
		m1.ModKey = "witch-atk"
		m1.Dur = c6WBDuration
		c.AddStatMod(m1)

		// --------------------------------------------------------------
		// Buff for the ACTIVE character (but does not apply to Fischl)
		// --------------------------------------------------------------
		for _, char := range c.Core.Player.Chars() {
			char := char
			m2 := character.StatMod{
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					vals := make([]float64, attributes.EndStatType)

					// Only if this is the active character AND it's not Fischl
					if char.Index() == c.Core.Player.Active() && char.Index() != c.Index() {
						atk := 0.225
						if c.Core.F < c6WBExpiry {
							atk *= 2
						}
						vals[attributes.ATKP] = atk
						return vals, true
					}

					return vals, false
				},
			}
			m2.ModKey = "witch-atk-team"
			m2.Dur = c6WBDuration
			c.AddStatMod(m2)
		}

		return false
	}, "fischl-witch-bonus-atk")

	// ------------------------------------------------
	//       ⚡ ElectroCharged → Elemental Mastery
	// ------------------------------------------------
	c.Core.Events.Subscribe(e.OnElectroCharged, func(args ...any) bool {
		// ---------------------------
		// Fischl always receives the EM buff
		// ---------------------------
		m1 := character.StatMod{
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				vals := make([]float64, attributes.EndStatType)
				em := 90.0

				// Doubled during C6 window
				if c.Core.F < c6WBExpiry {
					em *= 2
				}

				vals[attributes.EM] = em
				return vals, true
			},
		}
		m1.ModKey = "witch-em"
		m1.Dur = c6WBDuration
		c.AddStatMod(m1)

		// --------------------------------------------------------------
		// Buff for the ACTIVE character (but does not apply to Fischl)
		// --------------------------------------------------------------
		for _, char := range c.Core.Player.Chars() {
			char := char
			m2 := character.StatMod{
				AffectedStat: attributes.EM,
				Amount: func() ([]float64, bool) {
					vals := make([]float64, attributes.EndStatType)

					// Only active character, excluding Fischl
					if char.Index() == c.Core.Player.Active() && char.Index() != c.Index() {
						em := 90.0
						if c.Core.F < c6WBExpiry {
							em *= 2
						}
						vals[attributes.EM] = em
						return vals, true
					}

					return vals, false
				},
			}
			m2.ModKey = "witch-em-team"
			m2.Dur = c6WBDuration
			c.AddStatMod(m2)
		}

		return false
	}, "fischl-witch-bonus-em")
}
