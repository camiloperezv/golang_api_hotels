Programa para gestión de reservas hoteleras desarrollado en el lenguaje GO, bajo base de datos no relacionales Mongo DB.

Consulta de habitaciones disponibles para reserva

URL: https://udeain.herokuapp.com/api/v1/rooms?arrive_date=01-01-2017&leave_date=02-02-2017&city=05001&hosts=3&room_type=l

Realizar reserva de habitaciones

URL: https://udeain.herokuapp.com/api/v1/rooms/reserve

Se recibe una solicitud de tipo JSON con el siguiente formato:

{"arrive_date":"2017-10-26","leave_date":"2017-10-27","room_type":"l","capacity":1,
"beds":{"simple":1,"double":0},"hotel_id":"udeain_medellin",
"user":{"doc_type":"Cc","doc_id":"11521777","email":"cjmo@gmail.com","phone_number":"4448787"}}

Se retorna en caso de éxito una respuesta en formato JSON con id de reserva generada (alfanumérico): 

{"reservation_id":"ID RESERVA"}

En caso negativo se retorna un mensaje de error con el formato JSON:

{"message":"MENSAJE DE ERROR"}
