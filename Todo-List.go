package main

import (
	"encoding/json" // Für das Umwandeln zwischen JSON und Go-Datenstrukturen
	"fmt"           // Zum Ausgeben von Text in der Konsole
	"io/ioutil"     // Zum Lesen und Schreiben von Dateien
	"os"            // Zum Arbeiten mit Dateisystemen und Argumenten
	"strconv"       // Zum Umwandeln von Text in Zahlen
)

// Struktur für eine Aufgabe (Todo)
type Todo struct {
	ID        int    `json:"id"`        // Eine eindeutige ID für die Aufgabe
	Task      string `json:"task"`      // Die Beschreibung der Aufgabe
	Completed bool   `json:"completed"` // Ob die Aufgabe erledigt ist
}

func main() {
	// Hier lesen wir die Befehle aus, die nach dem Programmnamen eingegeben wurden
	args := os.Args[1:]

	// Wenn keine Befehle angegeben wurden, zeigen wir eine Anleitung
	if len(args) == 0 {
		printUsage()
		return
	}

	// Der erste Befehl entscheidet, was wir tun wollen
	command := args[0]

	switch command {
	case "add":
		// Wir brauchen mindestens 2 Argumente: den Befehl und die Aufgabe
		if len(args) < 2 {
			fmt.Println("Bitte eine Aufgabe eingeben, die hinzugefügt werden soll.")
			return
		}
		task := args[1]
		addTodoTask(task)
	case "update":
		// Hier brauchen wir die ID der Aufgabe und den neuen Text
		if len(args) < 3 {
			fmt.Println("Bitte eine ID und die neue Aufgabe eingeben.")
			return
		}
		// Die ID wird von Text in eine Zahl umgewandelt
		id, _ := strconv.Atoi(args[1])
		task := args[2]
		updateTodoTask(id, task)
	case "delete":
		// Zum Löschen brauchen wir nur die ID der Aufgabe
		if len(args) < 2 {
			fmt.Println("Bitte eine ID für die zu löschende Aufgabe eingeben.")
			return
		}
		id, _ := strconv.Atoi(args[1])
		deleteTodoTask(id)
	case "view":
		// Zeige die Liste der Aufgaben an
		viewTodoList()
	default:
		// Zeige die Anleitung, wenn der Befehl ungültig ist
		fmt.Println("Ungültiger Befehl. Verwende 'view', 'add', 'update' oder 'delete'.")
	}
}

// Funktion zum Hinzufügen einer neuen Aufgabe
func addTodoTask(task string) {
	// Lade die existierenden Aufgaben
	todos := loadTodos()
	// Finde die höchste ID und erhöhe sie um 1 für die neue Aufgabe
	newID := findHighestID(todos) + 1
	newTodo := Todo{ID: newID, Task: task, Completed: false}
	// Füge die neue Aufgabe der Liste hinzu
	todos = append(todos, newTodo)

	// Speichere die Aufgaben und zeige sie an
	saveTodos(todos)
	fmt.Println("Aufgabe hinzugefügt:", task)
}

// Funktion zum Speichern der Aufgaben in einer JSON-Datei
func saveTodos(todos []Todo) {
	data, err := json.MarshalIndent(todos, "", "  ") // JSON formatieren
	if err != nil {
		fmt.Println("Fehler beim Konvertieren in JSON:", err)
		return
	}

	err = ioutil.WriteFile("todos.json", data, 0644) // In eine Datei schreiben
	if err != nil {
		fmt.Println("Fehler beim Speichern der Datei:", err)
	}
}

// Funktion zum Anzeigen der Aufgaben
func viewTodoList() {
	todos := loadTodos() // Lade die Aufgaben
	fmt.Println("Aufgabenliste:")
	for _, todo := range todos {
		status := "offen"
		if todo.Completed {
			status = "erledigt"
		}
		fmt.Printf("ID: %d, Aufgabe: %s, Status: %s\n", todo.ID, todo.Task, status)
	}
}

// Funktion zum Aktualisieren einer Aufgabe
func updateTodoTask(id int, task string) {
	todos := loadTodos()
	found := false

	// Gehe durch alle Aufgaben und aktualisiere die richtige
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Task = task
			found = true
			break
		}
	}

	if found {
		saveTodos(todos)
		fmt.Println("Aufgabe aktualisiert:", task)
	} else {
		fmt.Println("Aufgabe mit ID", id, "nicht gefunden.")
	}
}

// Funktion zum Löschen einer Aufgabe
func deleteTodoTask(id int) {
	todos := loadTodos()
	newTodos := []Todo{}
	found := false

	// Füge nur die Aufgaben hinzu, die NICHT die zu löschende ID haben
	for _, todo := range todos {
		if todo.ID != id {
			newTodos = append(newTodos, todo)
		} else {
			found = true
		}
	}

	if found {
		saveTodos(newTodos)
		fmt.Println("Aufgabe gelöscht:", id)
	} else {
		fmt.Println("Aufgabe mit ID", id, "nicht gefunden.")
	}
}

// Funktion zum Laden der Aufgaben aus der Datei
func loadTodos() []Todo {
	file, err := ioutil.ReadFile("todos.json")
	if err != nil {
		// Wenn die Datei nicht existiert, fangen wir mit einer leeren Liste an
		fmt.Println("Keine Aufgaben gefunden. Beginne mit einer leeren Liste.")
		return []Todo{}
	}

	var todos []Todo
	err = json.Unmarshal(file, &todos) // JSON in Go-Objekte umwandeln
	if err != nil {
		fmt.Println("Fehler beim Laden der Aufgaben:", err)
		return []Todo{}
	}

	return todos
}

// Funktion zum Finden der höchsten ID
func findHighestID(todos []Todo) int {
	highestID := 0
	for _, todo := range todos {
		if todo.ID > highestID {
			highestID = todo.ID
		}
	}
	return highestID
}

// Anleitung für den Benutzer anzeigen
func printUsage() {
	fmt.Println("Verwendung:")
	fmt.Println("'add <Aufgabe>': Eine neue Aufgabe hinzufügen")
	fmt.Println("'update <ID> <neuer Text>': Eine vorhandene Aufgabe aktualisieren")
	fmt.Println("'delete <ID>': Eine Aufgabe löschen")
	fmt.Println("'view': Alle Aufgaben anzeigen")
}
