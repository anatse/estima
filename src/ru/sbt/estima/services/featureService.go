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

// Service retrieves comments for the feature. All the comments related for the whole feature not for feature text.
// This is done by this way because text can has many versions and comments may be lost.
// For paging results using next query (url) parameters
// offset - start offset for results
// pageSize - number of results
func (fs FeatureService) getComments (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]

	values := r.URL.Query()
	offset := model.GetInt(values, "offset", 0)
	pageSize := model.GetInt(values, "pageSize", 0)

	comments, err := fs.getDao().GetComments(feature, pageSize, offset)
	model.CheckErr(err)


	iComments := make ([]interface{}, len(comments))
	for idx, com := range comments {
		iComments[idx] = map[string]interface{} {
			"text": com.Comment.Text,
			"changeDate": com.Comment.CreateDate,
			"title": com.Comment.Title,
			"user": com.User,
		}
	}

	model.WriteAnyResponse (true, nil, iComments, w)
}

// Function adds comment to the feature
// This function attach comment directly to feature not for feature text
func (fs FeatureService) addComment (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]

	var com model.Comment
	model.ReadJsonBody(r, &com)

	user := model.GetUserFromRequest (w, r)
	comment := fs.getDao().AddComment(feature, com.Title, com.Text, user.GetCollection() + "/" + user.Key)

	model.WriteResponse(true, nil, comment, w)
}

// Service retrieve all text versions for the given feature
func (fs FeatureService) getTextVersionList(w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]

	versions, err := fs.getDao().GetTextVersionList(feature)
	model.CheckErr(err)

	model.WriteArrayResponse(true, nil, versions, w)
}

func (fs FeatureService) addFeature (w http.ResponseWriter, r *http.Request) {
	var process model.Process
	process.Key = mux.Vars(r)["id"]

	var feature model.Feature
	model.ReadJsonBody(r, &feature)

	featureId := fs.getDao().createAndConnectObjTx(feature, process, PRJ_EDGES, map[string]string{"label": "feature"})
	feature.Id = featureId

	model.WriteResponse(true, nil, feature, w)
}

func (fs FeatureService) removeFeature (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]
	fs.getDao().removeConnectedTx (feature.GetCollection(), PRJ_EDGES, feature.GetKey())
	model.WriteResponse(true, nil, feature, w)
}

func (fs FeatureService) updateFeature (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	model.ReadJsonBody(r, &feature)

	feature.Key = mux.Vars(r)["id"]
	entity, err := fs.getDao().Save(&feature)
	model.CheckErr(err)

	model.WriteResponse(true, nil, entity, w)
}

func (fs *FeatureService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle("/api/v.0.0.1/process/{id}/feature/list", handler(http.HandlerFunc(fs.findByProcess))).Methods("POST", "GET").Name("Features list for specified process")
	router.Handle("/api/v.0.0.1/process/{id}/feature/add", handler(http.HandlerFunc(fs.addFeature))).Methods("POST").Name("Add feature to process")
	router.Handle("/api/v.0.0.1/feature/{id}/update", handler(http.HandlerFunc(fs.updateFeature))).Methods("POST").Name("Update feature information")
	router.Handle("/api/v.0.0.1/feature/{id}", handler(http.HandlerFunc(fs.findById))).Methods("GET").Name("Get feature by id")
	router.Handle("/api/v.0.0.1/feature/{id}/remove", handler(http.HandlerFunc(fs.removeFeature))).Methods("POST", "DELETE").Name("Remove feature")
	router.Handle("/api/v.0.0.1/feature/{id}/text", handler(http.HandlerFunc(fs.getText))).Methods("POST", "GET").Name("Get text for feature. If version parameter is provided then return specified version of text")
	router.Handle("/api/v.0.0.1/feature/{id}/text/list", handler(http.HandlerFunc(fs.getTextVersionList))).Methods("POST", "GET").Name("Get all text versions for feature")
	router.Handle("/api/v.0.0.1/feature/{id}/addtext", handler(http.HandlerFunc(fs.addText))).Methods("POST").Name("Add text for feature")
	router.Handle("/api/v.0.0.1/feature/{id}/status", handler(http.HandlerFunc(fs.setStatus))).Methods("POST").Name("Set feature status")
	router.Handle("/api/v.0.0.1/feature/{id}/comments", handler(http.HandlerFunc(fs.getComments))).Methods("GET", "POST").Name("Get feature comments (use paging pageSize and offset URL parameters)")
	router.Handle("/api/v.0.0.1/feature/{id}/addcomment", handler(http.HandlerFunc(fs.addComment))).Methods("POST").Name("Add comment to feature")
}