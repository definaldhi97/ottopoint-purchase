package redeemtion

import (
	// lp "ottopoint-purchase/hosts/landing_page/host"
	"ottopoint-purchase/models"
	// "ottopoint-purchase/constants"
	// redeem "ottopoint-purchase/services/v2.1/Redeemtion"
)

func PaymentSplitBillServices(req models.PaymentSplitBillReq, param models.Params) models.Response {

	res := models.Response{}

	// printError := "\033[31m" // merah
	// printRes := "\033[34m"   // biru

	// fmt.Println(printRes, "[PaymentSplitBillServices]")

	// landingPage, errLP := lp.PaymentLandingPage("email", "firstname", "lastName", "phone", "merchantname", "trxId", 0)
	// if errLP != nil || landingPage.ResponseData.StatusCode != "00" {

	// 	fmt.Println(printError, "[PackageRedeemtion]-[PaymentSplitBillServices]")
	// 	fmt.Println(printError, fmt.Sprintf("[PaymentLandingPage]-[Error : %v]", errLP))

	// 	return res
	// }

	// resRedeem := models.Response{}

	// switch param.SupplierID {
	// // case constants.CODE_VENDOR_OTTOAG:
	// // 	fmt.Println(" [ Product OTTOAG ]")
	// // 	resRedeem = redeem.V21_VoucherOttoAg(req, param, header)
	// // case constants.CODE_VENDOR_UV:
	// // 	fmt.Println(" [ Product Ultra Voucher ]")
	// // 	resRedeem = redeem.V21_VoucherUV(req, param, header)
	// // case constants.CODE_VENDOR_SEPULSA:
	// // 	fmt.Println(" [ Product Sepulsa ]")
	// // 	resRedeem = redeem.V21_VoucherSepulsa(req, param, header)
	// case constants.CODE_VENDOR_AGREGATOR:
	// 	fmt.Println(" [ Product Agregator ]")
	// 	header.DeviceID = "H2H"
	// 	resRedeem = redeem.V21_VoucherAgServices.V21_VoucherAg(req, param, header)
	// }

	// if resRedeem.Meta.Code != 200 {
	// 	res = utils.GetMessageFailedErrorNew(res, constants.RC_FAILED_DECRYPT_VOUCHER, constants.RD_FAILED_DECRYPT_VOUCHER)
	// 	return res
	// }

	// res = models.Response {
	// 	Data : models.PaymentSplitBillResp{
	// 		Code      : "00",
	// 		Message   : "Success",
	// 		Success   : 1
	// 		Failed    : 0
	// 		Pending   : 0
	// 		UrlPayment: landingPage.ResponseData.EndpointURL,
	// 	},
	// 	Meta: utils.ResponseMetaOK(),
	// }

	return res

}
