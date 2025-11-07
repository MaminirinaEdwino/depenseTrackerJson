package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	colortext "github.com/maminirinaedwino/depenseTrackerJson/ColorText"
)

const Filename = "depense.foza"

// const osFlag = os.O_CREATE | os.O_RDONLY | os.O_TRUNC

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
	homeDir, err := os.UserHomeDir()
	ErrorChecker(err)
	fmt.Println(homeDir)
	file, err := os.Create(homeDir+"/"+Filename)
	ErrorChecker(err)
	file.WriteString(string(jsonData))
}

func (save *Save) WriteFile() {
	jsonData, err := json.MarshalIndent(save, "", "    ")
	ErrorChecker(err)
	homeDir, err := os.UserHomeDir()
	ErrorChecker(err)
	
	file, err := os.OpenFile(homeDir+"/"+Filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	// file, err := os.Open(Filename)
	ErrorChecker(err)
	file.Write(jsonData)
	fmt.Println("File Saved")
}
func (save *Save) ReadFile() {
	homeDir, err := os.UserHomeDir()
	ErrorChecker(err)
	fmt.Println("home ",homeDir)
	file, err := os.ReadFile(homeDir+"/"+Filename)
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
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Description : ")
	if scanner.Scan() {
		action.Description = scanner.Text()
		if action.Description == "" {
			action.GetDescription()
		}
	}
}
func (action *Action) GetValue(argent float32) {
	fmt.Print("Enter the value : ")
	fmt.Scan(&action.Value)
	if action.Value <= 0 {
		action.GetValue(argent)
	}
	if action.ActionType == 2 {
		if action.Value > argent {
			action.GetValue(argent)
		}
	}
}

func (save *Save) Addaction() {
	var action Action
	if len(save.Historique) == 0 {
		action.Id = 1
	} else {
		action.Id = save.Historique[len(save.Historique)-1].Id + 1
	}
	action.GetActionType()
	action.GetDescription()
	action.GetValue(float32(save.Argent))
	action.Date = time.Now().Format("02-01-2006 03:04:05")
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
func (action *Action) ReturnActionType() string {
	if action.ActionType == 1 {
		return colortext.GreenString("add")
	}
	return colortext.RedText("dépense")
}
func (action *Action) ShowAction(w *tabwriter.Writer) {
	fmt.Fprintf(w, "%d\t%s\t%s\t%.2f\t\t%s\n", action.Id, action.ReturnActionType(), action.Description, action.Value, action.Date)
}

func (save *Save) ListAllAction(actiontype int) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "ID\tAction\tDescription\tValue\t\tDate\t")

	for _, action := range save.Historique {
		if actiontype > 0 {
			if actiontype == action.ActionType {
				action.ShowAction(w)
			}
		}else{
			action.ShowAction(w)
		}
	}
	w.Flush()
	fmt.Println(colortext.GreenString("Argent : " + strconv.Itoa(save.Argent)))
}

func main() {
	fmt.Println(colortext.GreenString("Depense Tracker"))
	var save Save
	homeDir, err := os.UserHomeDir()
	ErrorChecker(err)

	if _, err := os.Stat(homeDir+"/"+Filename); err != nil {
		GenerateSaveFile()
	}
	addAction := flag.Bool("add", false, "Add a depense or money")
	deleteaction := flag.Bool("delete", false, "Delete an action")
	listAllAction := flag.Bool("list", false, "List all action")
	listAllAddAction := flag.Bool("list-add", false, "List all add action")
	lsitAllDepenseAction := flag.Bool("list-depense", false, "List all depense action")
	flag.Parse()

	switch {
	case *addAction:
		save.ReadFile()
		save.Addaction()
		save.WriteFile()
		save.ListAllAction(0)
	case *deleteaction:
		save.ReadFile()
		save.DeleteAction()
		save.WriteFile()
		save.ListAllAction(0)
	case *listAllAction:
		save.ReadFile()
		save.ListAllAction(0)
	case *listAllAddAction:
		save.ReadFile()
		save.ListAllAction(1)
	case *lsitAllDepenseAction:
		save.ReadFile()
		save.ListAllAction(2)
	default:
		fmt.Println(`
--add
		add an action
--delete
		delete an action
--list 
		list all action
--list-add
		list all add action
--list-depense
		list all depense action
		`)
	}
}
