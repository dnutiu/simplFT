package server

import (
	"testing"
	"math/rand"
)

func TestMakeStringStack(t *testing.T) {
	st := MakeStringStack(10)

	if st == nil {
		t.Errorf("MakeStringStack returned null!")
	}

	if st.Capacity() != 10 {
		t.Errorf("StringStack: Stack capacity is not ok! Want: 10 Got: %d", st.Capacity())
	}

	if st.IsEmpty() != true {
		t.Errorf("StringStack: Newly created stack is not empty!")
	}
}

func TestStringStack_CanPush(t *testing.T) {
	var st Stack = MakeStringStack(10)

	str := "Hello World"

	st.Push(str)

	if st.Top() != str {
		t.Errorf("StringStack: Push() failed. Want: %s Got: %s", str, st.Top())
	}

	if st.Size() != 1 {
		t.Errorf("StringStack: Size is not correct after one push. Want %d Got %d", 1, st.Size())
	}
}

func TestStringStack_StackOverflows(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("StringStack: Capacity of 0 doesn't overflow on Push()!")
		}
	}()

	st := MakeStringStack(0)
	st.Push(".")
}

func TestStringStack_InvalidType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("StringStack: Push() pushed a non-string type.")
		}
	}()

	st := MakeStringStack(1)
	st.Push(1)
}

func TestStringStack_Capacity(t *testing.T) {
	stCap := rand.Intn(100)
	st := MakeStringStack(stCap)

	if st.Capacity() != stCap {
		t.Errorf("StringStack: Invalid capacity! Want: %d Got: %d", stCap, st.Capacity())
	}
}

func TestStringStack_Size(t *testing.T) {
	pushes := rand.Intn(10)
	st := MakeStringStack(15)

	for i := 0; i < pushes; i++ {
		st.Push("a")
	}

	if st.Size() != pushes {
		t.Errorf("StringStack: Invalid size! Want: %d Got %d", pushes, st.Size())
	}
}

func TestStringStack_IsEmpty(t *testing.T) {
	st := MakeStringStack(10)

	if st.IsEmpty() == false {
		t.Errorf("StringStack: With no push, the stack is not empty!")
	}
}

func TestStringStack_IsEmptyAfterPushAndPop(t *testing.T) {
	st := MakeStringStack(10)
	pushes := rand.Intn(5)

	for i := 0; i < pushes; i++ {
		st.Push("a")
	}

	for i := 0; i < pushes; i++ {
		st.Pop()
	}

	if st.IsEmpty() == false {
		t.Errorf("StringStack: After push and pop, the stack is not empty!")
	}
}

func TestStringStack_Top(t *testing.T) {
	st := MakeStringStack(1)

	st.Push("A")

	if st.Top() != "A" {
		t.Errorf("StringStack: Top() returned invalid value!")
	}
}

func TestStringStack_TopUnderflow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("StringStack: Top() on empty stack didn't underflow.")
		}
	}()

	st := MakeStringStack(1)
	st.Top()
}

func TestStringStack_Pop(t *testing.T) {
	st := MakeStringStack(1)

	st.Push("A")

	if st.Pop() != "A" {
		t.Errorf("StringStack: Pop() returned invalid value!")
	}
}

func TestStringStack_Push(t *testing.T) {
	st := MakeStringStack(1)
	st.Push("A")
}

func TestStringStack_PushAndPop(t *testing.T) {
	st := MakeStringStack(12)
	characters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	pushes := rand.Intn(len(characters))

	for i := 0; i < pushes; i++ {
		st.Push(characters[i])
	}

	for i := pushes - 1; i >= 0; i-- {
		if val := st.Pop(); val != characters[i] {
			t.Errorf("StringStack: Pop() and Push() don't work correctly. Want %s Got %s", characters[i], val)
		}
	}
}
