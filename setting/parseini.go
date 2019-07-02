package setting

import (
	"flag"
	"log"
	"time"
	"github.com/go-ini/ini"
)

type App struct {
	Quizes 			string
}

var AppSetting = &App{}

type Server struct {
	RunMode      	string
	Port         	string
	ReadTimeout  	time.Duration
	WriteTimeout 	time.Duration
}

var ServerSetting = &Server{}


type Database struct {
	SqliteFile 		string
	TablePrefix 	string
}

var DatabaseSetting = &Database{}


var cfg *ini.File

func Setup() {
	var err error

	file := confFile()

	cfg, err = ini.Load(file)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second

}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}

func confFile() string {

	confPtr := flag.String("conf", "", "set private app.ini file")
	flag.Parse()

	if *confPtr != "" {
		return *confPtr
	} else {
		return "conf/app.ini"
	}
}
