# common
A common library includes some useful struct and utils across all services.
Import with 
```
go get github.com/cross-space-official/common
```
to get the latest version of common library.

## Setup middlewares for all!
Use this snippet in your `main.go`
```
common.SetupAll(engine)
```
Currently, this will set up 
1. the correlation injector 
2. error handling advisor which will convert the failure object to a standard `ErrorResponse`. 

## Configuration
Use the snippet in `main.go` to initialize the viper settings.
```
configuration.Init()
```
The function will first look at the local directory then use vault directory to override it.  

To get the config struct in runtime, e.g. with this sample `config.yaml`
```
data:
    test:
        field_a: "1"
        field_a: "2"
        field_a: "3"
```
Define your own config struct and build it with the following code
```
type TestConfig struct {
	FieldA string `mapstructure:"field_a"`
	FieldB string `mapstructure:"field_b"`
	FieldC string `mapstructure:"field_c"`
}

config := &TestConfig{}
Build("data.test", config)
```

## Logger
A common logger instance has been configured for the cross space microservices. 
Developers can easily obtain the logger with the `gin.Context` as input parameter.  
This context is mandatory to make metadata like `correlation id` logs out properly.  
```
import (
	"github.com/cross-space-official/common/logger"
)

logger.GetLoggerEntry(c).Info("test")
```
The logger entry can be reused.

## Http Client
A configured http client is provided to support customized features among microservices, such as `correlation id` 
propagation.
The embedded http client is `github.com/go-resty/resty/v2`, which is a wrapper over `http.client`, [WORM HOLE](https://github.com/go-resty/resty).  
It supports auto marshaling/unmarshalling and can be configured to be retryable.

A sample code of usage for inter microservices communication as below.
```
// Build a client entity which can be held at service level. This can be reused to execute multiple calls.
client := httpClient.GetHttpClientFactory().Build("http://localhost:8080", AuthPayload{ApiKey: "213"}, nil)

// Body can also pass in the struct, take your own try.
result := &CreateUploadURLsResponse{}
response, err := client.BuildRequest(context).
    SetBody(`{
        "files": [
            {
                "content_type": "image/png",
                "content_md5": "JBHQlZRJ5ov5jjB0+GweDQ==",
                "file_size": 56768
            }
        ]
    }`).
    SetResult(result).
    Post("/files?file_type=avatar")

if err != nil || response.IsError() {
    // Do your stuff
}
```
If the correlation id middleware is injected, it should wire the correlation id to the request header automatically. 

## Failure
Use failures to model the api-level errors.

1. Defines your own business service error or use existing error wrappers to carry stacktrace.
2. At your handler  generate `Failure` from the error.
3. Gin use the middleware.

See `handler_advisor_test.go` for usage. 

## Metrics
Metrics module provides a lazy initialized newrelic application. It reads config from path "data.newrelic".
Enable metrics by add the following middleware.
```	
r.Use(nrgin.Middleware(app))
```

### Report an issue 
Contact sidney@meta-bytes.io for any concerns, comments or report a bug.