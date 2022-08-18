package main

import "testing"

func Test_main(t *testing.T) {
	t.Skip("Only for debugging")
	t.Setenv("--web.port", "8083")
	main()
}
