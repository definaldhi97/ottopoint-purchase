package db

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models/dbmodels"
)

type Count struct {
	Count         int    `gorm:"column:count" json:"count"`
	AccountNumber string `gorm:"column:account_number"`
	CountFailed   int    `gorm:"column:count_failed"`
}

type MappingRc struct {
	InternalRc string `gorm:"column:internal_rc"`
	InternalRd string `gorm:"column:internal_rd"`
}

func GetCountInquiryGagal(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("SELECT count(*) as count FROM t_spending where cummulative_ref = ? and trans_type = 'Inquiry' and status = '01'", cummulative_ref).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[GetCountInquiryGagal]")
		fmt.Println("[Get Count Total Failed Inquiry]")
		fmt.Println(fmt.Sprintf("Failed to connect t_spending %v", err))

		return res, err
	}
	return res, nil

}

func GetCountSucc_Pyenment(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("SELECT count(*) as count FROM t_spending where cummulative_ref = ? and trans_type = ? and status = '00'", cummulative_ref, constants.CODE_TRANSTYPE_REDEMPTION).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[GetCountSucc_Pyenment]")
		fmt.Println("[Get Count Total Success Pyenment]")
		fmt.Println(fmt.Sprintf("Failed to connect t_spending %v", err))

		return res, err
	}
	return res, nil

}

func GetCountPending_Pyenment(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("select count(*) as count from public.t_spending where status = '09' and cummulative_ref = ? and trans_type = ?", cummulative_ref, constants.CODE_TRANSTYPE_REDEMPTION).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[GetCountPending_Pyenment]")
		fmt.Println("[Get Count Total Pending Pyenment]")
		fmt.Println(fmt.Sprintf("Failed to connect t_spending %v", err))

		return res, err
	}
	return res, nil

}

func GetCountFailedPyenment(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("select count(*) as count from public.t_spending where status = '01' and cummulative_ref = ? and trans_type = ?", cummulative_ref, constants.CODE_TRANSTYPE_REDEMPTION).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[GetCountFailedPyenment]")
		fmt.Println("[Get Count Total Failed Pyenment]")
		fmt.Println(fmt.Sprintf("Failed to connect t_spending %v", err))

		return res, err
	}
	return res, nil

}

func GetCountPyenment(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("select count(*) as count from public.t_spending where cummulative_ref = ? and trans_type = ?", cummulative_ref, constants.CODE_TRANSTYPE_REDEMPTION).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[GetCountPyenment]")
		fmt.Println("[Get Count Total Transaksi]")
		fmt.Println(fmt.Sprintf("Failed to connect t_spending %v", err))

		return res, err
	}
	return res, nil

}

func GetPyenmentFailed(cummulative_ref string) (Count, error) {
	res := Count{}

	// err := DbCon.Exec(`select * from users where phone = ?, status = true`, phone).Scan(&res).Error
	err := DbCon.Raw("select account_number, sum(point) as count, count(*) as count_failed from t_spending where status not in ('00','09','68') and cummulative_ref = ? and trans_type = ? group by account_number", cummulative_ref, constants.CODE_TRANSTYPE_REDEMPTION).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[GetPyenmentFailed]")
		fmt.Println("[Get Sum Pyenment Failed]")
		fmt.Println(fmt.Sprintf("Failed to connect database %v", err))
		return res, err
	}

	return res, nil
}

func GetResponseOttoag(issuer, rc string) (MappingRc, error) {
	res := MappingRc{}

	err := DbCon.Raw("select b.internal_rc,b.internal_rd from m_response_mapping a join m_response_internal b on (a.internal_rc=b.internal_rc) where a.institution_id = ? and a.institution_rc = ?", issuer, rc).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[GetResponseOttoag]")
		fmt.Println("[Get Data Mapping Response]")
		fmt.Println(fmt.Sprintf("Failed to connect t_spending %v", err))
		return res, err
	}

	return res, nil

}

func GetResponseCummulativeOttoAG(rc string) (dbmodels.MResponseInternal, error) {
	res := dbmodels.MResponseInternal{}

	err := DbCon.Raw("select * from m_response_internal where internal_rc = ?", rc).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[GetResponseOttoag]")
		fmt.Println("[Get Data Mapping Response]")
		fmt.Println(fmt.Sprintf("Failed to connect t_spending %v", err))
		return res, err
	}

	return res, nil

}
