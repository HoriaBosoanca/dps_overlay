package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func LoadItemData(summoners []Summoner) {
	for i := range summoners {
		for _, itemID := range summoners[i].ItemIDs {
			file, err := os.Open(fmt.Sprintf("../assets/items/%s.json", itemID.ItemID))
			if err != nil {
				fmt.Println("error opening item file for ", itemID.ItemID, ": ", err)
				continue
			}
			body, err := io.ReadAll(file)
			if err != nil {
				fmt.Println("error reading item file for ", itemID.ItemID, ": ", err)
				continue
			}
			file.Close()
			var wrapper struct {
				ImportedStats ItemImportedInfo `json:"stats"`
			}
			if err = json.Unmarshal(body, &wrapper); err != nil {
				fmt.Println("error parsing item data for ", itemID.ItemID, ": ", err)
				continue
			}
			summoners[i].ItemImportedInfo = append(summoners[i].ItemImportedInfo, wrapper.ImportedStats)
		}
	}
}

func LoadChampionData(summoners []Summoner) {
	for i := range summoners {
		fileName := strings.ReplaceAll(summoners[i].ChampionName, " ", "")
		fileName = strings.ReplaceAll(fileName, ".", "")
		fileName = strings.ReplaceAll(fileName, "'", "")
		file, err := os.Open(fmt.Sprintf("../assets/champions/%s.json", fileName))
		if err != nil {
			fmt.Println("error opening champion file for ", summoners[i].ChampionName, ": ", err)
			continue
		}
		body, err := io.ReadAll(file)
		if err != nil {
			file.Close()
			fmt.Println("error reading champion file for ", summoners[i].ChampionName, ": ", err)
			continue
		}
		file.Close()
		var wrapper struct {
			ImportedStats ChampionImportedInfo `json:"stats"`
		}
		if err = json.Unmarshal(body, &wrapper); err != nil {
			fmt.Println("error parsing champion data for ", summoners[i].ChampionName, ": ", err)
			continue
		}
		summoners[i].ChampionImportedInfo = wrapper.ImportedStats
	}
}
