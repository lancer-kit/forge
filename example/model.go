package main

type User struct {
	Name     string `db:"name" json:"name"`
	Age      int    `db:"age" json:"age"`
	FullName string `db:"full_name" json:"full_name"`
}
