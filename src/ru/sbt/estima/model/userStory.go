package model

import (
	ara "github.com/diegogub/aranGO"
)

// Structure represents User Story (smallest part of business requirement)
type UserStory struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty,required"`
	Description string `json:"description,omitempty"`
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
	return "stories"
}

func (us UserStory) GetError()(string, bool) {
	return us.Message, us.Error
}
