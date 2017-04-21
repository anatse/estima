package services

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/auth0/go-jwt-middleware"
	"math/rand"
	"fmt"
	"ru/sbt/estima/model"
	"strings"
	"ru/sbt/estima/conf"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// Function calulates random secret key for crypting auth coockie or header
func randStringBytesMaskImprSrc(n int) string {
	secret := conf.LoadConfig().Secret
	if secret != "" {
		return secret
	}

	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// Secret key, can be defined in confi.json file otherwize calculate using randomStringBuyesMaskImprSrc function
var mySigningKey = []byte(randStringBytesMaskImprSrc(64))
var config = conf.LoadConfig()

func createCookie (user model.EstimaUser, w http.ResponseWriter) {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["name"] = user.Name
	claims["displayName"] = user.DisplayName
	claims["mail"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["uid"] = user.Uid
	claims["roles"] = user.Roles

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)
	w.Header().Add("Authorization", "Bearer " + tokenString)
	http.SetCookie(w, &http.Cookie {
		config.Auth.CookieName,
		"Bearer " + tokenString,
		"/",
		"",
		time.Now().AddDate(1, 0, 0),
		"",
		config.Auth.MaxAge,
		false,
		false,
		"",
		nil})
}

func login (w http.ResponseWriter, username string, password string) {
	if username == "" || password == "" {
		panic("Username or/and password doesn't provided")
	}

	// Get user from LDAP
	var user *model.EstimaUser
	var err error
	if config.Ldap.Protocol == "fake" {
		user = model.NewUser(username, username + "@fake.com", password, username, "", nil)
	} else {
		user, err = model.FindUser(username, password)
		if err != nil {
			panic(err)
		}
	}

	// Try to find information from database
	dao := NewUserDao ()
	err = dao.FindOne(user)
	if err != nil {
		panic(err)
	}

	// If user not found
	if user.AraDoc().Id == "" {
		// Update user information in database
		_, err = NewUserDao().Save(user)
		if err != nil {
			panic(err)
		}
	}

	createCookie(*user, w)

	/* Finally, write the token to the browser window */
	model.WriteResponse (true, nil, *user, w)
}

var Login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	defer (func() {
		//if r := recover(); r != nil {
		//	fmt.Println("Recovered in Login:", r)
		//	model.WriteResponse (false, fmt.Sprint(r), nil, w)
		//}
	})()

	var li struct {
		Uname string `json:"uname"`
		Upass string `json:"upass"`
	}

	err := ReadJsonBodyAny(r, &li)
	if err != nil {
		panic(err)
	}

	username := li.Uname
	password := li.Upass
	login(w, username, password)

})

// Function generate new auth token (JWT) and store it in cookie. Also this function store user information in database
// if this user not exists yet
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	defer (func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetTokenHandler:", r)
			model.WriteResponse (false, fmt.Sprint(r), nil, w)
		}
	})()

	username := r.URL.Query().Get("uname")
	password := r.URL.Query().Get("upass")
	login(w, username, password)
})

//
// Using cookies to authenticate user
//
func FromAuthCookie(r *http.Request) (string, error) {
	var config = conf.LoadConfig()
	authHeader, err := r.Cookie(config.Auth.CookieName)
	if authHeader != nil {
		if authHeader.Value == "" {
			return "", nil // No error, just no token
		}

		authHeaderParts := strings.Split(authHeader.Value, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			return "", fmt.Errorf("Authorization cookie format must be Bearer {token}")
		}

		return authHeaderParts[1], nil
	}

	return "", err
}

var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
	Extractor: FromAuthCookie,
})

