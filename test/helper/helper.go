package helper

func DoPanicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
