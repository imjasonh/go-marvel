package main

import (
	"fmt"
)

func main() {
	r, err := NewClient("d96b5157cfc7a60cbfaa715dc23c3eb1", "ccbc72b222419e2a4e40b4027f3bcb356142651b").Series(2258, CommonRequest{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("%+v", r)
}
