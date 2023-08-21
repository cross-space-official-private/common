package errorhandle

import (
	"errors"
	"fmt"
	"github.com/joomcode/errorx"
	"log"
)

var testNamespace = errorx.NewNamespace("foo")
var testType = testNamespace.NewType("bar")
var testProperty0 = errorx.RegisterProperty("someKey")

func foofoo() error {
	return testType.Wrap(errors.New("Some Error"), "Additional Error Msg").WithProperty(testProperty0, "SomePropertyValue")
}

func ensureCallstackSample() *errorx.Error {
	err := errors.New("some error")
	return errorx.EnsureStackTrace(err)
}

func bar() {
	err := foofoo()
	if err != nil {
		panic(err) // Here convert err to a panic
	}
}

func Test1WithError() {
	// Try with error
	err, result := Try(func() int {
		bar()
		return 3
	}, func(err *errorx.Error) (bool, string) {
		if errorx.IsOfType(err, testType) {
			// handle
			log.Printf("Error: %+v", err)
		}

		value, _ := err.Property(testProperty0)
		return true, "SomeError. Property:" + value.(string)
	})

	fmt.Printf("Err: %+v, result: %+v", err, result)
}

func Test2WithNoError() {
	// Try with error
	err, result := Try(func() int {
		return 3
	}, func(err *errorx.Error) (bool, string) {
		if errorx.IsOfType(err, testType) {
			// handle
			log.Printf("Error: %+v", err)
		}
		return true, "SomeError"
	})

	fmt.Printf("Err: %+v, result: %+v", err, result)
}

func Test3WithUnrelatedPanic() {
	// Try with error
	err, result := Try(func() int {
		panic("Unrelated error")
		return 3
	}, func(err *errorx.Error) (bool, string) {
		if errorx.IsOfType(err, testType) {
			// handle
			log.Printf("Error: %+v", err)
		}
		return true, "SomeError"
	})

	fmt.Printf("Err: %+v, result: %+v", err, result)
}

func Test4WithPanicRethrow() {
	// Try with error
	err, result := Try(func() int {
		bar()
		return 3
	}, func(err *errorx.Error) (bool, string) {
		return false, "SomeError"
	})

	fmt.Printf("Err: %+v, result: %+v", err, result)
}
