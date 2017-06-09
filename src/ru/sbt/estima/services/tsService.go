package services

import (
	"ru/sbt/estima/model"
	"net/http"
	"github.com/gorilla/mux"
)

type TechStoryService struct {
}

func (ts *TechStoryService) addTechStory (w http.ResponseWriter, r *http.Request) {
	var story model.UserStory
	story.Key = mux.Vars(r)["id"]

	var techStory model.TechStory
	model.ReadJsonBody(r, &techStory)

	WithTechStoryDao(func(dao techStoryDao) {
		techStoryId := dao.createAndConnectObjTx(techStory, story, model.PRJ_EDGES, map[string]string{"label": "techStory"})
		techStory.Id = techStoryId
		model.WriteResponse(true, nil, techStory, w)
	})
}

func (ts *TechStoryService) removeTechStory (w http.ResponseWriter, r *http.Request) {
	var tss model.TechStory
	tss.Key = mux.Vars(r)["id"]

	WithTechStoryDao(func(dao techStoryDao) {
		dao.removeConnectedTx (tss.GetCollection(), model.PRJ_EDGES, tss.GetKey())
		model.WriteResponse(true, nil, tss, w)
	})
}

func (ts *TechStoryService) listTechStory (w http.ResponseWriter, r *http.Request) {
	var story model.UserStory
	story.Key = mux.Vars(r)["id"]

	WithTechStoryDao(func(dao techStoryDao) {
		uss, err := dao.FindByUS(story.GetCollection() + "/" + story.Key)
		model.CheckErr(err)
		model.WriteArrayResponse(true, nil, uss, w)
	})
}

// Service get text for given feature and, optional, text version
// If the version is not specified then service returns currently active text
func (ts *TechStoryService) getText (w http.ResponseWriter, r *http.Request) {
	var techStory model.TechStory
	techStory.Key = mux.Vars(r)["id"]

	// Check for version parameter
	values := r.URL.Query()
	version := model.GetInt (values, "version", -1)
	var found *model.VersionedText
	var err error

	WithTechStoryDao(func(dao techStoryDao) {
		if version != -1 {
			found, err = dao.GetTextByVersion(techStory, version)
		} else {
			found, err = dao.GetActiveText(techStory)
		}
		model.CheckErr(err)
		model.WriteResponse(true, nil, found, w)
	})
}

// Service find and retrieve feature by id (key)
func (ts *TechStoryService) findById (w http.ResponseWriter, r *http.Request) {
	var techStory model.TechStory
	techStory.Key = mux.Vars(r)["id"]
	WithTechStoryDao(func(dao techStoryDao) {
		err := dao.FindById(&techStory)
		model.CheckErr(err)
		model.WriteResponse(true, nil, techStory, w)
	})
}

// Service add text to feature
func (ts *TechStoryService) addText (w http.ResponseWriter, r *http.Request) {
	var techStory model.TechStory
	techStory.Key = mux.Vars(r)["id"]

	var text model.VersionedText
	model.ReadJsonBody(r, &text)
	WithTechStoryDao(func(dao techStoryDao) {
		vtext := dao.AddText(techStory, text.Text)
		model.WriteResponse(true, nil, vtext, w)
	})
}

// Service add text to feature
func (ts *TechStoryService) setStatus (w http.ResponseWriter, r *http.Request) {
	var techStory model.TechStory
	model.ReadJsonBody(r, &techStory)
	techStory.Key = mux.Vars(r)["id"]

	WithTechStoryDao(func(dao techStoryDao) {
		user := model.GetUserFromRequest(w, r)
		entity, err := dao.SetStatus(&techStory, techStory.Status, user)
		model.CheckErr(err)
		model.WriteResponse(true, nil, entity, w)
	})
}

