package model

import (
	ara "github.com/diegogub/aranGO"
	"math/big"
	"time"
)

// Structure represents User Story (smallest part of business requirement)
type TechStoryPrice struct {
	ara.Document `json:-`
	Group string `json:"group,omitempty,required"`
	StoryPoints big.Float `json:"storyPoints,omitempty"`
	Calculated time.Time `json:"calcDate,omitempty"`
}

func (us TechStoryPrice) Entity() interface{} {
	return struct{
		*TechStoryPrice

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

func (us TechStoryPrice) AraDoc() (ara.Document) {
	return us.Document
}

func (us TechStoryPrice)GetKey() string {
	return us.Key
}

func (us TechStoryPrice) GetCollection() string {
	return "tsprices"
}

func (us TechStoryPrice) GetError()(string, bool) {
	return us.Message, us.Error
}
