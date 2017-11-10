package main

/*
"strings"*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	//"reflect"
	"strconv"
	"strings"
	"testing"
)

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("Ha ocurrido un error. %v", err)
	}
}

func Test1(t *testing.T) {
	expected := "udeain"
	actual := "udeain"
	if actual != expected {
		t.Error("Test Fallido")
	}
}

func Test2(t *testing.T) {

	// formato de json enviado por medio de formulario web
	jsonDatos := []byte(`{"arrive_date":"2017-10-25","leave_date":"2017-10-26","room_type":"s","capacity":1,"beds":{"simple":1,"double":0},"hotel_id":"udeain_medellin",
		"user":{"doc_type":"CC","doc_id":"11521777","email":"cjmo@gmail.com","phone_number":"4448787"}}`)

	// Procesar datos recibidos
	var raw map[string]interface{}
	json.Unmarshal(jsonDatos, &raw)

	salida, _ := json.Marshal(raw["arrive_date"])

	expected := string("2017-10-25")

	arrive_date := string(salida)
	arrive_date = strings.Replace(arrive_date, "\"", "", -1)

	//fmt.Println(arrive_date)

	// prueba de verificación de fecha de llegada
	if arrive_date != expected {
		t.Error("Test Fallido para dato arrive_date")
	} else {
		fmt.Println("Test 2.1 de obtención de dato en respuesta aprobado")
	}

	// prueba de verificación de fechas (que la de llegada sea inferior a la de salida)
	salida, _ = json.Marshal(raw["leave_date"])
	leave_date := string(salida)
	leave_date = strings.Replace(leave_date, "\"", "", -1)

	if arrive_date > leave_date {
		t.Error("Test Fallido: fecha de llegada es posterior a la de salida")
	} else {
		fmt.Println("Test 2.2 de obtención de dato en respuesta aprobado")
	}

	// prueba de tipo de habitación
	salida, _ = json.Marshal(raw["room_type"])
	room_type := string(salida)
	room_type = strings.Replace(room_type, "\"", "", -1)
	room_type = strings.ToUpper(room_type)
	if room_type != "S" && room_type != "L" {
		t.Error("Test Fallido: Tipo de habitación distinta a las soportadas: S y L")
	} else {
		fmt.Println("Test 2.3 de obtención de dato en respuesta aprobado")
	}

	// prueba de capacidad de habitación (personas)
	salida, _ = json.Marshal(raw["capacity"])
	capacity := string(salida)
	capacity = strings.Replace(capacity, "\"", "", -1)
	capacity_number, err := strconv.Atoi(capacity)
	//fmt.Println(err)
	if capacity_number < 1 || err != nil {
		t.Error("Test Fallido: La capacidad de la habitación debe ser mínimo de 1 persona")
	} else {
		fmt.Println("Test 2.4 de obtención de dato en respuesta aprobado")
	}

}

func Test3(t *testing.T) {
	// prueba de puerto de comunicación
	puerto := os.Getenv("PORT")
	fmt.Println("Puerto: " + puerto)

	puerto_esperado := "8080"
	if puerto != puerto_esperado {
		t.Error("Test Fallido: Puerto de comunicación del servidor diferente al esperado")
	}
}

// pruebas de conexión y respuestas HTTP
func Test4(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	checkError(err, t)

	resp, err := http.Get("https://udeain.herokuapp.com/api/v1/rooms?arrive_date=01-01-2017&leave_date=02-02-2017&city=05001&hosts=3&room_type=l")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var raw map[string]interface{}
	json.Unmarshal(body, &raw)

	//fmt.Println("Tipo: ", reflect.TypeOf(raw))

	// obtener tipo de respuesta de petición HTTP
	switch interface{}(raw).(type) {
	case map[string]interface{}:
		fmt.Println("Test de Tipo de respuesta a petición HTTP aprobado")
	default:
		t.Error("Test Fallido: Tipo de respuesta a petición HTTP no esperado")
	}

	/*var tipo, errtipo = fmt.Printf("%T", raw)
	var tipo_string = string(tipo)
	var tipo_esperado = "map[string]interface {}"
	fmt.Println(tipo_string)
	fmt.Println(errtipo)

	if tipo_string != tipo_esperado {
		t.Error("Test Fallido: Tipo de respuesta no esperado")
	}*/

	// prueba a función interna
	rr := httptest.NewRecorder()
	http.HandlerFunc(getRoomsAvailable).ServeHTTP(rr, req)

	//Confirmar estado de código de respuesta
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Test Fallido: Código de estado HTTP erróneo. Esperado %d y Obtenido %d ", http.StatusOK, status)
	} else {
		fmt.Println("Test de estado de respuesta en servidor aprobado")
	}

	//fmt.Println(rr.Body.String())

}
