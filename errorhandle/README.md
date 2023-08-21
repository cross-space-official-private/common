# Usage
## Errors with callstack
```go
testType.Wrap(errors.New("Some Error"), "Additional Error Msg").WithProperty(testProperty0, "SomePropertyValue")
```

Can use wrap or decorate or EnsureStack to package a go standard error to a errorx error with call stack information

## Try catch

See four examples for usage of the try catch library.

The return value of a try block will always be the err, result. Err is the converted error after handling when the error is thrown and handled. result is the normal returned result if no errors have happened.

The handler's parameter will be the panic type that will be caught and if type is not matched the panic will propogated out of the try block. 
If a type has a match and you want to rethrow, just return false in the handler and the panic will be rethrown.

The callstack will be perserved for all the panics.

When combine with errorx, this will be powerful enough to get details all the way down to the source.

