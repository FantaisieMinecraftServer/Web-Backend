package models

type Category struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	DisplayItem DisplayItem `json:"display_item"`
	Items       []Item      `json:"items"`
}

type Item struct {
	ID              string      `bson:"_id,omitempty" json:"id"`
	Type            string      `bson:"type" json:"type"`
	Name            string      `bson:"name" json:"name"`
	Lore            []string    `bson:"lore" json:"lore"`
	Rarity          int         `bson:"rarity" json:"rarity"`
	MaxStackSize    int         `bson:"max_stack_size" json:"max_stack_size"`
	ItemID          string      `bson:"item_id" json:"item_id"`
	CustomModelData int         `bson:"custom_model_data" json:"custom_model_data"`
	Prices          PriceData   `bson:"prices" json:"prices"`
	Data            interface{} `bson:"data" json:"data"`
}

// 個別

type DisplayItem struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	CustomModelData int    `json:"custom_model_data"`
}

type PriceData struct {
	Purchase   int  `bson:"purchase" json:"purchase"`
	Selling    int  `bson:"selling" json:"selling"`
	CanSelling bool `bson:"can_selling" json:"can_selling"`
}

type Effect struct {
	Name      string `bson:"name" json:"name"`
	Duration  int    `bson:"duration" json:"duration"`
	Amplifier int    `bson:"amplifier" json:"amplifier"`
}

// 食料

type FoodData struct {
	Nutrition    int      `bson:"nutrition" json:"nutrition"`
	Saturation   int      `bson:"saturation" json:"saturation"`
	CanAlwaysEat bool     `bson:"can_always_eat" json:"can_always_eat"`
	EatSeconds   float64  `bson:"eat_seconds" json:"eat_seconds"`
	Effects      []Effect `bson:"effects" json:"effects"`
}

// 武器

type WeaponData struct {
	Group    string         `bson:"group" json:"group"`
	Crafting []CraftingData `bson:"crafting" json:"crafting"`
	Upgrades []UpgradeData  `bson:"upgrades" json:"upgrades"`
}

type CraftingData struct {
	Amount    int        `bson:"amount" json:"amount"`
	Materials []Material `bson:"materials" json:"materials"`
}

type UpgradeData struct {
	Level     int        `bson:"level" json:"level"`
	Cost      int        `bson:"cost" json:"cost"`
	Materials []Material `bson:"materials" json:"materials"`
}

type Material struct {
	Amount int  `bson:"amount" json:"amount"`
	Data   Item `bson:"data" json:"data"`
}

// 素材

type MaterialData struct {
}
