package data

type Summoner struct {
	SummonerInfo       SummonerInfo
	StaticChampionInfo StaticChampionInfo
	LivePlayerStats    LivePlayerStats
	LiveChampionInfo   LiveChampionInfo
	CalculatedStats    CalculatedStats
}

type SummonerInfo struct {
	RiotID       string `json:"riotId"`
	Team         string `json:"team"`
	ChampionName string `json:"championName"`
}

type StaticChampionInfo struct {
	LoadedSuccess bool
	Stats         struct {
		Armor           ChampionStatInfo `json:"armor"`
		MagicResistance ChampionStatInfo `json:"magicResistance"`
		AttackDamage    ChampionStatInfo `json:"attackDamage"`
	} `json:"stats"`
	Abilities struct {
		Q []Ability `json:"Q"`
		W []Ability `json:"W"`
		E []Ability `json:"E"`
		R []Ability `json:"R"`
	} `json:"abilities"`
}
type ChampionStatInfo struct {
	Flat     float32 `json:"flat"`
	PerLevel float32 `json:"perLevel"`
}
type Ability struct {
	Effects []struct {
		Leveling []struct {
			Modifiers []struct {
				Values []float32 `json:"values"`
				Units  []string  `json:"units"`
			} `json:"modifiers"`
		} `json:"leveling"`
	} `json:"effects"`
	// PHYSICAL_DAMAGE/MAGIC_DAMAGE/TRUE_DAMAGE/OTHER_DAMAGE/null
	DamageType string `json:"damageType"`
}

type LivePlayerStats struct {
	Stats struct {
		AbilityHaste                 float32 `json:"abilityHaste"`
		AbilityPower                 float32 `json:"abilityPower"`
		Armor                        float32 `json:"armor"`
		ArmorPenetrationFlat         float32 `json:"armorPenetrationFlat"`
		ArmorPenetrationPercent      float32 `json:"armorPenetrationPercent"`
		AttackDamage                 float32 `json:"attackDamage"`
		AttackRange                  float32 `json:"attackRange"`
		AttackSpeed                  float32 `json:"attackSpeed"`
		BonusArmorPenetrationPercent float32 `json:"bonusArmorPenetrationPercent"`
		BonusMagicPenetrationPercent float32 `json:"bonusMagicPenetrationPercent"`
		CritChance                   float32 `json:"critChance"`
		CritDamage                   float32 `json:"critDamage"`
		CurrentHealth                float32 `json:"currentHealth"`
		HealShieldPower              float32 `json:"healShieldPower"`
		HealthRegenRate              float32 `json:"healthRegenRate"`
		LifeSteal                    float32 `json:"lifeSteal"`
		MagicLethality               float32 `json:"magicLethality"`
		MagicPenetrationFlat         float32 `json:"magicPenetrationFlat"`
		MagicPenetrationPercent      float32 `json:"magicPenetrationPercent"`
		MagicResist                  float32 `json:"magicResist"`
		MaxHealth                    float32 `json:"maxHealth"`
		MoveSpeed                    float32 `json:"moveSpeed"`
		Omnivamp                     float32 `json:"omnivamp"`
		PhysicalLethality            float32 `json:"physicalLethality"`
		PhysicalVamp                 float32 `json:"physicalVamp"`
		ResourceMax                  float32 `json:"resourceMax"`
		ResourceRegenRate            float32 `json:"resourceRegenRate"`
		ResourceType                 string  `json:"resourceType"`
		ResourceValue                float32 `json:"resourceValue"`
		SpellVamp                    float32 `json:"spellVamp"`
		Tenacity                     float32 `json:"tenacity"`
	} `json:"championStats"`
	AbilityLevels struct {
		Q AbilityLevel `json:"Q"`
		W AbilityLevel `json:"W"`
		E AbilityLevel `json:"E"`
		R AbilityLevel `json:"R"`
	} `json:"abilities"`
}
type AbilityLevel struct {
	AbilityLevel int `json:"abilityLevel"`
}

