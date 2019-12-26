package ottopointmodels

type GetBalance struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	CustomerID string `json:"customerId"`
	Points     string `json:"points"`
}
