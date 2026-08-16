package data

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}
var liveClientDataUrl = "https://127.0.0.1:2999/liveclientdata"

func getPLayerRiotID() (string, error) {
	res, err := client.Get(liveClientDataUrl + "/activeplayer")
	if err != nil {
		return "", fmt.Errorf("error accessing endpoint: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}
	var wrapper struct {
		RiotID string `json:"riotId"`
	}
	if err = json.Unmarshal(body, &wrapper); err != nil {
		return "", fmt.Errorf("error parsing summoner stats: %w", err)
	}
	return wrapper.RiotID, nil
}

func loadSummoners() ([]Summoner, error) {
	res, err := client.Get(liveClientDataUrl + "/playerlist")
	if err != nil {
		return nil, fmt.Errorf("error accessing endpoint: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var summoners []Summoner
	if err = json.Unmarshal(body, &summoners); err != nil {
		return nil, fmt.Errorf("error parsing summoner list: %w", err)
	}
	loadChampionInfo(summoners)
	return summoners, nil
}

func findTeams(playerRiotID string, allSummoners []Summoner) (player *Summoner, enemies []*Summoner, err error) {
	for i := range allSummoners {
		if allSummoners[i].RiotID == playerRiotID {
			player = &allSummoners[i]
		}
	}
	if player == nil {
		return nil, nil, fmt.Errorf("player with RiotID %s not found", playerRiotID)
	}
	for i := range allSummoners {
		if allSummoners[i].Team != (*player).Team {
			enemies = append(enemies, &allSummoners[i])
		}
	}
	return player, enemies, nil
}

func updateSummoners(oldSummoners []*Summoner) error {
	res, err := client.Get(liveClientDataUrl + "/playerlist")
	if err != nil {
		return fmt.Errorf("error accessing endpoint: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}

	var newSummoners []Summoner
	if err = json.Unmarshal(body, &newSummoners); err != nil {
		return fmt.Errorf("error parsing summoner list: %w", err)
	}
	for i := range oldSummoners {
		for _, newSummoner := range newSummoners {
			if oldSummoners[i].RiotID == newSummoner.RiotID {
				oldSummoners[i].Level = newSummoner.Level
				oldSummoners[i].ItemIDs = newSummoner.ItemIDs
				break
			}
		}
	}
	return nil
}

func updatePlayerStats(player *Summoner) error {
	res, err := client.Get(liveClientDataUrl + "/activeplayer")
	if err != nil {
		return fmt.Errorf("error accessing endpoint: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}
	var wrapper struct {
		LiveStats LiveStats `json:"championStats"`
	}
	if err = json.Unmarshal(body, &wrapper); err != nil {
		return fmt.Errorf("error parsing summoner stats: %w", err)
	}
	player.LiveStats = wrapper.LiveStats
	return nil
}
