package constants

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
	KeyResponsePending = "pending"
	KeyResponseDefault = "default"
	KeyResponseFailed  = "failed"
	KeyResponseSucceed = "succeed"
)

const (
	CategoryPulsa        = "pulsa"
	CategoryPaketData    = "paket_data"
	CategoryFreeFire     = "free fire"
	CategoryMobileLegend = "mobile legends"
	CategoryPLN          = "pln"
	CategoryGame         = "game"
)

const (
	Success = "00 (Success)"
	Pending = "09 (Pending)"
	Failed  = "01 (Failed)"
)

const (
	UV     = "Ultra Voucher"
	OttoAG = "OttoAG"
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
