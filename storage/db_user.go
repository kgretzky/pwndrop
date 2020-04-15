package storage

import (
	"strings"
)

type DbUser struct {
	ID         int    `json:"id" storm:"id,increment"`
	Name       string `json:"name"`
	SearchName string `json:"search_name" storm:"unique"`
	Password   string `json:"password"`
}

func UserCreate(o *DbUser) (*DbUser, error) {
	o.SearchName = strings.ToLower(o.Name)
	err := db.Save(o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func UserList() ([]DbUser, error) {
	var os []DbUser
	err := db.All(&os)
	if err != nil {
		return nil, err
	}
	return os, nil
}

func UserGet(id int) (*DbUser, error) {
	var o DbUser
	err := db.One("ID", id, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func UserGetByName(username string) (*DbUser, error) {
	var o DbUser
	err := db.One("SearchName", strings.ToLower(username), &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func UserDelete(id int) error {
	o := &DbUser{
		ID: id,
	}
	err := db.DeleteStruct(o)
	if err != nil {
		return err
	}
	return nil
}
