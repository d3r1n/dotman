package main

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"bufio"
	"strconv"
	"os/exec"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"text/template"

	"github.com/5elenay/ezcli"
	"github.com/fatih/color"
	"github.com/tidwall/pretty"
)

/* ==== Globals ==== */

// colored output
var red 	= color.New(color.FgHiRed)
var green 	= color.New(color.FgHiGreen)
var blue 	= color.New(color.FgHiCyan).Add(color.Bold)

// Formatted time.Now()
var time_now = time.Now().Format("3:4:5 PM 2006-01-02")

/* ================ */

/* === FUNCTIONS === */

// check if error is not nil
func check_Error(err error) {
	if err != nil {
		panic(err)
	}
}

// check the file type of given file 
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

func Remove_Command_Array(slice []string, s int) []string {
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
func newGitConfig(RemoteName string, Branch string, RemoteUrl string) *GitConfig {
	return &GitConfig {
		RemoteName: 	RemoteName,
		Branch:    		Branch,
		RemoteUrl: 		RemoteUrl,
	} 
}

// Generate a new Config object.
func newConfig(Name string, Description string, InstallPath string, Git bool, Repository GitConfig, Dotfiles []Dotfile, Commands []string) *Config {
	return &Config {
		Name:        	Name,
		Description: 	Description,
		InstallPath: 	InstallPath,
		Git: 			Git,
		Repository: 	Repository,
		Dotfiles: 		Dotfiles,
		Commands: Commands,
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

func check_install_path_exist() bool {
	conf := read_config_file()

	if conf.InstallPath == "" {
		return false
	} else {
		return true
	}
}

// check if path exists
func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// read the dotman.json config file
func read_config_file() Config {

	// Get the file contents of the dotman.json file
	file_contents, err := ioutil.ReadFile("dotman.json")
	check_Error(err)

	// Convert contents to an Array of Dotfile objects
	conf := from_JSON(file_contents)
	return conf
}

// Add Dotfile Object to the dotman.json file
func add_dotfile(d Dotfile) {

	// Open the dotman.json file in current working directory as Read & Write mode
	file, err := os.OpenFile("dotman.json", os.O_CREATE | os.O_RDWR, 0644)
	check_Error(err)

	conf := read_config_file()

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
	wBytes, err := file.Write(to_JSON(conf))
	check_Error(err)

	_ = wBytes

	// flush the io writer and close the file
	file.Sync()
	file.Close()
}

// Remove Dotfile Object from the dotman.json file
func remove_dotfile(name string) []string {

	// Open the dotman.json file in current working directory as Read & Write mode
	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

	// if dotfile is exists set isFound to true and Remove the Dotfile Object from the array
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
		_, err := writer.Write(to_JSON(conf))
		check_Error(err)

		// flush the io writer and close the file
		writer.Flush()
		file.Close()

		return []string{Found.Location, Found.Type}
	}
}

// Update a Dotfile Object in the dotman.json file
func update_dotfile(name string) string {

	// Open the dotman.json file in current working directory as Write only mode
	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

	// if dotfile is exists set isFound to true and Update the Dotfile Object from the array
	// otherwise set isFound to false and continue until whole loop is finished
	var isFound bool
	var Location string
	for _, obj := range conf.Dotfiles {
		if obj.Name == name {
			green.Printf("\nFound \"%s\" DotFile\n", obj.Name)
			blue.Println("Updating...\n")
			isFound = true
			Location = obj.Location
			obj.LastUpdate = time_now
			copyFileOrDir(obj.Location, FILES_PATH)
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
		_, err := writer.Write(to_JSON(conf))
		check_Error(err)

		// flush the io writer and close the file
		writer.Flush()
		file.Close()

		// git stuff
		git_add(FILES_PATH + filepath.Base(Location))
		git_commit(fmt.Sprintf("Dotman: Updated the \"%s\" dotfile | " + time_now, name))
		git_push(conf.Repository.RemoteName, conf.Repository.Branch)

		return Location
	}
}

// Update all the Dotfile Objects in the dotman.json file
func update_all_dotfiles() {

	// Open the dotman.json file in current working directory as Write only mode
	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

	green.Printf("\nFound %d dotfiles in dotman.json\n", len(conf.Dotfiles))

	// if dotfile is already exists set isFound to true and Remove the Dotfile Object from the array
	// otherwise set isFound to false and continue until whole loop is finished
	for index, obj := range conf.Dotfiles {
		blue.Printf("Updating: %s...\n",  obj.Name)
		conf.Dotfiles[index].LastUpdate = time_now
		copyFileOrDir(obj.Location, FILES_PATH)
		
		// add file to git repository | executed for every dotfile in dotman.json
		git_add(FILES_PATH + filepath.Base(obj.Location)) 
	}

	// git stuff
	git_commit("Dotman: Updated all the dotfiles. | " + time_now)
	git_push(conf.Repository.RemoteName, conf.Repository.Branch)

	// Remove All the contents of the dotman.json file
	os.Truncate("dotman.json", 0)

	// new buffer writer for the file
	writer := bufio.NewWriter(file)
	out, err := writer.Write(to_JSON(conf))
	check_Error(err)
	_ = out

	// flush the io writer and close the file
	writer.Flush()
	file.Close()
}

// installer functions:

// Generate installer Bash File
func CreateInstaller(input []byte) {
	file, err := os.Create("./installer.sh")
	check_Error(err)

	writer := bufio.NewWriter(file)

	wBytes, err := writer.Write(input)
	check_Error(err)

	_ = wBytes

	writer.Flush()
	file.Close()
}

// Add Dotfile Object to the dotman.json file
func add_command(c string) {

	// Open the dotman.json file in current working directory as Read & Write mode
	file, err := os.OpenFile("dotman.json", os.O_CREATE | os.O_RDWR, 0644)
	check_Error(err)

	conf := read_config_file()

	// append the dotfile array
	conf.Commands = append(conf.Commands, c)

	// new buffer writer for the file
	wBytes, err := file.Write(to_JSON(conf))
	check_Error(err)

	_ = wBytes

	// flush the io writer and close the file
	file.Sync()
	file.Close()

	// git stuff
	gitAdd("./dotman.json", "Dotman: Added 1 Command", conf.Repository.RemoteName, conf.Repository.Branch)
}

// Remove Dotfile Object from the dotman.json file
func remove_command(c string) {

	// Open the dotman.json file in current working directory as Read & Write mode
	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

	// if dotfile is exists set isFound to true and Remove the Dotfile Object from the array
	// otherwise set isFound to false and continue until whole loop is finished
	var isFound bool
	for index, obj := range conf.Commands {
		if obj == c {
			green.Printf("Found \"%s\" Command\n", obj)
			blue.Println("Removing...")
			isFound = true
			conf.Commands = Remove_Command_Array(conf.Commands, index)
			break
		} else {
			isFound = false
			continue
		}
	}

	if !isFound {
		red.Printf("Can't find \"%s\" in Commands", c)
		os.Exit(1)

	} else {
		// Remove All the contents of the dotman.json file
		os.Truncate("dotman.json", 0)

		// new buffer writer for the file
		writer := bufio.NewWriter(file)
		_, err := writer.Write(to_JSON(conf))
		check_Error(err)

		// flush the io writer and close the file
		writer.Flush()
		file.Close()

		// git stuff
		gitRemove("./dotman.json", "Dotman: Removed 1 Command", conf.Repository.RemoteName, conf.Repository.Branch)
	}
}

func generate_installer() {
	
	var parsed_byte bytes.Buffer
	conf := read_config_file()
	
	tmpl := template.Must(template.New("installer").Parse(installer_template))

	err := tmpl.Execute(&parsed_byte, conf)
	check_Error(err)

	CreateInstaller(parsed_byte.Bytes())

	gitAdd("./installer.sh", "Dotman: Generated Installer Script", conf.Repository.RemoteName, conf.Repository.Branch)
}

// File Utils:

// Copy a file or a directory
func copyFileOrDir(from_location string, to_location string) {
	cmd := exec.Command("cp", "-f", "-r", from_location, to_location)
	err := cmd.Run()
	check_Error(err)
}

// Delete a file or a directory
func deleteFileOrDir(file string, file_type string) {
	if file_type == "file" {
		cmd := exec.Command("rm", file)
		cmd.Run()
	} else if file_type == "directory" {
		cmd := exec.Command("rm", "-d", file)
		err := cmd.Run()
		check_Error(err)
	}

}

/* ======================= */

/* ==== GIT UTILITIES ==== */

// git init
func git_init() {
	cmd := exec.Command("git", "init")
	cmd.Run()
}

// git remote add 
func git_remote_add(name string, url string) {
	cmd := exec.Command("git", "remote", "add", name, url)
	cmd.Run()
}

func git_branch_add(name string) {
	cmd := exec.Command("git", "branch", "-M", name)
	cmd.Run()
}

// git add
func git_add(path string) {
	cmd := exec.Command("git", "add", path)
	cmd.Run()
}

// git remove
func git_remove(path string) {
	cmd := exec.Command("git", "rm", path, "--cached")
	cmd.Run()
}

// git commit
func git_commit(message string) {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Run()
}

// git checkout
func git_checkout(branch string) {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Run()
}

// check the current working branch
func check_git_branch() string {
	cmd := exec.Command("git", "branch")
	out, err := cmd.Output()
	check_Error(err)
	return string(out)
}

// git push
func git_push(remote string, branch string) {
	cmd := exec.Command("git", "push", remote, branch)
	cmd.Run()
}

// All in one commands:

	// add dotfile
func gitAdd(add_path string, commit_message string, remote_name string, remote_branch string) {
	git_add(add_path)
	git_commit(commit_message)
	git_push(remote_name, remote_branch)
}
	// remove dotfile
func gitRemove(remove_path string, commit_message string, remote_name string, remote_branch string) {
	git_remove(remove_path)
	git_commit(commit_message)
	git_push(remote_name, remote_branch)
}

/* ======================= */

/* ==== Main Functions ==== */

func Init(c *ezcli.Command) {

	// variables ._.
	var conf Config

	var name 			string = ""
	var description 	string = "No description specified"
	var install_path 	string = ""
	var git 			bool   = false
	var remote_url 		string
	var branch_name 	string

	var repository 		GitConfig
	var dotfiles 		[]Dotfile = []Dotfile{}
	var commands		[]string = []string{}

	// Parse the positional argument "<name>"
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

	// if user specified a description, install path or git enable, update the variables
	for _, option := range c.CommandData.Options {
		switch option.Name {
			case "description":
				description = option.Value

			case "install_path":
				install_path = option.Value
				
			case "git":
				git = true

			case "remote_url":
				remote_url = option.Value
				
			case "branch_name":
				branch_name = option.Value
			}
		}
	}

	// if "git" is set and if remote url and branch name set too, new GitConfig Object will be created
	if git {
		if remote_url != "" && branch_name != "" {
			repository = *newGitConfig("origin", branch_name, remote_url)
			blue.Println("Git Enabled.\n")
		} else {
			red.Println("Please set both the Branch name and the Remote URL to use git repositories!\n")
			os.Exit(1)
		}
	}

	conf = *newConfig(name, description, install_path, git, repository, dotfiles, commands)
	json := to_JSON(conf)
	CreateAndWrite_JSON(json)

	// git stuff
	git_init()
	git_add("./dotman.json")
	git_commit("Dotman: Initialized Repository")
	git_branch_add(branch_name)
	git_checkout(branch_name)
	git_remote_add("origin", remote_url)
	git_push("origin", branch_name)

	green.Printf("Initialized Dotfile Project \"%s\" successfuly.", name)
}

func Add(c *ezcli.Command) {

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

	// Create a new Dotfile Object and add it to the dotman.json
	var inputDotfile Dotfile = *newDotfile(name, description, location, Type, priority, lastupdate)
	add_dotfile(inputDotfile)

	// if ./files/ exist then copy the file inside it, if not make a directory named ./files/ then copy the file inside it
	if exists(FILES_PATH) {
		copyFileOrDir(inputDotfile.Location, FILES_PATH)
	} else {
		os.Mkdir(FILES_PATH, os.FileMode(0755))
		copyFileOrDir(inputDotfile.Location, FILES_PATH)
	}

	conf := read_config_file()

	gitAdd(FILES_PATH + filepath.Base(inputDotfile.Location), fmt.Sprintf("Dotman: Added a \"%s\" dotfile | " + inputDotfile.LastUpdate, name), conf.Repository.RemoteName, conf.Repository.Branch)

	green.Printf("\nAdded \"%s\" to Dotfiles.\n", name)
}

func Remove(c *ezcli.Command) {	

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	// variables ._.
	var name string = ""

	// Parse the positional argument "<name>"
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

	conf := read_config_file()

	// remove dotfile from dotman.json and return the Location of the removed Dotfile
	path := remove_dotfile(name)
	file_name := FILES_PATH + filepath.Base(path[0]) // just filename.extension

	gitRemove(file_name, fmt.Sprintf("Dotman: Removed the \"%s\" dotfile | " + time_now, name), conf.Repository.RemoteName, conf.Repository.Branch) // remove dotfile from git repository
	deleteFileOrDir(file_name, path[1]) // delete the file or directory dotfile linked to

	green.Printf("\"%s\" successfuly removed.\n", name)
}

func Update(c *ezcli.Command) {

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	// variables ._.
	var name string = ""

	// Parse the positional argument "<name>"
	for i, v := range c.CommandData.Arguments {
		switch i {
			case 0:
				name = v
		}
	}
	
	// Wrong usage output
	if name == ""{
		red.Println("Please input the \"Name\" of the Dotfile you want to update.")
		blue.Println("Use \"@a\" to update all the dotfiles.")
		os.Exit(1)
	}

	// Special handle named "@a" to update all the dotfiles ( yes its inspired by minecraft :) )
	if name == "@a" {
		update_all_dotfiles()
		green.Println("\nSuccessfuly updated all the Dotfiles.")
	} else { // otherwise update single dotfile
		dotfile_path := update_dotfile(name)
		copyFileOrDir(dotfile_path, FILES_PATH)
		green.Printf("\"%s\" successfuly updated.\n", name)
	}
}

func Command(c *ezcli.Command) {

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	// variables ._.
	var method 	string = ""
	var command string = ""

	// Parse the positional arguments "<method> and <command>"
	for i, v := range c.CommandData.Arguments {
		switch i {
			case 0:
				method = v
			case 1:
				command = v
		}
	}
	
	// Wrong usage output
	if method == "" {
		red.Println("Please input the method you want to use.")
		blue.Println("Use \"add\" to add a command")
		blue.Println("Use \"remove\" to remove a command")
		os.Exit(1)
	} else if method != "add" && method != "remove"{
		red.Println("Invalid method!")
		blue.Println("Use \"add\" to add a command")
		blue.Println("Use \"remove\" to remove a command")
		os.Exit(1)
	} else if command == ""{
		red.Println("Command cannot be empty!")
		fmt.Printf("Try typing %s \n", color.CyanString("\"echo Hello World!\""))
		os.Exit(1)
	}

	if method == "add" {

		add_command(command)
		green.Println("\nAdded " + command + " to Commands.")

	} else if method == "remove" {
		remove_command(command)
		green.Println("\nRemoved " + command + " from Commands.")
	}

}

func Installer(c *ezcli.Command) {

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	// check if installation path is specified
	if !check_install_path_exist() {
		red.Println("Can't find Installation Path in dotman.json file.")
		blue.Println("Please initialize the Installation Path in the configuration.")
		os.Exit(1)
	}

	generate_installer()

	cmd := exec.Command("sudo chmod +x ./installer.sh")
	cmd.Run()

	green.Println("\nGenerated Installer!")
	blue.Println("\nPlease run following command to make installer usable:")
	blue.Println("\tsudo chmod +x ./installer.sh")
	blue.Println("\nCheck out ./installer.sh")

}
/* ======================= */