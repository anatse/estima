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

	var unitWithWeight struct {
		model.CalcUnit
		Weight float64
	}

	model.ReadJsonBody(r, &unitWithWeight)
	unitWithWeight.Id = unitWithWeight.GetCollection() + "/" + unitWithWeight.GetKey()

	WithCuDao(func(dao cuDao) {
		err := dao.connectCalcUnit(entityId, unitWithWeight.CalcUnit, unitWithWeight.Weight)
		model.CheckErr(err)
		model.WriteResponse(true, nil, unitWithWeight, w)
	})
}

func (us CuService) disconnectCalcUnit (w http.ResponseWriter, r *http.Request) {
	entity := model.TechStory{Document: ara.Document{Key: mux.Vars(r)["id"]}}
	entityId := entity.GetCollection() + "/" + entity.GetKey()

	var calcUnit model.CalcUnit
	model.ReadJsonBody(r, &calcUnit)
	calcUnit.Id = calcUnit.GetCollection() + "/" + calcUnit.GetKey()

	WithCuDao(func(dao cuDao) {
		dao.disconnectCalcUnit(entityId, calcUnit)
		model.WriteResponse(true, nil, calcUnit, w)
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

func (us *CuService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle("/api/v.0.0.1/cu/list", handler(http.HandlerFunc(us.listCalcUnit))).Methods("POST", "GET").Name("Calc units list")
	router.Handle("/api/v.0.0.1/cu/add", handler(http.HandlerFunc(us.addCalcUnit))).Methods("POST").Name("Add calc unit")
	router.Handle("/api/v.0.0.1/cu/{id}/update", handler(http.HandlerFunc(us.updateCalcUnit))).Methods("POST").Name("Update cal;c unit")
	router.Handle("/api/v.0.0.1/cu/{id}/remove", handler(http.HandlerFunc(us.removeCalcUnit))).Methods("POST", "DELETE").Name("Remove calc unit")
	router.Handle("/api/v.0.0.1/cu/{id}/addPrice", handler(http.HandlerFunc(us.addPrice))).Methods("POST").Name("Attach price to calc unit")
	router.Handle("/api/v.0.0.1/cu/{id}/listPrice", handler(http.HandlerFunc(us.findConnectedPrice))).Methods("POST", "GET").Name("List connected to calc unit prices")
	router.Handle("/api/v.0.0.1/cuprice/{id}/remove", handler(http.HandlerFunc(us.removePrice))).Methods("POST", "DELETE").Name("Remove calc unit price")
	router.Handle("/api/v.0.0.1/cuprice/{id}/update", handler(http.HandlerFunc(us.updatePrice))).Methods("POST").Name("Update calc unit price")
}