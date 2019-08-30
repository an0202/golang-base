package main

// "net"
// "time"

type scannedObject struct {
	host      string
	startPort int32
	endPort   int32
}

// func scanPort(startPort int32, endPort int32) {
// 	for i := startPort; i <= endPort; i++ {
// 		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", i), 3*time.Second)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		defer conn.Close()
// 		fmt.Println("access ok")
// 	}
// }

func main() {
	// conn, err := net.DialTimeout("tcp", "127.0.0.1:1080", 3*time.Second)
	// if err != nil {
	// 	fmt.Print(err, "\naaaa")
	// 	return
	// }
	// defer conn.Close()
	// fmt.Println("access ok")
	var scanner *scannedObject
}
