package main

import (
	"fmt"
	"testing"
)

func TestStuff(t *testing.T) {
	for i:=0;i<100;i++{
		fmt.Println(plusMinusOne())
	}
}
