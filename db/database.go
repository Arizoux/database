// handles the database functionality. The NewDatabase function creates a new Database instance, and the SaveToFile method saves the database as a binary file in the cdatabases directory. The ReadFromFile function reads the database from a binary file for debugging purposes.
// more functoinality's pending

package db

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Database struct to represent a database
type Database struct {
	Name string
}

// NewDatabase creates a new Database instance
func NewDatabase(name string) *Database {
	return &Database{Name: name}
}

// SaveToFile saves the database as a binary file in the 'cdatabases' directory
func (db *Database) SaveToFile() error {
	dir := "cdatabases"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' does not exist", dir)
	}

	filePath := dir + "/" + db.Name + ".db"
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("database file '%s' already exists", filePath)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create database file: %w", err)
	}
	defer file.Close()

	// Write a simple binary header: Magic number + Database name length + Database name
	// https://en.wikipedia.org/wiki/List_of_file_signatures .db file signature maybe for cross useability in future
	magicNumber := []byte("MYDB\000")
	if _, err = file.Write(magicNumber); err != nil {
		return fmt.Errorf("failed to write magic number: %w", err)
	}

	// Write the length of the database name
	nameLength := uint8(len(db.Name))
	if err = binary.Write(file, binary.LittleEndian, nameLength); err != nil {
		return fmt.Errorf("failed to write name length: %w", err)
	}

	// Write the database name
	if _, err = file.Write([]byte(db.Name)); err != nil {
		return fmt.Errorf("failed to write database name: %w", err)
	}

	return nil
}

// ReadFromFile reads the database from a binary file for debugging
func ReadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read and verify the magic number
	magic := make([]byte, 5)
	if _, err := file.Read(magic); err != nil {
		return fmt.Errorf("failed to read magic number: %w", err)
	}
	if string(magic) != "MYDB\000" {
		return fmt.Errorf("invalid database file format")
	}
	fmt.Printf("Magic Number: %s\n", string(magic))

	// Read database name length
	var nameLength uint8
	if err := binary.Read(file, binary.LittleEndian, &nameLength); err != nil {
		return fmt.Errorf("failed to read name length: %w", err)
	}
	fmt.Printf("Database Name Length: %d\n", nameLength)

	// Read database name
	name := make([]byte, nameLength)
	if _, err := file.Read(name); err != nil {
		return fmt.Errorf("failed to read database name: %w", err)
	}
	fmt.Printf("Database Name: %s\n", string(name))

	return nil
}
