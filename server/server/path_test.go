package server

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestMakePathFromStringStack(t *testing.T) {
	st := MakeStringStack(5)

	st.Push("first")
	st.Push("folder two")
	st.Push("trinity")

	path := MakePathFromStringStack(st)
	expected := fmt.Sprintf("%s%s/%s/%s/", BasePath, "first", "folder two", "trinity")

	if path != expected {
		t.Errorf("TestMakePathFromStringStack: Returned an invalid path! Want %s Got %s", expected, path)
	}

}

func TestMakePathFromStringStack_StackIsEmpty(t *testing.T) {
	st := MakeStringStack(1)
	path := MakePathFromStringStack(st)

	expected := BasePath
	if BasePath[len(BasePath)-1] != '/' {
		expected += "/"
	}

	if path != expected {
		t.Errorf("TestMakePathFromStringStack: Returned an invalid path!")
	}
}

func TestChangeDirectory(t *testing.T) {
	st := MakeStringStack(1)
	err := os.Chdir(BasePath)
	if err != nil {
		t.Errorf("TestChangeDirectory: Step1: %s", err.Error())
	}

	dirName := string(rand.Intn(10000))

	err = os.Mkdir(dirName, os.FileMode(0522))
	if err != nil {
		t.Errorf("TestChangeDirectory: Step2: %s", err.Error())
	}

	err = ChangeDirectory(st, dirName)
	if err != nil {
		t.Errorf("TestChangeDirectory: Step3: %s", err.Error())
	}

	err = os.Chdir(BasePath)
	if err != nil {
		t.Errorf("TestChangeDirectory: Step4: %s", err.Error())
	}

	err = os.Remove(dirName)
	if err != nil {
		t.Errorf("TestChangeDirectory: Step5: %s", err.Error())
	}
}

func TestChangeDirectory_InvalidDirectoryIsNotInStack(t *testing.T) {
	st := MakeStringStack(1)
	err := os.Chdir(BasePath)
	if err != nil {
		t.Errorf("TestChangeDirectory_InvalidDirectoryIsNotInStack: Step1: %s", err.Error())
	}

	dirName := string(rand.Intn(10000))

	// ignore no such directory error
	ChangeDirectory(st, dirName)
	if !st.IsEmpty() && st.Top() == dirName {
		t.Errorf("TestChangeDirectory: Stack is corrupted because invalid directory remained in stack")
	}

}

func TestChangeDirectory_InvalidDirectoryName(t *testing.T) {
	st := MakeStringStack(1)
	err := ChangeDirectory(st, "some/not/cool/directory/")
	if err != ErrInvalidDirectoryName {
		t.Error("TestChangeDirectory: Changed directory to something containing the '/' character!")
	}
}

func TestChangeDirectoryToPrevious_StackIsEmpty(t *testing.T) {
	st := MakeStringStack(1)
	err := ChangeDirectoryToPrevious(st)
	if err != ErrAlreadyAtBaseDirectory {
		t.Error("TestChangeDirectoryToPrevious: Stack is empty and no error occured!")
	}
}
