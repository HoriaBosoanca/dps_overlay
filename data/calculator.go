package data

type Summoner struct {
	RiotID               string `json:"riotId"`
	Team                 string `json:"team"`
	LiveStats            LiveStats
	ChampionName         string `json:"championName"`
	ChampionImportedInfo ChampionImportedInfo
	Level                int `json:"level"`
	ItemIDs              []struct {
		ItemID int `json:"itemID"`
	} `json:"items"`
	CalculatedStats CalculatedStats
}
type ChampionImportedInfo struct {
	LoadedSuccess   bool
	Armor           ChampionImportedStatInfo `json:"armor"`
	MagicResistance ChampionImportedStatInfo `json:"magicResistance"`
}
type ChampionImportedStatInfo struct {
	Flat     float32 `json:"flat"`
	PerLevel float32 `json:"perLevel"`
}
type ItemImportedInfo struct {
	LoadSuccess     bool
	Armor           ItemImportedStatInfo `json:"armor"`
	MagicResistance ItemImportedStatInfo `json:"magicResistance"`
}
type ItemImportedStatInfo struct {
	Flat float32 `json:"flat"`
}
type CalculatedStats struct {
	Armor           float32 `json:"armor"`
	MagicResistance float32 `json:"magicResistance"`
}
type LiveStats struct {
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
}

// CalculateResistances https://wiki.leagueoflegends.com/en-us/Champion_statistic
func calculateResistances(summoner []*Summoner) {
	for _, enemy := range summoner {
		enemy.CalculatedStats.Armor = enemy.ChampionImportedInfo.Armor.Flat
		enemy.CalculatedStats.MagicResistance = enemy.ChampionImportedInfo.MagicResistance.Flat
		for _, itemID := range enemy.ItemIDs {
			item := ItemInfoCache[itemID.ItemID]
			enemy.CalculatedStats.Armor += item.Armor.Flat
			enemy.CalculatedStats.MagicResistance += item.MagicResistance.Flat
		}
		enemy.CalculatedStats.Armor += enemy.ChampionImportedInfo.Armor.PerLevel * float32(enemy.Level-1) * (0.7025 + 0.0175*float32(enemy.Level-1))
		enemy.CalculatedStats.MagicResistance += enemy.ChampionImportedInfo.MagicResistance.PerLevel * float32(enemy.Level-1) * (0.7025 + 0.0175*float32(enemy.Level-1))
	}
}

// autoAttackDamage https://wiki.leagueoflegends.com/en-us/Armor
func autoAttackDamage(attacker Summoner, victim Summoner) float32 {
	return attacker.LiveStats.AttackDamage / (1 + victim.CalculatedStats.Armor/100)
}
