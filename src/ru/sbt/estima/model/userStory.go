package model

import (
	ara "github.com/diegogub/aranGO"
	"time"
)

// Structure represents User Story (smallest part of business requirement)
type UserStory struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty,required"`
	Description string `json:"description,omitempty"`
	// Who - owner of this story
	Who string `json:"who,omitempty"`
	// What - subject of the story
	What string `json:"what,omitempty"`
	// Why - reason
	Why string `json:"why,omitempty"`
	// Status of the story
	Status Status `json:"status,omitempty"`
	CreateDate time.Time `json:"createDate,omitempty"`
	Serial int `json:"serial"`
}

func (us UserStory) Entity() interface{} {
	return struct{
		*UserStory

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&us,
		nil,
		nil,
		nil,
		nil,
	}
}

func (us UserStory) AraDoc() (ara.Document) {
	return us.Document
}

func (us UserStory)GetKey() string {
	return us.Key
}

func (us UserStory) GetCollection() string {
	return "ustories"
}

func (us UserStory) GetError()(string, bool) {
	return us.Message, us.Error
}

func (us UserStory) CopyChanged (entity Entity) Entity {
	// Status cannot be changed by saving all object, it needs to use separate function setStatus
	panic("Not allowed to change user story")
}
