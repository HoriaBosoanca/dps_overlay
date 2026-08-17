package data

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Summoners []Summoner
var Player *Summoner
var Enemies []*Summoner
var ItemInfoCache map[int]ItemImportedInfo

func initGlobals() {
	Summoners = nil
	Player = nil
	Enemies = nil
	ItemInfoCache = make(map[int]ItemImportedInfo)
}

func InitGame() {
	initGlobals()
	playerRiotID, err := getPLayerRiotID()
	if err != nil {
		panic(fmt.Errorf("error getting player RiotID: %w", err))
	}
	Summoners, err = loadSummoners()
	if err != nil {
		panic(fmt.Errorf("error getting player list: %w", err))
		return
	}
	Player, Enemies, err = findTeams(playerRiotID, Summoners)
	if err != nil {
		panic(fmt.Errorf("error finding teams: %w", err))
		return
	}
	if Player == nil {
		panic(fmt.Errorf("error finding player: %w", err))
	}
}

func LoopGame() {
	err := updatePlayerStats(Player)
	if err != nil {
		fmt.Println("error getting active player stats: ", err)
		return
	}
	err = updateSummoners(Enemies)
	if err != nil {
		fmt.Println("error getting player list stats: ", err)
		return
	}
	loadItemInfo(Summoners)
	calculateResistances(Enemies)
	// AA damage
	drawText("AA damage:", rl.Yellow)
	for _, enemy := range Enemies {
		if !checkLoadSuccess(enemy) {
			drawText(fmt.Sprintf("error loading info for %s's champion/items", enemy.ChampionName), rl.Red)
			continue
		}
		text := fmt.Sprintf("%s: %.1f", enemy.ChampionName, autoAttackDamage(*Player, *enemy))
		drawText(text, rl.Red)
	}
	resetTextPos()
}

var Font rl.Font
var pos rl.Vector2

func drawText(text string, color rl.Color) {
	pos.Y += 20
	rl.DrawTextEx(Font, text, pos, 20, 0.5, color)
}
func resetTextPos() {
	pos = rl.Vector2{X: 20, Y: 0}
}
