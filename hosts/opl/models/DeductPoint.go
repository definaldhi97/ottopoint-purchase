package models

type PointResponse struct {
	// Utama
	PointsTransferId string `json:"pointsTransferId"`
	// Error point 0 / kosong, CustID not exist / Kosong
	Form struct {
		Children struct {
			Customer struct {
				Errors []string `json:"errors"` //  CustID not exist / Kosong
			} `json:"customer"`
			Points struct {
				Errors []string `json:"errors"` // Error point 0 / kosong
			} `json:"points"`
			ValidityDuration struct {
			} `json:"validityDuration"`
			Comment struct {
			} `json:"comment"`
		} `json:"children"`
	} `json:"form"`
	// Error Invalid JWT Token (Token Expiresd), Invalid credentials (Token kosong)
	Code    int    `json:"code"`
	Message string `json:"message"`
	// Error Forbidden (salah token)
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
