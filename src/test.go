package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Users struct which contains
// an array of users
type records struct {
	Users []User `json:"users"`
}

// User struct which contains a name
// a type and a list of social links
type User struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Age    int    `json:"Age"`
	Social Social `json:"social"`
}

// Social struct which contains a
// list of links
type Social struct {
	Facebook string `json:"facebook"`
	Twitter  string `json:"twitter"`
}

func main() {

	// read file
	data, err := ioutil.ReadFile("user.json")
	if err != nil {
		fmt.Print(err)
	}

	var obj records
	// fmt.Printf("%s", data)
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}

	for _, v := range obj.Users {
		fmt.Println(v)
		u := User{}
		user, _ := json.Marshal(u)
		fmt.Printf("%s\n", user)
	}

	// // define data structure
	// type DayPrice struct {
	// 	USD float32
	// 	EUR float32
	// 	GBP float32
	// }

	// // json data
	// var obj DayPrice

	// // unmarshall it
	// err = json.Unmarshal(data, &obj)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }

	// // can access using struct now
	// fmt.Printf("USD : %.2f\n", obj.USD)
	// fmt.Printf("EUR : %.2f\n", obj.EUR)
	// fmt.Printf("GBP : %.2f\n", obj.GBP)

}
