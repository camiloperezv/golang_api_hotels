package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"os"
	// "log"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	// // Firebase sdk
	// firebase "firebase.google.com/go"
	// context "golang.org/x/net/context"
	// "google.golang.org/api/option"

	firebase "github.com/wuman/firebase-server-sdk-go"
)

// var opt = option.WithCredentialsFile("firebaseAdminCredentials.json")
// var app, err = firebase.NewApp(context2.Background(), nil, opt)

// // if err != nil {
// //   return nil, fmt.Errorf("error initializing app: %v", err)
// // }

type Room struct {
	Hotel_id        string `json:"hotel_id"`
	Hotel_name      string `json:"hotel_name"`
	Hotel_thumbnail string `json:"hotel_thumbnail"`
	check_in        string `json:"check_in"`
}

type RoomInfo struct {
	Id         string
	start_date string
	end_date   string
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
	// var user User
	// _ = json.NewDecoder(req.Body).Decode(&user)
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"username": user.Username,
	// 	"password": user.Password,
	// })
	// tokenString, error := token.SignedString([]byte("secret"))
	// if error != nil {
	// 	fmt.Println(error)
	// }
	// json.NewEncoder(w).Encode(JwtToken{Token: tokenString})

	auth, _ := firebase.GetAuth()
	token, err := auth.CreateCustomToken("FJNWvK2wbrhA2XHhWSQuiLVVFHp2", nil)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: token})

	// opt := option.WithCredentialsFile("firebaseAdminCredentials.json")
	// app, err := firebase.NewApp(context2.Background(), nil, opt)

	// client, err := app.Auth()
	// if err != nil {
	// 	return nil, fmt.Errorf("error getting Auth client: %v", err)
	// }

	// token, err := client.VerifyIDToken("")
	// if err != nil {
	// 	return nil, fmt.Errorf("error verifying ID token: %v", err)
	// }

	// fmt.Printf("Verified ID token: %v\n", token)

	// client, err := app.Auth()
	// if err != nil {
	// 	return err
	// }

	// claims := map[string]interface{}{
	// 	"package": "gold",
	// }
	// token, err := client.CustomToken("",)
	// token, err := client.CustomToken("",)
	// if err != nil {
	// 	return nil, fmt.Errorf("error creando la token", err)
	// }

	// json.NewEncoder(w).Encode(JwtToken{Token: token})

}

