package test2

import "fmt"

func MapTest() {

	var map1 = make(map[string]string)

	map1["aaa"] = "111"
	map1["bbb"] = "222"
	map1["aaa"] = "333"

	for s := range map1 {
		fmt.Printf("%s => %s \n", s, map1[s])
	}
}
