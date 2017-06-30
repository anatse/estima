package conf

import (
	"encoding/json"
	"os"
	"log"
	"github.com/bradfitz/gomemcache/memcache"
	"fmt"
)

type Ldap struct {
	Protocol string
	Host string
	Port uint32
	Dn string
	BindDN string `json:"bindDN"`
	BindPwd string `json:"bindPwd"`
	PoolSize uint `json:"poolSize"`
}

type Database struct {
	Url	string
	User	string
	Password string
	Log bool
	Name string
	PoolSize uint `json:"poolSize"`
}

type Auth struct {
	CookieName string
	MaxAge int
}

type Memcached struct {
	Machines []struct {
		Host string
		Port int
	}

	cache *memcache.Client
}

type Profile struct {
	Name string
	Secret string
	Port int
	Ldap	Ldap
	Database Database
	Auth Auth
	Memcached Memcached
}

type Configuration struct {
	Active   string
	Profiles []Profile
}

func (cf Configuration) ActiveProfile ()(Profile) {
	var profile Profile
	for i:=0; i<len(cf.Profiles); i++ {
		if cf.Profiles[i].Name == cf.Active {
			profile = cf.Profiles[i]
		}
	}

	return profile
}

func (mem *Memcached) create () {
	if mem.Machines != nil && len(mem.Machines) > 0 {
		var memHosts []string = make ([]string, len(mem.Machines))
		for i, host := range mem.Machines {
			memHosts[i] = fmt.Sprintf("%s:%d", host.Host, host.Port)
		}

		mem.cache = memcache.New (memHosts...)
	}
}


func (profile Profile) Cache () *memcache.Client {
	GetLog().Println(profile.Memcached)

	if profile.Memcached.cache == nil {
		profile.Memcached.create()
	}

	return profile.Memcached.cache
}

var config Configuration
func LoadConfig() (Profile) {
	if config.Active != "" {
		return config.ActiveProfile()
	}

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.json"
	}

	GetLog().Printf("Config loaded from %s", cfgPath)

	file, _ := os.Open(cfgPath)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	GetLog().Println(config.ActiveProfile())
	return config.ActiveProfile()
}

var fileLog *log.Logger
var f *os.File
func GetLog () *log.Logger {
	if fileLog  == nil {
		f, err := os.OpenFile("estima.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			GetLog().Fatalf("error opening file: %v", err)
		}

		fileLog = log.New(f, "", log.LstdFlags)
	}

	f.Sync()
	return fileLog
}

func GetLogFile () *os.File {
	return f
}