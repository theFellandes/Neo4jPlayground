package models

type Model interface {
	ToMap() map[string]interface{}
}

type Task struct {
	Title       string
	Description string
}

func (t Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":       t.Title,
		"description": t.Description,
	}
}

type Person struct {
	Name string
}

func (p Person) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name": p.Name,
	}
}
