package g

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"sync"
)

var (
	//config *GlobalConfig
	lock = new(sync.RWMutex)
)

//
//func Content() *GlobalConfig {
//	lock.RLock()
//	defer lock.RUnlock()
//	return config
//}
//
//type LogConfig struct {
//	Path  string `json:"path"`
//	Level string `json:"level"`
//}
//
//type SentryConfig struct {
//	Dsn string `json:"dsn"`
//}
//
//type DBConfig struct {
//	DSN string `json:"dsn"`
//	//Addr         string `json:"addr"`
//	//User         string `json:"user"`
//	//Password     string `json:"password"`
//	//UserName         string `json:"name"`
//	//Charset      string `json:"charset"`
//	//ServerId     int    `json:"server_id"`
//	//MaxIdleConns int    `json:"max_idle_conns"`
//	//MaxOpenConns int    `json:"max_open_conns"`
//}
//
//type DBPoolConfig struct {
//	Uic       DBConfig `json:"uic"`
//	DashBoard DBConfig `json:"dashboard"`
//}
//
//type ServConfig struct {
//	Addr          string `json:"addr"`
//	AccessControl bool   `json:"access_control"`
//}
//
//type GlobalConfig struct {
//	Log    *LogConfig    `json:"log"`
//	Sentry *SentryConfig `json:"sentry"`
//	//DB     *DBConfig     `json:"db"`
//	DB   *DBPoolConfig `json:"db"`
//	Serv *ServConfig   `json:"serv"`
//}

func ParseConfig(cfg string) {

	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Infoln(err)
		flag.Usage()
		return
	}

	viper.SetConfigType("json")
	viper.SetConfigFile(cfg)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infoln("config changed：", e.Name)
	})

	if err := viper.ReadInConfig(); err != nil {
		log.Errorln(err)
		flag.Usage()
		return
	}

	lock.Lock()
	defer lock.Unlock()

	log.Infoln("read config file:", cfg, "successfully")
}
