# One Night Werewolf Rules

## Overview
One Night Werewolf is a fast-paced social deduction game where players have secret roles and must work together to identify the werewolves before it's too late.

## Game Setup

### Role Cards
**Critical Rule**: There must be **exactly 3 MORE role cards than players**.
- 3 players = 6 role cards (3 to players, 3 in center)
- 4 players = 7 role cards (4 to players, 3 in center)
- 5 players = 8 role cards (5 to players, 3 in center)
- etc.

The 3 extra cards are placed face-down in the "center" and are not assigned to any player initially.

### Recommended Role Distributions

**3-4 Players:**
- 1-2 Werewolves
- 1 Seer
- 1 Robber or Troublemaker
- Villagers to fill
- 3 center cards

**5-6 Players:**
- 2 Werewolves
- 1 Seer
- 1 Robber
- 1 Troublemaker
- 2 Masons
- Villagers to fill
- 3 center cards

**7-10 Players:**
- 2-3 Werewolves
- 1 Seer
- 1 Robber
- 1 Troublemaker
- 2 Masons
- 1 Drunk
- 1 Insomniac
- 1 Minion (optional)
- 1 Tanner (optional for chaos)
- Villagers to fill
- 3 center cards

## Game Phases

### 1. Role Reveal (5 seconds)
- Each player secretly views their role card
- Players should memorize their role

### 2. Night Phase
Players close their eyes and certain roles "wake up" in sequence:

**Order of Night Actions:**
1. **Werewolves** (open eyes, see each other)
2. **Minion** (open eyes, sees werewolves, but they don't see them)
3. **Masons** (open eyes, see each other)
4. **Seer** (choose: view one player's card OR two center cards)
5. **Robber** (swap cards with another player, view new card)
6. **Troublemaker** (swap two other players' cards, don't look)
7. **Drunk** (swap card with center card, don't look at new card)
8. **Insomniac** (view own card to see if it changed)

**Important**: Roles act in order, and later actions can affect earlier ones. Players do NOT know if their card was swapped unless they're the Insomniac.

### 3. Day Phase (Discussion & Voting)
- All players open eyes
- Discuss who might be werewolves (3-5 minutes typical)
- Vote simultaneously on who to eliminate
- Most votes = eliminated

## Win Conditions

### Village Team Wins If:
- At least one Werewolf is eliminated, OR
- No Werewolves are alive and no one is eliminated

### Werewolf Team Wins If:
- No Werewolves are eliminated AND someone is eliminated

### Tanner Wins If:
- The Tanner is eliminated (Tanner wins alone, all others lose)

### Special Cases:
- If ALL Werewolves are in the center (no player Werewolves), Village wins if no one dies
- Minion is on Werewolf team even if all Werewolves are in center

## Role Abilities

### Werewolf Team
- **Werewolf**: Sees other werewolves. Tries to avoid detection.
- **Minion**: Sees werewolves but they don't see them. Tries to protect werewolves.

### Village Team
- **Seer**: Views one player card OR two center cards
- **Robber**: Swaps card with another player, becomes their role
- **Troublemaker**: Swaps two other players' cards
- **Mason**: Sees other mason(s)
- **Drunk**: Swaps with center, doesn't know new role
- **Insomniac**: Checks if their role changed at end of night
- **Villager**: No special ability, must use deduction

### Solo
- **Tanner**: Wants to be eliminated to win alone

## Implementation Notes for Digital Version

### Phase Timers
- Role Reveal: 5 seconds (auto-advance)
- Night Phase: 30-60 seconds (role actions can be taken anytime)
- Day Phase: 2-5 minutes (with vote prompt when time low)

### Center Cards
- Must be implemented even for digital version
- Some roles interact with center cards (Seer, Robber, Drunk)
- Center cards are critical to the deduction logic

### Role Actions
- All night actions should be optional (player can skip)
- Actions should be taken secretly (private to each player)
- Backend must track original roles vs. final roles

