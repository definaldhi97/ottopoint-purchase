package routers

import (
	"net/http"

	"fmt"
	"ottopoint-purchase/controllers"
	checkStatus "ottopoint-purchase/controllers/v2/check_status"
	earning "ottopoint-purchase/controllers/v2/earnings"
	vouchers "ottopoint-purchase/controllers/v2/vouchers"
	callbacks "ottopoint-purchase/controllers/v2/vouchers/callbacks"
	use_vouchers "ottopoint-purchase/controllers/v2/vouchers/use_vouchers"

	callback_v21 "ottopoint-purchase/controllers/v2.1/vouchers/callbacks"
	redeemv21 "ottopoint-purchase/controllers/v2.1/vouchers/redeemtion"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func rootFunc(c *gin.Context) {
	fmt.Println("Welcome to Ottopoint-Purchase api")

	c.JSON(http.StatusOK, "Welcome to Ottopoint-Purchase api")
}

func Server(portStr string) error {
	router := gin.New()

	// set permissions
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"},
		AllowHeaders:     []string{"Origin", "authorization", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "AppsId", "ChannelId", "InstitutionId", "DeviceId", "Geolocation", "Signature", "Timestamp"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	// declare controllers
	// HealthCheck := new(controllers.HealthCheckService)

	// routes
	apiRoot := router.Group("/transaction")
	{
		v2Root := apiRoot.Group("/v2")
		{
			// transfer.GET("/status", checkStatusControllers.GetCheckStatus)
			v2Root.GET("/healthcheck", controllers.HealthCheckService)
			// v2Root.GET("/redeem", controllers.VoucherRedeemController)
			v2Root.POST("/usevoucher", use_vouchers.UseVouchersControllers)
			v2Root.POST("/earningpoint", earning.EarningsPointController)
			v2Root.POST("/check-status-earning", earning.CheckStatusEarningController)
			v2Root.POST("/check-status-scheduler", checkStatus.SchedulerCheckStatusController)
			v2Root.GET("/getEarning", earning.GetEarningRuleController)
			v2Root.POST("/redeempoint", redeemv21.RedeemtionControllerV21)
			v2Root.POST("/status/sepulsa", callback_v21.CallBackSepulsa_V21_Controller)
			v2Root.POST("/redeem/voucherag", callbacks.CallbackVoucherAggController)
			v2Root.POST("/usevoucher_uv", callbacks.CallBackUVController)
		}

		v21Root := apiRoot.Group("/v2.1")
		{
			v21Root.GET("/voucher/view", vouchers.ViewVoucherController)
			v21Root.GET("/usevoucher/vidio", use_vouchers.UseVoucherVidioController)
		}

	}

	err := router.Run(":" + portStr)
	if err != nil {
		fmt.Println(err)
	}

	return err
}
