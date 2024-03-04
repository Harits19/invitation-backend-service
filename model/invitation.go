package model

import "time"

type Invitation struct {
	Id      string     `json:"id" bson:"_id"`
	Name    string     `json:"name"`
	Music   string     `json:"music"`
	Initial string     `json:"initial"`
	Date    time.Time  `json:"date"`
	Groom   BrideGroom `json:"groom"`
	Bride   BrideGroom `json:"bride"`
	Address Address    `json:"address"`
}

type BrideGroom struct {
	Name       string `json:"name"`
	FatherName string `json:"father_name" bson:"father_name"`
	MotherName string `json:"mother_name" bson:"mother_name"`
	Address    string `json:"address"`
	Instagram  string `json:"instagram"`
}

type Address struct {
	Detail    string `json:"detail"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
