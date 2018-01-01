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

// MakeStringStack initializes a new StringStack pointer.
func MakeStringStack(capacity int) *StringStack {
	st := StringStack{}

	st.capacity = capacity
	st.index = 0
	st.items = make([]string, capacity, capacity)

	return &st
}

// Push pushes an item of type string to the stack.
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

// Pop returns the last pushed item from the stack and removes it
func (st *StringStack) Pop() interface{} {
	if st.Size() == 0 {
		panic(StackUnderflowError)
	}
	st.index--
	return st.items[st.index]
}

// Top return the last pushed item from the stack without removing it.
func (st *StringStack) Top() interface{} {
	if st.Size() == 0 {
		panic(StackUnderflowError)
	}
	return st.items[st.index-1]
}

// IsEmpty returns true if the stack contains no items.
func (st *StringStack) IsEmpty() bool {
	return st.Size() == 0
}

// Capacity returns the maximum items the stack can hold.
func (st *StringStack) Capacity() int {
	return st.capacity
}

// Size returns number of items the stack currently holds.
func (st *StringStack) Size() int {
	return st.index
}

// Items returns an array of the stack items as a copy.
func (st StringStack) Items() []string {
	return st.items
}
