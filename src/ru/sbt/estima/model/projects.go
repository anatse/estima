package model

import (
	ara "github.com/diegogub/aranGO"
	"time"
)

type Project struct {
	ara.Document `json:-`
	Number      string `json:"number,omitempty", unique:"projects"`
	Name string `json:"name",omitempty`
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

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&prj,
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
	return prj.Key
}

func (prj Project) GetCollection() string {
	return "projects"
}

func (prj Project) GetError()(string, bool){
	// default error bool and messages. Could be any kind of error
	return prj.Message, prj.Error
}

func (prj Project) CopyChanged (entity Entity) Entity {
	newPrj := entity.(Project)
	emptyTime := time.Time{}
	if newPrj.Name != "" {prj.Name = newPrj.Name}
	if newPrj.Status != "" {prj.Status = newPrj.Status}
	if newPrj.Description != "" {prj.Description = newPrj.Description}
	if newPrj.StartDate != emptyTime {prj.StartDate = newPrj.StartDate}
	if newPrj.EndDate != emptyTime {prj.EndDate = newPrj.EndDate}
	return prj
}

type Stage struct {
	ara.Document `json:-`
	Name     string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status string `json:"status"`
	StartDate time.Time `json:"startDate"`
	EndDate time.Time `json:"endDate"`
}

func NewStage (name string) Stage {
	var prj Stage
	prj.Name = name
	prj.StartDate = time.Now()
	return prj
}

func (stage Stage) Entity() interface{} {
	return struct{
		*Stage

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&stage,
		nil,
		nil,
		nil,
		nil,
	}
}

func (stage Stage) PrjEntity(prjKey string) interface{} {
	return struct{
		*Stage
		ProjectKey string `json:"projectKey,omitempty"`

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&stage,
		prjKey,
		nil,
		nil,
		nil,
		nil,
	}
}

func (stage Stage) AraDoc() (ara.Document) {
	return stage.Document
}

func (stage Stage)GetKey() string {
	return stage.Key
}

func (stage Stage) GetCollection() string {
	return "stages"
}

func (stage Stage) GetError()(string, bool){
	// default error bool and messages. Could be any kind of error
	return stage.Message, stage.Error
}

func (stage Stage) CopyChanged (entity Entity) Entity {
	newStage := entity.(Stage)
	emptyTime := time.Time{}
	if newStage.Name != "" {stage.Name = newStage.Name}
	if newStage.Status != "" {stage.Status = newStage.Status}
	if newStage.Description != "" {stage.Description = newStage.Description}
	if newStage.StartDate != emptyTime {stage.StartDate = newStage.StartDate}
	if newStage.EndDate != emptyTime {stage.EndDate = newStage.EndDate}
	return stage
}