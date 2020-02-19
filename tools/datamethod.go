/**
 * @Author: jie.an
 * @Description:
 * @File:  datamethod.go
 * @Version: 1.0.0
 * @Date: 2019/11/27 10:52
 */
package tools

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

//todo sort and remove duplicate
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