func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	token, _ := jwt.Parse(params["token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user User
		mapstructure.Decode(claims, &user)
		json.NewEncoder(w).Encode(user)
	} else {
		json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
	}
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				auth, _ := firebase.GetAuth()
				decodedToken, err := auth.VerifyIDToken(bearerToken[1])
				// println("hola")
				if err == nil {
					uid, found := decodedToken.UID()
					println("uid", uid)
					println("found", found)
					context.Set(req, "decoded", uid)
					next(w, req)

				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("udeain"))
}
func splitDate(date string) (retDate map[string]string) {
	dateArray := strings.Split(date, "-")
	//implementar validaciones para saber si existen los 3 atributos de la fecha
	retDate = map[string]string{
		"day":   dateArray[0],
		"month": dateArray[1],
		"year":  dateArray[2],
	}
	return
}
func getRooms(w http.ResponseWriter, r *http.Request) {

	/*
		vars := mux.Vars(r)

		//split the arrive date to get all the info
		arriveDate := vars["arriveDate"]
		arriveDateObj := splitDate(arriveDate)

		//split the arrive date to get all the info
		leaveDate := vars["leaveDate"]
		leaveDateObj := splitDate(leaveDate)

		city := vars["city"]
		hosts := vars["hosts"]
		roomType := vars["roomType"] */

	//city := "05001"

	arriveDate := r.URL.Query().Get("arrive_date")
	arriveDateObj := splitDate(arriveDate)
	leaveDate := r.URL.Query().Get("leave_date")
	city := r.URL.Query().Get("city")
	hosts := r.URL.Query().Get("hosts")
	roomType := r.URL.Query().Get("room_type")
	roomType = strings.ToUpper(roomType)

	println("searching.--.....----.")
	println("arriveDate", arriveDateObj["year"], arriveDateObj["month"], arriveDateObj["day"])
	//println("leaveDate",leaveDateObj["year"],leaveDateObj["month"],leaveDateObj["day"])
	println("leaveDate", leaveDate)
	println("city", city)
	println("hosts", hosts)
	println("roomType", roomType)

	session, err := mgo.Dial("mongodb://udeain:udeainmongodb@ds157444.mlab.com:57444/heroku_4r2js6cs")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("heroku_4r2js6cs").C("rooms")

	// result := Room{}
	var roomsObj []bson.M

	///////////////////////////////////////
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"room_type": roomType}},
		bson.M{"$match": bson.M{"city": city}},

		bson.M{"$lookup": bson.M{"from": "reservation", "localField": "id", "foreignField": "room_id", "as": "reservation"}},
		{"$unwind": bson.M{"path": "$reservation", "preserveNullAndEmptyArrays": false}},
		/* realizar filtro de fechas de reserva */
		bson.M{"$match": bson.M{"$or": []bson.M{bson.M{"reservation.start_date": bson.M{"$gte": leaveDate}}, bson.M{"reservation.end_date": bson.M{"$lte": arriveDate}}, bson.M{"reservation": bson.M{"$eq": nil}}}}},
		// omitir datos de reservas
		bson.M{"$project": bson.M{"reservation": 0}},
	}

	pipe := c.Pipe(pipeline)
	//resp := []bson.M{}
	err = pipe.All(&roomsObj)
	//respuesta, err :=  json.Marshal(resp)

	///////////////////////////////////////

	//err = c.Find(bson.M{"room_type": roomType, "city":city, "available":true}).All(&roomsObj)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("not found"))
		return
	}
	headerJson := []byte(`{`)
	if city == "05001" {
		cityInfo := []byte(`"hotel_id":"udeain_medellin","hotel_name":"udeain medellin", "hotel_location":{"address":"Cl. 5 Sur #42-2 a 42-70", "lat":"6.1992463", "long":"-75.5747155"},"hotel_thumbnail":"https://media-cdn.tripadvisor.com/media/photo-s/06/35/93/c2/hotel-el-deportista.jpg","check_in":"15:00","check_out":"13:00","hotel_website":"https://udeain.herokuapp.com", "rooms":`)
		headerJson = append(headerJson[:], cityInfo...)
	} else {
		cityInfo := []byte(`"hotel_id":"udeain_bogota","hotel_name":"udeain bogota", "hotel_location":{"address":"Cra. 14 #82-2 a 82-98", "lat":"4.667662", "long":"-74.0574518"},"hotel_thumbnail":"https://media-cdn.tripadvisor.com/media/photo-s/06/35/93/c2/hotel-el-deportista.jpg","check_in":"15:00","check_out":"13:00","hotel_website":"https://udeain.herokuapp.com", "rooms":`)
		headerJson = append(headerJson[:], cityInfo...)
	}
	respuesta, err := json.Marshal(roomsObj)

	jsonEnd := []byte(`}`)
	if string(respuesta) == "null" {
		respuesta = []byte(`[]`)
	}
	finalRes := append(headerJson[:], respuesta...)
	finalRes = append(finalRes[:], jsonEnd...)

	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("unable to get room"))
		return
	}

	/*var n Room
	if err2 := json.Unmarshal(finalRes, &n); err != nil {
		panic(err2)
	}
	fmt.Printf("%#v\n", n)
	fmt.Println( string(finalRes) )*/

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(finalRes)
	//w.Write( []byte(fmt.Sprintf("%v", n)) )

}

