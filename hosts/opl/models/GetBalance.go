package models

type BalanceResponse struct {
	FirstName                              string      `json:"firstName"`
	LastName                               string      `json:"lastName"`
	CustomerID                             string      `json:"customerId"`
	Points                                 float64     `json:"points"`
	P2PPoints                              float64     `json:"p2pPoints"`
	TotalEarnedPoints                      float64     `json:"totalEarnedPoints"`
	UsedPoints                             float64     `json:"usedPoints"`
	ExpiredPoints                          float64     `json:"expiredPoints"`
	LockedPoints                           float64     `json:"lockedPoints"`
	Level                                  string      `json:"level"`
	LevelName                              string      `json:"levelName"`
	LevelConditionValue                    float64     `json:"levelConditionValue"`
	NextLevel                              string      `json:"nextLevel"`
	NextLevelName                          string      `json:"nextLevelName"`
	NextLevelConditionValue                int         `json:"nextLevelConditionValue"`
	TransactionsAmountWithoutDeliveryCosts int         `json:"transactionsAmountWithoutDeliveryCosts"`
	TransactionsAmountToNextLevel          int         `json:"transactionsAmountToNextLevel"`
	AverageTransactionsAmount              string      `json:"averageTransactionsAmount"`
	TransactionsCount                      int         `json:"transactionsCount"`
	TransactionsAmount                     int         `json:"transactionsAmount"`
	Currency                               string      `json:"currency"`
	PointsExpiringNextMonth                float64     `json:"pointsExpiringNextMonth"`
	PointsExpiringBreakdown                interface{} `json:"pointsExpiringBreakdown"`
}
