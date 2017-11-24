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

	////
	/*collection := session.DB("heroku_4r2js6cs").C("reservation")
	pipeline := []bson.M{
		//bson.M{"$match": bson.M{"$or": []bson.M{bson.M{"start_date": bson.M{"$gte": fecha_fin}}, bson.M{"end_date": bson.M{"$lte": fecha_inicio}}}}},

		//Realizar 'Join' con documentos adicionales de hotel y datos de habitaciones//
		//bson.M{"$lookup": bson.M{"from": "rooms", "localField": "room_id", "foreignField": "id", "as": "rooms"}},
		bson.M{"$lookup": bson.M{"from": "hotel", "localField": "hotel_id", "foreignField": "hotel_id", "as": "hotel_details"}},
		// Realizar filtrado por tipo de habitación y ciudad //
		//{"$unwind": "$rooms"},
	}
	pipe := collection.Pipe(pipeline)
	resp := []bson.M{}
	err = pipe.All(&resp)*/
	////

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("not found"))
		return
	}

	//respuesta, err := json.Marshal(resp)
	respuesta, err := json.Marshal(reservasObj)
	longitud := len(reservasObj)
	fmt.Println(longitud)

	headerJsonPadre := []byte(`{"reservations": [`)
	///
	var headerJson []byte
	var cityInfo []byte
	var jsonEnd []byte
	var finalRes []byte

	for i := 0; i < longitud; i++ {
		fmt.Println(i)
		headerJson = []byte(`{`)
			
		cityInfo = []byte(`"hotel_id":"udeain_medellin","hotel_name":"udeain Medellín", "hotel_location":{"address":"Cra. 14 #82-2 a 82-98", "lat":"4.667662", "long":"-74.0574518"},"hotel_thumbnail":"https://media-cdn.tripadvisor.com/media/photo-s/06/35/93/c2/hotel-el-deportista.jpg","check_in":"15:00","check_out":"13:00","hotel_website":"https://udeain.herokuapp.com", "reservation":`)
		headerJson = append(headerJson[:], cityInfo...)

		jsonEnd = []byte(`}`)
		if string(respuesta) == "null" {
			respuesta = []byte(`[]`)
		}
		finalRes = append(headerJson[:], respuesta...)
		//finalRes := append(headerJson[:], datos[0]...)
		finalRes = append(finalRes[:], jsonEnd...)
	}
	///

	finalResPadre := append(headerJsonPadre[:], finalRes...)
	jsonEndPadre := []byte(`]}`)
	finalResPadre = append(finalResPadre[:], jsonEndPadre...)

	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("unable to get reservations"))
		return
	}

	// accesar a datos por índice
	/*var datos []bson.M
	err = json.Unmarshal(respuesta, &datos)
	fmt.Println("DATOS: ", datos[0])

	if err != nil {
		fmt.Println("error:", err)
	}*/

	///////////
	/*headerJson := []byte(`{`)
	
	cityInfo := []byte(`"reservations": [{ "hotel_id":"udeain_medellin","hotel_name":"udeain Medellín", "hotel_location":{"address":"Cra. 14 #82-2 a 82-98", "lat":"4.667662", "long":"-74.0574518"},"hotel_thumbnail":"https://media-cdn.tripadvisor.com/media/photo-s/06/35/93/c2/hotel-el-deportista.jpg","check_in":"15:00","check_out":"13:00","hotel_website":"https://udeain.herokuapp.com", "reservation":`)
	headerJson = append(headerJson[:], cityInfo...)

	jsonEnd := []byte(`}]}`)
	if string(respuesta) == "null" {
		respuesta = []byte(`[]`)
	}
	finalRes := append(headerJson[:], respuesta...)
	//finalRes := append(headerJson[:], datos[0]...)
	finalRes = append(finalRes[:], jsonEnd...)*/
	///////////

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	//w.Write(respuesta)
	//w.Write(finalRes)
	//w.Write(finalRes)
	w.Write(finalResPadre)
}