func getRoomsAvailable(w http.ResponseWriter, r *http.Request) {

	city := "05001"
	roomType := "s"
	fecha_inicio := "2017-10-28"
	fecha_fin := "2017-10-29"

	//roomType = r.Form.Get("room_type")
	roomType = r.URL.Query().Get("room_type")
	fecha_inicio = r.URL.Query().Get("arrive_date")
	fecha_fin = r.URL.Query().Get("leave_date")
	city = r.URL.Query().Get("city")

	fmt.Println(fecha_inicio + " " + fecha_fin + " " + city + " " + roomType)

	// establecer conexión
	session, err := mgo.Dial("mongodb://udeain:udeainmongodb@ds157444.mlab.com:57444/heroku_4r2js6cs")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection := session.DB("heroku_4r2js6cs").C("reservation")
	pipeline := []bson.M{
		/* filtro de fechas */

		bson.M{"$match": bson.M{"$or": []bson.M{bson.M{"start_date": bson.M{"$gte": fecha_fin}}, bson.M{"end_date": bson.M{"$lte": fecha_inicio}}}}},

		//bson.M{"$match": bson.M{"start_date": bson.M{"$lte": fecha_inicio} }},

		//bson.M{"$match": bson.M{ "$or": []bson.M{ bson.M{"hotel_id":"udeain_medellin"}, bson.M{"state": "awaiting"} } } },

		//bson.M{ "$or": []bson.M{ bson.M{"hotel_id":"udeain_medellin"}, bson.M{"state": "awaiting"} } },

		//bson.M{"$match": bson.M{"end_date": bson.M{"$lte": fecha_inicio} }},
		//bson.M{"$or": bson.M{ "start_date": bson.M{"$gte": fecha_fin}}},

		//bson.M{"$match": bson.M{"start_date": bson.M{"$gte": fecha_fin} }},
		//bson.M{"$match": {"$or": [{bson.M{"end_date": bson.M{"$lte": fecha_inicio} }},{bson.M{"start_date": bson.M{"$gte": fecha_fin} }}]} },
		//bson.M{"$match": bson.M{ "$or": [ bson.M{ "$lte": [ "end_date", fecha_inicio ] }, bson.M{ "$gte": [ "start_date", fecha_fin ] } ] } },
		//{ "$or": [ { "$lte": [ "end_date", fecha_inicio ] }, { "$gte": [ "start_date", fecha_fin ] } ] }

		/*Realizar 'Join' con documentos adicionales de hotel y datos de habitaciones*/
		bson.M{"$lookup": bson.M{"from": "rooms", "localField": "room_id", "foreignField": "id", "as": "rooms"}},
		bson.M{"$lookup": bson.M{"from": "hotel", "localField": "hotel_id", "foreignField": "hotel_id", "as": "hotel_details"}},
		/* Realizar filtrado por tipo de habitación y ciudad */
		{"$unwind": "$rooms"},
		bson.M{"$match": bson.M{"rooms.room_type": roomType}},
		bson.M{"$match": bson.M{"rooms.city": city}},
	}
	pipe := collection.Pipe(pipeline)
	resp := []bson.M{}
	err = pipe.All(&resp)

	respuesta, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("unable to get room"))
		return
	}

	//var a = jsonparser.GetInt(respuesta,"hotel_details", "[0]", "check_in")

	// asignar datos
	//jsonparser.Set(respuesta, []byte(hotel_name), "[0]", "hotel_name")
	//sjson.Set(`respuesta[0]`, "hotel_name", hotel_name)
	//respuesta[0].hotel_name = hotel_name;

	// asignar datos de acuerdo al formato
	// hotel_name := "";
	// hotel_thumb := "";
	// hotel_check_in := "";
	// hotel_check_out := "";
	// hotel_website := "";
	// hotel_address := "";
	// hotel_lat := "";
	// hotel_long := "";

	// hotel_name, err = jsonparser.GetString(respuesta, "[0]","hotel_details","[0]" ,"hotel_name")
	// hotel_thumb, err = jsonparser.GetString(respuesta, "[0]","hotel_details","[0]" ,"hotel_thumbnail")
	// hotel_check_in, err = jsonparser.GetString(respuesta, "[0]","hotel_details","[0]" ,"check_in")
	// hotel_check_out, err = jsonparser.GetString(respuesta, "[0]","hotel_details","[0]" ,"check_out")
	// hotel_website, err = jsonparser.GetString(respuesta, "[0]","hotel_details","[0]" ,"hotel_website")

	// hotel_address, err = jsonparser.GetString(respuesta, "[0]","hotel_details","[0]" ,"hotel_location", "address")
	// hotel_lat, err = jsonparser.GetString(respuesta, "[0]","hotel_details","[0]" ,"hotel_location", "lat")
	// hotel_long, err = jsonparser.GetString(respuesta, "[0]","hotel_details","[0]" ,"hotel_location", "long")
	// datos_hotel := map[string]string{"address": hotel_address, "lat": hotel_lat, "long": hotel_long}

	//fmt.Println( hotel_name )

	var datos []bson.M
	err = json.Unmarshal(respuesta, &datos)

	if err != nil {
		fmt.Println("error:", err)
	}

	//Asignar variables al Json
	// datos[0]["hotel_name"] = hotel_name;
	// datos[0]["hotel_thumbnail"] = hotel_thumb;
	// datos[0]["check_in"] = hotel_check_in;
	// datos[0]["check_out"] = hotel_check_out;
	// datos[0]["hotel_website"] = hotel_website;

	// datos[0]["hotel_location"] = datos_hotel

	/*var a = datos[0]["hotel_details"];
	md, ok := a.(map[string]interface{})
	fmt.Println( md["hotel_location"],ok )	*/

	// borrar datos reasignados
	// datos[0]["hotel_details"] = nil;

	// borrar datos adicionales temporalmente para retornar el formato establecido (falta hacer una operración para sacar habitaciones de estos que se borran)
	// respuesta, err =  json.Marshal(datos[0]) //////////
	// if err != nil {
	// 	w.WriteHeader(405)
	// 	w.Write([]byte("unable to get room"))
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respuesta)

}

