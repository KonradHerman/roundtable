package werewolf

import (
	"errors"
	"time"

	"github.com/yourusername/roundtable/internal/core"
)

// AdvanceToNight transitions from role reveal to night phase after all players acknowledge.
func (g *Game) AdvanceToNight() ([]core.GameEvent, error) {
	if g.phase != PhaseRoleReveal {
		return nil, errors.New("can only advance to night from role reveal phase")
	}

	// Check if all players have acknowledged
	if len(g.roleAcknowledgements) < len(g.players) {
		return nil, errors.New("not all players have acknowledged their roles")
	}

	g.phase = PhaseNight
	g.phaseStartedAt = time.Now()
	// Night phase has no automatic timer - host manually advances

	events := make([]core.GameEvent, 0)

	// Generate night script for the host
	allRoles := make([]RoleType, 0, len(g.roleAssignments)+len(g.centerCards))
	for _, role := range g.roleAssignments {
		allRoles = append(allRoles, role)
	}
	for _, role := range g.centerCards {
		allRoles = append(allRoles, role)
	}
	
	script := GenerateNightScript(allRoles)

	// Send script to host only
	hostID := ""
	for _, player := range g.players {
		// Find host - we'll need to track this, for now just use first player
		// In reality, the room knows who the host is
		hostID = player.ID
		break
	}

	scriptEvent, _ := core.NewPrivateEvent("night_script", "system", NightScriptPayload{
		Script: script,
	}, []string{hostID})
	events = append(events, scriptEvent)

	phaseEvent, _ := core.NewPublicEvent(core.EventPhaseChanged, "system", core.PhaseChangedPayload{
		Phase: core.GamePhase{
			Name:    string(PhaseNight),
			Message: "Night phase - everyone close your eyes",
		},
	})
	events = append(events, phaseEvent)

	return events, nil
}

// AdvanceToDay transitions the game from night to day phase.
func (g *Game) AdvanceToDay() ([]core.GameEvent, error) {
	if g.phase != PhaseNight {
		return nil, errors.New("can only advance to day from night phase")
	}

	g.phase = PhaseDay
	g.phaseStartedAt = time.Now()
	g.timerActive = false // Timer starts OFF

	events := make([]core.GameEvent, 0)

	phaseEvent, _ := core.NewPublicEvent(core.EventPhaseChanged, "system", core.PhaseChangedPayload{
		Phase: core.GamePhase{
			Name:    string(PhaseDay),
			Message: "Day phase - discuss and vote!",
		},
	})
	events = append(events, phaseEvent)

	return events, nil
}

// ToggleTimer turns the day phase timer on or off.
func (g *Game) ToggleTimer(enable bool, duration time.Duration) ([]core.GameEvent, error) {
	if g.phase != PhaseDay {
		return nil, errors.New("can only toggle timer during day phase")
	}

	g.timerActive = enable
	events := make([]core.GameEvent, 0)

	var phaseEndsAt *time.Time
	if enable {
		endTime := time.Now().Add(duration)
		g.phaseEndsAt = endTime
		phaseEndsAt = &endTime
	}

	timerEvent, _ := core.NewPublicEvent("timer_toggled", "system", TimerToggledPayload{
		Active:      enable,
		PhaseEndsAt: phaseEndsAt,
	})
	events = append(events, timerEvent)

	return events, nil
}

// ExtendTimer adds time to the day phase timer.
func (g *Game) ExtendTimer(seconds int) ([]core.GameEvent, error) {
	if g.phase != PhaseDay {
		return nil, errors.New("can only extend timer during day phase")
	}

	if !g.timerActive {
		return nil, errors.New("timer is not active")
	}

	g.phaseEndsAt = g.phaseEndsAt.Add(time.Duration(seconds) * time.Second)

	events := make([]core.GameEvent, 0)

	extendEvent, _ := core.NewPublicEvent("timer_extended", "system", TimerExtendedPayload{
		PhaseEndsAt: g.phaseEndsAt,
		ExtendedBy:  seconds,
	})
	events = append(events, extendEvent)

	return events, nil
}

// CheckPhaseTimeout checks if the current phase has expired and should advance.
func (g *Game) CheckPhaseTimeout() ([]core.GameEvent, error) {
	// Only auto-advance if timer is active during day phase
	if g.phase == PhaseDay && g.timerActive {
		now := time.Now()
		if now.After(g.phaseEndsAt) {
			// Timer expired but we don't auto-advance
			// Just turn off the timer
			g.timerActive = false
		}
	}

	return nil, nil
}
