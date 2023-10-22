package structs

import "github.com/google/uuid"

type GenderType string

const (
	Male   GenderType = "male"
	Female GenderType = "female"
)

type Person struct {
	ID         int        `json:"-" db:"id"`
	UUID       uuid.UUID  `json:"uuid" db:"uuid"`
	Name       string     `json:"name" db:"name"`
	Surname    string     `json:"surname" db:"surname"`
	Patronymic string     `json:"patronymic" db:"patronymic"`
	Age        uint       `json:"age" db:"age"`
	Gender     GenderType `json:"gender" db:"gender"`
	Nation     string     `json:"nation" db:"nation"`
}