func getReservationRequest(w http.ResponseWriter, r *http.Request) {
	decoded := context.Get(r, "decoded")
	var userJwt User
	mapstructure.Decode(decoded.(jwt.MapClaims), &userJwt)
	userId := userJwt.Username
	// establecer conexión con Base de Datos
	session, err := mgo.Dial("mongodb://udeain:udeainmongodb@ds157444.mlab.com:57444/heroku_4r2js6cs")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB("heroku_4r2js6cs").C("reservation")

	/*jsonDatos := []byte(`{"arrive_date":"2017-10-25","leave_date":"2017-10-26","room_type":"s","capacity":1,"beds":{"simple":1,"double":0},"hotel_id":"udeain_medellin",
	"user":{"doc_type":"CC","doc_id":"11521777","email":"cjmo@gmail.com","phone_number":"4448787"}}`)*/

	// Recibir datos Json enviados en solicitud POST
	jsonDatos, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	//println(string(jsonDatos))

	// Procesar datos recibidos
	var raw map[string]interface{}
	json.Unmarshal(jsonDatos, &raw)

	// obtener valor en particular  //arrive_date := raw["arrive_date"]
	salida, _ := json.Marshal(raw["arrive_date"])
	arrive_date := string(salida)
	arrive_date = strings.Replace(arrive_date, "\"", "", -1)
	salida, _ = json.Marshal(raw["leave_date"])
	leave_date := string(salida)
	leave_date = strings.Replace(leave_date, "\"", "", -1)
	salida, _ = json.Marshal(raw["room_type"])
	room_type := strings.ToUpper(string(salida))
	room_type = strings.Replace(room_type, "\"", "", -1)
	salida, _ = json.Marshal(raw["capacity"])
	capacity := string(salida)
	capacity_number, err := strconv.Atoi(capacity)
	salida, _ = json.Marshal(raw["hotel_id"])
	hotel_id := string(salida)
	hotel_id = strings.Replace(hotel_id, "\"", "", -1)
	salida, _ = json.Marshal(raw["beds"])
	beds := string(salida)
	salida, _ = json.Marshal(raw["user"])
	user := string(salida)

	// procesar subelemento 'beds'
	var rawBeds map[string]interface{}
	json.Unmarshal([]byte(beds), &rawBeds)
	salida, _ = json.Marshal(rawBeds["double"])
	beds_double := string(salida)
	beds_double = strings.Replace(beds_double, "\"", "", -1)
	salida, _ = json.Marshal(rawBeds["simple"])
	beds_simple := string(salida)
	beds_simple = strings.Replace(beds_simple, "\"", "", -1)

	// procesar subelemento 'user'
	var rawUser map[string]interface{}
	json.Unmarshal([]byte(user), &rawUser)
	salida, _ = json.Marshal(rawUser["doc_type"])
	doc_type := strings.ToLower(string(salida))
	doc_type = strings.Replace(doc_type, "\"", "", -1)
	salida, _ = json.Marshal(rawUser["doc_id"])
	doc_id := string(salida)
	doc_id = strings.Replace(doc_id, "\"", "", -1)
	salida, _ = json.Marshal(rawUser["email"])
	email := string(salida)
	email = strings.Replace(email, "\"", "", -1)
	salida, _ = json.Marshal(rawUser["phone_number"])
	phone_number := string(salida)
	phone_number = strings.Replace(phone_number, "\"", "", -1)

	// realizar búsqueda de habitaciones disponibles
	collection_rooms := session.DB("heroku_4r2js6cs").C("rooms")
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"room_type": room_type}},
		bson.M{"$match": bson.M{"hotel_id": hotel_id}},

		bson.M{"$lookup": bson.M{"from": "reservation", "localField": "id", "foreignField": "room_id", "as": "reservation"}},
		{"$unwind": bson.M{"path": "$reservation", "preserveNullAndEmptyArrays": true}},
		/* realizar filtro de fechas de reserva */
		bson.M{"$match": bson.M{"$or": []bson.M{bson.M{"reservation.start_date": bson.M{"$gte": leave_date}}, bson.M{"reservation.end_date": bson.M{"$lte": arrive_date}}, bson.M{"reservation": bson.M{"$eq": nil}}}}},
	}
	pipe := collection_rooms.Pipe(pipeline)
	resp := []bson.M{}
	err = pipe.All(&resp)
	respuesta, err := json.Marshal(resp)

	var datos []bson.M
	err = json.Unmarshal(respuesta, &datos)

	if err != nil {
		fmt.Println("error:", err)
	}

	//fmt.Println("habitaciones :", len(datos))

	// validar si hay habitaciones disponibles
	if len(datos) == 0 {
		w.WriteHeader(409)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message" : "No hay habitaciones disponibles para el rango de fechas especificado, intente de nuevo"}`))
		return
	}

	// obtener id de una habitación disponible
	var room_id = datos[0]["_id"]
	fmt.Println(room_id)

	// validar escogencia de habitación disponible
	if room_id == "" || room_id == nil {
		w.WriteHeader(409)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message" : "No hay habitaciones disponibles para el rango de fechas especificado, intente de nuevo"}`))
		return
	}

	// validación de errores de datos en json recibido "\"\""
	if arrive_date == "" || arrive_date == "null" || leave_date == "" || leave_date == "null" || room_type == "" ||
		room_type == "null" || capacity == "null" || capacity_number <= 0 || hotel_id == "null" || hotel_id == "" ||
		beds_double == "null" || beds_double == "" || beds_simple == "null" || beds_simple == "" ||
		doc_type == "null" || doc_type == "" || doc_id == "null" || doc_id == "" || email == "null" || email == "" ||
		phone_number == "null" || phone_number == "" {

		w.WriteHeader(409)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message" : "Los parámetros de la reserva no han sido especificados en su totalidad, o presentan errores de formato"}`))
		return
	} else if room_type != "S" && room_type != "L" {
		// validar tipo de habitación requerida*/
		w.WriteHeader(409)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message" : "Tipo de habitación no registrada"}`))
		return
	} else if doc_type != "cc" && doc_type != "pp" && doc_type != "ce" {
		// validar tipo de documento de identidad cc: cédula ciudananía; pp: pasaporte; ce: cédula extranjería
		w.WriteHeader(409)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message" : "Tipo de documento no válido. Tipos admitidos: 'cc': cédula ciudananía; 'pp': pasaporte; 'ce': cédula extranjería"}`))
		return
	}

	// guardar datos de reserva
	id_reserva := bson.NewObjectId().Hex()

	collection.Insert(bson.M{"userId": userId, "_id": id_reserva, "start_date": arrive_date, "end_date": leave_date, "state": "awaiting", "host_id": "0045123", "hotel_id": hotel_id,
		"room_type": room_type, "capacity": capacity_number, "beds_double": beds_double, "beds_simple": beds_simple, "doc_type": doc_type, "doc_id": doc_id,
		"email": email, "phone_number": phone_number, "room_id": room_id})
	println("ID reserva generada: " + id_reserva)

	// retornar respuesta de reserva gestionada
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	info_reserva := `{"reservation_id": "` + id_reserva + `"}`

	w.Write([]byte(info_reserva))

}

