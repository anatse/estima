package model

import (
	ara "github.com/diegogub/aranGO"
	"time"
)

type Project struct {
	ara.Document `json:-`
	Number      string `json:"number,omitempty", unique:"projects"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Flag string `json:"flag" enum:"G,Y,R,B"`
	StartDate   time.Time `json:"startDate,omitempty"`
	EndDate	time.Time `json:"endDate"`
}

func NewPrj (name string) Project {
	var prj Project
	prj.Number = name
	prj.SetKey(name)
	prj.StartDate = time.Now()
	return prj
}

func (prj Project) Entity() interface{} {
	return struct{
		*Project

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`
		OmitKey omit `json:"_key,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&prj,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
}

func (prj Project) AraDoc() (ara.Document) {
	return prj.Document
}

func (prj Project)GetKey() string {
	return prj.Number
}

func (prj Project) GetCollection() string {
	return "projects"
}

func (prj Project) GetError()(string, bool){
	// default error bool and messages. Could be any kind of error
	return prj.Message, prj.Error
}

type Stage struct {
	ara.Document `json:-`
	Name     string `json:"name,omitempty", unique:"projects"`
	Description string `json:"descriptio,omitempty"`
	Status string `json:"status"`
	StartDate time.Time `json:"startDate"`
	EndDate time.Time `json:"endDate"`
}

func NewStage (name string) Stage {
	var prj Stage
	prj.Name = name
	prj.SetKey(name)
	prj.StartDate = time.Now()
	return prj
}

func (stage Stage) Entity() interface{} {
	return struct{
		*Stage

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`
		OmitKey omit `json:"_key,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&stage,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
}

func (stage Stage) AraDoc() (ara.Document) {
	var p Stage
	p.Status = "10"

	return stage.Document
}

func (stage Stage)GetKey() string {
	return stage.Name
}

func (stage Stage) GetCollection() string {
	return "projects"
}

func (stage Stage) GetError()(string, bool){
	// default error bool and messages. Could be any kind of error
	return stage.Message, stage.Error
}