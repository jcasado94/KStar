package testutils

import (
	"io/ioutil"
	"os"
)

const (
	outExtension = ".out"
)

type testOutput interface {
	New(args ...interface{})
	Name() string
	Marshall() ([]byte, error)
	Unmarshal(bytes []byte) error
	TestFolderPath() string
}

func persist(to testOutput) {
	path := to.TestFolderPath() + to.Name() + outExtension
	bytes, err := to.Marshall()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, bytes, 0666)
}

func read(to testOutput, testName string) {
	bytes, err := ioutil.ReadFile(to.TestFolderPath() + testName + outExtension)
	if err != nil {
		panic(err)
	}
	err = to.Unmarshal(bytes)
	if err != nil {
		panic(err)
	}
}

// ReadTestOutput reads the output file from the file system, storing the result in *to and returning true.
// If the file does not exist, it creates the testOuput instance and returns false.
func ReadTestOutput(to testOutput, testName string, args ...interface{}) bool {
	if _, err := os.Stat(to.TestFolderPath() + testName + outExtension); os.IsNotExist(err) {
		to.New(args...)
		persist(to)
		return false
	}
	read(to, testName)
	return true
}
