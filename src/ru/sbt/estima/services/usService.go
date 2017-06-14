package services

import (
	"net/http"
	"ru/sbt/estima/model"
	"github.com/gorilla/mux"
	"log"
)

// Function for REST services
type UserStoryService struct {
}

func (us UserStoryService) addUserStory (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]

	var userStory model.UserStory
	model.ReadJsonBody(r, &userStory)

	WithUserStoryDao(func(dao userStoryDao) {
		userStoryId := dao.createAndConnectObjTx(userStory, feature, model.PRJ_EDGES, map[string]string{"label": "userStory"})
		userStory.Id = userStoryId
		model.WriteResponse(true, nil, userStory, w)
	})
}

func (us UserStoryService) removeUserStory (w http.ResponseWriter, r *http.Request) {
	var uss model.UserStory
	uss.Key = mux.Vars(r)["id"]

	WithUserStoryDao(func(dao userStoryDao) {
		dao.removeConnectedTx (uss.GetCollection(), model.PRJ_EDGES, uss.GetKey())
		model.WriteResponse(true, nil, uss, w)
	})
}

func (us UserStoryService) listUserStory (w http.ResponseWriter, r *http.Request) {
	var feature model.Feature
	feature.Key = mux.Vars(r)["id"]

	WithUserStoryDao(func(dao userStoryDao) {
		uss, err := dao.FindByFeature(feature.GetCollection() + "/" + feature.Key)
		model.CheckErr(err)
		model.WriteArrayResponse(true, nil, uss, w)
	})
}

// Service get text for given feature and, optional, text version
// If the version is not specified then service returns currently active text
func (fs *UserStoryService) getText (w http.ResponseWriter, r *http.Request) {
	var userStory model.UserStory
	userStory.Key = mux.Vars(r)["id"]

	// Check for version parameter
	values := r.URL.Query()
	version := model.GetInt (values, "version", -1)
	var found *model.VersionedText
	var err error

	WithUserStoryDao(func(dao userStoryDao) {
		if version != -1 {
			found, err = dao.GetTextByVersion(userStory, version)
		} else {
			found, err = dao.GetActiveText(userStory)
		}
		model.CheckErr(err)
		model.WriteResponse(true, nil, found, w)
	})
}

// Service find and retrieve feature by id (key)
func (fs UserStoryService) findById (w http.ResponseWriter, r *http.Request) {
	var userStory model.UserStory
	userStory.Key = mux.Vars(r)["id"]
	WithUserStoryDao(func(dao userStoryDao) {
		err := dao.FindById(&userStory)
		model.CheckErr(err)
		model.WriteResponse(true, nil, userStory, w)
	})
}

// Service add text to feature
func (fs UserStoryService) addText (w http.ResponseWriter, r *http.Request) {
	var userStory model.UserStory
	userStory.Key = mux.Vars(r)["id"]

	var text model.VersionedText
	model.ReadJsonBody(r, &text)
	WithUserStoryDao(func(dao userStoryDao) {
		vtext := dao.AddText(userStory, text.Text)
		model.WriteResponse(true, nil, vtext, w)
	})
}

// Service add text to feature
func (fs UserStoryService) setStatus (w http.ResponseWriter, r *http.Request) {
	var userStory model.UserStory
	model.ReadJsonBody(r, &userStory)
	userStory.Key = mux.Vars(r)["id"]

	WithUserStoryDao(func(dao userStoryDao) {
		user := model.GetUserFromRequest(w, r)
		entity, err := dao.SetStatus(&userStory, userStory.Status, user)
		model.CheckErr(err)
		model.WriteResponse(true, nil, entity, w)
	})
}

