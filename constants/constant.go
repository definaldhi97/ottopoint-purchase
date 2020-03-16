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
	CategoryFreeFire     = "free_fire"
	CategoryMobileLegend = "mobile_legend"
	CategoryToken        = "token"
)

const (
	Success = "(00 Success)"
	Pending = "(09 Pending)"
	Failed  = "(01 Failed)"
)

const (
	UV     = "Ultra Voucher"
	OttoAG = "OttoAG"
)
