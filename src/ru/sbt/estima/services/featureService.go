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

// Service get feature list for given process
func (fs FeatureService) findByProcess (w http.ResponseWriter, r *http.Request) {
	var process model.Process
	process.Key = mux.Vars(r)["id"]
	found, err := fs.getDao().FindByProcess(fmt.Sprintf("%s/%s", process.GetCollection(), process.GetKey()))
	model.CheckErr(err)
	model.WriteArrayResponse(true, nil, found, w)
}

// Service get text for given feature and, optional, text version
// If the version is not specified then service returns currently active text
func (fs *FeatureService) getText (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]

	// Check for version parameter
	values := r.URL.Query()
	version := model.GetInt (values, "version", -1)
	var found *model.VersionedText
	var err error
	if version != -1 {
		found, err = fs.getDao().GetTextByVersion(feature, version)
		model.CheckErr(err)
	} else {
		found, err = fs.getDao().GetActiveText(feature)
		model.CheckErr(err)
	}

	model.WriteResponse(true, nil, found, w)
}

// Service find and retrieve feature by id (key)
func (fs FeatureService) findById (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]
	err := fs.getDao().FindById(&feature)
	model.CheckErr(err)
	model.WriteResponse(true, nil, feature, w)
}

// Service add text to feature
func (fs FeatureService) addText (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]

	var text model.VersionedText
	model.ReadJsonBody(r, &text)

	vtext := fs.getDao().AddText(feature, text.Text)
	model.WriteResponse(true, nil, vtext, w)
}

// Service add text to feature
func (fs FeatureService) setStatus (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	model.ReadJsonBody(r, &feature)
	feature.Key = mux.Vars(r)["id"]

	entity, err := fs.getDao().SetStatus (&feature, feature.Status)
	model.CheckErr(err)
	model.WriteResponse(true, nil, entity, w)
}

func (fs *FeatureService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle("/api/v.0.0.1/process/{id}/feature/list", handler(http.HandlerFunc(fs.findByProcess))).Methods("POST", "GET").Name("Features list for specified process")
	router.Handle("/api/v.0.0.1/feature/{id}", handler(http.HandlerFunc(fs.findById))).Methods("GET").Name("Get feature by id")
	router.Handle("/api/v.0.0.1/feature/{id}/text", handler(http.HandlerFunc(fs.getText))).Methods("POST", "GET").Name("Get text for feature. If version parameter is provided then return specified version of text")
	router.Handle("/api/v.0.0.1/feature/{id}/addtext", handler(http.HandlerFunc(fs.addText))).Methods("POST").Name("Add text for feature")
	router.Handle("/api/v.0.0.1/feature/{id}/status", handler(http.HandlerFunc(fs.setStatus))).Methods("POST").Name("Set feature status")
}