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
	if summoner.StaticChampionInfo.LoadedSuccess == false {
		loadSuccess = false
	}
	for _, itemID := range summoner.LiveChampionInfo.ItemIDs {
		if ItemInfoCache[itemID.ItemID].LoadSuccess == false {
			loadSuccess = false
		}
	}
	return loadSuccess
}

func loadItemInfo(summoners []*Summoner) {
	for i := range summoners {
		for _, itemID := range summoners[i].LiveChampionInfo.ItemIDs {
			if _, exists := ItemInfoCache[itemID.ItemID]; exists {
				continue
			}
			ItemInfoCache[itemID.ItemID] = ItemInfo{LoadSuccess: false}
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
				ItemInfo ItemInfo `json:"stats"`
			}
			if err = json.Unmarshal(body, &wrapper); err != nil {
				fmt.Println("error parsing item data for ", itemID.ItemID, ": ", err)
				continue
			}
			wrapper.ItemInfo.LoadSuccess = true
			ItemInfoCache[itemID.ItemID] = wrapper.ItemInfo
		}
	}
}

func loadStaticChampionInfo(summoners []*Summoner) {
	for i := range summoners {
		summoners[i].StaticChampionInfo = StaticChampionInfo{LoadedSuccess: false}
		fileName := strings.ReplaceAll(summoners[i].SummonerInfo.ChampionName, " ", "")
		fileName = strings.ReplaceAll(fileName, ".", "")
		fileName = strings.ReplaceAll(fileName, "'", "")
		fileName = strings.ReplaceAll(fileName, "&", "")
		file, err := os.Open(fmt.Sprintf("../assets/champions/%s.json", fileName))
		if err != nil {
			fmt.Println("error opening champion file for ", summoners[i].SummonerInfo.ChampionName, ": ", err)
			continue
		}
		body, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			fmt.Println("error reading champion file for ", summoners[i].SummonerInfo.ChampionName, ": ", err)
			continue
		}
		if err = json.Unmarshal(body, &summoners[i].StaticChampionInfo); err != nil {
			fmt.Println("error parsing champion data for ", summoners[i].SummonerInfo.ChampionName, ": ", err)
			continue
		}
		summoners[i].StaticChampionInfo.LoadedSuccess = true
	}
}
