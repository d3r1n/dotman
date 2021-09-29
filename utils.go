package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/5elenay/ezcli"
	"github.com/fatih/color"
	"github.com/tidwall/pretty"
)

/* === FUNCTIONS === */

func check_Error(err error) {
	if err != nil {
		panic(err)
	}
}

func check_File_Type(path string) string {
	fileInfo, err := os.Stat(path)
	check_Error(err)

	if fileInfo.IsDir() {
		return "directory"
	} else {
		return "file"
	}
}

// go does not have a built in array remove function / method WTF?!?!?!?!?!?
func Remove_Dot_Array(slice []Dotfile, s int) []Dotfile {
	return append(slice[:s], slice[s+1:]...)
}

// Generate a new Dotfile object.
func newDotfile(Name string, Description string, Location string, Type string, Priority int64, LastUpdate string) *Dotfile {
	return &Dotfile{
		Name:        Name,
		Description: Description,
		Location:    Location,
		Type:        Type,
		Priority:    Priority,
		LastUpdate:  LastUpdate,
	}
}

// Generate a new GitConfig object.
func newGitConfig(Name string, Description string, Branch string, Origin string) *GitConfig {
	return &GitConfig {
		Name: 			Name,
		Description:    Description,
		Branch:    		Branch,
		Origin: 		Origin,
	} 
}

// Generate a new Config object.
func newConfig(Name string, Description string, InstallPath string, Git bool, Repository GitConfig, Dotfiles []Dotfile) *Config {
	return &Config {
		Name:        	Name,
		Description: 	Description,
		InstallPath: 	InstallPath,
		Git: 			Git,
		Repository: 	Repository,
		Dotfiles: 		Dotfiles,
	}
}

// Convert Dotfile Object to JSON Representation
func to_JSON(d Config) []byte {
	out, _ := json.Marshal(d)
	out = pretty.Pretty(out)
	return out
}

// Convert JSON Representation to Dotfile Object
func from_JSON(JSON []byte) Config {
	var config Config
	json.Unmarshal(JSON, &config)
	return config
}

// Generate JSON File
func CreateAndWrite_JSON(input []byte) {
	file, err := os.Create("./dotman.json")
	check_Error(err)

	writer := bufio.NewWriter(file)

	wBytes, err := writer.Write(input)
	check_Error(err)

	_ = wBytes

	writer.Flush()
	file.Close()
}

