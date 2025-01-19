package main

import (
	"database/db" // import the db package
	"fmt"         // format package for formatted I/O
	"os"          // os package for platform-independent interface to operating system functionality
)

func showHelp() {
	fmt.Println("Usage: go run main.go <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  create <database_name> 						 - Create a new database")
	fmt.Println("  list 									- Not implemented yet")
	fmt.Println("  delete <database_name> 						- delete a database")
	fmt.Println("  query  <database_name> <query> 						- Not implemented yet")
	fmt.Println("  addtable <database_name> <table> <columns>							- Not implemented yet")
	fmt.Println("  insert <db> <table> <data> 								- Not implemented yet")
	fmt.Println("  debug <database_name>  										- Read and debug a database file")
	fmt.Println("  help                    										- Show this help message")
}

func main() {
	// Check if a command was provided and argument len is greater than 2
	if len(os.Args) < 2 {
		fmt.Println("Error: No command provided.")
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	//Create a new database
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Error: No database name provided.")
			showHelp()
			return
		}
		name := os.Args[2]
		dbInstance := db.NewDatabase(name)
		if err := dbInstance.SaveToFile(); err != nil {
			fmt.Printf("Error saving database: %v\n", err)
			return
		}
		fmt.Printf("Database '%s' created successfully in 'cdatabases' directory\n", dbInstance.Name)

	// debug hexdump of database file
	case "debug":
		if len(os.Args) < 3 {
			fmt.Println("Error: No database name provided.")
			showHelp()
			return
		}
		filePath := "cdatabases/" + os.Args[2] + ".db"
		if err := db.ReadFromFile(filePath); err != nil { // call ReadFromFile is a function in db/database.go , nil for error handling
			fmt.Printf("Error reading database: %v\n", err)
		}
	case "help":
		showHelp()
	default:
		fmt.Printf("Error: Unknown command '%s'.\n", command)
		showHelp()
	}
}
