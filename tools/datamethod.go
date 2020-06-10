/**
 * @Author: jie.an
 * @Description:
 * @File:  datamethod.go
 * @Version: 1.0.0
 * @Date: 2019/11/27 10:52
 */
package tools

import (
	"fmt"
	"math/rand"
	"reflect"
)

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func StringFind(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

//https://www.golangprograms.com/remove-duplicate-values-from-slice.html
//func main() {
//	intSlice := []int{1,5,3,6,9,9,4,2,3,1,5}
//	fmt.Println(intSlice)
//	uniqueSlice := unique(intSlice)
//	fmt.Println(uniqueSlice)
//}
func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func UniqueStringList(StringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range StringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//randString
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.!@#$~%^&*(){}|")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//converting a []string to a []interface{}
func StrSliceToIfsSlice(StrSlice []string) (IfsSlice []interface{}) {
	IfsSlice = make([]interface{}, len(StrSlice))
	for i, v := range StrSlice {
		IfsSlice[i] = v
	}
	return IfsSlice
}

//converting a []interface{} to a []string
func IfsSliceToStrSlice(IfsSlice []interface{}) (StrSlice []string) {
	StrSlice = make([]string, len(IfsSlice))
	for i, v := range IfsSlice {
		StrSlice[i] = fmt.Sprint(v)
	}
	return StrSlice
}

// https://gist.github.com/bxcodec/c2a25cfc75f6b21a0492951706bc80b8
func StructToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
