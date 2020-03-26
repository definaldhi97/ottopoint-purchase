package dbmodels

import "time"

type MParameters struct {
	ID        string    `gorm:"column:id"`
	Group     string    `gorm:"column:"group`
	Code      string    `gorm:"column:code"`
	Desc      string    `gorm:"column:desc"`
	Value     string    `gorm:"column:value"`
	AddValue1 string    `gorm:"column:add_value1"`
	AddValue2 string    `gorm:"column:add_value2"`
	CreatedAt time.Time `gorm:"column:created_at"`
	CreatedBy string    `gorm:"column:created_by"`
	UpdatedAt time.Time `gorm:"column:update_by"`
	UpdateBy  string    `gorm:"column:update_by"`
}

func (t *MParameters) TableName() string {
	return "public.m_parameters"
}
