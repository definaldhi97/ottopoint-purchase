package models

type HealthcheckResponse struct {
	Redis    RedisHealthcheckResponse      `json:"redis"`
	Database DBHealthcheckResponse         `json:"database"`
	Services []ServicesHealthcheckResponse `json:"services"`
}

type RedisHealthcheckResponse struct {
	Type    string `json:"type"`
	Address string `json:"address"`
	Port    string `json:"port"`
}

type ServicesHealthcheckResponse struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type DBHealthcheckResponse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    string `json:"port"`
}
