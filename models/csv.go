package models

type CreateCSV struct {
	Tittle  string  `json:"tittle"`
	Messgae string  `json:"msg"`
	Data    DataCSV `json:"data"`
}

type DataCSV struct {
	Data1 string `json:"data1"`
	Data2 string `json:"data2"`
}
