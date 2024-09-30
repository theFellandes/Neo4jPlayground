package services

import (
	"Neo4jPlayground/internal/handlers"
	"Neo4jPlayground/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

// isRequestMethodValid checks whether the sent request method is valid or not.
func isRequestMethodValid(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// readBody reads the body of the request.
func readBody(w http.ResponseWriter, r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing request body: %v\n", err)
		}
	}(r.Body)

	return body
}

// parsePerson unmarshals json and returns the Person element
func parsePerson(w http.ResponseWriter, body []byte) models.Person {
	var data models.Person
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return models.Person{}
	}
	return data
}

// parseTask unmarshals json and returns the Task element
func parseTask(w http.ResponseWriter, body []byte) models.Task {
	var data models.Task
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return models.Task{}
	}
	return data
}

// createResponse sets the header and adds LoginRequest to body
func createResponse(model models.Model, w http.ResponseWriter) {
	responseData := model.ToMap()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// createResponseForList sets the header and adds list to body
func createResponseForList(models []models.Model, w http.ResponseWriter) {
	var responseData []map[string]interface{}
	for _, model := range models {
		responseData = append(responseData, model.ToMap())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleParameters returns the given parameter.
func handleParameters(r *http.Request, parameter string) string {
	vars := mux.Vars(r)
	param := vars[parameter]
	return param
}

func GetPersons(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodGet) {
		return
	}
	persons, err := handlers.GetPersons()
	if err != nil {
		http.Error(writer, "Failed to get person", http.StatusInternalServerError)
		return
	}

	// Convert []Person to []Model
	var modelList []models.Model
	for _, person := range persons {
		modelList = append(modelList, person)
	}

	createResponseForList(modelList, writer)
}

func AddPerson(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodPost) {
		return
	}
	person := parsePerson(writer, readBody(writer, request))
	err := handlers.CreatePerson(person)
	if err != nil {
		http.Error(writer, "Failed to create person", http.StatusInternalServerError)
	}
}

func UpdatePerson(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodPut) {
		return
	}
	person := parsePerson(writer, readBody(writer, request))
	// TODO: this will be changed, parameter is wrong.
	handlers.UpdatePerson(handleParameters(request, "name"), person.Name)
}

func DeletePerson(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodDelete) {
		return
	}
	handlers.DeletePerson(handleParameters(request, "name"))
}

func GetTasks(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodGet) {
		return
	}
	tasks, err := handlers.GetTasks()
	if err != nil {
		http.Error(writer, "Failed to get person", http.StatusInternalServerError)
		return
	}

	// Convert []Task to []Model
	var modelList []models.Model
	for _, task := range tasks {
		modelList = append(modelList, task)
	}

	createResponseForList(modelList, writer)
}

func AddTasks(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodPost) {
		return
	}
	task := parseTask(writer, readBody(writer, request))
	handlers.CreateTask(task)
}

func UpdateTasks(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodPut) {
		return
	}
	task := parseTask(writer, readBody(writer, request))
	handlers.UpdateTask(handleParameters(request, "title"), task)
}

func DeleteTasks(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodDelete) {
		return
	}
	task := parseTask(writer, readBody(writer, request))
	handlers.DeleteTask(task)
}

func AssignTask(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodPost) {
		return
	}
	person := parsePerson(writer, readBody(writer, request))
	task := parseTask(writer, readBody(writer, request))
	handlers.AssignPersonToTask(person, task)
}

func GetTasksForPerson(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodPost) {
		return
	}
	tasks, err := handlers.GetTasksForPerson(handleParameters(request, "name"))

	if err != nil {
		http.Error(writer, "Failed to get tasks", http.StatusInternalServerError)
		return
	}

	// Convert []Task to []Model
	var modelList []models.Model
	for _, task := range tasks {
		modelList = append(modelList, task)
	}

	createResponseForList(modelList, writer)
}

func GetPersonsForTask(writer http.ResponseWriter, request *http.Request) {
	if !isRequestMethodValid(writer, request, http.MethodPost) {
		return
	}
	tasks, err := handlers.GetPersonsForTask(handleParameters(request, "name"))

	if err != nil {
		http.Error(writer, "Failed to get tasks", http.StatusInternalServerError)
		return
	}

	// Convert []Task to []Model
	var modelList []models.Model
	for _, task := range tasks {
		modelList = append(modelList, task)
	}

	createResponseForList(modelList, writer)
}
