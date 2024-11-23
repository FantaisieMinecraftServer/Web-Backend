package models

const (
	TypeWeapon   = "weapon"
	TypeFood     = "food"
	TypeMaterial = "material"

	Dagger = "dagger"
	Sword  = "sword"
	Spear  = "spear"
	Hammer = "hammer"
	Wand   = "wand"
	Bow    = "bow"
)

var TypeNames = map[string]string{
	TypeWeapon:   "武器",
	TypeFood:     "食料",
	TypeMaterial: "素材",
}

var WeaponClasses = map[string]string{
	Dagger: "短剣",
	Sword:  "刀剣",
	Spear:  "槍",
	Hammer: "ハンマー",
	Wand:   "杖",
	Bow:    "弓",
}

var WeaponClassesCMD = map[string]int{
	Dagger: 11000,
	Sword:  11001,
	Spear:  11002,
	Hammer: 11003,
	Wand:   11004,
	Bow:    11005,
}
