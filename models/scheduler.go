package models

type SchedulerCheckStatusResp struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

type SchedulerCheckStatusData struct {
	Supplier string `json:"supplier"`
	Success  int    `json:"success"`
	Failed   int    `json:"failed"`
	Total    int    `json:"total"`
}

type SchedulerCheckStatusDataSupplier struct {
	Sepulsa       int `json:"sepulsa"`
	OttoAG        int `json:"ottoAG"`
	UltraVoucher  int `json:"ultraVoucher"`
	JempolKios    int `json:"jempolKios"`
	GudangVoucher int `json:"gudangVoucher"`
	VouicherAG    int `json:"voucherAG"`
}
