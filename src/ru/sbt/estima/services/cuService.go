package services

import (
	"ru/sbt/estima/model"
	"net/http"
	"github.com/gorilla/mux"
	ara "github.com/diegogub/aranGO"
	"time"
)

type CuService struct {}

func (cs CuService) listCalcUnit (w http.ResponseWriter, r *http.Request) {
	WithCuDao(func(dao cuDao) {
		values := r.URL.Query()
		offset := model.GetInt(values, "offset", 0)
		pageSize := model.GetInt(values, "pageSize", 0)
		cus := dao.FindAllUnit(offset, pageSize)
		model.WriteArrayResponse(true, nil, cus, w)
	})
}

func (cs CuService) listConnectedCalcUnit (w http.ResponseWriter, r *http.Request) {
	WithCuDao(func(dao cuDao) {
		entity := model.TechStory{Document: ara.Document{Key: mux.Vars(r)["id"]}}
		cus := dao.findConnectedCalcUnit(entity.GetCollection() + "/" + entity.GetKey())
		model.WriteArrayResponse(true, nil, cus, w)
	})
}

func (cs CuService) addCalcUnit (w http.ResponseWriter, r *http.Request) {
	var unit model.CalcUnit
	model.ReadJsonBody(r, &unit)
	unit.Changed = time.Now()

	WithCuDao(func(dao cuDao) {
		cmp, err := dao.Save(&unit)
		model.CheckErr(err)
		model.WriteResponse(true, nil, cmp, w)
	})
}

func (cs CuService) updateCalcUnit (w http.ResponseWriter, r *http.Request) {
	var unit model.CalcUnit
	model.ReadJsonBody(r, &unit)
	unit.Key = mux.Vars(r)["id"]
	unit.Changed = time.Now()

	WithCuDao(func(dao cuDao) {
		cmp, err := dao.Save(&unit)
		model.CheckErr(err)
		model.WriteResponse(true, nil, cmp, w)
	})
}

func (cs CuService) removeCalcUnit (w http.ResponseWriter, r *http.Request) {
	var cmp model.CalcUnit
	model.ReadJsonBody(r, &cmp)
	cmp.Key = mux.Vars(r)["id"]

	WithCuDao(func(dao cuDao) {
		err := dao.Database().Col(cmp.GetCollection()).Delete(cmp.GetKey())
		model.CheckErr(err)
		model.WriteResponse(true, nil, cmp, w)
	})
}

func (us CuService) connectCalcUnit (w http.ResponseWriter, r *http.Request) {
	entity := model.TechStory{Document: ara.Document{Key: mux.Vars(r)["id"]}}
	entityId := entity.GetCollection() + "/" + entity.GetKey()

	var unitEdge model.CalcUnitEdge
	model.ReadJsonBody(r, &unitEdge)
	unitEdge.Id = unitEdge.GetCollection() + "/" + unitEdge.GetKey()
	attrs := map[string]interface{} {
		"weight": unitEdge.Weight,
		"complexity": unitEdge.Complexity,
		"newFlag": unitEdge.NewFlag,
		"extCoef": unitEdge.ExtCoef,
	}

	WithCuDao(func(dao cuDao) {
		err := dao.connectCalcUnit(entityId, unitEdge.CalcUnit, attrs)
		model.CheckErr(err)
		model.WriteResponse(true, nil, unitEdge, w)
	})
}

func (us CuService) disconnectCalcUnit (w http.ResponseWriter, r *http.Request) {
	entity := model.TechStory{Document: ara.Document{Key: mux.Vars(r)["id"]}}
	entityId := entity.GetCollection() + "/" + entity.GetKey()

	calcUnit := model.CalcUnit{Document: ara.Document{Key: mux.Vars(r)["cuId"]}}
	calcUnit.Id = calcUnit.GetCollection() + "/" + calcUnit.GetKey()

	WithCuDao(func(dao cuDao) {
		dao.disconnectCalcUnit(entityId, calcUnit)
		model.WriteResponse(true, nil, calcUnit, w)
	})
}

func (ca CuService) updateCalcUnitEdge(w http.ResponseWriter, r *http.Request) {
	entity := model.TechStory{Document: ara.Document{Key: mux.Vars(r)["id"]}}
	entityId := entity.GetCollection() + "/" + entity.GetKey()

	var unitEdge model.CalcUnitEdge
	model.ReadJsonBody(r, &unitEdge)
	unitEdge.Id = unitEdge.GetCollection() + "/" + unitEdge.GetKey()

	WithCuDao(func(dao cuDao) {
		cu := dao.updateCalcUnitEdge(entityId, unitEdge)
		model.WriteResponse(true, nil, cu, w)
	})
}

