package main

import (
	"fmt"
	"github.com/mistral-hackathon/triageprof/internal/model"
)

func main() {
	var pb model.ProfileBundle
	var fb model.FindingsBundle
	fmt.Printf("ProfileBundle: %T, FindingsBundle: %T\n", pb, fb)
}