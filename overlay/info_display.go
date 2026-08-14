package overlay

import (
	"dps_overlay/data"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var summoners []data.Summoner
var player *data.Summoner
var enemies []*data.Summoner

func initDisplay() {
	summoners = nil
	player = nil
	enemies = nil
	playerRiotID, err := data.GetPLayerRiotID()
	if err != nil {
		fmt.Println("error getting player RiotID: ", err)
		return
	}
	summoners, err = data.LoadSummoners()
	if err != nil {
		fmt.Println("error getting summoners: ", err)
		return
	}
	player, enemies, err = data.FindTeams(playerRiotID, summoners)
	if err != nil {
		fmt.Println("error finding teams: ", err)
		return
	}
}

func loopDisplay() {
	err := data.UpdatePlayerStats(player)
	if err != nil {
		fmt.Println("error getting player stats: ", err)
		return
	}
	err = data.UpdateSummoners(summoners)
	if err != nil {
		fmt.Println("error getting summoner stats: ", err)
		return
	}
	data.CalculateResistances(enemies)
	// AA damage
	rl.DrawText("AA damage:", 20, 40, 20, rl.Yellow)
	y := int32(60)
	for _, enemy := range enemies {
		text := fmt.Sprintf("%s: %.1f", enemy.ChampionName, data.AAdamage(*player, *enemy))
		rl.DrawText(text, 20, y, 20, rl.Red)
		y += 20
	}
}
