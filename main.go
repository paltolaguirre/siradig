package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroserviciosiradig, router)
	fmt.Println("Microservicio de Siradig escuchando en el puerto: " + configuracion.Puertomicroserviciosiradig)
	log.Fatal(server)

}
