package web

import (
	"Neo4jPlayground/internal/config"
	"Neo4jPlayground/internal/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	PERSON_PATH      = "/neo/v1/persons"
	PERSON_WITH_NAME = "/neo/v1/persons/{name}"
	TASK_PATH        = "/neo/v1/tasks"
	TASK_WITH_TITLE  = "/neo/v1/tasks/{title}"
	ASSIGN_PATH      = "/neo/v1/assign"
)

func Routes(conf *config.Conf) {
	route := mux.NewRouter()
	route.HandleFunc(PERSON_PATH, services.GetPersons).Methods(http.MethodGet)
	route.HandleFunc(PERSON_PATH, services.AddPerson).Methods(http.MethodPost)
	route.HandleFunc(PERSON_WITH_NAME, services.UpdatePerson).Methods(http.MethodPut)
	route.HandleFunc(PERSON_WITH_NAME, services.DeletePerson).Methods(http.MethodDelete)

	route.HandleFunc(TASK_PATH, services.GetTasks).Methods(http.MethodGet)
	route.HandleFunc(TASK_PATH, services.AddTasks).Methods(http.MethodPost)
	route.HandleFunc(TASK_WITH_TITLE, services.UpdateTasks).Methods(http.MethodPut)
	route.HandleFunc(TASK_WITH_TITLE, services.DeleteTasks).Methods(http.MethodDelete)

	route.HandleFunc(ASSIGN_PATH, services.AssignTask).Methods(http.MethodPost)

	// Start the server and handle any potential errors
	server := &http.Server{
		Addr:    conf.App.Port, // Address to listen on
		Handler: route,         // The router to handle incoming requests
	}

	// Log fatal error if http.ListenAndServe fails
	log.Fatal(server.ListenAndServe())
}
