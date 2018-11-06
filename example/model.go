package main

//go:generate goplater model --type UserNewA  --suffix _q --tmpl ./q.tmpl

type UserNewA struct {
	Name     string `db:"name" json:"name"`
	Age      int    `db:"age" json:"age"`
	FullName string `db:"full_name" json:"full_name"`
}
