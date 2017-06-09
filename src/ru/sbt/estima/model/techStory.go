package model

import (
	ara "github.com/diegogub/aranGO"
)

type TechStory struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty,required"`
	Description string `json:"description,omitempty"`
	Status Status `json:"status,omitempty"`
}

func (cmp TechStory) Entity() interface{} {
	return struct{
		*TechStory

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&cmp,
		nil,
		nil,
		nil,
		nil,
	}
}

func (cmp TechStory) AraDoc() (ara.Document) {
	return cmp.Document
}

func (ts TechStory) CopyChanged (entity Entity) Entity {
	story := entity.(UserStory)
	if story.Name != "" {ts.Name = story.Name}
	if story.Description != "" {ts.Description = story.Description}
	return ts
}

func (cmp TechStory)GetKey() string {
	return cmp.Key
}

func (cmp TechStory) GetCollection() string {
	return "tstories"
}

func (cmp TechStory) GetError()(string, bool) {
	return cmp.Message, cmp.Error
}

type TechStoryWithText struct {
	UserStory
	Text string `json:"text"`
	Version int `json:"version"`
}

func (ts TechStoryWithText) Entity() interface{} {
	return struct{
		*TechStoryWithText

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&ts,
		nil,
		nil,
		nil,
		nil,
	}
}