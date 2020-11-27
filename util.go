package main

func panicIfError(e error) {
	if e != nil {
		panic(e)
	}
}
