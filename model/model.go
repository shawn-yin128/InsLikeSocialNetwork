package model

type Post struct {
	// struct tags used for json serialization if we want to insert into db as object
	// if we insert db with json string, no need for json serialization
	Id      string `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	Url     string `json:"url"`
	Type    string `json:"type"`
}
