package main
import "fmt"

func main() {
	fmt.Println("Go Interpreter Start！")
	// DebugScanner("let x = 5 ;note \n in -(x,3);note")
	// p := ScanAndParse("let x = 5 ;note \n in -(x,3);note")
	p := ScanAndParse("if zero?(-(x,11)) " +
		                   "then -(y,2) " +
		                   "else let x = 5 " +
		                          "in -(x,3)")
	fmt.Print(ExpPrettyStr(p.Exp1, 0))
}
