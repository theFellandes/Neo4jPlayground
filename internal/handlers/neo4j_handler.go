package handlers

import (
	"Neo4jPlayground/internal/config"
	"Neo4jPlayground/internal/models"
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

var driver neo4j.DriverWithContext

func InitDriver(c *config.Conf) {
	connectionString := fmt.Sprintf("neo4j://%s:%d", c.DB.Ip, c.DB.Port)
	var err error
	driver, err = neo4j.NewDriverWithContext(connectionString, neo4j.BasicAuth(c.DB.User, c.DB.Password, ""))
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}
}

func CloseDriver() {
	if driver != nil {
		err := driver.Close(context.Background())
		if err != nil {
			return
		}
	}
}

func withSession(fn func(neo4j.SessionWithContext) (interface{}, error)) (interface{}, error) {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	return fn(session)
}

func CreateTask(task models.Task) {
	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "CREATE (t:Task {title: $title, description: $description}) RETURN t"
			params := map[string]interface{}{"title": task.Title, "description": task.Description}
			_, err := tx.Run(context.Background(), query, params)
			return nil, err
		})
	})

	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	} else {
		log.Println("Task created successfully!")
	}
}

func CreatePerson(person models.Person) error {
	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "CREATE (p:Person {name: $name}) RETURN p"
			params := map[string]interface{}{"name": person.Name}
			_, err := tx.Run(context.Background(), query, params)
			return nil, err
		})
	})

	if err != nil {
		log.Fatalf("Failed to create person: %v", err)
		return err
	}
	fmt.Println("Person created successfully!")
	return nil
}

func GetPersons() ([]models.Person, error) {
	var persons []models.Person

	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "MATCH (p:Person) RETURN p.name"
			result, err := tx.Run(context.Background(), query, nil)
			if err != nil {
				return nil, err
			}
			for result.Next(context.Background()) {
				name, _ := result.Record().Get("p.name")
				person := models.Person{
					Name: name.(string),
				}
				persons = append(persons, person)
			}
			return nil, result.Err()
		})
	})

	if err != nil {
		log.Fatalf("Failed to get persons: %v", err)
		return nil, err
	}

	return persons, nil
}

func UpdatePerson(oldName, newName string) {
	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := `
				MATCH (p:Person {name: $oldName})
				SET p.name = $newName
				RETURN p`
			params := map[string]interface{}{"oldName": oldName, "newName": newName}
			_, err := tx.Run(context.Background(), query, params)
			return nil, err
		})
	})

	if err != nil {
		log.Fatalf("Failed to update person: %v", err)
	} else {
		log.Println("Person updated successfully!")
	}
}

func DeletePerson(name string) {
	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "MATCH (p:Person {name: $name}) DETACH DELETE p"
			params := map[string]interface{}{"name": name}
			_, err := tx.Run(context.Background(), query, params)
			return nil, err
		})
	})

	if err != nil {
		log.Fatalf("Failed to delete person: %v", err)
	} else {
		log.Println("Person deleted successfully!")
	}
}

func AssignPersonToTask(person models.Person, task models.Task) {
	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := `
			MATCH (p:Person {name: $personName}), (t:Task {title: $taskTitle})
			MERGE (p)-[:ASSIGNED_TO]->(t)
			RETURN p, t`
			params := map[string]interface{}{"personName": person.Name, "taskTitle": task.Title}
			_, err := tx.Run(context.Background(), query, params)
			return nil, err
		})
	})

	if err != nil {
		log.Fatalf("Failed to assign person to task: %v", err)
	} else {
		fmt.Println("Person assigned to task successfully!")
	}
}

func GetTasksForPerson(personName string) ([]models.Task, error) {
	var tasks []models.Task

	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := `
				MATCH (p:Person {name: $name})-[:ASSIGNED_TO]->(t:Task)
				RETURN t.title, t.description`
			params := map[string]interface{}{"name": personName}
			result, err := tx.Run(context.Background(), query, params)
			if err != nil {
				return nil, err
			}
			for result.Next(context.Background()) {
				title, _ := result.Record().Get("t.title")
				description, _ := result.Record().Get("t.description")
				task := models.Task{
					Title:       title.(string),
					Description: description.(string),
				}
				tasks = append(tasks, task)
			}
			return nil, result.Err()
		})
	})

	if err != nil {
		log.Fatalf("Failed to get tasks for person '%s': %v", personName, err)
		return nil, err
	}

	return tasks, nil
}

func GetPersonsForTask(taskTitle string) ([]models.Person, error) {
	var persons []models.Person

	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := `
				MATCH (t:Task {title: $title})<-[:ASSIGNED_TO]-(p:Person)
				RETURN p.name`
			params := map[string]interface{}{"title": taskTitle}
			result, err := tx.Run(context.Background(), query, params)
			if err != nil {
				return nil, err
			}
			for result.Next(context.Background()) {
				name, _ := result.Record().Get("p.name")
				person := models.Person{
					Name: name.(string),
				}
				persons = append(persons, person)
			}
			return nil, result.Err()
		})
	})

	if err != nil {
		log.Fatalf("Failed to get persons for task '%s': %v", taskTitle, err)
		return nil, err
	}

	return persons, nil
}

func GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "MATCH (t:Task) RETURN t.title, t.description"
			result, err := tx.Run(context.Background(), query, nil)
			if err != nil {
				return nil, err
			}
			for result.Next(context.Background()) {
				title, _ := result.Record().Get("t.title")
				description, _ := result.Record().Get("t.description")
				task := models.Task{
					Title:       title.(string),
					Description: description.(string),
				}
				tasks = append(tasks, task)
			}
			return nil, result.Err()
		})
	})

	if err != nil {
		log.Fatalf("Failed to get tasks: %v", err)
		return nil, err
	}

	return tasks, nil
}

func UpdateTask(oldTitle string, updatedTask models.Task) {
	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := `
			MATCH (t:Task {title: $oldTitle})
			SET t.title = $newTitle, t.description = $newDescription
			RETURN t`
			params := map[string]interface{}{"oldTitle": oldTitle, "newTitle": updatedTask.Title, "newDescription": updatedTask.Description}
			_, err := tx.Run(context.Background(), query, params)
			return nil, err
		})
	})

	if err != nil {
		log.Fatalf("Failed to update task: %v", err)
	} else {
		fmt.Println("Task updated successfully!")
	}
}

func DeleteTask(task models.Task) {
	_, err := withSession(func(session neo4j.SessionWithContext) (interface{}, error) {
		return session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "MATCH (t:Task {title: $title}) DETACH DELETE t"
			params := map[string]interface{}{"title": task.Title}
			_, err := tx.Run(context.Background(), query, params)
			return nil, err
		})
	})

	if err != nil {
		log.Fatalf("Failed to delete task: %v", err)
	} else {
		fmt.Println("Task deleted successfully!")
	}
}
