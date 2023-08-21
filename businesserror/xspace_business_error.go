package businesserror

type XSpaceBusinessError interface {
	Error() string
	Stacktrace() string
	Message() string
}
