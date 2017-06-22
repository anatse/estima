package model

import (
	ara "github.com/diegogub/aranGO"
	"time"
)

type CalcUnit struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty,required"`
	Description string `json:"description,omitempty"`
	Version string `json:"version,omitempty"`
	Status int `json:"status,omitempty"`
	Changed time.Time `json:"changed,omitempty"`
}

func (cmp CalcUnit) Entity() interface{} {
	return struct{
		*CalcUnit

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

func (cmp CalcUnit) AraDoc() (ara.Document) {
	return cmp.Document
}

func (cmp CalcUnit)GetKey() string {
	return cmp.Key
}

func (cmp CalcUnit) GetCollection() string {
	return "calcunits"
}

func (cmp CalcUnit) GetError()(string, bool) {
	return cmp.Message, cmp.Error
}

func (cmp CalcUnit) CopyChanged (entity Entity) Entity {
	emptyTime := time.Time{}
	unit := entity.(CalcUnit)
	if unit.Name != "" {cmp.Name = unit.Name}
	if unit.Description != "" {cmp.Description = unit.Description}
	if unit.Version != "" {cmp.Version = unit.Version}
	if unit.Changed != emptyTime {cmp.Changed = unit.Changed}
	return cmp
}

// Information will be copied to edge and participated after in calculation process
type CalcUnitEdge struct {
	CalcUnit
	Weight float64 `json:"weight,omitempty,required"`
	Complexity float64 `json:"complexity,omitempty,required"`
	NewFlag float64 `json:"newFlag,omitempty,required"`
	ExtCoef float64 `json:"extCoef,omitempty,required"`
}

func (cmp CalcUnitEdge) Entity() interface{} {
	return struct{
		*CalcUnitEdge

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