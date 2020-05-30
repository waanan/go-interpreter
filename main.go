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
	//pStr := "let x = 7" +
	//	    "in let y = 2" +
	//	         "in let y = let x = -(x,1) in -(x,y)" +
	//	         "in minus(-(-(x,8),y))"
	//pStr := "minus(-(minus(5),9))"
	//pStr := "equal?(3,4)"
	//pStr :=  "let x=4 in null?(cdr( cons(x, cons(cons(-(x,1),emptylist),emptylist )) ))"
	//pStr := "let x=4 in list(x,-(x,1),-(x,3))"
	//pStr := "let x=4 in cond {zero?(-(x,1)) > 5}{zero?(-(x,4)) > 6}end"
	//pStr := "if 3 then 2 else 1"
	//pStr := "let x=30 in let x=-(x,1) y=-(x,2) in -(x,y)"
	pStr := "let u=7  in unpack x y = cons (u,cons(3,emptylist)) in -(x,y)"
	v := Run(pStr)
	fmt.Print(v.GetPrettyStr())
}
