// Utility functions used for various purposes in the project

package main

func panicIfError(e error) {
	if e != nil {
		panic(e)
	}
}
