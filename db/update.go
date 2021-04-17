package db

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"time"

	"github.com/sirupsen/logrus"
)

func UpdateVoucher(use time.Time, couponId string) (dbmodels.TSpending, error) {
	res := dbmodels.TSpending{}

	err := DbCon.Raw(`update t_spending set is_used = true, used_at = ? where coupon_id = ?`, use, couponId).Scan(&res).Error
	if err != nil {
		logrus.Info("Failed to UpdateVoucher from database", err)
		return res, err
	}
	logrus.Info("Update Voucher :", res)

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
		logrus.Info("Failed to UpdateVoucherSepulsa from database", err)
		return res, err
	}
	logrus.Info("Update Voucher :", res)

	return res, nil
}

func UpdateTSchedulerRetry(transactionID string) (dbmodels.TSchedulerRetry, error) {
	res := dbmodels.TSchedulerRetry{}

	err := DbCon.Raw(`update t_scheduler_retry set is_done = true where transaction_id = ?`, transactionID).Scan(&res).Error
	if err != nil {
		logrus.Info("Failed to UpdateTSchedulerRetry from database", err)
		return res, err
	}
	logrus.Info("Update Scheduler Retry :", res)
	return res, nil
}

func UpdateTEarning(pointId, id string) (dbmodels.TEarning, error) {
	res := dbmodels.TEarning{}

	err := DbCon.Exec(`update t_earning set points_transfer_id = ? where id = ?`, pointId, id).Scan(&res).Error
	if err != nil {
		logrus.Info("Failed to UpdateVoucher from database", err)

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
		logrus.Info("Failed to UpdateVoucher from database", err)
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
		logrus.Info("Failed to UpdateVoucher from database", err)
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
		logrus.Info("Failed to UpdateVoucherSepulsa from database", err)
		return res, err
	}
	logrus.Info("Update Voucher :", res)

	return res, nil
}

type VoucherTypeDB struct {
	VoucherType int

	OrderId      string
	ResponseCode string
	ResponseDesc string

	VoucherId    string
	VoucherCode  string
	VoucherName  string
	IsRedeemed   bool
	RedeemedDate string
}

func UpdateVoucherbyVoucherType(req VoucherTypeDB, trxId string) error {
	res := dbmodels.TSpending{}
	logrus.Println(">>> UpdateVoucherbyVoucherType <<<")

	var status string
	var err error

	switch req.ResponseCode {
	case "00":
		status = "00"
	case "09", "68":
		status = "09"
	default:
		status = "01"
	}

	if trxId == "" {

		logrus.Error("[PackageDB]-[UpdateVoucherbyVoucherType]")
		logrus.Error(fmt.Sprintf("[TrxID Kosong]", trxId))

		return err
	}

	// PPOB
	if req.VoucherType == 1 {

		logrus.Println(">>> VoucherType PPOB <<<")

		if req.OrderId != "" {

			err = DbCon.Raw(
				"update t_spending set responder_rc = ?, responder_rd = ?, status = ?, rrn = ?, updated_at = ? where cummulative_ref = ?",
				req.ResponseCode, req.ResponseDesc, status, req.OrderId, time.Now(), trxId).Scan(&res).Error

		} else {
			err = DbCon.Raw(
				"update t_spending set responder_rc = ?, responder_rd = ?, status = ?, updated_at = ? where cummulative_ref = ?",
				req.ResponseCode, req.ResponseDesc, status, time.Now(), trxId).Scan(&res).Error
		}

	}

	// Voucher Code
	if req.VoucherType == 2 {
		err = DbCon.Raw(
			"update t_spending set is_used = ?, rrn = ?, voucher_code = ?, used_at = ?, updated_at = ? where cummulative_ref = ?",
			req.IsRedeemed, req.OrderId, req.VoucherCode, req.RedeemedDate, time.Now(), trxId).Scan(&res).Error
	}

	if err != nil {

		logrus.Error("[PackageDB]-[UpdateVoucherbyVoucherType]")
		logrus.Error(fmt.Sprintf("[Failed get Data by TrxID : %v from TSpending]-[Error : %v]", trxId, err))

		return err
	}

	return nil
}

func UpdateTrxVoucher(param models.Params, trxId, status string) error {
	res := dbmodels.TSpending{}

	logrus.Println(fmt.Sprintf("[Start]-[UpdateTrxVoucherAG]-[%v]", trxId))

	err := DbCon.Exec(`update t_spending set responder_data = ?, requestor_data = ?, responder_rc = ?, responder_rd, status = ? where transaction_id = ?`,
		param.DataSupplier.Response, param.DataSupplier.Request, param.DataSupplier.Rc, param.DataSupplier.Rd, status, trxId).Scan(&res).Error
	if err != nil {

		logrus.Error("[PackageDB]-[UpdateTrxVoucherAG]")
		logrus.Error(fmt.Sprintf("[Error : %v]-[TrxID : %v]", err, trxId))

		return err
	}

	return nil
}
