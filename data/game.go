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
		fmt.Println("error getting player RiotID: ", err)
		return
	}
	Summoners, err = loadSummoners()
	if err != nil {
		fmt.Println("error getting player list: ", err)
		return
	}
	Player, Enemies, err = findTeams(playerRiotID, Summoners)
	if err != nil {
		fmt.Println("error finding teams: ", err)
		return
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
	rl.DrawText("AA damage:", 20, 40, 20, rl.Yellow)
	y := int32(40)
	for _, enemy := range Enemies {
		y += 20
		if !checkLoadSuccess(enemy) {
			rl.DrawText(fmt.Sprintf("error loading info for %s's champion/items", enemy.ChampionName), 20, y, 20, rl.Red)
			continue
		}
		text := fmt.Sprintf("%s: %.1f", enemy.ChampionName, autoAttackDamage(*Player, *enemy))
		rl.DrawText(text, 20, y, 20, rl.Red)
	}
}
