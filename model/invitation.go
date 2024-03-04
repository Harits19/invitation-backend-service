package model

type Invitation struct {
	Id   string `json:"id"`
	Name      string     `json:"name"`
	Music     string     `json:"music"`
	Initial   string     `json:"initial"`
	Date      string     `json:"date"`
	Groom     BrideGroom `json:"groom"`
	Bride     BrideGroom `json:"bride"`
	Address   Address    `json:"address"`
	Instagram string     `json:"instagram"`
}

type BrideGroom struct {
	Name       string `json:"name"`
	FatherName string `json:"father_name"`
	MotherName string `json:"mother_name"`
	Address    string `json:"address"`
}

type Address struct {
	Detail    string `json:"detail"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
