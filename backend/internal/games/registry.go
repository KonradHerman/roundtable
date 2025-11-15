package games

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/yourusername/roundtable/internal/core"
	"github.com/yourusername/roundtable/internal/games/werewolf"
)

// GameFactory creates a new game instance.
type GameFactory func() core.Game

// ConfigParser parses raw JSON config into game-specific config.
type ConfigParser func([]byte) (core.GameConfig, error)

// Registry manages available game types.
type Registry struct {
	factories map[string]GameFactory
	parsers   map[string]ConfigParser
}

// NewRegistry creates a new game registry with all available games.
func NewRegistry() *Registry {
	r := &Registry{
		factories: make(map[string]GameFactory),
		parsers:   make(map[string]ConfigParser),
	}

	// Register games
	r.Register("werewolf", werewolf.NewGame, werewolf.ParseConfig)
	// Future: r.Register("avalon", avalon.NewGame, avalon.ParseConfig)
	// Future: r.Register("bohnanza", bohnanza.NewGame, bohnanza.ParseConfig)

	return r
}

// Register adds a game type to the registry.
func (r *Registry) Register(gameType string, factory GameFactory, parser ConfigParser) {
	r.factories[gameType] = factory
	r.parsers[gameType] = parser
}

// CreateGame creates a new game instance by type.
func (r *Registry) CreateGame(gameType string) (core.Game, error) {
	factory, exists := r.factories[gameType]
	if !exists {
		return nil, fmt.Errorf("unknown game type: %s", gameType)
	}

	return factory(), nil
}

// ParseConfig parses raw JSON config for a game type.
func (r *Registry) ParseConfig(gameType string, data json.RawMessage) (core.GameConfig, error) {
	parser, exists := r.parsers[gameType]
	if !exists {
		return nil, fmt.Errorf("unknown game type: %s", gameType)
	}

	return parser(data)
}

// IsRegistered checks if a game type is available.
func (r *Registry) IsRegistered(gameType string) bool {
	_, exists := r.factories[gameType]
	return exists
}

// ListGames returns all registered game types.
func (r *Registry) ListGames() []string {
	games := make([]string, 0, len(r.factories))
	for gameType := range r.factories {
		games = append(games, gameType)
	}
	return games
}

// ValidateConfig validates a config without creating a game.
func (r *Registry) ValidateConfig(gameType string, data json.RawMessage) error {
	config, err := r.ParseConfig(gameType, data)
	if err != nil {
		return err
	}

	if err := config.Validate(); err != nil {
		return errors.New("invalid configuration: " + err.Error())
	}

	return nil
}
