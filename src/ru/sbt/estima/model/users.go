package model

import (
	ara "github.com/diegogub/aranGO"
	"fmt"
	"crypto/sha1"
	"encoding/base64"
	"gopkg.in/ldap.v2"
	"ru/sbt/estima/conf"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"log"
	"sync"
)

type EstimaUser struct {
	ara.Document `json:-`

	Name     string `json:"name,omitempty", unique:"users"`
	Email    string `json:"email,omitempty, unique:"users""`
	Password string  `json:"-"`
	DisplayName  string `json:"displayName,omitempty"`
	Uid string `json:"uid,omitempty", unique:"users"`
	Roles	[]string `json:"roles,omitempty"`
}

func NewUser (name string, email string, password string, displayName string, uid string, roles []string, key string) *EstimaUser {
	var pwd string

	if password != "" {
		pwd = base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(password)))
	}

	var user EstimaUser
	user.Name = name
	user.Email = email
	user.DisplayName = displayName
	user.Password = pwd
	user.Uid = uid
	user.Roles = roles
	user.Key = key

	return &user
}

type omit *struct{}
func (user EstimaUser) Entity() interface{} {
	return struct{
		*EstimaUser

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&user,
		nil,
		nil,
		nil,
		nil,
	}
}

func (user EstimaUser) AraDoc() (ara.Document) {
	return user.Document
}

func (user EstimaUser)GetKey() string {
	return user.Key
}

func (user EstimaUser) GetCollection() string {
	return "users"
}

func (user EstimaUser) GetError()(string, bool){
	// default error bool and messages. Could be any kind of error
	return user.Message, user.Error
}

func (user EstimaUser) CopyChanged (entity Entity) Entity {
	newUser := entity.(EstimaUser)
	if newUser.Name != "" {user.Name = newUser.Name}
	if newUser.Email != "" {user.Email = newUser.Email}
	if newUser.DisplayName != "" {user.DisplayName = newUser.DisplayName}
	if newUser.Password != "" {user.Password = newUser.Password}
	if newUser.Uid != "" {user.Uid = newUser.Uid}
	if newUser.Roles != nil {user.Roles = newUser.Roles}

	return user
}

func printLdapAttrs (sr *ldap.SearchResult) {
	for i:=0;i<len(sr.Entries);i++ {
		entry := sr.Entries[i]
		for a:=0;a<len(entry.Attributes);a++ {
			attr := entry.Attributes[a]
			fmt.Println(attr.Name + " = " + attr.Values[0])
		}
		fmt.Println("----------------------------------------")
	}
}

// Global LDAP connections pool
var ldapConnectionsPool []*ldap.Conn
type chanUser struct  {
	user string
	password string
}

var ldapJobs chan chanUser
type LockedUser struct {
	sync.WaitGroup
	user *EstimaUser
	err error
}

type JobMap struct {
	sync.RWMutex
	jobMap map[string]*LockedUser
}

var jobMap JobMap = JobMap {
	jobMap: make(map[string]*LockedUser),
}

func writeMap (key string, user *LockedUser) {
	jobMap.Lock()
	defer jobMap.Unlock()

	jobMap.jobMap[key] = user
}

func readMap (key string) *LockedUser {
	jobMap.Lock()
	defer jobMap.Unlock()

	return jobMap.jobMap[key]
}

func SendToAuth (username string, password string) {
	ldapJobs <- chanUser{username, password}
}

// Function create and filled LDAP connection pool
func InitLdapPool (poolSize int) {
	if ldapConnectionsPool == nil {
		config := conf.LoadConfig()
		if config.Ldap.Protocol == "fake" {
			return
		}

		ldapConnectionsPool = make ([]*ldap.Conn, poolSize)
		for index := range ldapConnectionsPool {
			// Connect to LDAP catalog
			conn, err := ldap.Dial(config.Ldap.Protocol, fmt.Sprintf("%s:%d", config.Ldap.Host, config.Ldap.Port))
			CheckErr(err)

			// Bind to main user
			CheckErr(conn.Bind(config.Ldap.BindDN, config.Ldap.BindPwd))

			// Store connections to LDAP connections pool
			ldapConnectionsPool[index] = conn
		}

		// Make channels
		ldapJobs = make(chan chanUser)

		// Start worker routines
		for i:=0; i< poolSize; i++ {
			go worker(i, config.Ldap, ldapConnectionsPool[i], ldapJobs)
		}
	}
}

// Function close all LDAP Connections and clear connection pool
func FinishLdapPool () {
	if ldapConnectionsPool != nil {
		for _, con := range ldapConnectionsPool {
			defer con.Close()
		}

		ldapConnectionsPool = nil
		close (ldapJobs)
	}
}

func worker(id int, config conf.Ldap, conn *ldap.Conn, jobs chan chanUser) {
	for job := range jobs {
		fmt.Printf("+++ Start processing worker(%d) for job (%v)\n", id, job)

		defer CheckErr(conn.Bind(config.BindDN, config.BindPwd))

		user, err := auth (conn, config, job.user, job.password)
		luser := readMap(job.user)
		luser.user = user
		luser.err = err
		luser.Done()
	}
}

func auth (conn *ldap.Conn, config conf.Ldap, username string, password string) (retUser *EstimaUser, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in FindUser:", r)
			retErr = fmt.Errorf ("%v", r)
			fmt.Printf("Error is: %v\n", retErr)
		}
	}()

	searchRequest := ldap.NewSearchRequest (
		config.Dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		10,
		0,
		false,
		fmt.Sprintf("(&(objectClass=person)(sAMAccountName=%s))", username),
		[]string{"dn", "uid", "givenName", "displayName", "sn", "cn", "name", "mail", "sAMAccountName"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	CheckErr (err)

	if len(sr.Entries) == 0 {
		log.Panicf("Authentication for user %v failed", username)
	}

	// Get user info
	entry := sr.Entries[0]

	// printLdapAttrs (sr)

	// Trying to authenticate using found user
	CheckErr(conn.Bind(entry.GetAttributeValue("cn"), password))

	retUser = NewUser(
		username,
		entry.GetAttributeValue("mail"),
		password,
		entry.GetAttributeValue("displayName"),
		entry.GetAttributeValue("uid"),
		nil,
		"")

	return retUser, retErr
}

/**
 Function used to check prc name and password through the LDAP.
 Additional documentation about LDAP library located here https://godoc.org/gopkg.in/ldap.v2
 This configuration only for ActiveDirectory installation
 */
func FindUser (username string, password string) (retUser *EstimaUser, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in FindUser:", r)
			retErr = fmt.Errorf("%v", r)
		}
	}()

	lu := LockedUser{}
	lu.Add(1)

	writeMap(username, &lu)
	SendToAuth(username, password)

	lu.Wait()

	return lu.user, lu.err
}

func GetUserFromRequest (w http.ResponseWriter, r *http.Request) (*EstimaUser) {
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	var roles []string
	rolesClaim := claims["roles"]
	if rolesClaim != nil {
		rolesInterface := rolesClaim.([]interface{})
		roles = make([]string, len(rolesInterface))
		for i := range rolesInterface {
			roles[i] = rolesInterface[i].(string)
		}
	}

	return NewUser(
		claims["name"].(string),
		claims["mail"].(string),
		"",
		claims["displayName"].(string),
		claims["uid"].(string),
		roles,
		claims["key"].(string),
	)
}