type LiveChampionInfo struct {
	RiotID  string `json:"riotId"`
	Level   int    `json:"level"`
	ItemIDs []struct {
		ItemID int `json:"itemID"`
	} `json:"items"`
}
type ItemInfo struct {
	LoadSuccess     bool
	Armor           ItemStatInfo `json:"armor"`
	MagicResistance ItemStatInfo `json:"magicResistance"`
}
type ItemStatInfo struct {
	Flat float32 `json:"flat"`
}

type CalculatedStats struct {
	Armor           float32 `json:"armor"`
	MagicResistance float32 `json:"magicResistance"`
}

// https://wiki.leagueoflegends.com/en-us/Champion_statistic
func scaleStat(growth float32, level float32) float32 {
	return growth * (level - 1.0) * (0.7025 + 0.0175*(level-1.0))
}

// https://wiki.leagueoflegends.com/en-us/Armor
func damageMitigation(initialDamage float32, damageType string, victim Summoner) float32 {
	switch damageType {
	case "PHYSICAL_DAMAGE":
		return initialDamage / (1 + victim.CalculatedStats.Armor/100)
	case "MAGIC_DAMAGE":
		return initialDamage / (1 + victim.CalculatedStats.MagicResistance/100)
	case "TRUE_DAMAGE":
		return initialDamage
	case "OTHER_DAMAGE":
		return 0.0
	default:
		return 0.0
	}
}

func bonusAD(summoner Summoner) float32 {
	startingAD := summoner.StaticChampionInfo.Stats.AttackDamage.Flat +
		scaleStat(summoner.StaticChampionInfo.Stats.AttackDamage.PerLevel, float32(summoner.LiveChampionInfo.Level))
	return summoner.LivePlayerStats.Stats.AttackDamage - startingAD
}

func calculateResistances(summoner []*Summoner) {
	for _, s := range summoner {
		s.CalculatedStats.Armor = s.StaticChampionInfo.Stats.Armor.Flat
		s.CalculatedStats.MagicResistance = s.StaticChampionInfo.Stats.MagicResistance.Flat
		for _, itemID := range s.LiveChampionInfo.ItemIDs {
			item := ItemInfoCache[itemID.ItemID]
			s.CalculatedStats.Armor += item.Armor.Flat
			s.CalculatedStats.MagicResistance += item.MagicResistance.Flat
		}
		s.CalculatedStats.Armor +=
			scaleStat(s.StaticChampionInfo.Stats.Armor.PerLevel, float32(s.LiveChampionInfo.Level))
		s.CalculatedStats.MagicResistance +=
			scaleStat(s.StaticChampionInfo.Stats.MagicResistance.PerLevel, float32(s.LiveChampionInfo.Level))
	}
}

func autoAttackDamage(attacker Summoner, victim Summoner) float32 {
	return damageMitigation(attacker.LivePlayerStats.Stats.AttackDamage, "PHYSICAL_DAMAGE", victim)
}

func abilityDamage(attacker Summoner, ability Ability, abilityLevel int, victim Summoner) float32 {
	var totalDamage float32
	if abilityLevel == 0 {
		return 0.0
	}
	for _, effect := range ability.Effects {
		for _, leveling := range effect.Leveling {
			for _, modifier := range leveling.Modifiers {
				if len(modifier.Values) >= abilityLevel && len(modifier.Units) >= abilityLevel {
					switch modifier.Units[abilityLevel-1] {
					case "":
						totalDamage += modifier.Values[abilityLevel-1]
					case "% AD":
						totalDamage += attacker.LivePlayerStats.Stats.AttackDamage * (modifier.Values[abilityLevel-1] / 100.0)
					case "% AP":
						totalDamage += attacker.LivePlayerStats.Stats.AbilityPower * (modifier.Values[abilityLevel-1] / 100.0)
					case "% bonus AD":
						totalDamage += bonusAD(attacker) * (modifier.Values[abilityLevel-1] / 100.0)
					case "%":
					default:
						return 0.0
					}
				} else {
					// return 0 damage to symbolize error if ability rank is higher than expected
					return 0.0
				}
			}
		}
	}
	return damageMitigation(totalDamage, ability.DamageType, victim)
}
