package services

import (
	"github.com/gorilla/mux"
	"net/http"
	"ru/sbt/estima/model"
	"log"
)

// Function for REST services
type ComponentService struct {
}

func (cs ComponentService) listComponent (w http.ResponseWriter, r *http.Request) {
	WithComponentDao(func(dao componentDao) {
		cmps := dao.FindAllComponents()
		model.WriteArrayResponse(true, nil, cmps, w)
	})
}

func (cs ComponentService) addComponent (w http.ResponseWriter, r *http.Request) {
	var cmp model.Component
	model.ReadJsonBody(r, &cmp)

	WithComponentDao(func(dao componentDao) {
		cmp, err := dao.Save(cmp)
		model.CheckErr(err)
		model.WriteResponse(true, nil, cmp, w)
	})
}

func (cs ComponentService) updateComponent (w http.ResponseWriter, r *http.Request) {
	var cmp model.Component
	model.ReadJsonBody(r, &cmp)

	cmp.Key = mux.Vars(r)["id"]

	WithComponentDao(func(dao componentDao) {
		cmp, err := dao.Save(cmp)
		model.CheckErr(err)
		model.WriteResponse(true, nil, cmp, w)
	})
}

func (cs ComponentService) removeComponent (w http.ResponseWriter, r *http.Request) {
	var cmp model.Component
	model.ReadJsonBody(r, &cmp)
	cmp.Key = mux.Vars(r)["id"]

	WithComponentDao(func(dao componentDao) {
		err := dao.Database().Col(cmp.GetCollection()).Delete(cmp.GetKey())
		model.CheckErr(err)
		model.WriteResponse(true, nil, cmp, w)
	})
}

func (us ComponentService) addComponentToFeature (w http.ResponseWriter, r *http.Request) {
	feature := model.Feature{}
	feature.Key = mux.Vars(r)["id"]
	feature.Id = feature.GetCollection() + "/" + feature.GetKey()

	var cmp model.Component
	model.ReadJsonBody(r, &cmp)
	cmp.Id = cmp.GetCollection() + "/" + cmp.GetKey()

	log.Printf("feature.Id: %v, component: %v", feature.Id, cmp.Id)

	WithComponentDao(func(dao componentDao) {
		err := dao.ConnectComponent(cmp, feature)
		model.CheckErr(err)
		model.WriteResponse(true, nil, cmp, w)
	})
}

func (us ComponentService) listComponentConsumer (w http.ResponseWriter, r *http.Request) {
	var cmp model.Component
	cmp.Key =  mux.Vars(r)["id"]
	cmp.Id = cmp.GetCollection() + "/" +cmp.GetKey()

	WithComponentDao(func(dao componentDao) {
		consumers := dao.FindAllConsumer(cmp)
		model.WriteAnyResponse(true, nil, consumers, w)
	})
}

func (us ComponentService) removeComponentFromFeature (w http.ResponseWriter, r *http.Request) {
	log.Printf("removeComponentFromFeature")

	var cmp model.Component
	cmp.Id = cmp.GetCollection() + "/" + mux.Vars(r)["cmpId"]

	var feature model.Feature
	feature.Id  = feature.GetCollection() + "/" + mux.Vars(r)["id"]

	WithComponentDao(func(dao componentDao) {
		dao.DisconnectComponent(cmp, feature)
		model.WriteResponse(true, nil, nil, w)
	})
}

func (us ComponentService) listComponentForFeature (w http.ResponseWriter, r *http.Request) {
	feature := model.Feature{}
	feature.Key = mux.Vars(r)["id"]
	feature.Id = feature.GetCollection() + "/" + feature.GetKey()

	WithComponentDao(func(dao componentDao) {
		cmps := dao.FindConnectedComponents(feature)
		model.WriteArrayResponse(true, nil, cmps, w)
	})
}

func (us *ComponentService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle("/api/v.0.0.1/component/list", handler(http.HandlerFunc(us.listComponent))).Methods("POST", "GET").Name("Components list")
	router.Handle("/api/v.0.0.1/component/add", handler(http.HandlerFunc(us.addComponent))).Methods("POST").Name("Add component")
	router.Handle("/api/v.0.0.1/component/{id}/update", handler(http.HandlerFunc(us.updateComponent))).Methods("POST").Name("Update component")
	router.Handle("/api/v.0.0.1/component/{id}/remove", handler(http.HandlerFunc(us.removeComponent))).Methods("POST", "DELETE").Name("Remove component")
	router.Handle("/api/v.0.0.1/feature/{id}/addComponent", handler(http.HandlerFunc(us.addComponentToFeature))).Methods("POST", "GET").Name("Attach component to feature")
	router.Handle("/api/v.0.0.1/feature/{id}/listComponent", handler(http.HandlerFunc(us.listComponentForFeature))).Methods("POST", "GET").Name("List components for feature")
	router.Handle("/api/v.0.0.1/feature/{id}/component/{cmpId}/remove", handler(http.HandlerFunc(us.removeComponentFromFeature))).Methods("POST", "DELETE").Name("Attach component to feature")
	router.Handle("/api/v.0.0.1/component/{id}/listConsumer", handler(http.HandlerFunc(us.listComponentConsumer))).Methods("POST", "GET").Name("Get list of component consumers")
}