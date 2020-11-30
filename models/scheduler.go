package models

type SchedulerCheckStatusResp struct {
	Data  []SchedulerCheckStatusData
	Total int `json:"total"`
}

type SchedulerCheckStatusData struct {
	Supplier string `json:"supplier"`
	Success  int    `json:"success"`
	Failed   int    `json:"failed"`
	Total    int    `json:"total"`
}
