package conf

import (
	"encoding/json"
	"os"
	"fmt"
	"log"
	"github.com/bradfitz/gomemcache/memcache"
)

type Ldap struct {
	Protocol string
	Host string
	Port uint32
	Dn string
	BindDN string `json:"bindDN"`
	BindPwd string `json:"bindPwd"`
}

type Database struct {
	Url	string
	User	string
	Password string
	Log bool
	Name string
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

		log.Printf("hosts: %v", memHosts)
		mem.cache = memcache.New (memHosts...)
		log.Printf("client: %v", mem.cache)
	}
}


func (profile Profile) Cache () *memcache.Client {
	log.Println(profile.Memcached)

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

	log.Printf("Config loaded from %s", cfgPath)

	file, _ := os.Open(cfgPath)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	fmt.Println(config.ActiveProfile())
	return config.ActiveProfile()
}
