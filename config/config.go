package config

import (
	goflag "flag"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// MFLIntegratorConf is the global conf
var MFLIntegratorConf Config

func init() {
	// ./mfld --config-file /etc/mflintegrator/mfld.yml

	configFile := flag.String("config-file", "",
		"The path to the configuration file of the application")

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	viper.SetConfigName("mfld")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/mflintegrator/")

	if len(*configFile) > 0 {
		viper.SetConfigFile(*configFile)
		log.Printf("Config File %v", *configFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// log.Fatalf("Configuration File: %v Not Found", *configFile, err)
			panic(fmt.Errorf("Fatal Error %w \n", err))

		} else {
			log.Fatalf("Error Reading Config: %v", err)

		}
	}

	err := viper.Unmarshal(&MFLIntegratorConf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Printf(">>>>> %v", MFLIntegratorConf.Database.DBUsername)
	log.Printf(">>>>> %v", MFLIntegratorConf.Database.DBPassword)
	log.Printf(">>>>>++ %v", MFLIntegratorConf.Server.MaxRetries)
	log.Printf(">>>>>++ %v", MFLIntegratorConf.Server)

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatalf("unable to reread configuration into global conf: %v", err)
		}
		_ = viper.Unmarshal(&MFLIntegratorConf)
		log.Printf(">>>>>+++++++++ %v", viper.GetInt("server.request_process_interval"))
	})
	viper.WatchConfig()
}

// Config is the top level cofiguration object
type Config struct {
	Database struct {
		URI        string `mapstructure:"uri" env:"MFLINTEGRATOR_DB" env-default:"postgres://postgres:postgres@localhost/mfldb?sslmode=disable"`
		DBHost     string `mapstructure:"db_host" env:"MFLINTEGRATOR_DB_HOST" env-default:"localhost"`
		DBPort     string `mapstructure:"db_port" env:"MFLINTEGRATOR_DB_PORT" env-default:"5432"`
		DBUsername string `mapstructure:"db_username" env:"MFLINTEGRATOR_USER" env-description:"API user name"`
		DBPassword string `mapstructure:"db_password" env:"MFLINTEGRATOR_PASSWORD" env-description:"API user password"`
	} `yaml:"database"`

	Server struct {
		Host                    string `mapstructure:"host" env:"MFLINTEGRATOR_HOST" env-default:"localhost"`
		Port                    string `mapstructure:"http_port" env:"MFLINTEGRATOR_SERVER_PORT" env-description:"Server port" env-default:"9090"`
		ProxyPort               string `mapstructure:"proxy_port" env:"MFLINTEGRATOR_PROXY_PORT" env-description:"Server port" env-default:"9191"`
		MaxRetries              int    `mapstructure:"max_retries" env:"MFLINTEGRATOR_MAX_RETRIES" env-default:"3"`
		StartOfSubmissionPeriod string `mapstructure:"start_submission_period" env:"START_SUBMISSION_PERIOD" env-default:"18"`
		EndOfSubmissionPeriod   string `mapstructure:"end_submission_period" env:"END_SUBMISSION_PERIOD" env-default:"24"`
		MaxConcurrent           int    `mapstructure:"max_concurrent" env:"MFLINTEGRATOR_MAX_CONCURRENT" env-default:"5"`
		RequestProcessInterval  int    `mapstructure:"request_process_interval" env:"REQUEST_PROCESS_INTERVAL" env-default:"4"`
		LogDirectory            string `mapstructure:"logdir" env:"MFLINTEGRATOR_LOGDIR" env-default:"/var/log/dispatcher2"`
		UseSSL                  string `mapstructure:"use_ssl" env:"MFLINTEGRATOR_USE_SSL" env-default:""`
		SSLClientCertKeyFile    string `mapstructure:"ssl_client_certkey_file" env:"SSL_CLIENT_CERTKEY_FILE" env-default:""`
		SSLServerCertKeyFile    string `mapstructure:"ssl_server_certkey_file" env:"SSL_SERVER_CERTKEY_FILE" env-default:""`
		SSLTrustedCAFile        string `mapstructure:"ssl_trusted_cafile" env:"SSL_TRUSTED_CA_FILE" env-default:""`
	} `yaml:"server"`

	API struct {
		// BaseMFLURL string ``
		Email     string `yaml:"email" env:"MFLINTEGRATOR_EMAIL" env-description:"API user email address"`
		AuthToken string `yaml:"authtoken" env:"RAPIDPRO_AUTH_TOKEN" env-description:"API JWT authorization token"`
		SmsURL    string `yaml:"smsurl" env:"MFLINTEGRATOR_SMS_URL" env-description:"API SMS endpoint"`
	} `yaml:"api"`
}
