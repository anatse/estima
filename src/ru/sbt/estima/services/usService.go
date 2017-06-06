package services

import (
	"net/http"
	"ru/sbt/estima/model"
	"github.com/gorilla/mux"
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

func (us *UserStoryService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle("/api/v.0.0.1/feature/{id}/userstory/list", handler(http.HandlerFunc(us.listUserStory))).Methods("POST", "GET").Name("User story list for specified feature")
	router.Handle("/api/v.0.0.1/feature/{id}/userstory/add", handler(http.HandlerFunc(us.addUserStory))).Methods("POST").Name("Add user story to feature")
	router.Handle("/api/v.0.0.1/userstory/{id}/remove", handler(http.HandlerFunc(us.removeUserStory))).Methods("POST", "DELETE").Name("Remove user story")
}