package model

import (
	ara "github.com/diegogub/aranGO"
	"fmt"
	"log"
	"crypto/sha1"
	"encoding/json"
	"encoding/base64"
	"gopkg.in/ldap.v2"
	"ru/sbt/estima/conf"
)

type EstimaUser struct {
	ara.Document
	Name     string
	Email    string
	Password string
	DisplayName  string
	Uid string
	Roles	[]string
}

func NewUser (name string, email string, password string, displayName string, uid string) *EstimaUser {
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

	var entity Entity = user
	println(entity)

	return &user
}

func (user EstimaUser) AraDoc() (ara.Document) {
	return user.Document
}

func (user EstimaUser) ToJson () ([]byte, error) {
	return json.Marshal(user)
}

func (user EstimaUser) Copy (entity Entity) {
	var from EstimaUser = entity.(EstimaUser)
	user.Name = from.Name
	user.Email = from.Email
	user.DisplayName = from.DisplayName
	user.Password = from.Password
	user.Document = from.Document
	user.Uid = from.Uid
}

func (user EstimaUser) FromJson (jsUser []byte) (error) {
	var retUser EstimaUser
	err := json.Unmarshal(jsUser, &retUser)
	if err == nil {
		user.Copy(retUser)
	}

	return err
}


/**
 Function used to check user name and password through the LDAP.
 Additional documentation about LDAP library located here https://godoc.org/gopkg.in/ldap.v2
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
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer l.Close()

	cn := fmt.Sprintf("cn=%s,%s", username, config.Ldap.Dn)
	println ("cn: " + cn)
	println ("user: " + username)
	// Authenticate using given username and password
	err = l.Bind(cn, password)
	if err != nil {
		log.Panic(err)
		println ("error occurred")
		return nil, err
	}

	// Search for the uswr details
	searchRequest := ldap.NewSearchRequest (
		config.Ldap.Dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		20,
		0,
		false,
		fmt.Sprintf("(&(objectClass=person)(cn=%s))", username),
		[]string{"dn", "uid", "givenName", "displayName", "sn", "cn", "name", "mail", "sAMAccountName"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Print(err)
	}

	for i:=0;i<len(sr.Entries);i++ {
		entry := sr.Entries[i]
		for a:=0;a<len(entry.Attributes);a++ {
			attr := entry.Attributes[a]
			fmt.Println(attr.Name + " = " + attr.Values[0])
		}
		fmt.Println("----------------------------------------")
	}

	entry := sr.Entries[0]
	retUser = NewUser(
		username,
		entry.GetAttributeValue("mail"),
		password,
		entry.GetAttributeValue("displayName"),
		entry.GetAttributeValue("uid"))

	return retUser, retErr
}
