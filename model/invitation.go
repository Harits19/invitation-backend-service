package model

import "time"

type Invitation struct {
	Id      *string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string     `json:"name"`
	Music   *string     `json:"music"`
	Initial *string    `json:"initial"`
	Date    time.Time  `json:"date"`
	Groom   BrideGroom `json:"groom"`
	Bride   BrideGroom `json:"bride"`
	Address Address    `json:"address"`
	Photo   Photo      `json:"photo"`
}

type Photo struct {
	Cover string `json:"cover"`
	Side  struct {
		Top    string `json:"top"`
		Bottom string `json:"bottom"`
	}
	Background string   `json:"background"`
	Slide      []string `json:"slide"`
	Divider    string   `json:"divider"`
	Gallery    []string `json:"gallery"`
}

type BrideGroom struct {
	Name       *string `json:"name"`
	FatherName string  `json:"father_name" bson:"father_name"`
	MotherName string  `json:"mother_name" bson:"mother_name"`
	Address    string  `json:"address"`
	Instagram  string  `json:"instagram"`
	Photo      *string  `json:"photo"`
}

type Address struct {
	Detail    string `json:"detail"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
