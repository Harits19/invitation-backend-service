package model

import "time"

type Invitation struct {
	Id      *string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string     `json:"name" bson:"name,omitempty"`
	Music   *string    `json:"music" bson:"music,omitempty"`
	Initial *string    `json:"initial" bson:"initial,omitempty"`
	Date    time.Time  `json:"date" bson:"date,omitempty"`
	Groom   BrideGroom `json:"groom" bson:"groom,omitempty"`
	Bride   BrideGroom `json:"bride" bson:"bride,omitempty"`
	Address Address    `json:"address" bson:"address,omitempty"`
	Photo   *Photo     `json:"photo" bson:"photo,omitempty"`
}

type Photo struct {
	Cover *string `json:"cover" bson:"cover,omitempty"`
	Side  struct {
		Top    *string `json:"top" bson:"top,omitempty"`
		Bottom *string `json:"bottom" bson:"bottom,omitempty"`
	} `json:"side" bson:"side,omitempty"`
	Background *string   `json:"background" bson:"background,omitempty"`
	Slide      *[]string `json:"slide" bson:"slide,omitempty"`
	Divider    *string   `json:"divider" bson:"divider,omitempty"`
	Gallery    *[]string `json:"gallery" bson:"gallery,omitempty"`
}

type BrideGroom struct {
	Name       *string `json:"name" bson:"name,omitempty"`
	FatherName string  `json:"father_name" bson:"father_name,omitempty"`
	MotherName string  `json:"mother_name" bson:"mother_name,omitempty"`
	Address    string  `json:"address" bson:"address,omitempty"`
	Instagram  string  `json:"instagram" bson:"instagram,omitempty"`
	Photo      *string `json:"photo" bson:"photo,omitempty"`
}

type Address struct {
	Detail    string `json:"detail" bson:"detail,omitempty"`
	Latitude  string `json:"latitude" bson:"latitude,omitempty"`
	Longitude string `json:"longitude" bson:"longitude,omitempty"`
}
