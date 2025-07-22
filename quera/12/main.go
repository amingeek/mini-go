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
}

type Map struct {
	players map[string]bool
	chanel  chan string
	mu      sync.Mutex
}

type Game struct {
	players map[string]*Player
	maps    map[int]*Map
	mu      sync.Mutex
}

func NewGame(mapIds []int) (*Game, error) {
	maps := make(map[int]*Map)

	for _, id := range mapIds {
		if id <= 0 {
			return nil, errors.New("map ID must be greater than 0")
		}

		m := &Map{
			players: make(map[string]bool),
			chanel:  make(chan string, 100),
		}

		go m.FanOutMessages()

		maps[id] = m
	}

	return &Game{
		players: make(map[string]*Player),
		maps:    maps,
	}, nil
}

func (g *Game) ConnectPlayer(name string) error {
	normalized := strings.ToLower(name)
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.players[normalized]; exists {
		return errors.New("player already exists")
	}

	g.players[normalized] = &Player{
		name:     name,
		location: 0,
		chanel:   make(chan string, 100),
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

	oldMap, oldExists := g.maps[player.location]

	player.location = mapId

	if oldExists {
		oldMap.chanel <- fmt.Sprintf("%s left the map", player.name)
	}

	newMap.chanel <- fmt.Sprintf("%s entered the map", player.name)

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
}

func (p *Player) GetChannel() <-chan string {
	return nil
}

func (p *Player) SendMessage(msg string) error {
	return nil
}

func (p *Player) GetName() string {
	return ""
}
