package models

import "time"

// RedisClusterEnv ..
type RedisClusterEnv struct {
	AddressCluster1 string `json:"address_cluster_1"`
	AddressCluster2 string `json:"address_cluster_2"`
	AddressCluster3 string `json:"address_cluster_3"`
}

// RedisNonClusterEnv ..
type RedisNonClusterEnv struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

// DatabaseEnv ..
type DatabaseEnv struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
}

// ServiceEnv ..
type ServiceEnv struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	HealthCheckKey string `json:"health_check_key"`
}

// HealthCheckResponse ..
type HealthCheckResponse struct {
	Redis    []RedisHealthCheck    `json:"redis,omitempty"`
	Database []DatabaseHealthCheck `json:"database,omitempty"`
	Kafka    []KafkaHealthCheck    `json:"kafka,omitempty"`
	Service  []ServiceHealthCheck  `json:"service,omitempty"`
}

// RedisHealthCheck ..
type RedisHealthCheck struct {
	Address     string `json:"address"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// DatabaseHealthCheck ..
type DatabaseHealthCheck struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// KafkaHealthCheck ..
type KafkaHealthCheck struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// ServiceHealthCheck ..
type ServiceHealthCheck struct {
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}
