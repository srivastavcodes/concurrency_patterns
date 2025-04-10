package main

func race2() {
	var data int
	go func() { data++ }()
}
