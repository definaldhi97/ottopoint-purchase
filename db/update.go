package db

import (
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models/dbmodels"
	"time"

	"github.com/astaxie/beego/logs"
)

func UpdateVoucher(use time.Time, couponId string) (dbmodels.TSpending, error) {
	res := dbmodels.TSpending{}

	err := DbCon.Raw(`update t_spending set is_used = true, used_at = ? where coupon_id = ?`, use, couponId).Scan(&res).Error
	if err != nil {
		logs.Info("Failed to UpdateVoucher from database", err)
		return res, err
	}
	logs.Info("Update Voucher :", res)

	return res, nil
}

func UpdateVoucherSepulsa(status, respDesc, reqData, transactionID, orderID string) (dbmodels.TSpending, error) {
	res := dbmodels.TSpending{}
	var internalCode string
	switch status {
	case "Success":
		internalCode = "00"
	case "Pending":
		internalCode = "09"
	default:
		internalCode = "01"
	}

	err := DbCon.Raw(`update t_spending set responder_rd = ?, responder_rc = ?, responder_data = ?, status = ? where rrn = ? and transaction_id = ?`, status, respDesc, reqData, internalCode, transactionID, orderID).Scan(&res).Error
	if err != nil {
		logs.Info("Failed to UpdateVoucherSepulsa from database", err)
		return res, err
	}
	logs.Info("Update Voucher :", res)

	return res, nil
}

func UpdateTSchedulerRetry(transactionID string) (dbmodels.TSchedulerRetry, error) {
	res := dbmodels.TSchedulerRetry{}

	err := DbCon.Raw(`update t_scheduler_retry set is_done = true where transaction_id = ?`, transactionID).Scan(&res).Error
	if err != nil {
		logs.Info("Failed to UpdateTSchedulerRetry from database", err)
		return res, err
	}
	logs.Info("Update Scheduler Retry :", res)
	return res, nil
}

func UpdateTEarning(pointId, id string) (dbmodels.TEarning, error) {
	res := dbmodels.TEarning{}

	err := DbCon.Exec(`update t_earning set points_transfer_id = ? where id = ?`, pointId, id).Scan(&res).Error
	if err != nil {
		logs.Info("Failed to UpdateVoucher from database", err)

		return res, err
	}

	return res, nil
}

func UpdateVoucherAg(redeemDate, usedDate, spendingID string) (dbmodels.TSpending, error) {
	res := dbmodels.TSpending{}

	err := DbCon.Raw(
		"update t_spending set redeem_at = ?, used_at = ?, is_used = true where id = ?",
		redeemDate, usedDate, spendingID,
	).Scan(&res).Error

	if err != nil {
		logs.Info("Failed to UpdateVoucher from database", err)
		return res, err
	}

	return res, nil

}

func UpdateTSchedulerVoucherAG(transactionID string) error {

	err := DbCon.Raw(
		"update t_scheduler_retry set is_done = true where transaction_id = ? and code = ?",
		transactionID, constants.CodeSchedulerVoucherAG,
	).Error

	if err != nil {
		logs.Info("Failed to UpdateVoucher from database", err)
		return err
	}

	return nil
}

func UpdateVoucherAgSecond(status, respDesc, tspendingID string) (dbmodels.TSpending, error) {
	res := dbmodels.TSpending{}
	var internalCode string
	switch status {
	case "Success":
		internalCode = "00"
	case "Pending":
		internalCode = "09"
	default:
		internalCode = "01"
	}

	err := DbCon.Raw(`update t_spending set responder_rd = ?, responder_rc = ?, status = ? where id = ?`, status, respDesc, internalCode, tspendingID).Scan(&res).Error
	if err != nil {
		logs.Info("Failed to UpdateVoucherSepulsa from database", err)
		return res, err
	}
	logs.Info("Update Voucher :", res)

	return res, nil
}
