package routers

import (
	"fmt"
	"io"
	"os"
	"ottopoint-purchase/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"ottodigital.id/library/httpserver/ginserver"
	ottologer "ottodigital.id/library/logger"
	"ottodigital.id/library/ottotracing"
	"ottodigital.id/library/utils"
)

var (
	redeem             string
	use_voucher        string
	deductPoint        string
	paymentQR          string
	reversePoint       string
	healthcheck        string
	earningPoint       string
	splitbill          string
	comulative         string
	usevoucher_uv      string
	checkStatusEarning string

	nameservice     string
	agentracinghost string

	debugmode string
	readto    int
	writeto   int
)

func init() {
	//TODO pls change UPERCASE & _ Not using dot

	healthcheck = utils.GetEnv("healthcheck", "/transaction/v2/healthcheck")
	redeem = utils.GetEnv("redeem", "/transaction/v2/redeem")
	use_voucher = utils.GetEnv("use_voucher", "/transaction/v2/usevoucher")
	comulative = utils.GetEnv("comulative", "/transaction/v2/redeempoint")
	deductPoint = utils.GetEnv("deduct_point", "/transaction/v2/deduct")
	reversePoint = utils.GetEnv("reverse_point", "/transaction/v2/reversal")
	earningPoint = utils.GetEnv("earning_point", "/transaction/v2/earningpoint")
	splitbill = utils.GetEnv("splitbill", "/transaction/v2/splitbill")
	usevoucher_uv = utils.GetEnv("usevoucher_uv", "/transaction/v2/usevoucher_uv")
	checkStatusEarning = utils.GetEnv("usevoucher_uv", "/check-status-earning")

	debugmode = utils.GetEnv("apps.debug", "debug")

	nameservice = utils.GetEnv("OTTOPOINT_PURCHASE_NAMESERVICE", "ottopoint-purchase")

	agentracinghost = utils.GetEnv("AGENT_TRACING_HOST_OTTOPOINT_PURCHASE", "13.250.21.165:5775")
	// readto = utils.GetEnv("server.readtimeout", 30)
	// writeto = utils.GetEnv("server.writetimeout", 30)

}

func Server(listenAddr string) error {

	ottoRouter := OttoRouter{}
	ottoRouter.InitTracing()
	ottoRouter.Routers()
	defer ottoRouter.Close()

	err := ginserver.GinServerUp(listenAddr, ottoRouter.Router)

	if err != nil {
		fmt.Println("Error : ", err)

		return err
	}
	sugarLogger := ottologer.GetLogger()
	sugarLogger.Info("Server UP ", zap.String("Address", listenAddr))

	return err

}

type OttoRouter struct {
	Tracer   opentracing.Tracer
	Reporter jaeger.Reporter
	Closer   io.Closer
	Err      error
	Ginfunc  gin.HandlerFunc
	Router   *gin.Engine
}

func (ottoRouter *OttoRouter) Routers() {
	gin.SetMode(debugmode)

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"},
		AllowHeaders:     []string{"Origin", "authorization", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		//AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge: 86400,
	}))

	router.Use(ottoRouter.Ginfunc)
	router.Use(gin.Recovery())

	// router.GET(cashbackbyproduct, controllers.InquiryController)
	router.GET(healthcheck, controllers.HealthCheckService)
	router.POST(redeem, controllers.VoucherRedeemController)
	router.POST(comulative, controllers.VoucherComulativeController)
	router.POST(use_voucher, controllers.UseVouhcerController)
	router.POST(deductPoint, controllers.PointController)
	router.POST(reversePoint, controllers.ReversePointController)
	router.POST(earningPoint, controllers.EarningsPointController)
	router.POST(splitbill, controllers.DeductSplitBillController)
	router.POST(usevoucher_uv, controllers.UseVouhcerUVController)
	router.POST(checkStatusEarning, controllers.CheckStatusEarningController)

	ottoRouter.Router = router

}

func (ottoRouter *OttoRouter) InitTracing() {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "PROD"
	}
	tracer, reporter, closer, err := ottotracing.InitTracing(fmt.Sprintf("%s::%s", nameservice, hostName), agentracinghost, ottotracing.WithEnableInfoLog(true))
	if err != nil {
		fmt.Println("Error :", err)
	}
	opentracing.SetGlobalTracer(tracer)

	ottoRouter.Closer = closer
	ottoRouter.Reporter = reporter
	ottoRouter.Tracer = tracer
	ottoRouter.Err = err
	ottoRouter.Ginfunc = ottotracing.OpenTracer([]byte("api-request-"))

}

func (ottoRouter *OttoRouter) Close() {
	ottoRouter.Closer.Close()
	ottoRouter.Reporter.Close()
}
