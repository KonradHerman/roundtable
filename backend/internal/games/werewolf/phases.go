package werewolf

import (
	"errors"
	"time"

	"github.com/KonradHerman/roundtable/internal/core"
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

	// Phase change event (public)
	phaseEvent, _ := core.NewPublicEvent(core.EventPhaseChanged, "system", core.PhaseChangedPayload{
		Phase: core.GamePhase{
			Name:    string(PhaseNight),
			Message: "Night phase - everyone close your eyes",
		},
	})
	events = append(events, phaseEvent)

	// Send role-specific wakeup events

	// 1. Werewolves wake up and see each other
	werewolves := g.getPlayersByRole(RoleWerewolf)
	if len(werewolves) > 0 {
		for _, ww := range werewolves {
			otherWerewolves := make([]string, 0)
			for _, other := range werewolves {
				if other != ww {
					otherWerewolves = append(otherWerewolves, other)
				}
			}
			wakeupEvent, _ := core.NewPrivateEvent("werewolf_wakeup", "system", WerewolfWakeupPayload{
				OtherWerewolves: otherWerewolves,
			}, []string{ww})
			events = append(events, wakeupEvent)
		}
	}

	// 2. Masons wake up and see each other
	masons := g.getPlayersByRole(RoleMason)
	if len(masons) > 0 {
		for _, mason := range masons {
			otherMasons := make([]string, 0)
			for _, other := range masons {
				if other != mason {
					otherMasons = append(otherMasons, other)
				}
			}
			wakeupEvent, _ := core.NewPrivateEvent("mason_wakeup", "system", MasonWakeupPayload{
				OtherMasons: otherMasons,
			}, []string{mason})
			events = append(events, wakeupEvent)
		}
	}

	// Generate night script for the host
	allRoles := make([]RoleType, 0, len(g.roleAssignments)+len(g.centerCards))
	for _, role := range g.roleAssignments {
		allRoles = append(allRoles, role)
	}
	allRoles = append(allRoles, g.centerCards...)

	script := GenerateNightScript(allRoles)

	// Send script to host only
	if g.hostID != "" {
		scriptEvent, _ := core.NewPrivateEvent("night_script", "system", NightScriptPayload{
			Script: script,
		}, []string{g.hostID})
		events = append(events, scriptEvent)
	}

	return events, nil
}

// AdvanceToDay transitions the game from night to day phase.
func (g *Game) AdvanceToDay() ([]core.GameEvent, error) {
	if g.phase != PhaseNight {
		return nil, errors.New("can only advance to day from night phase")
	}

	events := make([]core.GameEvent, 0)

	// Insomniac wakes up and sees their final role (after all swaps)
	for playerID, role := range g.roleAssignments {
		if g.originalRoles[playerID] == RoleInsomniac {
			// Send the Insomniac their final role
			insomniacEvent, _ := core.NewPrivateEvent("insomniac_result", "system", InsomniacResultPayload{
				FinalRole: role,
			}, []string{playerID})
			events = append(events, insomniacEvent)
			g.nightActionsComplete[RoleInsomniac] = true
		}
	}

	g.phase = PhaseDay
	g.phaseStartedAt = time.Now()
	g.timerActive = false // Timer starts OFF

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
