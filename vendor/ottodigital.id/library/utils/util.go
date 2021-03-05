package utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/envconfig"

	//"strconv"

	"time"
)

var (
	readtimeout       string
	readheadertimeout string
	writetimeout      string
	idletimeout       string
	maxheaderbyte     string

	serverfileca         string
	serverfileprivatekey string
	serverfilepubkey     string
	servertls12client    string

	consuladdres             string
	consulregid              string
	consulregname            string
	consulregserver          string
	consulregport            string
	consulhealtCheckHttp     string
	consulhealtcheckInterval string
	consulhealtcheckTimeout  string
	consulonof               string

	kafkaBrokerUrl            string
	kafkaClient               string
	kafkaProducerTimeout      string
	kafkaProducerDialTimeout  string
	kafkaProducerReadTimeout  string
	kafkaProducerWriteTimeout string
	kafkaProducerMaxmsgbyte   string
)

//$ export  httpserver_tls12=OFF

func init() {

	/*

		consul registration
		export CONSULREG_HOSTADDRES="13.250.21.165:8500"
		export CONSULREG_IDSERVER="OTTOCASHBACK-01"
		export CONSULREG_NAME="OTTOCASHBACK-01"
		export CONSULREG_SERVER="OTTOCASHBACK-01"
		export CONSULREG_PORT="8997"

		#consul healthcheck
		export CONSULREG_HCURL="http://127.0.0.1:8997/healtcheck"
		export CONSULREG_HCINTERVAL="10s"
		export CONSULREG_HCTIMEOUT="3s"


		#	kafka.BROKERSURL", "13.250.26.210:9092,13.250.26.210:9092,13.250.26.210:9092")
		#	KAFKAPRO_CLIENTID", "CLIENTID")
		#	KAFKAPRO_TO", "10s")
		#	KAFKAPRO_DIALTO", "10s")
		#	KAFKAPRO_READTO", "10s")
		#	KAFKAPRO_WRITETO", "10s")
		#	KAFKAPRO_MAXMSGBYTE", "50000000")

		#SERVER LOCAL
		export SERVER_READTIMEOUT="60s"
		export SERVER_RHTIMEOUT="10s"
		export SERVER_WRITETIMEOUT="10s"
		export SERVER_IDLETIMEOUT="0"
		export SERVER_MAXBYTES="0"
		export SERVER_FILECA="keys/cert.pem"
		export SERVER_PRIVATEKEY="keys/key.pem"
		export SERVER_PUBLICKEY="keys/cert.pem"
		export SERVER_TLS12STATUS="OFF"
	*/

}

const (
	CONSULREGPREFIX     = "consulreg"
	SERVERPREFIX        = "server"
	KAFKAPRODUCERPREFIX = "kafkapro"
	KAFKACONSUMERPREFIX = "kafkaconsumer"
	OTTOHTTPCLIENT      = "ottohttpclient"
)

type EnvHttpClient struct {
	HttpClientTimeout string `envconfig:"timeout"`
	HttpClientRetry   int    `envconfig:"retry"`
	HttpClientTracing bool   `envconfig:"enabletracing"`
	HttpClientDebug   bool   `envconfig:"debug"`
}

type EnvKafkaProcedureConfig struct {
	KafkaBrokerUrl            string        `envconfig:"brokerurl"`
	KafkaClient               string        `envconfig:"clientid"`
	KafkaProducerTimeout      time.Duration `envconfig:"to"`
	KafkaProducerDialTimeout  time.Duration `envconfig:"dialto"`
	KafkaProducerReadTimeout  time.Duration `envconfig:"readto"`
	KafkaProducerWriteTimeout time.Duration `envconfig:"writeto"`
	KafkaProducerMaxmsgbyte   int           `envconfig:"maxmsgbyte"`
}

type EnvKafkaConsumerConfig struct {
	KafkaZookeeper string `envconfig:"kafkazookeeper"`
	KafkaBroker    string `envconfig:"kafkabroker"`
}
type ConsulConfig struct {
	HostAddres string `envconfig:"hostaddres"`
	IdServer   string
	Name       string
	Server     string
	Port       int
	Hcurl      string
	HcInterval string
	HcTimeout  string
	Status     bool
}

type ServerConfig struct {

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body.
	ReadHeaderTimeout time.Duration `envconfig:"rhtimeout"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.
	IdleTimeout time.Duration

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int `envconfig:"maxbytes"`

	Serverfileca         string `envconfig:"FILECA"`
	Serverfileprivatekey string `envconfig:"PRIVATEKEY"`
	Serverfilepubkey     string `envconfig:"PUBLICKEY"`
	Servertls12client    string `envconfig:"TLS12STATUS"`
}

func GetServerConfig() *ServerConfig {
	var serverCfg ServerConfig
	err := envconfig.Process(SERVERPREFIX, &serverCfg)
	fmt.Println("Error Config Consul : ", err)
	return &serverCfg
}

func GetConsulConfig() ConsulConfig {
	var consulcfg ConsulConfig
	err := envconfig.Process(CONSULREGPREFIX, &consulcfg)
	fmt.Println("Error Config Consul : ", err)
	return consulcfg
}

func GetEnvConfigConsumerKafka() EnvKafkaConsumerConfig {
	var cfg EnvKafkaConsumerConfig
	err := envconfig.Process(KAFKACONSUMERPREFIX, &cfg)
	fmt.Println("Error Config Kafka Consumer : ", err)
	return cfg
}

func GetEnvConfigProcedurKafka() EnvKafkaProcedureConfig {
	var cfg EnvKafkaProcedureConfig
	err := envconfig.Process(KAFKAPRODUCERPREFIX, &cfg)
	fmt.Println("Error Config Consul : ", err)
	return cfg

}

// Simple helper function to read an environment or return a default value
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func GetServerTlsConfig() *tls.Config {
	if servertls12client == "ON" {

		caCert, err := ioutil.ReadFile(serverfileca)
		if err != nil {
			fmt.Println("Error : ", err)

		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		tlsConfig.BuildNameToCertificate()
		return tlsConfig
	}
	return &tls.Config{InsecureSkipVerify: true}

}

func StrucToMap(in interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}

func GetEnvConfigOttoHttpReq() EnvHttpClient {
	var cfg EnvHttpClient
	err := envconfig.Process(OTTOHTTPCLIENT, &cfg)
	fmt.Println("Error Config HttpClient : ", err)
	return cfg
}
