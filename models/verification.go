package module_models

type Code struct {
	Value     string `gorm:"unique"`
	ExpiresAt int64  `gorm:"autoCreateTime"`
}
