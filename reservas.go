package main

import (
	"fmt"
	//"github.com/gorilla/mux"
	//"os"
	"encoding/json"
	//"github.com/gorilla/handlers"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"io/ioutil"	
	//"strconv"
	"strings"
	"net/http"
	firebase "github.com/wuman/firebase-server-sdk-go"
	"github.com/gorilla/context"
)

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {

	auth, _ := firebase.GetAuth()
	token, err := auth.CreateCustomToken("FJNWvK2wbrhA2XHhWSQuiLVVFHp2", nil)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: token})
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			authorizationHeader := req.Header.Get("authorization")
			if authorizationHeader != "" {
				bearerToken := strings.Split(authorizationHeader, " ")
				if len(bearerToken) == 2 {
					auth, _ := firebase.GetAuth()
					decodedToken, err := auth.VerifyIDToken(bearerToken[1])
					// uid, found := decodedToken.UID()
					// println("uid", uid)
					// println("found", found)
					if err == nil {
						uid, found := decodedToken.UID()
						println("uid", uid)
						println("found", found)
						// context.Set(req, "decoded", uid)
						context.Set(req, "uid", uid)
						next(w, req)
					} else {
						fmt.Println(err)
						json.NewEncoder(w).Encode(Exception{Message: "Invalid token"})
					}
				}
			} else {
				json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
			}
		})
	}
	
	func TestEndpoint(w http.ResponseWriter, req *http.Request) {
		decoded := context.Get(req, "decoded")
		// // var user User
		// mapstructure.Decode(decoded.(jwt.MapClaims), &user)
		// json.NewEncoder(w).Encode(user)
		json.NewEncoder(w).Encode(decoded)
	}

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