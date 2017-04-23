package model

import (
	ara "github.com/diegogub/aranGO"
	"time"
)

type Entity interface {
	AraDoc() (ara.Document)
	Entity() interface{}
	GetKey() string
	GetCollection() string
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
