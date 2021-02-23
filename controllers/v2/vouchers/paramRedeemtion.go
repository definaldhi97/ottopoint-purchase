package vouchers

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/controllers"
	"ottopoint-purchase/models"
	"strings"
)

func ParamRedeemtion(custId string, data models.VoucherDetailsManagement) models.Params {

	fmt.Println("[Start]-[ParamRedeemtion]")

	var result models.Params

	var category string
	t := strings.ToLower(data.BrandName)
	switch t {
	case constants.CategoryFreeFire, constants.CategoryMobileLegend:
		category = "Game"
	default:
		category = data.BrandName
	}

	var categoriesID *string
	if len(data.CategoriesID) == 0 {
		fmt.Println("Not CategoriesID")
		categoriesID = nil
	} else {
		categoriesID = &data.CategoriesID[0]
	}

	if data.CodeSuplier == constants.CODE_VENDOR_OTTOAG {
		validateVerfix := controllers.ValidatePerfix(custId, data.ExternalProductCode, category)
		if validateVerfix == false {

			fmt.Println("Invalid verfix")
			result.ResponseCode = 500

			// result = utils.GetMessageResponse(result, 500, false, errors.New("Nomor akun ini tidak terdafatr"))
			// ctx.JSON(http.StatusOK, res)
			return result
		}

		if category == "" {

			fmt.Println("Invalid Category")
			result.ResponseCode = 500

			// result = utils.GetMessageResponse(result, 500, false, errors.New("Invalid BrandName"))
			// ctx.JSON(http.StatusOK, res)
			return result
		}
	}

	result = models.Params{
		// AccountNumber:       dataToken.Data,
		// MerchantID:          dataUser.MerchantID,
		ResponseCode:        200,
		ProductType:         category,
		ProductCode:         data.ExternalProductCode,
		CouponCode:          data.ExternalProductCode,
		SupplierID:          data.CodeSuplier,
		NamaVoucher:         data.VoucherName,
		Point:               int(data.CostPoints),
		Category:            strings.ToLower(category),
		ExpDate:             data.ActivityActiveTo,
		CategoryID:          categoriesID,
		UsageLimitVoucher:   data.UsageLimit,
		ProductCodeInternal: data.InternalProductCode,
		RewardID:            data.RewardID,
		ProductID:           data.ProductID,

		// InstitutionID:       header.InstitutionID,
		// AccountId: custId,
		// CampaignID:          campaignId,
	}

	return result
}
