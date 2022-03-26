package data_test

import "fmt"

func GetCmd20(whereArg string) string {

	return fmt.Sprintf("SELECT name FROM Clients WHERE %s", whereArg)
}

func GetCmd21(whereArg string) string {
	return fmt.Sprintln("SELECT name FROM Clients WHERE s=1")
}

func GetCmd22(whereArg []byte) string {
	return fmt.Sprintf("SELECT name FROM User WHERE id=%s", string(whereArg))
}
