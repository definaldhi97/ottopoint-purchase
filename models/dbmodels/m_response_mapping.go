package dbmodels

type MResponseMapping struct {
	InternalRc    string `gorm:"column:internal_rc"`
	InstitutionId string `gorm:"column:institution_id"`
	InstitutionRc string `gorm:"column:institution_id"`
	InstitutionRd string `gorm:"column:institution_id"`
}

func (t *MResponseMapping) TableName() string {
	return "public.m_response_mapping"
}
