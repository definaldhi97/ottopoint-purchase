package dbmodels

type User struct {
	ID         int    `gorm:"id";pk json:"id"`
	Nama       string `gorm:"nama" json:"nama"`
	LastName   string `gorm:"last_name"`
	Phone      string `gorm:"phone" json:"phone"`
	CustID     string `gorm:"cust_id" json:"cust_id"`
	Email      string `gorm:"email" json:"email"`
	MerchantID string `gorm:"merchant_id" json:"merchant_id"`
	Password   string `gorm:"password" json:"password"`
	Status     bool   `gorm:"status" json:"status"`
	// CreatedAT time.Time `gorm:"created_at"`
	// UpdateAT  time.Time `gorm:"updated_at"`
}

func (t *User) TableName() string {
	return "public.users"
}
