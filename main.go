package main
import "fmt"

func main() {
	fmt.Println("Go Interpreter Start！")
	// DebugScanner("let x = 5 ;note \n in -(x,3);note")
	// p := ScanAndParse("let x = 5 ;note \n in -(x,3);note")
	//p := ScanAndParse("if zero?(-(x,11)) " +
	//	                   "then -(y,2) " +
	//	                   "else let x = 5 " +
	//	                          "in -(x,3)")
	pStr := "let x = 7" +
		    "in let y = 2" +
		         "in let y = let x = -(x,1) in -(x,y)" +
		         "in -(-(x,8),y)"
	v := Run(pStr)
	fmt.Print(v.GetPrettyStr())
}
