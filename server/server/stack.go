package server

// Stack interface
type Stack interface {
	// Push pushes an item to the stack. Returns an error if it fails.
	Push(item interface{})
	// Pop retrieves an item from the stack and removes it.
	Pop() interface{}
	// Top peeks at an item from the stack.
	Top() interface{}
	// IsEmpty returns bool indicating if the stack is empty.
	IsEmpty() bool
	// Capacity returns the capacity of the stack.
	Capacity() int
	// Size returns the size of the stack
	Size() int
}

// StringStack is a stack that holds string objects
type StringStack struct {
	index    int
	capacity int
	items    []string
}

func MakeStringStack(capacity int) *StringStack {
	st := StringStack{}

	st.capacity = capacity
	st.index = 0
	st.items = make([]string, capacity, capacity)

	return &st
}

func (st *StringStack) Push(item interface{}) {
	if st.index == st.Capacity() {
		panic(StackOverflowError)
	}

	value, ok := item.(string)
	if ok == false {
		panic(StackInvalidTypeError)
	}

	st.items[st.index] = value
	st.index++
}

func (st *StringStack) Pop() interface{} {
	if st.Size() == 0 {
		panic(StackUnderflowError)
	}
	st.index--
	return st.items[st.index]
}

func (st *StringStack) Top() interface{} {
	if st.Size() == 0 {
		panic(StackUnderflowError)
	}
	return st.items[st.index-1]
}

func (st *StringStack) IsEmpty() bool {
	return st.Size() == 0
}

func (st *StringStack) Capacity() int {
	return st.capacity
}

func (st *StringStack) Size() int {
	return st.index
}