// Service retrieves comments for the feature. All the comments related for the whole feature not for feature text.
// This is done by this way because text can has many versions and comments may be lost.
// For paging results using next query (url) parameters
// offset - start offset for results
// pageSize - number of results
func (ts *TechStoryService) getComments (w http.ResponseWriter, r *http.Request) {
	var techStory model.TechStory
	techStory.Key = mux.Vars(r)["id"]

	values := r.URL.Query()
	offset := model.GetInt(values, "offset", 0)
	pageSize := model.GetInt(values, "pageSize", 0)

	WithTechStoryDao(func(dao techStoryDao) {
		comments, err := dao.GetComments(techStory, pageSize, offset)
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
func (ts *TechStoryService) addComment (w http.ResponseWriter, r *http.Request) {
	var techStory model.TechStory
	techStory.Key = mux.Vars(r)["id"]

	var com model.Comment
	model.ReadJsonBody(r, &com)

	WithTechStoryDao(func(dao techStoryDao) {
		user := model.GetUserFromRequest(w, r)
		comment := dao.AddComment(techStory, com.Title, com.Text, user.GetCollection()+"/"+user.Key)

		model.WriteResponse(true, nil, comment, w)
	})
}

// Service retrieve all text versions for the given feature
func (ts TechStoryService) getTextVersionList(w http.ResponseWriter, r *http.Request) {
	var techStory model.TechStory
	techStory.Key = mux.Vars(r)["id"]

	WithTechStoryDao(func(dao techStoryDao) {
		versions, err := dao.GetTextVersionList(techStory)
		model.CheckErr(err)
		model.WriteArrayResponse(true, nil, versions, w)
	})
}

func (ts TechStoryService) updateTechStory (w http.ResponseWriter, r *http.Request) {
	var techStory model.TechStory
	model.ReadJsonBody(r, &techStory)

	WithTechStoryDao(func(dao techStoryDao) {
		techStory.Key = mux.Vars(r)["id"]
		entity, err := dao.Save(&techStory)
		model.CheckErr(err)
		model.WriteResponse(true, nil, entity, w)
	})
}

func (ts *TechStoryService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle("/api/v.0.0.1/userstory/{id}/tstory/list", handler(http.HandlerFunc(ts.listTechStory))).Methods("POST", "GET").Name("Tech story list for specified US")
	router.Handle("/api/v.0.0.1/userstory/{id}/tstory/add", handler(http.HandlerFunc(ts.addTechStory))).Methods("POST").Name("Add tech story to US")
	router.Handle("/api/v.0.0.1/tstory/{id}/update", handler(http.HandlerFunc(ts.updateTechStory))).Methods("POST").Name("Update tech story")
	router.Handle("/api/v.0.0.1/tstory/{id}/remove", handler(http.HandlerFunc(ts.removeTechStory))).Methods("POST", "DELETE").Name("Remove tech story")
	router.Handle("/api/v.0.0.1/tstory/{id}/get", handler(http.HandlerFunc(ts.findById))).Methods("POST", "GET").Name("Get tech story by key")
	router.Handle("/api/v.0.0.1/tstory/{id}/status", handler(http.HandlerFunc(ts.setStatus))).Methods("POST").Name("Change status of tech story")
	router.Handle("/api/v.0.0.1/tstory/{id}/addcomment", handler(http.HandlerFunc(ts.addComment))).Methods("POST").Name("Add comment to tech story")
	router.Handle("/api/v.0.0.1/tstory/{id}/comment", handler(http.HandlerFunc(ts.getComments))).Methods("POST", "GET").Name("Get comments of tech story")
	router.Handle("/api/v.0.0.1/tstory/{id}/addtext", handler(http.HandlerFunc(ts.addText))).Methods("POST").Name("Add text to tech story")
	router.Handle("/api/v.0.0.1/tstory/{id}/text", handler(http.HandlerFunc(ts.getText))).Methods("POST", "GET").Name("Get active text for tech story")
	router.Handle("/api/v.0.0.1/tstory/{id}/text/list", handler(http.HandlerFunc(ts.getTextVersionList))).Methods("POST", "GET").Name("Get all text versions for tech story")
}
