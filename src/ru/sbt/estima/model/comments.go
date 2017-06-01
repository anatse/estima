package model

import (
	ara "github.com/diegogub/aranGO"
	"time"
)

type Comment struct {
	ara.Document `json:-`
	Title string `json:"name,omitempty"`
	Text string `json:"text,omitempty,required"`
	CreateDate time.Time `json:"createDate,omitempty"`
}

type CommentWithUser struct {
	Comment Comment `json:"comment"`
	User EstimaUser `json:"user"`
}

func (com Comment) Entity() interface{} {
	return struct{
		*Comment

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&com,
		nil,
		nil,
		nil,
		nil,
	}
}

func (com Comment) CopyChanged (entity Entity) Entity {
	panic("Comment cannot be changed")
}

func (com Comment) AraDoc() (ara.Document) {
	return com.Document
}

func (com Comment)GetKey() string {
	return com.Key
}

func (com Comment) GetCollection() string {
	return "comments"
}

func (com Comment) GetError()(string, bool) {
	// default error bool and messages. Could be any kind of error
	return com.Message, com.Error
}
