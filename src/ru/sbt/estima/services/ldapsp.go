package services

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/auth0/go-jwt-middleware"
	"math/rand"
	"encoding/json"
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

var src = rand.NewSource(time.Now().UnixNano())
func randStringBytesMaskImprSrc(n int) string {
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

var mySigningKey = []byte(randStringBytesMaskImprSrc(64))
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	defer (func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetTokenHandler:", r)
			var err model.ErrorObj = model.ErrorObj{
				"get-token",
				fmt.Sprint(r),
				"001",
			}
			js, _ := json.Marshal(err)
			w.Header().Set("Content-Type", "application/json;utf-8")
			w.Write([]byte(js))
		}
	})()

	username := r.URL.Query().Get("uname")
	password := r.URL.Query().Get("upass")
	if username == "" || password == "" {
		panic("Username or/and password doesn't provided")
	}

	// Get user from LDAP
	user, err := model.FindUser(username, password)
	if err != nil {
		panic(err)
	}

	// Update user information in database
	userEntity, err := NewUserDao ().Save(*user)
	if err != nil {
		panic(err)
	}

	*user = userEntity.(model.EstimaUser)

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

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	w.Header().Add("Authorization", "Bearer " + tokenString)
	var config = conf.LoadConfig()
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

	/* Finally, write the token to the browser window */
	w.Header().Set("Content-Type", "application/json;utf-8")
	w.Write([]byte("{success: true}"))
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

		// TODO: Make this a bit more robust, parsing-wise
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

