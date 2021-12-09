package main

import (
	"fmt"
	"strconv"
)

var grammatical bool = true
var finished bool = false
var initialized bool = false
var inputArray = []string{"id", "+", "id", "*", "id"}
var parStack parseStack
var inQueue inputQueue
var treeStack treeNodeStack
//The grammar provided by Dr. Anthony Maida at the University of Louisiana at Lafayette.
//The assignment was to construct a functioning SR Parser that used this grammar in tandem with the action table and goto table.
//As such, aTable and gTable were also provided by Dr. Maida.
var grammar = [6][]string{
	{"E", "->", "E", "+", "T"},
	{"E", "->", "T"},
	{"T", "->", "T", "*", "F"},
	{"T", "->", "F"},
	{"F", "->", "(", "E", ")"},
	{"F", "->", "id"}}

var aTable = [12][6]string{ // action table
	{"S5", "", "", "S4", "", ""},     // 0
	{"", "S6", "", "", "", "accept"}, // 1
	{"", "R2", "S7", "", "R2", "R2"}, // 2
	{"", "R4", "R4", "", "R4", "R4"}, // 3
	{"S5", "", "", "S4", "", ""},     // 4
	{"", "R6", "R6", "", "R6", "R6"}, // 5
	{"S5", "", "", "S4", "", ""},     // 6
	{"S5", "", "", "S4", "", ""},     // 7
	{"", "S6", "", "", "S11", ""},    // 8
	{"", "R1", "S7", "", "R1", "R1"}, // 9
	{"", "R3", "R3", "", "R3", "R3"}, // 10
	{"", "R5", "R5", "", "R5", "R5"}, // 11
}

var gTable = [12][3]string{
	{"1", "2", "3"}, // 0
	{"", "", ""},    // 1
	{"", "", ""},    // 2
	{"", "", ""},    // 3
	{"8", "2", "3"}, // 4
	{"", "", ""},    // 5
	{"", "9", "3"},  // 6
	{"", "", "10"},  // 7
	{"", "", ""},    // 8
	{"", "", ""},    // 9
	{"", "", ""},    // 10
	{"", "", ""},    // 11
}

