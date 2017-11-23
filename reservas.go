package main

import (
	//"fmt"
	//"github.com/gorilla/mux"
	//"os"
	"encoding/json"
	//"github.com/gorilla/handlers"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"io/ioutil"	
	//"strconv"
	//"strings"
	"net/http"
)

func getReservations(w http.ResponseWriter, r *http.Request){

	session, err := mgo.Dial("mongodb://udeain:udeainmongodb@ds157444.mlab.com:57444/heroku_4r2js6cs")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("heroku_4r2js6cs").C("reservation")

	var reservasObj []bson.M

	err = c.Find(nil).All(&reservasObj)

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("not found"))
		return
	}

	respuesta, err := json.Marshal(reservasObj)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("unable to get reservations"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respuesta)
}