package model

import (
	ara "github.com/diegogub/aranGO"
	"time"
)

type VersionedText struct {
	ara.Document `json:-`
	Text string `json:"text,omitempty,required"`
	Version int `json:"version,required"`
	CreateDate time.Time `json:"createDate,omitempty"`
	Active bool `json:"active,required"`
}

func (txt VersionedText) Entity() interface{} {
	return struct{
		*VersionedText

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&txt,
		nil,
		nil,
		nil,
		nil,
	}
}

func (txt VersionedText) AraDoc() (ara.Document) {
	return txt.Document
}

func (txt VersionedText)GetKey() string {
	return txt.Key
}

func (txt VersionedText) GetCollection() string {
	return "verstext"
}

func (txt VersionedText) GetError()(string, bool) {
	return txt.Message, txt.Error
}

func (txt VersionedText) CopyChanged (entity Entity) Entity {
	// Can't be changed
	return entity
}