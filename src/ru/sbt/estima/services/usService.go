package services

import (
	"net/http"
	"ru/sbt/estima/model"
	"github.com/gorilla/mux"
)

// Function for REST services
type UserStoryService struct {
	dao *userStoryDao
}

func (fs UserStoryService) getDao() userStoryDao {
	if fs.dao == nil {
		fs.dao = NewUserStoryDao()
	}

	return *fs.dao
}

func (fs UserStoryService) addUserStory (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]

	var userStory model.UserStory
	model.ReadJsonBody(r, &userStory)

	featureId := fs.getDao().createAndConnectObjTx(userStory, feature, model.PRJ_EDGES, map[string]string{"label": "feature"})
	feature.Id = featureId

	model.WriteResponse(true, nil, feature, w)
}

//func (fs FeatureService) removeFeature (w http.ResponseWriter, r *http.Request) {
//	var feature model.Feature
//	feature.Key = mux.Vars(r)["id"]
//	fs.getDao().removeConnectedTx (feature.GetCollection(), PRJ_EDGES, feature.GetKey())
//	model.WriteResponse(true, nil, feature, w)
//}
//
//func (fs FeatureService) updateFeature (w http.ResponseWriter, r *http.Request) {
//	var feature model.Feature
//	model.ReadJsonBody(r, &feature)
//
//	feature.Key = mux.Vars(r)["id"]
//	entity, err := fs.getDao().Save(&feature)
//	model.CheckErr(err)
//
//	model.WriteResponse(true, nil, entity, w)
//}

func (fs *UserStoryService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	//router.Handle("/api/v.0.0.1/process/{id}/feature/list", handler(http.HandlerFunc(fs.findByProcess))).Methods("POST", "GET").Name("Features list for specified process")
	//router.Handle("/api/v.0.0.1/process/{id}/feature/add", handler(http.HandlerFunc(fs.addFeature))).Methods("POST").Name("Add feature to process")
	//router.Handle("/api/v.0.0.1/feature/{id}/update", handler(http.HandlerFunc(fs.updateFeature))).Methods("POST").Name("Update feature information")
	//router.Handle("/api/v.0.0.1/feature/{id}", handler(http.HandlerFunc(fs.findById))).Methods("GET").Name("Get feature by id")
	//router.Handle("/api/v.0.0.1/feature/{id}/remove", handler(http.HandlerFunc(fs.removeFeature))).Methods("POST", "DELETE").Name("Remove feature")
	//router.Handle("/api/v.0.0.1/feature/{id}/text", handler(http.HandlerFunc(fs.getText))).Methods("POST", "GET").Name("Get text for feature. If version parameter is provided then return specified version of text")
	//router.Handle("/api/v.0.0.1/feature/{id}/text/list", handler(http.HandlerFunc(fs.getTextVersionList))).Methods("POST", "GET").Name("Get all text versions for feature")
	//router.Handle("/api/v.0.0.1/feature/{id}/addtext", handler(http.HandlerFunc(fs.addText))).Methods("POST").Name("Add text for feature")
	//router.Handle("/api/v.0.0.1/feature/{id}/status", handler(http.HandlerFunc(fs.setStatus))).Methods("POST").Name("Set feature status")
	//router.Handle("/api/v.0.0.1/feature/{id}/comments", handler(http.HandlerFunc(fs.getComments))).Methods("GET", "POST").Name("Get feature comments (use paging pageSize and offset URL parameters)")
	//router.Handle("/api/v.0.0.1/feature/{id}/addcomment", handler(http.HandlerFunc(fs.addComment))).Methods("POST").Name("Add comment to feature")
}