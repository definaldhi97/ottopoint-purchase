package vouchers

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	"strings"
)

func SwitchDataVoucher(data models.VoucherDetailsManagement) models.Params {
	var result models.Params

	var producrType string
	t := strings.ToLower(data.BrandName)
	switch t {
	case constants.CategoryFreeFire, constants.CategoryMobileLegend:
		producrType = "Game"
	default:
		producrType = data.BrandName
	}

	var categoriesID *string
	if len(data.CategoriesID) == 0 {
		fmt.Println("Not CategoriesID")
		categoriesID = nil
	} else {
		categoriesID = &data.CategoriesID[0]
	}
	result = models.Params{
		ProductType:         producrType,
		ProductCode:         data.ExternalProductCode,
		CouponCode:          data.ExternalProductCode,
		SupplierID:          data.CodeSuplier,
		NamaVoucher:         data.VoucherName,
		Point:               int(data.CostPoints),
		Category:            strings.ToLower(producrType),
		ExpDate:             data.ActivityActiveTo,
		CategoryID:          categoriesID,
		UsageLimitVoucher:   data.UsageLimit,
		ProductCodeInternal: data.InternalProductCode,
		RewardID:            data.RewardID,
		ProductID:           data.ProductID,
	}

	return result
}