// Service retrieves comments for the feature. All the comments related for the whole feature not for feature text.
// This is done by this way because text can has many versions and comments may be lost.
// For paging results using next query (url) parameters
// offset - start offset for results
// pageSize - number of results
func (fs UserStoryService) getComments (w http.ResponseWriter, r *http.Request) {
	var userStory model.UserStory
	userStory.Key = mux.Vars(r)["id"]

	values := r.URL.Query()
	offset := model.GetInt(values, "offset", 0)
	pageSize := model.GetInt(values, "pageSize", 0)

	WithUserStoryDao(func(dao userStoryDao) {
		comments, err := dao.GetComments(userStory, pageSize, offset)
		model.CheckErr(err)

		iComments := make([]interface{}, len(comments))
		for idx, com := range comments {
			iComments[idx] = map[string]interface{}{
				"text":       com.Comment.Text,
				"createDate": com.Comment.CreateDate,
				"title":      com.Comment.Title,
				"user":       com.User,
			}
		}

		model.WriteAnyResponse(true, nil, iComments, w)
	})
}

// Function adds comment to the feature
// This function attach comment directly to feature not for feature text
func (fs UserStoryService) addComment (w http.ResponseWriter, r *http.Request) {
	var userStory model.UserStory
	userStory.Key = mux.Vars(r)["id"]

	var com model.Comment
	model.ReadJsonBody(r, &com)

	WithUserStoryDao(func(dao userStoryDao) {
		user := model.GetUserFromRequest(w, r)
		comment := dao.AddComment(userStory, com.Title, com.Text, user.GetCollection()+"/"+user.Key)

		model.WriteResponse(true, nil, comment, w)
	})
}

// Service retrieve all text versions for the given feature
func (fs UserStoryService) getTextVersionList(w http.ResponseWriter, r *http.Request) {
	var userStory model.UserStory
	userStory.Key = mux.Vars(r)["id"]

	WithUserStoryDao(func(dao userStoryDao) {
		versions, err := dao.GetTextVersionList(userStory)
		model.CheckErr(err)
		model.WriteArrayResponse(true, nil, versions, w)
	})
}

func (fs UserStoryService) updateUserStory (w http.ResponseWriter, r *http.Request) {
	var userStory model.UserStory
	model.ReadJsonBody(r, &userStory)
	log.Printf("Trying to change US.status = %v", userStory.Status)

	WithUserStoryDao(func(dao userStoryDao) {
		userStory.Key = mux.Vars(r)["id"]
		entity, err := dao.Save(&userStory)
		model.CheckErr(err)
		model.WriteResponse(true, nil, entity, w)
	})
}

func (us *UserStoryService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle("/api/v.0.0.1/feature/{id}/userstory/list", handler(http.HandlerFunc(us.listUserStory))).Methods("POST", "GET").Name("User story list for specified feature")
	router.Handle("/api/v.0.0.1/feature/{id}/userstory/add", handler(http.HandlerFunc(us.addUserStory))).Methods("POST").Name("Add user story to feature")
	router.Handle("/api/v.0.0.1/userstory/{id}/update", handler(http.HandlerFunc(us.updateUserStory))).Methods("POST").Name("Update user story")
	router.Handle("/api/v.0.0.1/userstory/{id}/remove", handler(http.HandlerFunc(us.removeUserStory))).Methods("POST", "DELETE").Name("Remove user story")
	router.Handle("/api/v.0.0.1/userstory/{id}/get", handler(http.HandlerFunc(us.findById))).Methods("POST", "GET").Name("Get user story by key")
	router.Handle("/api/v.0.0.1/userstory/{id}/status", handler(http.HandlerFunc(us.setStatus))).Methods("POST").Name("Change status of user story")
	router.Handle("/api/v.0.0.1/userstory/{id}/addcomment", handler(http.HandlerFunc(us.addComment))).Methods("POST").Name("Add comment to user story")
	router.Handle("/api/v.0.0.1/userstory/{id}/comment", handler(http.HandlerFunc(us.getComments))).Methods("POST", "GET").Name("Get comments of user story")
	router.Handle("/api/v.0.0.1/userstory/{id}/addtext", handler(http.HandlerFunc(us.addText))).Methods("POST").Name("Add text to user story")
	router.Handle("/api/v.0.0.1/userstory/{id}/text", handler(http.HandlerFunc(us.getText))).Methods("POST", "GET").Name("Get active text for user story")
	router.Handle("/api/v.0.0.1/userstory/{id}/text/list", handler(http.HandlerFunc(us.getTextVersionList))).Methods("POST", "GET").Name("Get all text versions for user story")
}