func (ca CuService) findConnectedCalcUnits(w http.ResponseWriter, r *http.Request) {
	entity := model.TechStory{Document: ara.Document{Key: mux.Vars(r)["id"]}}
	entityId := entity.GetCollection() + "/" + entity.GetKey()

	WithCuDao(func(dao cuDao) {
		cmp := dao.findConnectedCalcUnit(entityId)
		model.WriteArrayResponse(true, nil, cmp, w)
	})
}

func (cs CuService) addPrice (w http.ResponseWriter, r *http.Request) {
	var unit model.CalcUnit
	unit.Key = mux.Vars(r)["id"]

	var price model.CalcUnitPrice
	model.ReadJsonBody(r, &price)
	price.Changed = time.Now()

	WithCuDao(func(dao cuDao) {
		price = dao.addPrice(unit, price)
		model.WriteResponse(true, nil, price, w)
	})
}

func (cs CuService) removePrice (w http.ResponseWriter, r *http.Request) {
	var price model.CalcUnitPrice
	price.Key = mux.Vars(r)["id"]

	WithCuDao(func(dao cuDao) {
		dao.removePrice(price)
		model.WriteResponse(true, nil, nil, w)
	})
}

func (cs CuService) updatePrice (w http.ResponseWriter, r *http.Request) {
	var unitPrice model.CalcUnitPrice
	model.ReadJsonBody(r, &unitPrice)
	unitPrice.Key = mux.Vars(r)["id"]
	unitPrice.Changed = time.Now()

	WithCuDao(func(dao cuDao) {
		cmp, err := dao.Save(&unitPrice)
		model.CheckErr(err)
		model.WriteResponse(true, nil, cmp, w)
	})
}

func (ca CuService) findConnectedPrice(w http.ResponseWriter, r *http.Request) {
	var unit model.CalcUnit
	unit.Key = mux.Vars(r)["id"]
	unit.Id = unit.GetCollection() + "/" + unit.GetKey()

	WithCuDao(func(dao cuDao) {
		cmp := dao.findConnectedPrice(unit)
		model.WriteArrayResponse(true, nil, cmp, w)
	})
}

func (ca CuService) calculateProjectPrice(w http.ResponseWriter, r *http.Request) {
	projectKey := mux.Vars(r)["id"]

	WithCuDao(func(dao cuDao) {
		res := dao.calculateProjectPricePlain(projectKey)
		model.WriteAnyResponse(true, nil, res, w)
	})
}

func (us *CuService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle("/api/v.0.0.1/cu/list", handler(http.HandlerFunc(us.listCalcUnit))).Methods("POST", "GET").Name("Calc units list")
	router.Handle("/api/v.0.0.1/cu/add", handler(http.HandlerFunc(us.addCalcUnit))).Methods("POST").Name("Add calc unit")
	router.Handle("/api/v.0.0.1/ts/{id}/listCu", handler(http.HandlerFunc(us.findConnectedCalcUnits))).Methods("POST", "GET").Name("List calc units connected to tech story")
	router.Handle("/api/v.0.0.1/ts/{id}/addCu", handler(http.HandlerFunc(us.connectCalcUnit))).Methods("POST").Name("Connect calc unit to tech story")
	router.Handle("/api/v.0.0.1/ts/{id}/cu/{cuId}/remove", handler(http.HandlerFunc(us.disconnectCalcUnit))).Methods("POST", "DELETE").Name("Disconnect calc unit from tech story")
	router.Handle("/api/v.0.0.1/ts/{id}/cu/{cuId}/update", handler(http.HandlerFunc(us.updateCalcUnitEdge))).Methods("POST", "DELETE").Name("Disconnect calc unit from tech story")
	router.Handle("/api/v.0.0.1/cu/{id}/update", handler(http.HandlerFunc(us.updateCalcUnit))).Methods("POST").Name("Update calc unit")
	router.Handle("/api/v.0.0.1/cu/{id}/remove", handler(http.HandlerFunc(us.removeCalcUnit))).Methods("POST", "DELETE").Name("Remove calc unit")
	router.Handle("/api/v.0.0.1/cu/{id}/addPrice", handler(http.HandlerFunc(us.addPrice))).Methods("POST").Name("Attach price to calc unit")
	router.Handle("/api/v.0.0.1/cu/{id}/listPrice", handler(http.HandlerFunc(us.findConnectedPrice))).Methods("POST", "GET").Name("List connected to calc unit prices")
	router.Handle("/api/v.0.0.1/cuprice/{id}/remove", handler(http.HandlerFunc(us.removePrice))).Methods("POST", "DELETE").Name("Remove calc unit price")
	router.Handle("/api/v.0.0.1/cuprice/{id}/update", handler(http.HandlerFunc(us.updatePrice))).Methods("POST").Name("Update calc unit price")
	router.Handle("/api/v.0.0.1/calcproject/{id}", handler(http.HandlerFunc(us.calculateProjectPrice))).Methods("POST", "GET").Name("Calculate projects price")
}