//mongodb://udeain:udeainmongodb@ds157444.mlab.com:57444/heroku_4r2js6cs
//http://localhost:8080/api/v1/rooms/arrive_date/01-01-2017/leave_date/02-02-2017/city/05001/hosts/3/room_type/l
func main() {
	fmt.Println("start server 8080")
	// opt := option.WithCredentialsFile("firebaseAdminCredentials.json")
	// app, err := firebase.NewApp(context.Background(), nil, opt)

	// if err != nil {
	// 	return nil, fmt.Errorf("error inicializando la aplicación", err)
	// }

	firebase.InitializeApp(&firebase.Options{
		ServiceAccountPath: "firebaseAdminCredentials.json",
	})

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")

	//r.HandleFunc("/api/v1/rooms/arrive_date/{arriveDate}/leave_date/{leaveDate}/city/{city}/hosts/{hosts}/room_type/{roomType}", getRooms).Methods("GET")
	r.HandleFunc("/api/v1/rooms", getRooms).Methods("GET")
	r.HandleFunc("/api/v1/rooms_info", getRoomsAvailable).Methods("GET")
	r.HandleFunc("/api/v1/rooms/reserve", ValidateMiddleware(getReservationRequest)).Methods("POST")

	r.HandleFunc("/authenticate", CreateTokenEndpoint).Methods("POST")
	r.HandleFunc("/protected", ProtectedEndpoint).Methods("GET")
	r.HandleFunc("/test", ValidateMiddleware(TestEndpoint)).Methods("GET")

	http.Handle("/", r)
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	//http.ListenAndServe("0.0.0.0:"+port, nil)
	corsObj := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":"+port, handlers.CORS(corsObj)(r))
	// http.ListenAndServe(":"+port, nil)
}
