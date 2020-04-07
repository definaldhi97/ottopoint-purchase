package db

import (
	"fmt"
)

type Count struct {
	Count         int    `gorm:"count" json:"count"`
	AccountNumber string `gorm:"account_number"`
}

func GetCountInquiryGagal(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("SELECT count(*) as count FROM redeem_transactions where cummulative_ref = ? and trans_type = 'Inquiry' and responder_data = '01'", cummulative_ref).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[Get Count Total Failed Pyenment]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return res, err
	}
	return res, nil

}

func GetCountSucc_Pyenment(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("SELECT count(*) as count FROM redeem_transactions where cummulative_ref = ? and trans_type = 'Payment' and responder_data = '00'", cummulative_ref).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[Get Count Total Failed Pyenment]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return res, err
	}
	return res, nil

}

func GetCountPending_Pyenment(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("select count(*) as count from public.redeem_transactions where responder_data = '09' and cummulative_ref = ? and trans_type = 'Payment'", cummulative_ref).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[Get Count Total Failed Pyenment]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return res, err
	}
	return res, nil

}

func GetCountFailedPyenment(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("select count(*) as count from public.redeem_transactions where responder_data = '01' and cummulative_ref = ? and trans_type = 'Payment'", cummulative_ref).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[Get Count Total Failed Pyenment]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return res, err
	}
	return res, nil

}

func GetCountPyenment(cummulative_ref string) (Count, error) {
	res := Count{}
	err := DbCon.Raw("select count(*) as count from public.redeem_transactions where cummulative_ref = ?", cummulative_ref).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[Get Count Total Failed Pyenment]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return res, err
	}
	return res, nil

}

func GetPyenmentFailed(cummulative_ref string) (Count, error) {
	res := Count{}

	// err := DbCon.Exec(`select * from users where phone = ?, status = true`, phone).Scan(&res).Error
	err := DbCon.Raw("select account_number, sum(point) as count from redeem_transactions where responder_data not in ('00','09','68') and cummulative_ref = ? and trans_type ='Payment' group by account_number", cummulative_ref).Scan(&res).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[Get Failed Pyenment]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))
		return res, err
	}

	return res, nil
}
