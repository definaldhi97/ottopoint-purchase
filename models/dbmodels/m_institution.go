package dbmodels

import "time"

type MInstution struct {
	ID               int       `gorm:"primary_key" json:"id"`
	PartnerID        string    `gorm:"column:partner_id" json:"partnerID"`
	Name             string    `gorm:"column:name" json:"name"`
	BrandName        string    `gorm:"column:brand_name" json:"brandName"`
	Email            string    `gorm:"column:email" json:"email"`
	UserType         int       `gorm:"column:user_type" json:"UserType"`
	Address          string    `gorm:"column:address" json:"address"`
	ResidenceAddress string    `gorm:"column:residence_address" json:"residenceAddress"`
	BusinessType     string    `gorm:"column:business_type" json:"businessType"`
	TaxNumber        string    `gorm:"column:tax_number" json:"taxNumber"`
	Phone            string    `gorm:"column:phone" json:"phone"`
	PicName          string    `gorm:"column:pic_name" json:"picName"`
	PicEmail         string    `gorm:"column:pic_email" json:"picEmail"`
	PicPhone         string    `gorm:"column:pic_phone" json:"picPhone"`
	Status           string    `gorm:"column:status" json:"status"` // draft, waiting_approve, approved
	ApproveDate      time.Time `gorm:"column:approve_date" json:"approveDate"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"createdAt"`
	CreatedBy        string    `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updatedAt"`
	UpdatedBy        string    `gorm:"column:updated_by" json:"updatedBy"`
	IsActive         bool      `gorm:"column:is_active" json:"isActive"`
	CallBackUrl      string    `gorm:"column:callback_url"`
	NOtificationID   int       `gorm:"column:notification_id"`
}

func (t *MInstution) TableName() string {
	return "public.m_institution"
}
