/**
 * @Author: jie.an
 * @Description:
 * @File:  excel_var.go
 * @Version: 1.0.0
 * @Date: 2020/02/17 10:55
 */

package excel

import "golang-base/tools"

//HeaderLine Last Position
var LastPosition = map[int]string{
	1: "A1",
	2: "B1",
	3: "C1",
	4: "D1",
	5: "E1",
	6: "F1",
	7: "G1",
	8: "H1",
	9: "I1",
	10: "J1",
	11: "K1",
	12: "L1",
	13: "M1",
	14: "N1",
	15: "O1",
	16: "P1",
	17: "Q1",
	18: "R1",
	19: "S1",
	20: "T1",
	21: "U1",
	22: "V1",
	23: "W1",
	24: "X1",
	25: "Y1",
	26: "Z1",
	27: "AA1",
}

// DescribeLastPosition return last position for headline
func DescribeLastPosition(len int) string {
	if s, ok := LastPosition[len]; ok {
		return s
	}
	tools.ErrorLogger.Fatalln("Length Of HeaderLine Out Of Defined (Length: 27 , Column: AA1)")
	return ""
}
