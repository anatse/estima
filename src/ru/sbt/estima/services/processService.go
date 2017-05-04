package services

import (
	"github.com/gorilla/mux"
	"net/http"
	"ru/sbt/estima/model"
)

// Function for REST services
type ProcessService struct {
	dao *processDao
}

func (ps ProcessService) getDao() processDao {
	if ps.dao == nil {
		ps.dao = NewProcessDao()
	}

	return *ps.dao
}

func (ps ProcessService) findByStage (w http.ResponseWriter, r *http.Request) {
	stageId := mux.Vars(r)["id"]

	processList, err := ps.getDao().FindByStage(stageId)
	model.CheckErr (err)

	model.WriteArrayResponse(true, nil, processList, w)
}

func (ps ProcessService) create (w http.ResponseWriter, r *http.Request) {
	stageId := mux.Vars(r)["id"]
	var stage model.Stage
	stage.SetKey(stageId)

	var prc model.Process
	model.ReadJsonBody(r, &prc)

	prc.Id = ps.getDao().createAndConnectObjTx(
		prc,
		stage,
		PRJ_EDGES,
		nil)

	err := ps.getDao().FindById(&prc)
	model.CheckErr(err)

	model.WriteResponse(true, nil, prc, w)
}

func (ps ProcessService) setStatus (w http.ResponseWriter, r *http.Request) {
	var prc model.Process
	prc.Name = mux.Vars(r)["id"]
	var status struct {
		Name string `json:"name"`
		Status string `json:"status"`
	}

	model.ReadJsonBodyAny(r, &status)
	ps.getDao().SetStatus(&prc, status.Status)
	model.WriteResponse(true, nil, prc, w)
}

func (ps ProcessService) remove (w http.ResponseWriter, r *http.Request) {
	var prc model.Process
	prc.Name = mux.Vars(r)["id"]
	ps.getDao().removeConnectedTx (prc.GetCollection(), PRJ_EDGES, prc.GetKey())
	model.WriteResponse(true, nil, prc, w)
}

func (ps ProcessService) findOne (w http.ResponseWriter, r *http.Request) {
	var p model.Process
	p.Name = mux.Vars(r)["id"] // Number field used as identifier

	err := ps.getDao().FindById(&p)
	model.CheckErr (err)

	model.WriteResponse(true, nil, p, w)
}

func (ps ProcessService) updateProcess (w http.ResponseWriter, r *http.Request) {
	var p model.Process
	p.Name = mux.Vars(r)["id"] // Number field used as identifier
	err := ps.getDao().FindById(&p)
	model.CheckErr (err)

	var prj model.Process
	model.ReadJsonBody (r, &prj)
	prj.Id = p.Id

	pe, err := ps.getDao().Save(prj)
	model.CheckErr (err)

	model.WriteResponse(true, nil, pe, w)
}

func (ps *ProcessService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle ("/stage/{id}/process/list", handler(http.HandlerFunc(ps.findByStage))).Methods("POST", "GET").Name("Process list for selected project stage")
	router.Handle ("/stage/{id}/process/create", handler(http.HandlerFunc(ps.create))).Methods("POST").Name("Create and app process to projects stage")
	router.Handle ("/process/{id}/status", handler(http.HandlerFunc(ps.setStatus))).Methods("POST").Name("Set process status")
	router.Handle ("/process/{id}/remove", handler(http.HandlerFunc(ps.remove))).Methods("POST", "DELETE").Name("Remove process")
	router.Handle ("/process/{id}", handler(http.HandlerFunc(ps.findOne))).Methods("GET").Name("GET: Retrieve information about selected process")
	router.Handle ("/process/{id}", handler(http.HandlerFunc(ps.updateProcess))).Methods("POST").Name("POST: Update selected process")
}
