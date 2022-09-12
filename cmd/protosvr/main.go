package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/julienschmidt/httprouter"
	"jcc.dev/protosvr/internal"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Test struct {
	Id   int
	Name string
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func HandleError(err error, w http.ResponseWriter) {
	log.Printf("\n%v", internal.PrettyPrintError(internal.WrapErrorN(err, 2)))

	// todo handle other error types
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Internal Server Error.")
}

func ListEntities(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	values, err := internal.GetEntityMetaList()
	if err != nil {
		HandleError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(values)
	if err != nil {
		HandleError(err, w)
		return
	}
}

func EntityPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var values []*Test
	err := pgxscan.Select(context.Background(), internal.Db, &values, "select * from test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(w, "EntityPost, %s\n", params.ByName("entityName"))
}

func EntityUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "EntityUpdate, %s %s\n", params.ByName("entityName"), params.ByName("id"))
}

func EntityGetAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	entityName, err := internal.GetCanonicalEntityName(params.ByName("entityName"))
	if err != nil {
		HandleError(err, w)
		return
	}

	queryParams := r.URL.Query()
	takeStr := queryParams.Get("take")
	var take int64
	if len(strings.Trim(takeStr, " ")) > 0 {
		take, err = strconv.ParseInt(takeStr, 10, 64)
		if err != nil {
			HandleError(err, w)
			return
		}
	}

	skipStr := queryParams.Get("skip")
	var skip int64
	if len(strings.Trim(skipStr, " ")) > 0 {
		skip, err = strconv.ParseInt(skipStr, 10, 64)
		if err != nil {
			HandleError(err, w)
			return
		}
	}

	options := internal.GetAllOptions{
		Take: take,
		Skip: skip,
	}

	res, err := internal.GetAllEntities(entityName, &options)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(err, w)
		return
	}
}

func EntityGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	entityName, err := internal.GetCanonicalEntityName(params.ByName("entityName"))
	if err != nil {
		HandleError(err, w)
		return
	}

	idStr := params.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		HandleError(err, w)
		return
	}

	res, err := internal.GetSingleEntity(entityName, id)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(err, w)
		return
	}
}

func main() {
	internal.SetupDb()
	defer internal.Db.Close()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/_entities", ListEntities)
	router.POST("/api/e/:entityName", EntityPost)
	router.PUT("/api/e/:entityName/:id", EntityUpdate)
	router.GET("/api/e/:entityName", EntityGetAll)
	router.GET("/api/e/:entityName/:id", EntityGet)

	log.Fatal(http.ListenAndServe(":8080", router))
}
