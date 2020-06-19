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
	//p := ScanAndParse("((proc(x) proc(y) -(x,y) 2) 3)")

	//fmt.Print(ExpPrettyStr(p.Exp1,0))
	//pStr := "let x = 7" +
	//	    "in let y = 2" +
	//	         "in let y = let x = -(x,1) in -(x,y)" +
	//	         "in -(-(x,8),y)"
	//pStr := "let x = 200" +
	//	    "in let f = proc(z) -(z,x)" +
	//	       "in let x = 100" +
	//	          "in let g = proc(z) -(z,x)" +
	//	             "in -((f 1), (g 1))"
	//pStr := "letrec double(x) = if zero?(x) " +
	//	"                       then 0 " +
	//	"                       else -((double -(x,1)),2)" +
	//	"    in (double 3)"
	//pStr := "let makemult = proc (maker) " +
	//	"proc (x)" +
	//	"if zero?(x)" +
	//	"then 0" +
	//	"else -(((maker maker) -(x,1)),4)" +
	//	"in let times = proc(x) ((makemult makemult) x)" +
	//	"in (times 3)"

	pStr := "let even = proc (maker)" +
		"proc (x)" +
		"if zero?(x)" +
		"then 1" +
		"else " +
	    	"if zero?(-(x,1)) then 0 " +
		     "else ((maker maker) -(x,2))" +
		"in let event = proc(x) ((even even) x)" +
		"in (event 41)"
	v := Run(pStr)
	fmt.Print(v.GetPrettyStr())
}
