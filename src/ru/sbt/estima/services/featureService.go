package services

import (
	"github.com/gorilla/mux"
	"net/http"
	"ru/sbt/estima/model"
	"fmt"
)

// Function for REST services
type FeatureService struct {
	dao *featureDao
}

func (fs FeatureService) getDao() featureDao {
	if fs.dao == nil {
		fs.dao = NewFeatureDao()
	}

	return *fs.dao
}

func (fs FeatureService) findByProcess (w http.ResponseWriter, r *http.Request) {
	found, err := fs.getDao().FindByProcess(fmt.Sprintf("%s/%s", model.Process{}.GetCollection(), mux.Vars(r)["id"]))
	model.CheckErr(err)
	model.WriteArrayResponse(true, nil, found, w)
}

func (fs *FeatureService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle("/process/{id}/feature/list", handler(http.HandlerFunc(fs.findByProcess))).Methods("POST", "GET").Name("Features list for specified process")
}

func (fs *FeatureService) getText () {

}