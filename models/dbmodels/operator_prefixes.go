package dbmodels

type OperatorPrefixes struct {
	ID           int    `gorm:"column:id;primary_key"`
	Prefix       string `gorm:"prefix"`
	OperatorCode int    `gorm:"operator_code"`
	CreatedAt    string `gorm:"created_at"`
	UpdatedAt    string `gorm:"updated_at"`
}

func (t *OperatorPrefixes) TableName() string {
	return "operator_prefixes"
}
