package db

import (
	"ottopoint-purchase/models/dbmodels"

	"github.com/astaxie/beego/logs"
)

func UpdateVoucher(use, couponId string) (dbmodels.TransaksiRedeem, error) {
	res := dbmodels.TransaksiRedeem{}

	err := DbCon.Raw(`update redeem_transactions set trans_type = 'Used', used_at = ? where coupon_id = ?`, use, couponId).Scan(&res).Error
	if err != nil {
		logs.Info("Failed to UpdateVoucher from database", err)
		return res, err
	}
	logs.Info("Update Voucher :", res)

	return res, nil
}
