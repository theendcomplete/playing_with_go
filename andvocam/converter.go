package main

import (
	"fmt"
	"log"
	"os/exec"
)

// Конвертирует файл с использованием проприетарного по
func convertFile(dirname string) {
	fmt.Println(dirname + "App_E_Dog.exe")
	cmd := exec.Command("powershell", dirname+"App_E_Dog.exe")
	cmd.Dir = dirname
	log.Printf("Running command and waiting for it to finish...")
	error := cmd.Run()
	if error != nil {
		fmt.Println("Error launching:", error.Error())
		log.Printf("Command finished with error: %v", error)
	}
}
