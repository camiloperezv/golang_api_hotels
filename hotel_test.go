package main

/*
"strings"*/

import (
	"encoding/json"
	"fmt"
	//"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	//"reflect"
	"bytes"
	"net/url"
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

func TestPuerto(t *testing.T) {
	// prueba de puerto de comunicación
	puerto := os.Getenv("PORT")
	fmt.Println("Puerto: " + puerto)

	puerto_esperado := "8080"
	if puerto != puerto_esperado {
		t.Error("Test Fallido: Puerto de comunicación del servidor diferente al esperado")
	}
}

// pruebas de conexión y respuestas HTTP
func TestHTTP(t *testing.T) {
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

// Prueba de procesamiento de fechas en formato especificado
func TestFechas(t *testing.T) {

	date := "2017-11-05"
	dateObj := splitDate(date)
	fecha_esperada := "2017"
	//fmt.Println(dateObj["year"])

	// prueba de año
	if dateObj["year"] != fecha_esperada {
		t.Errorf("Test Fallido: Formato de fecha incorrecto. Año esperado %s y Obtenido %s ", fecha_esperada, dateObj["year"])
	}

	// prueba de mes
	fecha_esperada = "11"
	if dateObj["month"] != fecha_esperada {
		t.Errorf("Test Fallido: Formato de fecha incorrecto. Mes esperado %s y Obtenido %s ", fecha_esperada, dateObj["month"])
	}

	// prueba de día
	fecha_esperada = "05"
	if dateObj["day"] != fecha_esperada {
		t.Errorf("Test Fallido: Formato de fecha incorrecto. Día esperado %s y Obtenido %s ", fecha_esperada, dateObj["day"])
	}
}

func TestEstructuras(t *testing.T) {
	room := Room{Hotel_id: "udeain_medellin", Hotel_name: "Udea IN", Hotel_thumbnail: "N/A", check_in: "11:00 am"}

	if room.Hotel_id != "udeain_medellin" {
		t.Errorf("Test Fallido: Datos en estructura distintos a los especificados. Esperado %s y Obtenido %s ", room.Hotel_id, "udeain_medellin")
	}
}

func TestReserva(t *testing.T) {

	form := url.Values{
		"arrive_date": {"2017-11-26"},
		"leave_date":  {"2017-11-27"},
		"room_type":   {"l"},
		"capacity":    {"1"},
		"hotel_id":    {"udeain_medellin"},
		/*"beds":        {"simple": {"1"}, "double": {"0"}},
		"user":        {"doc_type": {"Ccu"}, "doc_id": {"11521777"}, "email": {"cjmo@gmail.com"}, "phone_number": {"4448787"}},*/
	}
	body_request := bytes.NewBufferString(form.Encode())

	resp, err := http.Post("https://udeain.herokuapp.com/api/v1/rooms/reserve", "application/json", body_request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// obtener mensaje de respuesta para solicitud de reserva
	var raw map[string]interface{}
	json.Unmarshal(body, &raw)
	salida, _ := json.Marshal(raw["message"])
	mensaje := string(salida)
	mensaje = strings.Replace(mensaje, "\"", "", -1)
	valor_esperado := "No hay habitaciones disponibles para el rango de fechas especificado, intente de nuevo"

	if err != nil {
		fmt.Println("error:", err)
	}

	//fmt.Println(string(salida))
	if mensaje != valor_esperado {
		t.Errorf("Test Fallido: Mensaje de respuesta para reserva erróneo. Esperado: %s y Obtenido: %s ", valor_esperado, mensaje)
	} else {
		fmt.Println("Test de Mensaje de respuesta para reserva aprobado")
	}

	reader := strings.NewReader("")
	req, _ := http.NewRequest("POST", "/api/v1/rooms/reserve", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// probar código de respuesta
	if w.Code != 200 {
		t.Errorf("Test Fallido: Código de respuesta HTTP para reservas erróneo. Esperado: %d y Obtenido: %d ", 200, w.Code)
	} else {
		fmt.Println("Test de Código de respuesta HTTP para reservas aprobado")
	}

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Header.Get("Content-Type"))

}

func TestPublish(t *testing.T) {

}
