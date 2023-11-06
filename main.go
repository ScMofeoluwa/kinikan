/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	// "kinikan/cmd"
	"kinikan/platform/heroku"
	"kinikan/utils"
)

func main() {
	// cmd.Execute()
	platform := heroku.New()
	serviceImages, _ := utils.ExtractImagesFromCompose("")

	if err := platform.CreateAddOns(serviceImages); err != nil {
		fmt.Printf("error creating add-ons: %v\n", err)
	}
}