func parse1step() {
	//put inits before switch
	if !initialized {

		parStack = newParseStack()
		inQueue = newInputQueue()
		treeStack = newTreeStack()

		parStack.push(pstackItem{"", "0"})
		for i := 0; i < 5; i++ {
			inQueue.enqueue(inputArray[i])
		}
		inQueue.enqueue("$") //End of input symbol
		initialized = true
	}
	//Print Parse Stack
	fmt.Printf("%-17s", parStack.String())

	//Print Input Queue
	fmt.Printf("%-17s", inQueue.String())

	var column int
	switch inQueue.next() {
	case "id":
		column = 0
	case "+":
		column = 1
	case "*":
		column = 2
	case "(":
		column = 3
	case ")":
		column = 4
	case "$":
		column = 5
	}
	/*switch parStack.top().stateSym {
	case "0":
		row = 0
	case "1":
		row = 1
	case "2":
		row = 2
	case "3":
		row = 3
	case "4":
		row = 4
	case "5":
		row = 5
	case "6":
		row = 6
	case "7":
		row = 7
	case "8":
		row = 8
	case "9":
		row = 9
	case "10":
		row = 10
	case "11":
		row = 11
	}*/

	row := parStack.top().stateSym
	//These two statements effectively replace the commented switch statement above.
	//Print Action Lookup
	actLook := ("[" + row + "," + inQueue.next() + "]")
	fmt.Printf("%-17s", actLook)

	choice := "ungrammatical"
	rowNum, err := strconv.Atoi(row)
	if err != nil {
		fmt.Println("ERROR")
	}
	action := aTable[rowNum][column]
	switch string(action[0]) {
	case "S":
		choice = "shift"
		//Print Action Value -needs adjustment
		fmt.Printf("%-98s", action)
	case "R":
		choice = "reduce"
		//Print Action Value
		fmt.Printf("%-17s", action)
	case "a":
		choice = "accept"
		//Print Action Value -needs adjustment
		fmt.Printf("%-109s", action)
	}

	//At this point, we should have our initial parse stack and input queue.
	switch choice {

	//On accept, we are done and grammatical = true.
	case "accept":
		finished = true

	//On error, we are done and grammatical = false.
	case "ungrammatical":
		grammatical = false
		finished = true

	//On shift, push token to parse stack in pair with action lookup number. Dequeue token.
	case "shift":
		pushing := pstackItem{inQueue.next(), string(action[1])}

		//Print Stack Action
		fmt.Printf("%-11s", "push "+pushing.String())

		parStack.push(pushing)
		//If id, push onto parse tree stack
		if inQueue.next() == "id" {
			//init queue
			subtree := newTreeQueue()
			//push node to queue
			subtree.push(treeNode{"", "id", nil})
			//push new queue to stack
			treeStack.push(subtree)
		}
		inQueue.dequeue()

	//On reduce, call popnum, then assign goto based on remaining stack. Push LHS and goto val.
	case "reduce":
		rule := action[1] - 1 - '0'
		lhs := string(grammar[rule][0]) //put into pair as X

		//Print LHS
		fmt.Printf("%-17s", lhs)

		lenOfRhs := len(grammar[rule]) - 2
		//Print len(RHS)
		fmt.Printf("%-17d", lenOfRhs)

		switch lenOfRhs {
		case 1:
			parStack.popNum(lenOfRhs)
			//pop queue on top of stack
			temp := treeStack.pop()
			var tempPtr *treeNodeQueue = &temp
			//push new node to queue
			newTree := newTreeQueue()
			newTree.push(treeNode{lhs, "", tempPtr})
			//push queue back to stack
			treeStack.push(newTree)
		case 3:
			var op string
			for i := len(parStack.String()) - 1; i > len(parStack.String())-7; i-- {
				switch string(parStack.String()[i]) {
				case "+":
					op = "+"
				case "*":
					op = "*"
				case "(":
					op = "("
				case ")":
					op = ")"
				}
			}
			parStack.popNum(lenOfRhs)
			//pop two items from parse stack
			que1 := treeStack.pop()
			que2 := treeStack.pop()
			//insert operator
			opNode := treeNode{"", op, nil}
			//order: 1st queue nodes -> opnode -> 2nd queue nodes
			//push new subtree to tree stack
			newQueue := newTreeQueue()
			for i := 0; i < que2.len(); i++ {
				newNode := que2.pop()
				newQueue.push(newNode)
			}
			newQueue.push(opNode)
			for i := 0; i < que1.len(); i++ {
				newNode := que1.pop()
				newQueue.push(newNode)
			}
			treeStack.push(newQueue)
			//pop queue on top of stack
			temp := treeStack.pop()
			var tempPtr *treeNodeQueue = &temp
			//push new node to queue
			newTree := newTreeQueue()
			newTree.push(treeNode{lhs, "", tempPtr})
			//push queue back to stack
			treeStack.push(newTree)
		}

		//Print Temp Stack
		fmt.Printf("%-17s", parStack.String())

		lookupStr := parStack.top().stateSym
		lookup, err := strconv.Atoi(lookupStr)
		if err != nil {
			fmt.Print("ERROR")
		}
		//Print Goto Lookup
		gotoLook := ("[" + lookupStr + "," + lhs + "]")
		fmt.Printf("%-13s", gotoLook)
		//using strconv package here to avoid a long switch statement
		var gotoVal string

		//Row based on lookup, column based on lhs
		switch lhs {
		case "E":
			gotoVal = gTable[lookup][0]
		case "T":
			gotoVal = gTable[lookup][1]
		case "F":
			gotoVal = gTable[lookup][2]
		}

		//Print Goto Value
		fmt.Printf("%-17s", gotoVal)
		//Simply push lhs:gotoVal pair to parStack.
		parStack.push(pstackItem{lhs, gotoVal})
		fmt.Printf("%-11s", "push "+parStack.top().String())

	}
	//Finally, print parse tree stack
	fmt.Println(treeStack.String())
}
//main function just formats the data and presents it in a readable way.
func main() {
	//fmt.Println("Stack		Input		Act.Lookup		Act.Val		LHS		RHS length		temp stack		goto lookup		goto value		stack action")
	//fmt.Println("_______________________________________________________________________________________________________________________________________")
	fmt.Printf("%-17s", "Stack")
	fmt.Printf("%-17s", "Input")
	fmt.Printf("%-17s", "Act.Lookup")
	fmt.Printf("%-17s", "Act.Val")
	fmt.Printf("%-17s", "LHS")
	fmt.Printf("%-17s", "RHS Length")
	fmt.Printf("%-17s", "Temp Stack")
	fmt.Printf("%-13s", "Goto Lookup")
	fmt.Printf("%-17s", "Goto Value")
	//fmt.Println("Stack Action")
	fmt.Printf("%-15s", "Stack Action")
	fmt.Println("Parse Tree Stack")
	fmt.Println("______________________________________________________________________________________________________________________________________________________________________________________")
	for !finished {
		parse1step()
	}
	switch grammatical {
	case true:
		indCount := 1
		opCount := 0
		var indention int
		for i := 0; i < len(treeStack.String()); i++ {
			switch string(treeStack.String()[i]) {
			case "[":
			case "]":
			case "i":
				fmt.Print(string(treeStack.String()[i]))
			case "d":
				fmt.Println(string(treeStack.String()[i]))
				opCount++
				indCount = opCount
				//indention := opCount * 5
				//fmt.Printf("%-"+strconv.Itoa(indention)+"s", "")
			default:
				switch indCount == opCount {
				case true:
					indention = indCount * 5
					fmt.Printf("%-"+strconv.Itoa(indention)+"s", "")
					fmt.Println(string(treeStack.String()[i]))
					fmt.Printf("%-"+strconv.Itoa(indention)+"s", "")
					indCount++
				case false:
					fmt.Println(string(treeStack.String()[i]))
					indention = indCount * 5
					fmt.Printf("%-"+strconv.Itoa(indention)+"s", "")
					indCount++
				}

			}
		}
	case false:
		fmt.Println("This is ungrammatical.")
	}
	//if !grammatical {
	//fmt.Println("This expression is ungrammatical.")
	//}
	//si0 := pstackItem{"", "S0"} // random stack entries (stack item)
	//si1 := pstackItem{"E", "5"}
	//si2 := pstackItem{"T", "6"}
	//pstack := newParseStack() // new empty stack
	//pstack.push(si0)
	//pstack.push(si1)
	//pstack.push(si2)
	//fmt.Println(pstack.top()) // prints "T6", no side-effect
	//fmt.Println(pstack.pop()) // prints "T6", and pops stack
	//fmt.Println(pstack.len()) // prints "2", because one item has been popped
	//fmt.Println(pstack)       // prints "S0E5", no spaces b/c goes into a table
}
