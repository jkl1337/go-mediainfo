package mediainfo

import (
	"fmt"
)

func ExampleCount() {
	mi := New()
	err := mi.Open("cmd/Example.ogg")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mi.Count(StreamVideo), mi.Count(StreamAudio))
	// Output: 0 1
}

func ExampleGet() {
	mi := New()
	err := mi.Open("cmd/Example.ogg")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mi.Get(StreamGeneral, 0, "Format"))
	// Output: OGG
}
