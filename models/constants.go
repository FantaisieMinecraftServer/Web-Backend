package models

const (
	CategoryWeapon   = "weapon"
	CategoryFood     = "food"
	CategoryMaterial = "material"

	Dagger = "dagger"
	Sword  = "sword"
	Spear  = "spear"
	Hammer = "hammer"
	Wand   = "wand"
	Bow    = "bow"
)

var CategoryNames = map[string]string{
	CategoryWeapon:   "武器",
	CategoryFood:     "食料",
	CategoryMaterial: "素材",
}

var WeaponsGroup = map[string]string{
	Dagger: "短剣",
	Sword:  "刀剣",
	Spear:  "槍",
	Hammer: "ハンマー",
	Wand:   "杖",
	Bow:    "弓",
}
