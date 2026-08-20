package data

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Summoners []Summoner
var Player *Summoner
var Enemies []*Summoner
var ItemInfoCache map[int]ItemInfo

func initGlobals() {
	Summoners = nil
	Player = nil
	Enemies = nil
	ItemInfoCache = make(map[int]ItemInfo)
}

func InitGame() error {
	initGlobals()
	playerRiotID, err := getPlayerRiotID()
	if err != nil {
		return fmt.Errorf("error getting player RiotID: %w", err)
	}
	Summoners, err = loadSummonerInfo()
	if err != nil {
		return fmt.Errorf("error loading summoner info: %w", err)
	}
	Player, Enemies, err = findTeams(playerRiotID, Summoners)
	if err != nil {
		return fmt.Errorf("error splitting summoners into teams: %w", err)
	}
	loadStaticChampionInfo([]*Summoner{Player})
	loadStaticChampionInfo(Enemies)
	return nil
}

func LoopGame() {
	err := updateLivePlayerStats(Player)
	if err != nil {
		fmt.Println("error updating live player stats: ", err)
	}
	err = updateLiveChampionInfo(Enemies)
	if err != nil {
		fmt.Println("error updating live champion info: ", err)
	}
	loadItemInfo(Enemies)
	calculateResistances(Enemies)
	for _, enemy := range Enemies {
		if !checkLoadSuccess(enemy) {
			drawText(fmt.Sprintf("error loading info for %s's champion/items/abilities", enemy.SummonerInfo.ChampionName), rl.Red)
			continue
		}
		drawText(fmt.Sprintf("%s:", enemy.SummonerInfo.ChampionName), rl.Blue)
		drawText(fmt.Sprintf("AA: %.0f", autoAttackDamage(*Player, *enemy)), rl.Yellow)
		if Player.LivePlayerStats.Stats.CritChance > 0 {
			drawText(fmt.Sprintf("Crit: %.0f", critDamage(*Player, *enemy)), rl.Yellow)
		}
		drawText(fmt.Sprintf("Q: %.0f", abilityDamage(*Player,
			Player.StaticChampionInfo.Abilities.Q[0],
			Player.LivePlayerStats.AbilityLevels.Q.AbilityLevel,
			*enemy)), rl.Yellow)
		drawText(fmt.Sprintf("W: %.0f", abilityDamage(*Player,
			Player.StaticChampionInfo.Abilities.W[0],
			Player.LivePlayerStats.AbilityLevels.W.AbilityLevel,
			*enemy)), rl.Yellow)
		drawText(fmt.Sprintf("E: %.0f", abilityDamage(*Player,
			Player.StaticChampionInfo.Abilities.E[0],
			Player.LivePlayerStats.AbilityLevels.E.AbilityLevel,
			*enemy)), rl.Yellow)
		drawText(fmt.Sprintf("R: %.0f", abilityDamage(*Player,
			Player.StaticChampionInfo.Abilities.R[0],
			Player.LivePlayerStats.AbilityLevels.R.AbilityLevel,
			*enemy)), rl.Yellow)
	}
	resetTextPos()
}

var Font rl.Font
var pos = rl.Vector2{X: 20, Y: 0}

func drawText(text string, color rl.Color) {
	pos.Y += 20
	rl.DrawTextEx(Font, text, pos, 20, 0.5, color)
}
func resetTextPos() {
	pos = rl.Vector2{X: 20, Y: 0}
}
