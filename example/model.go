package main

//go:generate goplater model --type UserSomethingNew  --suffix _q --tmpl ./q.tmpl

type UserSomethingNew struct {
	Name     string `db:"name" json:"name"`
	Age      int    `db:"age" json:"age"`
	FullName string `db:"full_name" json:"full_name"`
}
