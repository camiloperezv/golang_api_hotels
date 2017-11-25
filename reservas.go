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
	// "strings"
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

		//c.Update({"reserve_id": id_reserva },{"$set" : {"state": "C"}})

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
	respuesta, err := json.Marshal(reservasObj[0])
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
