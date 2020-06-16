package dbmodels

type MResponseInternal struct {
	InternalRc string `gorm:"column:internal_rc"`
	InternalRd string `gorm:"column:internal_rd"`
}

func (t *MResponseInternal) TableName() string {
	return "public.m_response_internal"
}
