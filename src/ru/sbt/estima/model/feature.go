package model

import (
	ara "github.com/diegogub/aranGO"
)

// Structure represents Feature (atomic business part of end to end process)
type Feature struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty,required"`
	Description string `json:"description,omitempty"`
	Status Status `json:"status,omitempty"`
	// Важность - величина обратная приоритету, используется для формирования последовательности фич бэклога
	Importance int  `json:"importance,omitempty"`
}

func (fea Feature) Entity() interface{} {
	return struct{
		*Feature

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&fea,
		nil,
		nil,
		nil,
		nil,
	}
}

func (fea Feature) AraDoc() (ara.Document) {
	return fea.Document
}

func (fea Feature)GetKey() string {
	return fea.Key
}

func (fea Feature) GetCollection() string {
	return "features"
}

func (fea Feature) GetError()(string, bool) {
	return fea.Message, fea.Error
}

func (fea Feature) CopyChanged (entity Entity) Entity {
	newFea := entity.(Feature)
	if newFea.Name != "" {fea.Name = newFea.Name}
	// Status cannot be changed by saving all object, it needs to use separate function setStatus
	//if newFea.Status != -1 {fea.Status = newFea.Status}
	if newFea.Description != "" {fea.Description = newFea.Description}
	return fea
}