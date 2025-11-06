package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

const Filename = "depense.json"
const osFlag = os.O_CREATE | os.O_RDONLY | os.O_TRUNC

type Action struct {
	Id          int     `json:"id"`
	ActionType  int     `json:"actiontype"`
	Description string  `json:"description"`
	Value       float32 `json:"value"`
	Date        string  `json:"date"`
}

type Save struct {
	Argent     int      `json:"argent"`
	Historique []Action `json:"historique"`
}

func ErrorChecker(err error) {
	if err != nil {
		panic(err)
	}
}

func GenerateSaveFile() {
	var save Save
	save.Argent = 0
	jsonData, err := json.MarshalIndent(save, "", "    ")
	ErrorChecker(err)
	// file, err := os.OpenFile(Filename, os.O_CREATE | os.O_RDONLY | os.O_TRUNC, os.ModePerm)
	file, err := os.Create(Filename)
	ErrorChecker(err)
	file.WriteString(string(jsonData))
	fmt.Println("File Generated")
}

func (save *Save) WriteFile() {
	jsonData, err := json.MarshalIndent(save, "", "    ")
	ErrorChecker(err)
	file, err := os.OpenFile(Filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	// file, err := os.Open(Filename)
	ErrorChecker(err)
	file.Write(jsonData)
	fmt.Println(save)
	fmt.Println("File Saved")
}
func (save *Save) ReadFile() {
	file, err := os.ReadFile(Filename)
	ErrorChecker(err)

	err = json.Unmarshal(file, &save)
	ErrorChecker(err)
}

func (action *Action) GetActionType() {
	fmt.Println("Choose the action type : \n1=> add\n2=> Depense")
	fmt.Scanln(&action.ActionType)
	if action.ActionType != 2 && action.ActionType != 1 {
		action.GetActionType()
	}
}
func (action *Action) GetDescription() {
	fmt.Println("Enter the Description : ")
	fmt.Scanln(&action.Description)
}
func (action *Action) GetValue() {
	fmt.Print("Enter the value : ")
	fmt.Scan(&action.Value)
	if action.Value <= 0 {
		action.GetValue()
	}
}

func (save *Save) Addaction() {
	var action Action
	if len(save.Historique) == 0 {
		action.Id = 1
	} else {
		action.Id = save.Historique[len(save.Historique) -1].Id +1
	}
	action.GetActionType()
	action.GetDescription()
	action.GetValue()
	action.Date = time.Now().GoString()
	save.Historique = append(save.Historique, action)
	if action.ActionType == 1 {
		save.Argent += int(action.Value)
	} else {
		save.Argent -= int(action.Value)
	}
}

func (action *Action) GetId() {
	fmt.Print("Enter the ID ")
	fmt.Scan(&action.Id)
	if action.Id <= 0 {
		action.GetId()
	}
}

func (save *Save) DeleteAction() {
	var act Action
	act.GetId()
	for idx, action := range save.Historique {
		if action.Id == act.Id {
			if action.ActionType == 1 {
				save.Argent -= int(act.Value)
			} else {
				save.Argent += int(action.Value)
			}
			save.Historique = append(save.Historique[:idx], save.Historique[idx+1:]...)
		}
	}
}

func main() {
	fmt.Println("Depense Tracker")
	
	if _, err := os.Stat(Filename); err != nil {
		GenerateSaveFile()
	}
	addAction := flag.Bool("add", false, "Add a depense or money")
	deleteaction := flag.Bool("delete", false, "Delete an action")

	flag.Parse()

	switch {
	case *addAction:

		var save Save
		save.ReadFile()
		save.Addaction()
		save.WriteFile()
	case *deleteaction:

		var save Save
		save.ReadFile()
		save.DeleteAction()
		save.WriteFile()
	default:
		fmt.Println(`
--add
		add an action
--delete
		delete an action
--list 
		list all action

		`)
	}
}
