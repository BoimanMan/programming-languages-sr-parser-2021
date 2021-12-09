package main

/* Implements datatypes:
pstackItem           An item on the parse stack for an SR parser
parseStack           A parse stack for an SR parser.
*/

/* Stack operations implemented:
1. function newParseStack()         Returns new empty parse stack.
2. method push(itm pstackItem)      Pushes item on stack.
3. method top()                     Returns the top of the stack. No side-effect.
4. method pop()                     Pops the stack and returns the popped item.
5. method popnum(n int)             Pops n items from stack. No return value.
6. method String()                  Returns a string representation of the stack.

NOTE: All of the methods and datatypes listed above were provided by Dr. Anthony Maida at the University of Louisiana at Lafayette.
All other methods not listed above were created by me using Dr. Maida's existing datatypes and methods as a reference.
*/

import (
	"container/list"
)

/* pstackItem
==============*/
type pstackItem struct {
	grammarSym, stateSym string
}

/* String()
Implement a String() method with this exact signature.
The print statements in the "fmt" package will understand
this and use it to print instances of this datatype.
*/
func (se pstackItem) String() string {
	return se.grammarSym + se.stateSym
}

/* parseStack: list implementation
==================================*/
type parseStack struct {
	stack *list.List
}

/* inputQueue ***NEW***: list implementation*/
type inputQueue struct {
	queue *list.List
}

/*********************************/

/* treeNode ***NEW***: For Parse Tree*/
type treeNode struct {
	parent   string //LHSsym
	termsym  string //terminal symbol
	children *treeNodeQueue
}
type treeNodeQueue struct {
	queue *list.List
}
type treeNodeStack struct {
	stack *list.List
}

/**********************************/
/* newParseStack()
Creates and returns a new empty stack
*/
func newParseStack() parseStack {
	ps := parseStack{}
	ps.stack = list.New()
	return ps
}

/* newInputQueue() ***NEW***
Creates and returns a new input queue
*/
func newInputQueue() inputQueue {
	iq := inputQueue{}
	iq.queue = list.New()
	return iq
}

/*********************************/

/* newTreeStack() **NEW** */
func newTreeStack() treeNodeStack {
	ts := treeNodeStack{}
	ts.stack = list.New()
	return ts
}

/**************************/

func newTreeQueue() treeNodeQueue {
	tq := treeNodeQueue{}
	tq.queue = list.New()
	return tq
}

/***************************/
func (stk parseStack) push(itm pstackItem) {
	stk.stack.PushFront(itm)
}

/***NEW***/
func (que inputQueue) enqueue(tok string) {
	que.queue.PushBack(tok)
}

//We want to push back instead of front since this is a queue.
/*********************************/

func (que treeNodeQueue) push(nd treeNode) {
	que.queue.PushBack(nd)
}
func (que treeNodeQueue) topPush(nd treeNode) {
	que.queue.PushFront(nd)
}

func (stk treeNodeStack) push(que treeNodeQueue) {
	stk.stack.PushFront(que)
}

/* top()
Returns top of the stack. No side-effect.
*/
func (stk parseStack) top() pstackItem {
	e := stk.stack.Front()
	return e.Value.(pstackItem)
}

/***NEW***/
func (que inputQueue) next() string {
	top := que.queue.Front()
	return top.Value.(string)
}

/*********************************/

func (stk parseStack) pop() pstackItem {
	e := stk.stack.Front()
	if e != nil {
		stk.stack.Remove(e)
		return e.Value.(pstackItem)
	}
	return pstackItem{"", ""}
}

/***NEW***/
func (que inputQueue) dequeue() string {
	e := que.queue.Front()
	if e != nil {
		que.queue.Remove(e)
		return e.Value.(string)
	}
	return ""
}

/***************************************/
func (stk treeNodeStack) pop() treeNodeQueue {
	e := stk.stack.Front()
	if e != nil {
		stk.stack.Remove(e)
		return e.Value.(treeNodeQueue)
	}
	return treeNodeQueue{}
}

/*********************************/

func (que treeNodeQueue) pop() treeNode {
	e := que.queue.Front()
	if e != nil {
		que.queue.Remove(e)
		return e.Value.(treeNode)
	}
	return treeNode{}
}
func (stk parseStack) popNum(n int) {
	for i := 1; i <= n; i++ {
		stk.pop()
	}
}

//Should not need an equivalent to popNum for the input queue.
//Input queue only lets one token go at a time.
func (stk parseStack) len() int {
	return stk.stack.Len()
}

/***NEW***/
func (que inputQueue) len() int {
	return que.queue.Len()
}

func (que treeNodeQueue) len() int {
	return que.queue.Len()
}

func (stk treeNodeStack) len() int {
	return stk.stack.Len()
}

//Just in case.
/**********************************/

func (stk parseStack) String() string {
	str := ""
	if stk.stack.Len() > 0 {
		e := stk.stack.Back()
		str = e.Value.(pstackItem).String()
		for e.Prev() != nil {
			str += e.Prev().Value.(pstackItem).String()
			e = e.Prev()
		}
	}
	return str
}

func (que inputQueue) String() string {
	str := ""
	if que.queue.Len() > 0 {
		e := que.queue.Front()
		str = e.Value.(string)
		for e.Next() != nil {
			str += e.Next().Value.(string)
			e = e.Next()
		}
	}
	return str
}

//Converts ONE tree node and it's children to string
func (node treeNode) NodeString(e *list.Element) string {
	str := ""
	switch e.Value.(treeNode).parent {
	case "":
		str = e.Value.(treeNode).termsym
	default:
		str = "[" + e.Value.(treeNode).parent + e.Value.(treeNode).children.String() + "]"
	}
	return str
}
func (que treeNodeQueue) String() string {
	str := ""
	if que.queue.Len() > 0 {
		e := que.queue.Front()
		str = e.Value.(treeNode).NodeString(e)
		for e.Next() != nil {
			str += e.Next().Value.(treeNode).NodeString(e.Next())
			e = e.Next()
		}
	}
	return str
}
func (stk treeNodeStack) String() string {
	str := ""
	if stk.stack.Len() > 0 {
		e := stk.stack.Front()
		str = e.Value.(treeNodeQueue).String()
		for e.Next() != nil {
			str += e.Next().Value.(treeNodeQueue).String()
			e = e.Next()
		}
	}
	return str
}

/*********************************************/
