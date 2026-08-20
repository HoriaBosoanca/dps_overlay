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

func getPlayerRiotID() (string, error) {
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

func loadSummonerInfo() ([]Summoner, error) {
	res, err := client.Get(liveClientDataUrl + "/playerlist")
	if err != nil {
		return nil, fmt.Errorf("error accessing endpoint: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var summonerInfo []SummonerInfo
	if err = json.Unmarshal(body, &summonerInfo); err != nil {
		return nil, fmt.Errorf("error parsing summoner list: %w", err)
	}
	var summoners []Summoner
	for _, info := range summonerInfo {
		summoners = append(summoners, Summoner{SummonerInfo: info})
	}
	return summoners, nil
}

func findTeams(playerRiotID string, allSummoners []Summoner) (player *Summoner, enemies []*Summoner, err error) {
	for i := range allSummoners {
		if allSummoners[i].SummonerInfo.RiotID == playerRiotID {
			player = &allSummoners[i]
		}
	}
	if player == nil {
		return nil, nil, fmt.Errorf("player with RiotID %s not found", playerRiotID)
	}
	for i := range allSummoners {
		if allSummoners[i].SummonerInfo.Team != (*player).SummonerInfo.Team {
			enemies = append(enemies, &allSummoners[i])
		}
	}
	return player, enemies, nil
}

func updateLiveChampionInfo(summoners []*Summoner) error {
	res, err := client.Get(liveClientDataUrl + "/playerlist")
	if err != nil {
		return fmt.Errorf("error accessing endpoint: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}

	var liveChampionInfo []LiveChampionInfo
	if err = json.Unmarshal(body, &liveChampionInfo); err != nil {
		return fmt.Errorf("error parsing summoner list: %w", err)
	}
	for i := range summoners {
		for _, champInfo := range liveChampionInfo {
			if summoners[i].SummonerInfo.RiotID == champInfo.RiotID {
				summoners[i].LiveChampionInfo.Level = champInfo.Level
				summoners[i].LiveChampionInfo.ItemIDs = champInfo.ItemIDs
				break
			}
		}
	}
	return nil
}

func updateLivePlayerStats(player *Summoner) error {
	res, err := client.Get(liveClientDataUrl + "/activeplayer")
	if err != nil {
		return fmt.Errorf("error accessing endpoint: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}
	if err = json.Unmarshal(body, &player.LivePlayerStats); err != nil {
		return fmt.Errorf("error parsing summoner stats: %w", err)
	}
	return nil
}
