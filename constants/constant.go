package constants

const (
	FullPointMethod = 0
	SplitBillMethod = 1
)

const (
	FieldNoHP        = "001"
	FieldIDPelanggan = "002"
	FieldIDServer    = "003"
	FieldNoKartu     = "004"
)

const (
	NotifTypeCashPayment = "0"
	NotifTypeQRPayment   = "1"
)

const (
	OTTOPAY   = "OTTOPAY"
	INDOMARCO = "INDOMARCO"
	BOGASARI  = "BOGASARI"
	PEDE      = "PEDE"
	OTTOSG    = "OTTOSG"
)

const (
	TokenKeyRedis       = "OTTOPAY-SESSION-TOKEN:%s"
	FdsErrorCodeSession = "IB-1010"
)

const (
	MAXUDP = 60000
)

const (
	SplitBill = "split_bill"
	Point     = "point"
)

const (
	Transfer = "transfer"
	Spend    = "spend"
)

const (
	CodeSplitBill     = "TSP06"
	CodeRedeemtion    = "TSP02"
	CodeReversal      = "TAD04"
	PaymentSplitBill  = "PCS01"
	SpendingSplitBill = "TSP05"
)

const (
	KeyResponsePending = "pending"
	KeyResponseDefault = "default"
	KeyResponseFailed  = "failed"
	KeyResponseSucceed = "succeed"
)

const (
	CategoryPulsa        = "pulsa"
	CategoryPaketData    = "paket data"
	CategoryFreeFire     = "free fire"
	CategoryMobileLegend = "mobile legends"
	CategoryPLN          = "pln"
	CategoryGame         = "game"
	CategoryVidio        = "vidio"
)

const (
	Success = "00"
	Pending = "09"
	Failed  = "01"
	TimeOut = "68"
)

const (
	UV            = "Ultra Voucher"
	OttoAG        = "OttoAG"
	Sepulsa       = "Sepulsa"
	VoucherAg     = "Voucher Aggregator"
	JempolKios    = "Jempol Kios"
	GudangVoucher = "Gudang Voucher"
	SecurePage    = "SecurePage"
)

const (
	GeneralSpending  = "GSR"
	Multiply         = "MEP"
	InstantReward    = "IRR"
	CustomeEventRule = "CER"
	EventRule        = "ERC"
	CustomerReferral = "CRR"
)

const (
	RC_ERROR_HEADER_MANDATORY = 61
	RD_ERROR_HEADER_MANDATORY = "Invalid header mandatory"

	RC_ERROR_USER_INACTIVE = 66
	RD_ERROR_USER_INACTIVE = "Account is Inactive"

	RC_ERROR_INVALID_TOKEN = 60
	RD_ERROR_INVALID_TOKEN = "Token or Session Expired Please Login Again"

	RC_ERROR_USER_LINKED_INACTIVE = 67
	RD_ERROR_USER_LINKED_INACTIVE = "User Linked is Inactive"

	RC_ERROR_INVALID_SIGNATURE = 81
	RD_ERROR_INVALID_SIGNATURE = "Signature mismatched"

	RC_ERROR_VOUCHER_NOTFOUND = 162
	RD_ERROR_VOUCHER_NOTFOUND = "Voucher Not Found"

	RC_ERROR_ACC_NOT_ELIGIBLE = 72
	RD_ERROR_ACC_NOT_ELIGIBLE = "Nomor belum eligible"

	RC_ERROR_FAILED_GETBALANCE = 73
	RD_ERROR_FAILED_GETBALANCE = "Failed to Get Balance"

	RC_ERROR_DUPLICATE_TRXID = 74
	RD_ERROR_DUPLICATE_TRXID = "Duplicate TrxID"

	RC_ERROR_FAILED_TRANS_POINT = 80
	RD_ERROR_FAILED_TRANS_POINT = "Gagal Transfer Point"

	RC_ERROR_FAILED_REDEEM_VOUCHER = 86
	RD_ERROR_FAILED_REDEEM_VOUCHER = "Gagal Redeem Voucher"

	RC_ERROR_FAILED_MAX_BUY_VOUCHER = 87
	RD_ERROR_FAILED_MAX_BUY_VOUCHER = "Anda mencapai batas maksimal pembelian voucher"

	RC_ERROR_FAILED_GET_HISTORY_VOUCHER = 89
	RD_ERROR_FAILED_GET_HISTORY_VOUCHER = "Gagal Get History Voucher Customer"

	RC_ERROR_FAILED_TRANSACTION = 106
	RD_ERROR_FAILED_TRANSACTION = "Transaksi gagal (other error)"

	RC_ERROR_FAILED_REVERSAL_VOUCHER = 70
	RD_ERROR_FAILED_REVERSAL_VOUCHER = "Gagal Reversal Voucher"

	RC_ERROR_PENDING_TRANSACTION = 109
	RD_ERROR_PENDING_TRANSACTION = "Transaksi Pending"

	RC_ERROR_NOT_ENOUGH_BALANCE = 27
	RD_ERROR_NOT_ENOUGH_BALANCE = "Point Tidak Mencukupi"

	RC_ERROR_FAILED_GET_POINT = 107
	RD_ERROR_FAILED_GET_POINT = "Gagal Dapat Point"

	RC_PARAMETER_INVALID = 201
	RD_PARAMETER_INVALID = "Parameter Invalid"

	RC_VOUCHER_NOTFOUND = 422
	RD_VOUCHER_NOTFOUND = "Voucher Tidak Ditemukan"

	RC_VOUCHER_NOT_AVAILABLE = "208"
	RD_VOUCHER_NOT_AVAILABLE = "Voucher not available"

	RC_FAILED_DECRYPT_VOUCHER = 202
	RD_FAILED_DECRYPT_VOUCHER = "Failed to Decrypt Voucher"

	RC_LIMIT_VOUCHER_USER_NOT_AVAILABLE = "209"
	RD_LIMIT_VOUCHER_USER_NOT_AVAILABLE = "Payment count limit exceeded"

	// Nomor belum eligible

	// Failed to GetBalance

	// Duplicate TrxID

	// Gagal Transfer Point
	// Gagal Redeem Voucher

	// Anda mencapai batas maksimal pembelian voucher

	// Gagal Get History Voucher Customer

	// Transaksi Gagal

	// Gagal Reversal Voucher

	// Transaksi Pending

)

