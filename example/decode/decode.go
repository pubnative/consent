package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pubnative/consent"
)

func main() {
	str := "BOQ7WlgOQ7WlgABACDENABwAAABJOACgACAAQABA"
	if len(os.Args) == 2 {
		str = os.Args[1]
	}

	c, err := consent.ParseV1(str)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(c); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
