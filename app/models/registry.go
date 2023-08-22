package models

type Model struct {
	Model interface{}
}

func RegisterModels() []Model {
	return []Model{
		{Model: Employee{}},
		{Model: Attendance{}},
	}
}
