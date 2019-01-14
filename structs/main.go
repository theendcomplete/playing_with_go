package main

import "fmt"

type contactInfo struct {
	email   string
	zipCode int
}

type person struct {
	firstName string
	lastName  string
	contactInfo
}

func main() {
	var alex person
	alex.firstName = "Alex"
	alex.lastName = "Anderson"
	// alex := person{firstName: "Alex", lastName: "Anderson"}
	// fmt.Println(alex)
	// fmt.Printf("%+v", alex)
	jim := person{
		firstName: "Jim",
		lastName:  "Morrison",
		contactInfo: contactInfo{
			email:   "jimmy@jim.ru",
			zipCode: 17589,
		},
	}
	// jim.print()
	// jimPointer := &jim
	// fmt.Println(&jim)
	jim.updateName("Hui")
	jim.print()
}

func (p person) print() {
	fmt.Printf("%+v", p)
}
func (pointerToPerson *person) updateName(newFirstName string) {
	(*pointerToPerson).firstName = newFirstName
}
