package main

import (
	"PowerfulCLIAppsInGo/todo"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const todoFileNameEnvVar = "TODO_FILENAME"

// Default file name
var todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	// Parsing command line flags
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()
	// Check if the user defined the ENV VAR for a custom file name
	if os.Getenv(todoFileNameEnvVar) != "" {
		todoFileName = os.Getenv(todoFileNameEnvVar)
	}
	// Define an items list
	l := &todo.List{}
	// Use the Get method to read to do items from file
	if err := l.Get(todoFileName); err != nil {
		log.Fatalf("Error while loading the %s file: %v", todoFileName, err)
	}
	// Decide what to do based on the number of arguments provided
	switch {
	// For no extra arguments, print the list
	case *list:
		// List current to do items
		fmt.Print(l)
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			log.Fatalf("Cannot complete the item %d due to: %v", complete, err)
		}
		// Save the new list
		saveList(l, todoFileName)
	case *add:
		// When any arguments (excluding flags) are provided, they will be
		// used as the new task
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			log.Fatalf("%v", err)
		}
		l.Add(t)
		// Save the new list
		saveList(l, todoFileName)
	default:
		log.Fatal("Invalid option")
	}

}

func saveList(l *todo.List, todoFileName string) {
	if err := l.Save(todoFileName); err != nil {
		log.Fatalf("Couldn't save the todo to the file: %v", err)
	}
}

// getTask function decides where to get the description for a new
// task from: arguments or STDIN
func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}

	return s.Text(), nil
}
