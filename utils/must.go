package utils

func Must(err error) {
	if IsNil(err) {
		return
	}

	panic(err)
}
