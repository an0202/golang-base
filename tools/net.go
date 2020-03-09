/**
 * @Author: jie.an
 * @Description:
 * @File:  net.go
 * @Version: 1.0.0
 * @Date: 2020/03/04 13:26
 */

package tools

import (
	"fmt"
	"net"
)

//DNS Lookup
func Lookup() {
	cname, err := net.LookupIP("www.a.shifen.com")
	if err != nil {
		fmt.Println(nil)
	}
	fmt.Println(cname)
}