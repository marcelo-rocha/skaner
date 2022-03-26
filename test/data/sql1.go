package main

import "fmt"

func GetCmd(whereArg string) string {

	return fmt.Sprintf("SELECT name FROM Clients WHERE %s", whereArg)
}
