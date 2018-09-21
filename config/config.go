package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	AgentContractAddressKey    = "AGENT_CONTRACT_ADDRESS"
	AutoSSLDomainKey           = "AUTO_SSL_DOMAIN"
	AutoSSLCacheDirKey         = "AUTO_SSL_CACHE_DIR"
	BlockchainEnabledKey       = "BLOCKCHAIN_ENABLED"
	ConfigPathKey              = "CONFIG_PATH"
	DaemonListeningPortKey     = "DAEMON_LISTENING_PORT"
	DaemonTypeKey              = "DAEMON_TYPE"
	DbPathKey                  = "DB_PATH"
	EthereumJsonRpcEndpointKey = "ETHEREUM_JSON_RPC_ENDPOINT"
	ExecutablePathKey          = "EXECUTABLE_PATH"
	HdwalletIndexKey           = "HDWALLET_INDEX"
	HdwalletMnemonicKey        = "HDWALLET_MNEMONIC"
	LogLevelKey                = "LOG.LEVEL"
	LogFormatterKey            = "LOG.FORMATTER"
	PassthroughEnabledKey      = "PASSTHROUGH_ENABLED"
	PassthroughEndpointKey     = "PASSTHROUGH_ENDPOINT"
	PollSleepKey               = "POLL_SLEEP"
	PrivateKeyKey              = "PRIVATE_KEY"
	ServiceTypeKey             = "SERVICE_TYPE"
	SSLCertPathKey             = "SSL_CERT"
	SSLKeyPathKey              = "SSL_KEY"
	WireEncodingKey            = "WIRE_ENCODING"

	defaultConfigJson string = `
{
	"auto_ssl_cache_dir": ".certs",
	"auto_ssl_domain": "",
	"blockchain_enabled": true,
	"daemon_listening_port": 5000,
	"daemon_type": "grpc",
	"db_path": "snetd.db",
	"ethereum_json_rpc_endpoint": "http://127.0.0.1:8545",
	"hdwallet_index": 0,
	"hdwallet_mnemonic": "",
	"passthrough_enabled": false,
	"poll_sleep": "5s",
	"service_type": "grpc",
	"ssl_cert": "",
	"ssl_key": "",
	"wire_encoding": "proto",
	"log":  {
		"level": "info",
		"formatter": {
			"type": "json",
			"timezone": "UTC"
		},
		"output": {
			"type": "file",
			"file_pattern": "./snet-daemon.%Y%m%d.log",
			"current_link": "./snet-daemon.log",
			"clock_timezone": "UTC",
			"rotation_time_in_sec": 86400,
			"max_age_in_sec": 604800,
			"rotation_count": 0
		}
	}
}
`
)

var vip *viper.Viper

func init() {
	vip = viper.New()
	vip.SetEnvPrefix("SNET")
	vip.AutomaticEnv()

	setDefaultsFromJsonString(defaultConfigJson)

	vip.AddConfigPath(".")
}

func setDefaultsFromJsonString(data string) {
	var err error
	var temporaryConfig *viper.Viper = viper.New()

	temporaryConfig.SetConfigType("json")

	err = temporaryConfig.ReadConfig(strings.NewReader(data))
	if err != nil {
		panic(fmt.Sprintf("Cannot load default config: %v", err))
	}

	for key, value := range temporaryConfig.AllSettings() {
		vip.SetDefault(key, value)
	}
}

func Vip() *viper.Viper {
	return vip
}

func Validate() error {
	switch dType := vip.GetString(DaemonTypeKey); dType {
	case "grpc":
		switch sType := vip.GetString(ServiceTypeKey); sType {
		case "grpc":
		case "jsonrpc":
		case "process":
			if vip.GetString(ExecutablePathKey) == "" {
				return errors.New("EXECUTABLE required with SERVICE_TYPE 'process'")
			}
		default:
			return fmt.Errorf("unrecognized SERVICE_TYPE '%+v'", sType)
		}

		switch enc := vip.GetString(WireEncodingKey); enc {
		case "proto":
		case "json":
		default:
			return fmt.Errorf("unrecognized WIRE_ENCODING '%+v'", enc)
		}
	case "http":
	default:
		return fmt.Errorf("unrecognized DAEMON_TYPE '%+v'", dType)
	}

	if vip.GetBool(BlockchainEnabledKey) {
		if vip.GetString(PrivateKeyKey) == "" && vip.GetString(HdwalletMnemonicKey) == "" {
			return errors.New("either PRIVATE_KEY or HDWALLET_MNEMONIC are required")
		}
	}

	certPath, keyPath := vip.GetString(SSLCertPathKey), vip.GetString(SSLKeyPathKey)
	if (certPath != "" && keyPath == "") || (certPath == "" && keyPath != "") {
		return errors.New("SSL requires both key and certificate when enabled")
	}

	return nil
}

func LoadConfig(configFile string) error {
	vip.SetConfigFile(configFile)
	return vip.ReadInConfig()
}

func WriteConfig(configFile string) error {
	vip.SetConfigFile(configFile)
	return vip.WriteConfig()
}

func GetString(key string) string {
	return vip.GetString(key)
}

func GetInt(key string) int {
	return vip.GetInt(key)
}

func GetDuration(key string) time.Duration {
	return vip.GetDuration(key)
}

func GetBool(key string) bool {
	return vip.GetBool(key)
}
