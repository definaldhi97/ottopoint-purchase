package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"
)

func GetResposeInternal(rc string) (dbmodels.MResponseInternal, error) {
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
