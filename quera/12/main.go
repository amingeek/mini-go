package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Player struct {
	name     string
	location int
	chanel   chan string
	game     *Game
}

type Map struct {
	players map[string]bool
	chanel  chan string
	mu      sync.Mutex
	game    *Game
}

type Game struct {
	players map[string]*Player
	maps    map[int]*Map
	mu      sync.Mutex
}

func formatName(name string) string {
	if len(name) == 0 {
		return ""
	}
	return strings.ToUpper(string(name[0])) + strings.ToLower(name[1:])
}

func NewGame(mapIds []int) (*Game, error) {
	game := &Game{
		players: make(map[string]*Player),
		maps:    make(map[int]*Map),
	}

	for _, id := range mapIds {
		if id <= 0 {
			return nil, errors.New("map ID must be greater than 0")
		}

		m := &Map{
			players: make(map[string]bool),
			chanel:  make(chan string, 100),
			game:    game, // تنظیم ارجاع به بازی
		}

		game.maps[id] = m
	}

	// بعد از ساخت کامل بازی، FanOutMessages را برای هر نقشه اجرا کن
	for _, m := range game.maps {
		go m.FanOutMessages()
	}

	return game, nil
}

func (g *Game) ConnectPlayer(name string) error {
	normalized := strings.ToLower(name)
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.players[normalized]; exists {
		return errors.New("player already exists")
	}

	g.players[normalized] = &Player{
		name:     formatName(name),
		location: 0,
		chanel:   make(chan string, 100),
		game:     g,
	}

	return nil
}

func (g *Game) SwitchPlayerMap(name string, mapId int) error {
	normalized := strings.ToLower(name)

	g.mu.Lock()
	defer g.mu.Unlock()

	player, exists := g.players[normalized]
	if !exists {
		return errors.New("player not found")
	}

	newMap, mapExists := g.maps[mapId]
	if !mapExists {
		return errors.New("map not found")
	}

	if player.location == mapId {
		return errors.New("player already in this map")
	}

	oldMap, oldExists := g.maps[player.location]
	if oldExists {
		oldMap.mu.Lock()
		delete(oldMap.players, normalized)
		oldMap.mu.Unlock()
	}

	newMap.mu.Lock()
	newMap.players[normalized] = true
	newMap.mu.Unlock()

	player.location = mapId

	return nil
}

func (g *Game) GetPlayer(name string) (*Player, error) {
	normalized := strings.ToLower(name)

	g.mu.Lock()
	defer g.mu.Unlock()

	player, exists := g.players[normalized]
	if !exists {
		return nil, errors.New("player not found")
	}

	return player, nil
}

func (g *Game) GetMap(mapId int) (*Map, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	m, exists := g.maps[mapId]
	if !exists {
		return nil, errors.New("map not found")
	}
	return m, nil
}

func (m *Map) FanOutMessages() {
	for {
		msg := <-m.chanel

		m.mu.Lock()
		for playerName := range m.players {
			player := m.game.players[playerName]
			if player != nil {

				senderEndIndex := strings.Index(msg, " says:")
				if senderEndIndex != -1 {
					senderName := msg[:senderEndIndex]
					if strings.ToLower(senderName) == playerName {
						continue
					}
				}

				select {
				case player.chanel <- msg:
				default:
				}
			}
		}
		m.mu.Unlock()
	}
}

func (p *Player) GetChannel() <-chan string {
	return p.chanel
}

func (p *Player) SendMessage(msg string) error {
	if p.location == 0 {
		return errors.New("player is not in any map")
	}

	p.game.mu.Lock()
	defer p.game.mu.Unlock()

	m, exists := p.game.maps[p.location]
	if !exists {
		return errors.New("map not found")
	}

	formattedName := formatName(p.name)
	m.chanel <- fmt.Sprintf("%s says: %s", formattedName, msg)

	return nil
}

func (p *Player) GetName() string {
	return p.name
}