// Check if dotman.json exists in the current working directory
func check_config_exist() bool {
	if _, err := os.Stat("./dotman.json"); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Add Dotfile Object to the dotman.json file
func add_dotfile(d Dotfile) {

	red := color.New(color.FgRed)

	// Open the dotman.json file in current working directory as Read & Write mode
	file, err := os.OpenFile("dotman.json", os.O_CREATE | os.O_RDWR, 0644)
	check_Error(err)

	// Get the file contents of the dotman.json file
	file_contents, err := os.ReadFile("dotman.json")
	check_Error(err)

	// Convert contents to an Array of Dotfile objects
	conf := from_JSON(file_contents)

	// Check if the dotfile is already declared
	for i := 0; i < len(conf.Dotfiles); i++ {
		if d.Name == conf.Dotfiles[i].Name {
			red.Printf("\"%s\" already exists in dotman.json\n", d.Name)
			os.Exit(1)
		}
	}

	// append the dotfile array
	conf.Dotfiles = append(conf.Dotfiles, d)

	// new buffer writer for the file
	wBytes, _ := file.Write(to_JSON(conf))
	check_Error(err)

	_ = wBytes

	// flush the io writer and close the file
	file.Sync()
	file.Close()
}

// Remove Dotfile Object from the dotman.json file
func remove_dotfile(name string) []string {

	red := color.New(color.FgRed)
	green := color.New(color.FgHiGreen)
	blue := color.New(color.FgHiCyan)

	// Open the dotman.json file in current working directory as Read & Write mode
	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	// Get the file contents of the dotman.json file
	file_contents, err := ioutil.ReadFile("dotman.json")
	check_Error(err)

	// Convert contents to an Array of Dotfile objects
	conf := from_JSON(file_contents)

	// if dotfile is already exists set isFound to true and Remove the Dotfile Object from the array
	// otherwise set isFound to false and continue until whole loop is finished
	var isFound bool
	var Found Dotfile
	for index, obj := range conf.Dotfiles {
		if obj.Name == name {
			green.Printf("Found \"%s\" DotFile\n", obj.Name)
			blue.Println("Removing...")
			isFound = true
			Found = obj
			conf.Dotfiles = Remove_Dot_Array(conf.Dotfiles, index)
			break
		} else {
			isFound = false
			continue
		}
	}

	if !isFound {
		red.Printf("Can't find \"%s\" in dotman.json", name)
		os.Exit(1)

		return []string{}
	} else {
		// Remove All the contents of the dotman.json file
		os.Truncate("dotman.json", 0)

		// new buffer writer for the file
		writer := bufio.NewWriter(file)
		wBytes, _ := writer.Write(to_JSON(conf))
		check_Error(err)

		_ = wBytes

		// flush the io writer and close the file
		writer.Flush()
		file.Close()

		return []string{Found.Location, Found.Type}
	}
}

func update_dotfile(name string) string {

	red := color.New(color.FgRed)
	green := color.New(color.FgHiGreen)
	blue := color.New(color.FgHiCyan)

	// Open the dotman.json file in current working directory as Read & Write mode
	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	// Get the file contents of the dotman.json file
	file_contents, err := ioutil.ReadFile("dotman.json")
	check_Error(err)

	// Convert contents to an Array of Dotfile objects
	conf := from_JSON(file_contents)

	// if dotfile is already exists set isFound to true and Remove the Dotfile Object from the array
	// otherwise set isFound to false and continue until whole loop is finished
	var isFound bool
	var Found Dotfile
	for index, obj := range conf.Dotfiles {
		if obj.Name == name {
			green.Printf("Found \"%s\" DotFile\n", obj.Name)
			blue.Println("Updating...")
			isFound = true
			Found = obj
			conf.Dotfiles[index].LastUpdate = time.Now().Format("3:4:5 PM 2006-01-02")
			break
		} else {
			isFound = false
			continue
		}
	}

	if !isFound {
		red.Printf("Can't find \"%s\" in dotman.json", name)
		os.Exit(1)

		return ""
	} else {
		// Remove All the contents of the dotman.json file
		os.Truncate("dotman.json", 0)

		// new buffer writer for the file
		writer := bufio.NewWriter(file)
		wBytes, _ := writer.Write(to_JSON(conf))
		check_Error(err)

		_ = wBytes

		// flush the io writer and close the file
		writer.Flush()
		file.Close()

		return Found.Location
	}
}

func update_all_dotfiles() {

	green := color.New(color.FgHiGreen)
	blue := color.New(color.FgHiCyan)

	// Open the dotman.json file in current working directory as Read & Write mode
	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	// Get the file contents of the dotman.json file
	file_contents, err := ioutil.ReadFile("dotman.json")
	check_Error(err)

	// Convert contents to an Array of Dotfile objects
	conf := from_JSON(file_contents)

	green.Printf("Found %d dotfiles in dotman.json\n\n", len(conf.Dotfiles))

	// if dotfile is already exists set isFound to true and Remove the Dotfile Object from the array
	// otherwise set isFound to false and continue until whole loop is finished
	for index, obj := range conf.Dotfiles {
		blue.Printf("Updating: %s...\n",  obj.Name)
		conf.Dotfiles[index].LastUpdate = time.Now().Format("3:4:5 PM 2006-01-02")
		copyFileOrDir(obj.Location, FILES_PATH)
	}

	// Remove All the contents of the dotman.json file
	os.Truncate("dotman.json", 0)

	// new buffer writer for the file
	writer := bufio.NewWriter(file)
	wBytes, _ := writer.Write(to_JSON(conf))
	check_Error(err)

	_ = wBytes

	// flush the io writer and close the file
	writer.Flush()
	file.Close()
}

// Update Dotfile Object(s) in the dotman.json file

// Copy a file or a directory
func copyFileOrDir(from_location string, to_location string) {
	cmd := exec.Command("cp", "-f", from_location, to_location)
	err := cmd.Run()
	check_Error(err)
}

// Delete a file or a directory
func deleteFileOrDir(file string, file_type string) {
	if file_type == "file" {
		cmd := exec.Command("rm", file)
		err := cmd.Run()
		check_Error(err)
	} else if file_type == "directory" {
		cmd := exec.Command("rm", "-d", file)
		err := cmd.Run()
		check_Error(err)
	}

}

/* ======================= */

/* ==== GIT UTILITIES ==== */

// git add

// git commit

// git push

/* ======================= */

/* ==== Main Utilities ==== */

func Init(c *ezcli.Command) {

	red := color.New(color.FgRed)
	green := color.New(color.FgHiGreen)
	blue := color.New(color.FgHiCyan)

	// variables ._.
	var conf Config

	var name 			string = ""
	var description 	string = "No description specified"
	var install_path 	string = ""
	var git 			bool = false
	var repository 		GitConfig
	var dotfiles 		[]Dotfile = []Dotfile{}

	// Parse the positional arguments "<name>"
	for i, v := range c.CommandData.Arguments {
		switch i {
			case 0:
				name = v
	}

	// Wrong usage output
	if name == "" {
		red.Println("Please input the \"Name\" of the Dotman Project.")
		os.Exit(1)
	}

	// if user specified a description or a priority update the variables
	for _, option := range c.CommandData.Options {
		switch option.Name {
			case "description":
				description = option.Value

			// priority must be between or equal 1 and 3
			case "install_path":
				install_path = option.Value
			
			case "git":
				git = true
			}	
		}
	}

	// if "git" is set true then initialize a Git repository
	if git {
		repository = *newGitConfig("", "", "", "")
		blue.Println("Git Enabled.\n")
	}

	conf = *newConfig(name, description, install_path, git, repository, dotfiles)
	json := to_JSON(conf)
	CreateAndWrite_JSON(json)

	green.Printf("Initialized Dotfile Project \"%s\" successfuly.", name)
}

func Add(c *ezcli.Command) {

	red := color.New(color.FgRed)
	green := color.New(color.FgHiGreen)
	blue := color.New(color.FgHiCyan)

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	// variables ._.
	var name 		string = ""
	var location 	string = ""
	var description string = "No description specified"
	var Type 		string = ""
	var priority 	int64  = 1
	var lastupdate 	string = time.Now().Format("3:4:5 PM 2006-01-02")

	// Parse the positional arguments "<name>" and "<location>"
	for i, v := range c.CommandData.Arguments {
		switch i {
			case 0:
				name = v

			case 1:
				location = v
		}
	}

	// Wrong usage output
	if name == "" || location == "" {
		red.Println("Please input the \"Name\" and \"Location\" of the Dotfile.")
		os.Exit(1)
	} else if name == "@a" {
		red.Println("\"@a\" is a special name. You cant use it!")
		os.Exit(1)
	}

	// if user specified a description or a priority update the variables
	for _, option := range c.CommandData.Options {
		switch option.Name {
			case "description":
				description = option.Value

			// priority must be between or equal 1 and 3
			case "priority":
				if i, _ := strconv.ParseInt(option.Value, 0, 8); i > 3 {
					priority = 3
				} else if i, _ := strconv.ParseInt(option.Value, 0, 8); i < 0 {
					priority = 1
				} else {
					priority = 1
			}
		}
	}

	// determine the type of the file being added
	if check_File_Type(location) == "file" {
		Type = "file"
	} else if check_File_Type(location) == "directory" {
		Type = "directory"
	}

	var inputDotfile Dotfile = *newDotfile(name, description, location, Type, priority, lastupdate)
	add_dotfile(inputDotfile)

	if exists(FILES_PATH) {
		copyFileOrDir(inputDotfile.Location, FILES_PATH)
	} else {
		os.Mkdir(FILES_PATH, os.FileMode(0755))
		copyFileOrDir(inputDotfile.Location, FILES_PATH)
	}

	green.Printf("\nAdded \"%s\" to Dotfiles.\n\n", name)
}

func Remove(c *ezcli.Command) {

	red := color.New(color.FgRed)
	green := color.New(color.FgHiGreen)

	// variables ._.
	var name string = ""

	// Parse the positional arguments "<name>" and "<location>"
	for i, v := range c.CommandData.Arguments {
		switch i {
			case 0:
				name = v
		}
	}
	
	// Wrong usage output
	if name == ""{
		red.Println("Please input the \"Name\" of the Dotfile you want to remove.")
		os.Exit(1)
	}

	path := remove_dotfile(name)
	file_name := FILES_PATH + filepath.Base(path[0])
	deleteFileOrDir(file_name, path[1])

	green.Printf("\"%s\" successfuly removed.\n\n", name)
}

func Update(c *ezcli.Command) {

	red := color.New(color.FgRed)
	green := color.New(color.FgHiGreen)

	// variables ._.
	var name string = ""

	// Parse the positional arguments "<name>" and "<location>"
	for i, v := range c.CommandData.Arguments {
		switch i {
			case 0:
				name = v
		}
	}
	
	// Wrong usage output
	if name == ""{
		red.Println("Please input the \"Name\" of the Dotfile you want to remove.")
		os.Exit(1)
	}

	if name == "@a" {
		update_all_dotfiles()
		green.Println("Successfuly updated all the Dotfiles.\n\n")
	} else {
		dotfile_path := update_dotfile(name)
		copyFileOrDir(dotfile_path, FILES_PATH)
		green.Printf("\"%s\" successfuly updated.\n\n", name)
	}

	
}

/* ======================= */
