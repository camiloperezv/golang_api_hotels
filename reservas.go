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
	"github.com/gorilla/context"
	"net/http"
	"strings"
)

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

func TestEndpoint(w http.ResponseWriter, req *http.Request) {
	decoded := context.Get(req, "decoded")
	// // var user User
	// mapstructure.Decode(decoded.(jwt.MapClaims), &user)
	// json.NewEncoder(w).Encode(user)
	json.NewEncoder(w).Encode(decoded)
}

func deleteReservation(w http.ResponseWriter, r *http.Request, id_reserva string){

	session, err := mgo.Dial("mongodb://udeain:udeainmongodb@ds157444.mlab.com:57444/heroku_4r2js6cs")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Configurar sesión Mongo DB
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("heroku_4r2js6cs").C("reservation")

	var reservaObj bson.M

	// Buscar reserva recibida
	err = c.Find(bson.M{"reserve_id": id_reserva}).One(&reservaObj)

	reserva, err := json.Marshal(reservaObj)
	
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	if (reserva != nil && string(reserva) != "null" ){
		//w.Write(reserva)
		// Desactivar reserva

		c.Upsert(
			bson.M{"reserve_id": id_reserva},
			bson.M{"$set": bson.M{"state": "C"}},
		)

		mensaje := "La reserva con identificador " + id_reserva + " fue cancelada exitosamente!!"
		w.Write([]byte(`{"message" : "` + mensaje + `"}`))

	}else{
		mensaje := "La reserva con identificador " + id_reserva + " no se encuentra en nuestros registros, intente de nuevo"
		w.Write([]byte(`{"message" : "` + mensaje + `"}`))
	}

}

func getReservations(w http.ResponseWriter, r *http.Request) {

	// obtener parámetro de reserva
	id_reserva := r.URL.Query().Get("reserve_id")
	//fmt.Println("Reserva "+ id_reserva)
	if (id_reserva != ""){
		deleteReservation(w , r , id_reserva)			
		return	
	}

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

	//respuesta, err := json.Marshal(resp)
	respuesta, err := json.Marshal(reservasObj)
	longitud := len(reservasObj)	

	// Decodificar datos de JSON
	var raw []map[string]interface{}
	json.Unmarshal(respuesta, &raw)
		
	for i := 0; i < longitud; i++ {
		// agregar reserva a cada elemento
		raw[i]["reservation"] = reservasObj[i]

		// obtener id de hotel
		salida, _ := json.Marshal(raw[i]["hotel_id"])
		hotel := string(salida)
		hotel = strings.Replace(hotel, "\"", "", -1)
		fmt.Println( hotel )

		// agregar datos adicionales	
		if (hotel == "udeain_medellin"){
			raw[i]["hotel_id"] = hotel
			raw[i]["hotel_name"] = "udeain medellin"
			raw[i]["hotel_thumbnail"] = "http://www.kohler.com/common/images/global_accounts/Hotel-Indigo-Thumbnail.jpg"
			raw[i]["hotel_location"] = bson.M{"address":"Cl. 5 Sur #42-2 a 42-70", "lat":"6.1992463", "long":"-75.5747155"}
			raw[i]["check_in"] = "15:00"
			raw[i]["check_out"] = "13:00"
			raw[i]["hotel_website"] = "https://udeain.herokuapp.com"			

		}else if (hotel == "udeain_bogota"){
			raw[i]["hotel_id"] = hotel
			raw[i]["hotel_name"] = "udeain bogota"
			raw[i]["hotel_thumbnail"] = "http://images.citybreakcdn.com/image.aspx?ImageId=329020&width=300&height=300&crop=1"
			raw[i]["hotel_location"] = bson.M{"address":"Cra. 14 #82-2 a 82-98", "lat":"4.667662", "long":"-74.0574518"}
			raw[i]["check_in"] = "15:30"
			raw[i]["check_out"] = "13:30"
			raw[i]["hotel_website"] = "https://udeain.herokuapp.com"
		}
			
		// eliminar datos no requeridos
		delete(raw[i], "_id");
		delete(raw[i], "doc_id");
		delete(raw[i], "arrive_date");
		delete(raw[i], "doc_type");
		delete(raw[i], "email");
		delete(raw[i], "host_id");
		delete(raw[i], "leave_date");
		delete(raw[i], "reserve_id");
		delete(raw[i], "room_id");
		delete(raw[i], "state");
		delete(raw[i], "userId");		
	}
	
	
	// actualizar datos de respuesta
	respuesta, err = json.Marshal(raw)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respuesta)
}