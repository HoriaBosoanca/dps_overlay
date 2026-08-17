package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func checkLoadSuccess(summoner *Summoner) bool {
	loadSuccess := true
	if summoner.ChampionImportedInfo.LoadedSuccess == false {
		loadSuccess = false
	}
	for _, itemID := range summoner.ItemIDs {
		if ItemInfoCache[itemID.ItemID].LoadSuccess == false {
			loadSuccess = false
		}
	}
	return loadSuccess
}

func loadItemInfo(summoners []Summoner) {
	for i := range summoners {
		for _, itemID := range summoners[i].ItemIDs {
			if _, exists := ItemInfoCache[itemID.ItemID]; exists {
				continue
			}
			ItemInfoCache[itemID.ItemID] = ItemImportedInfo{LoadSuccess: false}
			file, err := os.Open(fmt.Sprintf("../assets/items/%d.json", itemID.ItemID))
			if err != nil {
				fmt.Println("error opening item file for ", itemID.ItemID, ": ", err)
				continue
			}
			body, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				fmt.Println("error reading item file for ", itemID.ItemID, ": ", err)
				continue
			}
			var wrapper struct {
				ImportedStats ItemImportedInfo `json:"stats"`
			}
			if err = json.Unmarshal(body, &wrapper); err != nil {
				fmt.Println("error parsing item data for ", itemID.ItemID, ": ", err)
				continue
			}
			wrapper.ImportedStats.LoadSuccess = true
			ItemInfoCache[itemID.ItemID] = wrapper.ImportedStats
		}
	}
}

func loadChampionInfo(summoners []Summoner) {
	for i := range summoners {
		summoners[i].ChampionImportedInfo = ChampionImportedInfo{LoadedSuccess: false}
		fileName := strings.ReplaceAll(summoners[i].ChampionName, " ", "")
		fileName = strings.ReplaceAll(fileName, ".", "")
		fileName = strings.ReplaceAll(fileName, "'", "")
		fileName = strings.ReplaceAll(fileName, "&", "")
		file, err := os.Open(fmt.Sprintf("../assets/champions/%s.json", fileName))
		if err != nil {
			fmt.Println("error opening champion file for ", summoners[i].ChampionName, ": ", err)
			continue
		}
		body, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			fmt.Println("error reading champion file for ", summoners[i].ChampionName, ": ", err)
			continue
		}
		var wrapper struct {
			ImportedStats ChampionImportedInfo `json:"stats"`
		}
		if err = json.Unmarshal(body, &wrapper); err != nil {
			fmt.Println("error parsing champion data for ", summoners[i].ChampionName, ": ", err)
			continue
		}
		wrapper.ImportedStats.LoadedSuccess = true
		summoners[i].ChampionImportedInfo = wrapper.ImportedStats
	}
}
