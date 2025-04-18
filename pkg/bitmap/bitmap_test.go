package bitmap

import (
	"fmt"
	"testing"
)

func TestBitmap_Load(t *testing.T) {
	b := NewBitmap(5) //5 * 8 = 40位

	b.Set("aaa")
	fmt.Println(hash("aaa"))
	b.Set("bbb")
	fmt.Println(hash("bbb"))
	b.Set("aaa")
	fmt.Println(hash("aaa"))
	b.Set("121")
	fmt.Println(hash("121"))
	b.Set("122")
	fmt.Println(hash("122"))
	b.Set("123")
	fmt.Println(hash("123"))

	for _, bit := range b.bits {
		t.Logf("%b %v", bit, bit)
	}
	fmt.Println(len(b.bits))
	fmt.Println(len([]byte{1}))
	fmt.Println([]byte{1}[0])

}
