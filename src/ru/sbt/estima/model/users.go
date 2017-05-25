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
	"github.com/gorilla/context"
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

	config := conf.LoadConfig()
	l, err := ldap.Dial(config.Ldap.Protocol, fmt.Sprintf("%s:%d", config.Ldap.Host, config.Ldap.Port))
	CheckErr (err)
	defer l.Close()

	// cn := fmt.Sprintf("cn=%s,%s", username, config.Ldap.Dn)
	// Authenticate using given username and password
	err = l.Bind(username, password)
	CheckErr (err)

	// Search for the uswr details
	searchRequest := ldap.NewSearchRequest (
		config.Ldap.Dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		10,
		0,
		false,
		fmt.Sprintf("(&(objectClass=person)(cn=%s))", username),
		[]string{"dn", "uid", "givenName", "displayName", "sn", "cn", "name", "mail", "sAMAccountName"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	CheckErr (err)

	entry := sr.Entries[0]
	retUser = NewUser(
		username,
		entry.GetAttributeValue("mail"),
		password,
		entry.GetAttributeValue("sn"),
		entry.GetAttributeValue("uid"),
		nil,
		"")

	return retUser, retErr
}

func GetUserFromRequest (w http.ResponseWriter, r *http.Request) (*EstimaUser) {
	user := context.Get(r, "user")
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