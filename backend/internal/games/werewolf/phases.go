package werewolf

import (
	"time"

	"github.com/yourusername/roundtable/internal/core"
)

// AdvanceToDay transitions the game from night to day phase.
// This should be called when the night phase timer expires.
func (g *Game) AdvanceToDay() ([]core.GameEvent, error) {
	if g.phase != PhaseNight {
		return nil, nil // Already past night phase
	}

	g.phase = PhaseDay
	g.phaseStartedAt = time.Now()
	g.phaseEndsAt = g.phaseStartedAt.Add(g.config.DayDuration)

	events := make([]core.GameEvent, 0)

	phaseEvent, _ := core.NewPublicEvent(core.EventPhaseChanged, "system", core.PhaseChangedPayload{
		Phase: core.GamePhase{
			Name:    string(PhaseDay),
			EndsAt:  &g.phaseEndsAt,
			Message: "Day phase - discuss and vote!",
		},
	})
	events = append(events, phaseEvent)

	return events, nil
}

// CheckPhaseTimeout checks if the current phase has expired and should advance.
func (g *Game) CheckPhaseTimeout() ([]core.GameEvent, error) {
	now := time.Now()

	if now.After(g.phaseEndsAt) {
		switch g.phase {
		case PhaseNight:
			return g.AdvanceToDay()
		case PhaseDay:
			// Day phase expired without all votes - could auto-end or extend
			// For now, we'll just let it stay in day phase
			return nil, nil
		}
	}

	return nil, nil
}