//Cons Code
var (
	CONS_LINK_ACCOUNT       = "link"
	CONS_UNLINK_ACCOUNT     = "unlink"
	CONS_USER_STATUS_ACTIVE = "active"
)

const (
	CODE_TRANSTYPE_REDEMPTION = "TSP02"
	CODE_TRANSTYPE_INQUERY    = "TSP01"
)

// push notif general
const (
	CODE_EARNING_POINT         = "earning_voucher"
	CODE_REVERSAL_POINT        = "reversal_point"
	CODE_REDEEM_PLN            = "redeem_pln"
	CODE_EARNING_VOUCHER       = "earning_voucher"
	CODE_VOUCHER_EXPIRED       = "voucher_expired"
	CODE_GIFT_POINT_ACTIVATION = "gift_point_activation"
	TOPIC_PUSHNOTIF_GENERAL    = "ottopoint-notification-topics"
	TOPIC_PUSHNOTIF_REVERSAL   = "ottopoint-notification-reversal"
	CODE_REDEEM_VIDIO          = "redeem_vidio"
	CODE_REDEEM_PLN_SMS        = "sms_pln"
)

const (
	MsgSuccess = "Success"
)

const (
	CodeScheduler              = "SC001" // retry
	CodeSchedulerSepulsa       = "SC003"
	CodeSchedulerSpending      = "SC006"
	CodeSchedulerUV            = "SC007"
	CodeSchedulerOttoAG        = "SC008"
	CodeSchedulerVoucherAG     = "SC009"
	CodeSchedulerJempolKios    = "SC010"
	CodeSchedulerGudangVoucher = "SC011"
)

const (
	CODE_APPS_NOTIF     = 1
	CODE_SMS_NOTIF      = 2
	CODE_SMS_APPS_NOTIF = 3
)

const (
	CODE_VENDOR_OTTOAG     = "A1OAG001"
	CODE_VENDOR_UV         = "A1UVO002"
	CODE_VENDOR_SEPULSA    = "A1SPS003"
	CODE_VENDOR_AGREGATOR  = "A1OPG004"
	CODE_VENDOR_DUMY       = "A1ODM005"
	CODE_VENDOR_JempolKios = "A1JPK005"
	CODE_VENDOR_GV         = "A1GDV006"
)

const (
	CODE_SUCCESS   = "00"
	CODE_FAILED    = "01"
	CODE_STATUS_TO = "68"
)

const (
	RECEIVER_NAME_UV = "OTTOPOINT"
)

var (
	EXPDATE_VOUCHER = 15
)

const (
	// PARAMETER CONFIG UV
	CODE_CONFIG_UV_GROUP   = "UVCONFIG"
	CODE_CONFIG_UV_NAME    = ""
	CODE_CONFIG_UV_EMAIL   = "UV_EMAIL_ORDER"
	CODE_CONFIG_UV_EXPIRED = "UV_EXPIRED_VOUCHER"
	CODE_CONFIG_UV_PHONE   = "UV_PHONE_ORDER"

	// PARAMETER CONFIG COREPOINT
	CODE_CONFIG_COREPOINT_GROUP              = "CORE_POINT"
	CODE_CONFIG_COREPOINT_EXPIRED_POINT_DAYS = "EXPIRED_POINT_DAYS"

	// PARAMETER CONFIG AGREGATOR
	CODE_CONFIG_AGG_GROUP = "VOUCHER_AGG"
	CODE_CONFIG_AGG_EMAIL = "VOUCHER_AG_EMAIL_ORDER"
	CODE_CONFIG_AGG_PHONE = "VOUCHER_AG_PHONE_ORDER"
	CODE_CONFIG_AGG_EXPD  = "VOUCHER_AG_EXPIRED_VOUCHER"
	CODE_CONFIG_AGG_NAME  = ""
)

const (
	CODE_NOMOR_TELP   = "001"
	CODE_ID_Pelanggan = "002"
	CODE_ID_Server    = "003"
	CODE_NOMOR_KARTU  = "004"
)

const (
	VoucherTypePPOB        = "ppob"
	VoucherTypeVoucherCode = "voucher code"
)

const (
	CreatedbySystem = "System"
)

const (
	TypeCash  = "cash"
	TypePoint = "point"
)

var Iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
