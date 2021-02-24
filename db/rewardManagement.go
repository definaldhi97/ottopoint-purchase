package db

import (
	"fmt"
	"log"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"strings"
)

func GetVoucherDetails(rewardID string) (models.VoucherDetailsManagement, error) {
	fmt.Println(">>>>>>>> Get Voucher/Reward Details <<<<<<<")
	var result models.VoucherDetailsManagement
	var dtaVoucher models.VoucherDetailsManagement1

	// err := DbCon.Raw(`select * from product.m_reward where id = ?`, rewardID).Scan(&result).Error
	err := DbCon.Raw(`select 
	a.id as reward_id, 
	a."name", 
	a.cost_in_points, 
	a.usage_limit, 
	d."name" as brand_name, 
	a.activity_active_from, 
	a.activity_active_to, 
	a.categories as categories_id, 
	c.code as code_suplier,  
	a.reward_codes, 
	b.external_code, 
	b.code as internal_code,
	b.id as m_product_id,
	array(
		select array_to_string(array_agg(mf.code), ',') from product.m_field_brand mfb
		left join product.m_field mf on mf.id = mfb.m_field_id
		where mfb.m_product_brand_id = b.m_product_brand_id
		group by mfb.sort_order
		order by mfb.sort_order asc
	) fields
	from product.m_reward a join product.m_product b on (a.m_product_id = b.id)
	join vendor.m_vendor c on (b.m_vendor_id = c.id)
	join product.m_product_brand d on (d.id = b.m_product_brand_id)  where a.id = ?`, rewardID).Scan(&dtaVoucher).Error
	if err != nil {
		fmt.Println("Failed Get Voucher/Reward Details : ", err)
		return result, err
	}

	replaceStr := strings.ReplaceAll(dtaVoucher.CategoriesID, "[", "")
	replaceStr = strings.ReplaceAll(replaceStr, "]", "")
	replaceStr = strings.ReplaceAll(replaceStr, "\"", "")
	replaceStr = strings.ReplaceAll(replaceStr, ",", " ")
	CategoriesIDArray := strings.Fields(replaceStr)

	result.RewardID = dtaVoucher.RewardID
	result.VoucherName = dtaVoucher.VoucherName
	result.CostPoints = dtaVoucher.CostPoints
	result.UsageLimit = dtaVoucher.UsageLimit
	result.BrandName = dtaVoucher.BrandName
	result.ActivityActiveFrom = dtaVoucher.ActivityActiveFrom
	result.ActivityActiveTo = dtaVoucher.ActivityActiveTo
	result.CategoriesID = CategoriesIDArray
	result.CodeSuplier = dtaVoucher.CodeSuplier
	result.RewardCodes = dtaVoucher.RewardCodes
	result.ExternalProductCode = dtaVoucher.ExternalProductCode
	result.InternalProductCode = dtaVoucher.InternalProductCode
	result.ProductID = dtaVoucher.ProductID

	fmt.Println(result)

	return result, nil

}

func UpdateUsageLimitVoucher(rewadID string, latestLimitVoucher int) error {

	var result dbmodels.MRewardModel
	err := DbCon.Model(result).Unscoped().Where("id = ?", rewadID).Update(map[string]interface{}{"usage_limit": latestLimitVoucher}).Error
	if err != nil {
		log.Print("Failed Update h_Bulk", err)
		return err
	}
	return nil
}

func Get_MReward(id string) (dbmodels.MRewardModel, error) {
	var result dbmodels.MRewardModel

	err := DbCon.Raw(`select * from product.m_reward where id = ?`, id).Scan(&result).Error
	if err != nil {
		fmt.Println("Failed Get Voucher/Reward : ", err)
		return result, err
	}

	return result, nil

}
