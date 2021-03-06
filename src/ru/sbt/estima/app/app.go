package app

import (
	"net/http"
	"ru/sbt/estima/model"
	"ru/sbt/estima/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"reflect"
	"unsafe"
	"encoding/json"
	"github.com/gorilla/handlers"
	"fmt"
	"compress/gzip"
	"ru/sbt/estima/conf"
)

func JwtHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer (func() {
			if r := recover(); r != nil {
				araErr := model.GetAraError (r)
				if araErr != nil {
					ae := araErr.(model.AraError)
					model.WriteResponse(false, fmt.Sprintf("%s", model.GetErrorText(ae)), nil, w)
				} else {
					model.WriteResponse(false, fmt.Sprint(r), nil, w)
				}
			}
		})()

		// Let secure process the request. If it returns an error,
		// that indicates the request should not continue.
		err := services.JwtMiddleware.CheckJWT(w, r)

		// If there was an error, do not continue.
		if err != nil {
			model.WriteResponse(false, fmt.Sprintf("Forbidden: %v", err), nil, w)
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func WebSocketHandler (h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}
		channel, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			model.CheckErr(err)
		}

		defer channel.Close()
		for {
			mt, message, err := channel.ReadMessage()
			if err != nil {
				conf.GetLog().Println("read:", err)
				break
			}
			conf.GetLog().Printf("recv: %s", message)
			err = channel.WriteMessage(mt, message)
			if err != nil {
				conf.GetLog().Println("write:", err)
				break
			}
		}
	})
}

// Router information struct
type routeInfo struct {
	Name string
	Path string
}

// Fill router information
func getRoutes (router *mux.Router) []routeInfo {
	v := reflect.ValueOf(router)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	routes := v.FieldByName("routes")
	len := routes.Len()

	ris := make ([]routeInfo, len)

	for i:=0;i<len;i++ {
		route := routes.Index(i)
		rp := (*mux.Route)(unsafe.Pointer (route.Pointer()))

		path, _ := rp.GetPathTemplate()
		rp.GetName()
		ris[i] = routeInfo{
			rp.GetName(),
			path,
		}
	}

	return ris
}

func instService (w http.ResponseWriter, req *http.Request) {
	collections := []string {
		"users",
		"projects",
		"stages",
		"processes",
		"features",
		"ustories",
		"verstext",
		"comments",
		"components",
		"cuprices",
		"calcunits",
		"tstories",
		"tsprices",
	}

	services.GetPool().Use(func(iDao interface{}) {
		dao := iDao.(*services.BaseDao)
		for _, col := range collections {
			dao.Col(col)
		}

		// create edge collection
		dao.EdgeCol(model.PRJ_EDGES)

		// Create unique index for edge collection
		dao.Database().Col(model.PRJ_EDGES).CreatePersistent(true, "_from", "_to", "label");

		// Write success flag
		model.WriteAnyResponse(true, nil, nil, w)
	})
}

func nextStatuses (w http.ResponseWriter, r *http.Request) {
	user := model.GetUserFromRequest (w, r)
	values := r.URL.Query()
	status := model.GetStatus(values, "status", model.STATUS_NEW)

	var ret []model.Status
	curStatus := model.FromStatus(status)
	nextStatuses := curStatus.NextStatuses()

	for _, nextStatus := range nextStatuses {
		if curStatus.CanMoveTo(nextStatus, user.Roles) {
			ret = append(ret, nextStatus)
		}
	}

	model.WriteAnyResponse(true, nil, ret, w)
}

func PrepareRoute () *mux.Router {
	var us services.UserService
	var ps services.ProjectService
	var pcs services.ProcessService
	var fs services.FeatureService
	var uss services.UserStoryService
	var tss services.TechStoryService
	var cs services.ComponentService
	var cus services.CuService

	model.RegisterService("user", us)
	model.RegisterService("project", ps)
	model.RegisterService("process", pcs)
	model.RegisterService("feature", fs)
	model.RegisterService("userStory", uss)
	model.RegisterService("techStory", tss)
	model.RegisterService("component", cs)
	model.RegisterService("calculate", cus)

	r := model.GetRouter()
	r.Handle("/api/v.0.0.1/get-token", services.GetTokenHandler).Methods("GET").Name("Login router (GET). Query parameters uname & upass. This router is deprecated and will be removed in next release")
	r.Handle("/api/v.0.0.1/login", services.Login).Methods("POST").Name("Login router (POST). Body: uname & upass")
	r.Handle("/api/v.0.0.1/init", http.HandlerFunc(instService)).Methods("POST", "GET").Name("Create database collections")
	r.Handle("/api/v.0.0.1/nextStatuses", JwtHandler(http.HandlerFunc(nextStatuses))).Methods("POST", "GET").Name("Get next statuses for current status and user")

	r.Handle("/echo", WebSocketHandler(nil))

	us.ConfigRoutes(r, JwtHandler)
	ps.ConfigRoutes(r, JwtHandler)
	pcs.ConfigRoutes(r, JwtHandler)
	fs.ConfigRoutes(r, JwtHandler)
	uss.ConfigRoutes(r, JwtHandler)
	tss.ConfigRoutes(r, JwtHandler)
	cs.ConfigRoutes(r, JwtHandler)
	cus.ConfigRoutes(r, JwtHandler)

	// Function build router for get router information
	var routesInformation = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		js, _ := json.Marshal(getRoutes(r))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(js))
	})

	r.Handle("/api/v.0.0.1/ri", routesInformation)

	// Add static router. Should be last in routes list
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./views")))
	return r
}

func AppRun () {
	conf.GetLog().Printf("Starting application...")

	// Init LDAP
	model.InitLdapPool(2)

	r := PrepareRoute()
	addr := fmt.Sprintf(":%d", conf.LoadConfig().Port)

	//err := http.ListenAndServeTLS(":9443", "server.crt", "server.key", handlers.CompressHandler(handlers.LoggingHandler(os.Stdout, r)))
	err := http.ListenAndServe(addr, handlers.CompressHandlerLevel(handlers.LoggingHandler(conf.GetLogFile(), r), gzip.BestCompression))
	if err != nil {
		conf.GetLog().Fatal("ListenAndServe: ", err)
	}

	model.FinishLdapPool()
}

