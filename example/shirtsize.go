// Copyright 2017 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

//go:generate goplater -type=ShirtSize

type ShirtSize byte

const (
	NA ShirtSize = iota
	XS
	S
	M
	L
	XL
)

//go:generate goplater -type=WeekDay

type WeekDay int

const (
	Monday WeekDay = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

var defWeekDayValueToName = map[WeekDay]string{
	Monday:    "Dilluns",
	Tuesday:   "Dimarts",
	Wednesday: "Dimecres",
	Thursday:  "Dijous",
	Friday:    "Divendres",
	Saturday:  "Dissabte",
	Sunday:    "Diumenge",
}

var defShirtSizeValueToName = map[ShirtSize]string{
	NA: "NA",
	XS: "XS",
	S:  "S",
	M:  "M",
	L:  "L",
	XL: "XL",
}

func main() {
	v := struct {
		Size ShirtSize
		Day  WeekDay
	}{M, Friday}
	if err := json.NewEncoder(os.Stdout).Encode(v); err != nil {
		log.Fatal(err)
	}

	input := `{"Size":"XL", "Day":"Dimarts"}`
	if err := json.NewDecoder(strings.NewReader(input)).Decode(&v); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decoded %s as %+v\n", input, v)
}
