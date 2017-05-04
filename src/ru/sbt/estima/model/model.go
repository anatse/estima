package model

import (
	ara "github.com/diegogub/aranGO"
	"time"
	"crypto/des"
	"bytes"
	"crypto/cipher"
	"fmt"
)

type Entity interface {
	AraDoc() (ara.Document)
	Entity() interface{}
	GetKey() string
	GetCollection() string
}

func byteToString (bts []byte) string {
	var buffer bytes.Buffer
	for i := 0; i< len(bts); i++ {
		buffer.WriteString(fmt.Sprintf("%d", bts[i]))
	}

	return buffer.String()
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func GetCryptedKey (key string, salt string) string {
	block, err := des.NewCipher([]byte(salt))
	//cipher, err := des.NewCipher()
	if err != nil {
		panic(err)
	}

	//var dest []byte = make([]byte, len(key))
	blockSize := block.BlockSize()
	fmt.Println(blockSize)

	origData := PKCS5Padding([]byte(key), blockSize)
	fmt.Println(origData)

	blockMode := cipher.NewCBCEncrypter(block, []byte(salt))
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)

	//fmt.Println(base64.URLEncoding.EncodeToString(cryted))

	data := make ([]byte, len(origData))
	dblockMode := cipher.NewCBCDecrypter(block, []byte(salt))
	dblockMode.CryptBlocks(data, cryted)

	data = PKCS5UnPadding (data)
	return byteToString(cryted)
}

//
// jsom samoples https://eager.io/blog/go-and-json/
// https://mholt.github.io/json-to-go/
//
type JSONTime time.Time

func (t JSONTime)MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := time.Now().String()
	//stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
	return []byte(stamp), nil
}
