package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"github.com/5elenay/ezcli"
	"github.com/fatih/color"
	"github.com/tidwall/pretty"
)

/* ==== Globals ==== */

// colored output
var red = color.New(color.FgHiRed)
var green = color.New(color.FgHiGreen)
var blue = color.New(color.FgHiCyan).Add(color.Bold)
var white = color.New(color.FgHiWhite).Add(color.Bold)
var yellow = color.New(color.FgHiYellow)

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
func removeDotfileArray(slice []Dotfile, s int) []Dotfile {
	return append(slice[:s], slice[s+1:]...)
}

func removeInstallerArray(slice []Installer, s int) []Installer {
	return append(slice[:s], slice[s+1:]...)
}

func removeCommandArray(slice []Command, s int) []Command {
	return append(slice[:s], slice[s+1:]...)
}

/* === Generators === */

// Generate a new Config struct.
func newConfig(Name string, Description string, InstallPath string, Git bool, Repository GitConfig, Dotfiles []Dotfile, Installers []Installer) *Config {
	return &Config{
		Name:        Name,
		Description: Description,
		InstallPath: InstallPath,
		Git:         Git,
		Repository:  Repository,
		Dotfiles:    Dotfiles,
		Installers:  Installers,
	}
}

// Generate a new Dotfile struct.
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

// Generate a new GitConfig struct.
func newGitConfig(RemoteName string, Branch string, RemoteUrl string) *GitConfig {
	return &GitConfig{
		RemoteName: RemoteName,
		Branch:     Branch,
		RemoteUrl:  RemoteUrl,
	}
}

// Generate a new Installer struct
func newInstaller(Name string, Description string, Commands []Command) *Installer {
	return &Installer{
		Name:        Name,
		Description: Description,
		Commands:    Commands,
	}
}

// Generate a new Command struct
func newCommand(Execute string, Sudo bool) *Command {
	return &Command{
		Execute: Execute,
		Sudo:    Sudo,
	}
}

/* ================ */

// Convert Dotfile struct to JSON Representation
func to_JSON(d Config) []byte {
	out, _ := json.Marshal(d)
	out = pretty.Pretty(out)
	return out
}

// Convert JSON Representation to Dotfile struct
func from_JSON(JSON []byte) Config {
	var config Config
	err := json.Unmarshal(JSON, &config)
	check_Error(err)
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

	file_contents, err := ioutil.ReadFile("dotman.json")
	check_Error(err)

	conf := from_JSON(file_contents)
	return conf
}

// Add Dotfile struct to the dotman.json file
func add_dotfile(d Dotfile) {

	file, err := os.OpenFile("dotman.json", os.O_CREATE|os.O_RDWR, 0644)
	check_Error(err)

	conf := read_config_file()

	// Check if the dotfile is already declared
	for i := 0; i < len(conf.Dotfiles); i++ {
		if d.Name == conf.Dotfiles[i].Name {
			red.Printf("\"%s\" already exists in dotman.json\n", d.Name)
			os.Exit(1)
		}
	}

	conf.Dotfiles = append(conf.Dotfiles, d)

	wBytes, err := file.Write(to_JSON(conf))
	check_Error(err)

	_ = wBytes

	err = file.Sync()
	check_Error(err)
	err = file.Close()
	check_Error(err)
}

// Remove Dotfile struct from the dotman.json file
func remove_dotfile(name string) []string {

	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

	// if dotfile is exists set isFound to true and Remove the Dotfile struct from the array
	// otherwise set isFound to false and continue until whole loop is finished
	var isFound bool
	var Found Dotfile
	for index, obj := range conf.Dotfiles {
		if obj.Name == name {
			green.Printf("Found \"%s\" DotFile\n", obj.Name)
			blue.Println("Removing...")
			isFound = true
			Found = obj
			conf.Dotfiles = removeDotfileArray(conf.Dotfiles, index)
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
		err_ := os.Truncate("dotman.json", 0)
		check_Error(err_)

		writer := bufio.NewWriter(file)
		_, err := writer.Write(to_JSON(conf))
		check_Error(err)

		writer.Flush()
		file.Close()

		return []string{Found.Location, Found.Type}
	}
}

// Update a Dotfile struct in the dotman.json file
func update_dotfile(name string) string {

	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

	// if dotfile is exists set isFound to true and Update the Dotfile struct from the array
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
		err_ := os.Truncate("dotman.json", 0)
		check_Error(err_)

		writer := bufio.NewWriter(file)
		_, err := writer.Write(to_JSON(conf))
		check_Error(err)

		writer.Flush()
		file.Close()

		// git stuff
		git_add(FILES_PATH + filepath.Base(Location))
		git_commit(fmt.Sprintf("Dotman: Updated the \"%s\" dotfile | "+time_now, name))
		git_push(conf.Repository.RemoteName, conf.Repository.Branch)

		return Location
	}
}

// Update all the Dotfile structs in the dotman.json file
func update_all_dotfiles() {

	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

	green.Printf("\nFound %d dotfiles in dotman.json\n", len(conf.Dotfiles))

	// iterate over all the dotfiles and updae them
	for index, obj := range conf.Dotfiles {
		blue.Printf("Updating: %s...\n", obj.Name)
		conf.Dotfiles[index].LastUpdate = time_now
		copyFileOrDir(obj.Location, FILES_PATH)

		git_add(FILES_PATH + filepath.Base(obj.Location))
	}

	// git stuff
	git_commit("Dotman: Updated all the dotfiles. | " + time_now)
	git_push(conf.Repository.RemoteName, conf.Repository.Branch)

	err_ := os.Truncate("dotman.json", 0)
	check_Error(err_)

	writer := bufio.NewWriter(file)
	out, err := writer.Write(to_JSON(conf))
	check_Error(err)
	_ = out

	writer.Flush()
	file.Close()
}


// === Installer Functions ===

// add a new installer to the installers in dotman.json
func add_installer(ins Installer) {

	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

    // check if the installer is already declared
	for _, obj := range conf.Installers {
		if ins.Name == obj.Name {
			red.Printf("Installer \"%s\" already exists in dotman.json\n", ins.Name)
			os.Exit(1)
		}
	}

	conf.Installers = append(conf.Installers, ins)

	wBytes, err := file.Write(to_JSON(conf))
	check_Error(err)

	_ = wBytes

	err = file.Sync()
	check_Error(err)
	err = file.Close()
	check_Error(err)
}

// remove an installer from the installers in dotman.json
func remove_installer(name string) {

	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

    // if installer is exists set isFound to true and Remove the Installer struct from the array
    // otherwise set isFound to false and continue until whole loop is finished
	var isFound bool
	for index, obj := range conf.Installers {
		if obj.Name == name {
			green.Printf("Found \"%s\" DotFile\n", obj.Name)
			blue.Println("Removing...")
			isFound = true
			conf.Installers = removeInstallerArray(conf.Installers, index)
			break
		} else {
			isFound = false
			continue
		}
	}

	if !isFound {
		red.Printf("Can't find \"%s\" in dotman.json", name)
		os.Exit(1)
	} else {
		err_ := os.Truncate("dotman.json", 0)
		check_Error(err_)
		writer := bufio.NewWriter(file)
		_, err := writer.Write(to_JSON(conf))
		check_Error(err)

		writer.Flush()
		file.Close()
	}
}

// Generate installer Bash File
func CreateInstallerScript(name string, input []byte) {
	file, err := os.Create(INSTALLERS_PATH + "installer[" + name + "].sh")
	check_Error(err)

	writer := bufio.NewWriter(file)

	wBytes, err := writer.Write(input)
	check_Error(err)

	_ = wBytes

	writer.Flush()
	file.Close()
}

// add a new command to a specific installer
func add_command(name string, c Command) {

	file, err := os.OpenFile("dotman.json", os.O_CREATE|os.O_RDWR, 0644)
	check_Error(err)

	conf := read_config_file()

    // if  Installer is exists set isFound to true and add a new command to the Installer.Commands array
    // otherwise set isFound to false and continue until whole loop is finished
	var isFound bool
	var indx int
	for index, obj := range conf.Installers {
		if obj.Name == name {
			green.Printf("Found \"%s\" Installer\n", obj.Name)
			blue.Println("Adding...")
			isFound = true
			conf.Installers[index].Commands = append(conf.Installers[index].Commands, c)
			indx = index
			break
		} else {
			isFound = false
			continue
		}
	}

	if !isFound {
		red.Printf("Can't find Installer \"%s\" in dotman.json", name)
		os.Exit(1)
	} else {
		err_ := os.Truncate("dotman.json", 0)
		check_Error(err_)

		writer := bufio.NewWriter(file)
		_, err := writer.Write(to_JSON(conf))
		check_Error(err)

		writer.Flush()
		file.Close()
	}

	// git stuff
	gitAdd("./dotman.json", "Dotman: Added 1 Command to "+conf.Installers[indx].Name, conf.Repository.RemoteName, conf.Repository.Branch)
}

// remove a command from a specific installer
func remove_command(name string, command string) {

	file, err := os.OpenFile("dotman.json", os.O_CREATE|os.O_RDWR, 0644)
	check_Error(err)

	conf := read_config_file()

    // if command is exists set isFound to true and remove the command from the Installer.Commands array
    // otherwise set isFound to false and continue until whole loop is finished
	var isFound bool
	var indx int
	for index, obj := range conf.Installers {
		if obj.Name == name {
			green.Printf("Found \"%s\" Installer\n", obj.Name)

			for i, o := range conf.Installers[index].Commands {
				if o.Execute == command {
					green.Printf("Found \"%s\" Command\n", o.Execute)
					blue.Println("Removing...")
					conf.Installers[index].Commands = removeCommandArray(conf.Installers[index].Commands, i)
					isFound = true
					break
				} else {
					isFound = false
					continue
				}
			}

			indx = index
			break
		}
	}

	if !isFound {
		red.Printf("Can't find Command \"%s\" in Installer \"%s\" dotman.json", command, name)
		os.Exit(1)
	} else {
		err_ := os.Truncate("dotman.json", 0)
		check_Error(err_)

		writer := bufio.NewWriter(file)
		_, err := writer.Write(to_JSON(conf))
		check_Error(err)

		writer.Flush()
		file.Close()
	}

	// git stuff
	gitAdd("./dotman.json", "Dotman: Added 1 Command to "+conf.Installers[indx].Name, conf.Repository.RemoteName, conf.Repository.Branch)
}

// Generate a new bash script for a specific installer
func generate_installer(name string) {

	file, err := os.OpenFile("dotman.json", os.O_WRONLY, 6666)
	check_Error(err)

	conf := read_config_file()

    if !exists(INSTALLERS_PATH) {
        err := os.Mkdir(INSTALLERS_PATH, 0777)
        check_Error(err)
    }

    // if installer is exists set isFound to true and generate a new bash script
    // otherwise set isFound to false and continue until whole loop is finished
	var isFound bool
	for index, obj := range conf.Installers {
		if obj.Name == name {
			isFound = true

			var parsed_byte bytes.Buffer
			conf := read_config_file()

			tmpl := template.Must(template.New("installer").Parse(installer_template))

			err := tmpl.Execute(&parsed_byte, Tmpl{
                InstallPath: conf.InstallPath,
                Installer: conf.Installers[index],
            })
            
			check_Error(err)

			CreateInstallerScript(name, parsed_byte.Bytes())

			break
		} else {
			isFound = false
			continue
		}
	}

	if !isFound {
		red.Printf("Can't find Installer \"%s\" in dotman.json", name)
		os.Exit(1)
	} else {
		// Remove All the contents of the dotman.json file
		err_ := os.Truncate("dotman.json", 0)
		check_Error(err_)

		// new buffer writer for the file
		writer := bufio.NewWriter(file)
		_, err := writer.Write(to_JSON(conf))
		check_Error(err)

		// flush the io writer and close the file
		writer.Flush()
		file.Close()
	}

	gitAdd(INSTALLERS_PATH+"installer["+name+"].sh", "Dotman: Generated Installer Script ["+name+"]", conf.Repository.RemoteName, conf.Repository.Branch)
}

// Show the status of the dotman.json file
func show_status(option string) {

	conf := read_config_file()

    // show basic info
	if option == "normal" {

		color.New(color.FgHiCyan).Add(color.Bold).Add(color.Italic).Add(color.Italic).Printf("\n\n+-- Dotman Status --+\n")
		fmt.Printf("\n")

		white.Printf("%s", "Name: ")
		color.New(color.FgHiMagenta).Printf("%s\n", conf.Name)

		white.Printf("%s", "Description: ")
		yellow.Printf("\"%s\"\n", conf.Description)

		white.Printf("%s", "Install Path: ")
		yellow.Printf("%s\n", conf.InstallPath)

		white.Printf("%s", "Git Enabled? ")
		color.New(color.FgHiMagenta).Printf("%t\n", conf.Git)

		white.Printf("\n%s\n", "Repository\n+--------+")

		white.Printf("* %s", "Name: ")
		yellow.Printf("%s\n", conf.Repository.RemoteName)
		white.Printf("* %s", "Branch: ")
		yellow.Printf("%s\n", conf.Repository.Branch)
		white.Printf("* %s", "URL: ")
		yellow.Printf("%s\n", conf.Repository.RemoteUrl)
		white.Printf("+--------+")

		white.Printf("\n\n%s", "Dotfiles: ")
		green.Printf("%d Dotfiles in this dot.\n", len(conf.Dotfiles))

		white.Printf("%s", "Installers: ")
		green.Printf("%d Custom installers in this dot.\n\n", len(conf.Installers))

		color.New(color.FgHiCyan).Add(color.Bold).Printf("+-------------------+\n")
    
    // show detailed info about dotfiles
	} else if option == "dotfiles" {

		color.New(color.FgHiCyan).Add(color.Bold).Add(color.Italic).Add(color.Italic).Printf("\n\n+-- Dotman Dotfiles --+\n")
		fmt.Printf("\n")

		for _, obj := range conf.Dotfiles {
			red.Add(color.Bold).Add(color.Italic).Add(color.Italic).Printf("+-------------+\n")

			white.Printf("Name: ")
			yellow.Printf("%s\n", obj.Name)

			white.Printf("Description: ")
			yellow.Printf("%s\n", obj.Description)

			white.Printf("Location: ")
			yellow.Printf("%s\n", obj.Location)

			white.Printf("Type: ")
			yellow.Printf("%s\n", obj.Type)

			white.Printf("Priority: ")
			yellow.Printf("%d\n", obj.Priority)

			white.Printf("Last Update: ")
			yellow.Printf("%s\n", obj.LastUpdate)

			red.Add(color.Bold).Add(color.Italic).Add(color.Italic).Printf("+-------------+\n\n")
		}

		color.New(color.FgHiCyan).Add(color.Bold).Printf("+-------------------+\n")
    
    // show detailed info about installers
	} else if option == "installers" {
		color.New(color.FgHiCyan).Add(color.Bold).Add(color.Italic).Add(color.Italic).Printf("\n\n+-- Dotman Commands --+\n")
		fmt.Printf("\n")

		for _, obj := range conf.Installers {
			red.Add(color.Bold).Add(color.Italic).Add(color.Italic).Printf("+-------------+\n")

			white.Printf("Installer: ")
			yellow.Printf("%s\n", obj.Name)

			white.Printf("Description: ")
			yellow.Printf("%s\n", obj.Description)

			white.Printf("\nCommands:\n+--------+\n")

			for i, cmd := range obj.Commands {
				white.Printf("* Command: %s\n", color.HiYellowString(cmd.Execute))
				white.Printf("* Is Sudo?  %s\n", color.HiMagentaString(strconv.FormatBool(cmd.Sudo)))
				if i != len(obj.Commands)-1 {
					white.Printf("+--------+\n")
				}
			}

			white.Printf("+--------+\n")

			red.Add(color.Bold).Add(color.Italic).Add(color.Italic).Printf("+-------------+\n\n")
		}

		color.New(color.FgHiCyan).Add(color.Bold).Printf("+--------------------+\n")
	}

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
/* func check_git_branch() string {
	cmd := exec.Command("git", "branch")
	out, err := cmd.Output()
	check_Error(err)
	return string(out)
} */

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

	var conf Config

	var name string = ""
	var description string = "No description specified"
	var install_path string = ""
	var git bool = false
	var remote_url string
	var branch_name string

	var repository GitConfig
	var dotfiles []Dotfile = []Dotfile{}
	var installers []Installer = []Installer{}

	// Parse the positional argument "<name>"
	for i, v := range c.CommandData.Arguments {
		switch i {
		case 0:
			name = v
		}
	}

	// Wrong usage output
	if name == "" {
		red.Println("Please input the \"Name\" of the Dotman Project.")
		os.Exit(1)
	}

	for _, option := range c.CommandData.Options {
		switch option.Name {
		case "description":
			description = option.Value

		case "install_path":
			install_path = option.Value

		case "git":
			git = true

		case "remote":
			remote_url = option.Value

		case "branch":
			branch_name = option.Value
		}
	}

	// if "git" is set and if remote url and branch name set too, new GitConfig struct will be created
	if git {
		if remote_url != "" && branch_name != "" {
			repository = *newGitConfig("origin", branch_name, remote_url)
			blue.Println("Git Enabled.\n")
		} else {
			red.Println("Please set both the Branch name and the Remote URL to use git repositories!\n")
			os.Exit(1)
		}
	}

	conf = *newConfig(name, description, install_path, git, repository, dotfiles, installers)
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

	var name string = ""
	var location string = ""
	var description string = "No description specified"
	var Type string = ""
	var priority int64 = 1
	var lastupdate string = time.Now().Format("3:4:5 PM 2006-01-02")

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

    if !exists(location) {
        red.Println("The location you specified does not exist.")
        os.Exit(1)
    }
    
	// determine the type of the file being added
	if check_File_Type(location) == "file" {
		Type = "file"
	} else if check_File_Type(location) == "directory" {
		Type = "directory"
	}

	var inputDotfile Dotfile = *newDotfile(name, description, location, Type, priority, lastupdate)
	add_dotfile(inputDotfile)

	// if ./files/ exist then copy the file inside it, if not make a directory named ./files/ then copy the file inside it
	if exists(FILES_PATH) {
		copyFileOrDir(inputDotfile.Location, FILES_PATH)
	} else {
		err := os.Mkdir(FILES_PATH, os.FileMode(0755))
		check_Error(err)
		copyFileOrDir(inputDotfile.Location, FILES_PATH)
	}

	conf := read_config_file()

	gitAdd(FILES_PATH+filepath.Base(inputDotfile.Location), fmt.Sprintf("Dotman: Added a \"%s\" dotfile | "+inputDotfile.LastUpdate, name), conf.Repository.RemoteName, conf.Repository.Branch)

	green.Printf("\nAdded \"%s\" to Dotfiles.\n", name)
}

func Remove(c *ezcli.Command) {

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	var name string = ""

	// Parse the positional argument "<name>"
	for i, v := range c.CommandData.Arguments {
		switch i {
		case 0:
			name = v
		}
	}

	// Wrong usage output
	if name == "" {
		red.Println("Please input the \"Name\" of the Dotfile you want to remove.")
		os.Exit(1)
	}

	conf := read_config_file()

	// remove dotfile from dotman.json and return the Location of the removed Dotfile
	path := remove_dotfile(name)
	file_name := FILES_PATH + filepath.Base(path[0]) // just filename.extension

	gitRemove(file_name, fmt.Sprintf("Dotman: Removed the \"%s\" dotfile | "+time_now, name), conf.Repository.RemoteName, conf.Repository.Branch) // remove dotfile from git repository
	deleteFileOrDir(file_name, path[1])                                                                                                           // delete the file or directory dotfile linked to

	green.Printf("\"%s\" successfuly removed.\n", name)
}

func Update(c *ezcli.Command) {

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	var name string = ""

	// Parse the positional argument "<name>"
	for i, v := range c.CommandData.Arguments {
		switch i {
		case 0:
			name = v
		}
	}

	// Wrong usage output
	if name == "" {
		red.Println("Please input the \"Name\" of the Dotfile you want to update.")
		blue.Println("Use \"@a\" to update all the dotfiles.")
		os.Exit(1)
	}

	// Special handle named "@a" to update all the dotfiles || yes its inspired by minecraft :)
	if name == "@a" {
		update_all_dotfiles()
		green.Println("\nSuccessfuly updated all the Dotfiles.")
	} else { // otherwise update single dotfile
		dotfile_path := update_dotfile(name)
		copyFileOrDir(dotfile_path, FILES_PATH)
		green.Printf("\"%s\" successfuly updated.\n", name)
	}
}

func CommandHandler(c *ezcli.Command) {

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	var method string = ""

	var name string = ""
	var command string = ""
	var sudo bool = false

	// Parse the positional arguments "<method>, <command>" and "<name>"
	for i, v := range c.CommandData.Arguments {
		switch i {
		case 0:
			method = v
		case 1:
			name = v
		case 2:
			command = v
		}
	}

	for _, option := range c.CommandData.Options {
		switch option.Name {
		case "sudo":
			sudo = true
		}
	}

	// Wrong usage output
	if method == "" {
		red.Println("Please input the method you want to use.")
		blue.Println("Use \"add\" to add a command")
		blue.Println("Use \"remove\" to remove a command")
		os.Exit(1)
	} else if method != "add" && method != "remove" {
		red.Println("Invalid method!")
		blue.Println("Use \"add\" to add a command")
		blue.Println("Use \"remove\" to remove a command")
		os.Exit(1)
	} else if method == "add" && (name == "" || command == "") {
		red.Println("Empty field(s)! Please input all the fields")
		fmt.Printf("Try typing \"hello %s\"\n", color.CyanString("\"echo Hello World!\""))
		os.Exit(1)
	} else if method == "remove" && (name == "" || command == "") {
		red.Println("Empty field(s)! Please input all the fields")
		os.Exit(1)
	}

	if method == "add" {

		Command := newCommand(command, sudo)

		add_command(name, *Command)
		green.Println("\nAdded " + name + " to Commands.")

	} else if method == "remove" {
		remove_command(name, command)
		green.Println("\nRemoved " + name + " from Commands.")
	}

}

func InstallerHandler(c *ezcli.Command) {

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	var method string = ""

	var name string = ""
	var description string = ""

	// Parse the positional arguments "<method> and <name>"
	for i, v := range c.CommandData.Arguments {
		switch i {
		case 0:
			method = v
		case 1:
			name = v
		}
	}

	for _, option := range c.CommandData.Options {
		switch option.Name {
		case "description":
			description = option.Value
		}
	}

	// Wrong usage output
	if method == "" {
		red.Println("Please input the method you want to use.")
		blue.Println("Use \"add\" to add a command")
		blue.Println("Use \"remove\" to remove a command")
		os.Exit(1)
	} else if method != "add" && method != "remove" {
		red.Println("Invalid method!")
		blue.Println("Use \"add\" to add a command")
		blue.Println("Use \"remove\" to remove a command")
		os.Exit(1)
	} else if method == "add" && name == "" {
		red.Println("Empty field(s)! Please input all the fields")
		fmt.Printf("Try typing \"hello %s\"\n", color.CyanString("\"echo Hello World!\""))
		os.Exit(1)
	} else if method == "remove" && name == "" {
		red.Println("Empty field(s)! Please input all the fields")
		os.Exit(1)
	}

	if method == "add" {

		Instaler := newInstaller(name, description, []Command{})

		add_installer(*Instaler)
		green.Println("\nAdded " + name + " to Installer.")

	} else if method == "remove" {
		remove_installer(name)
		green.Println("\nRemoved " + name + " from Installers.")
	}

}

func Install(c *ezcli.Command) {

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

	var name string = ""

	// Parse the positional arguments "<name>"
	for i, v := range c.CommandData.Arguments {
		switch i {
		case 0:
			name = v
		}
	}

	// Wrong usage output
	if name == "" {
		red.Println("Please input the \"Name\" of the Installer you want to add.")
		os.Exit(1)
	}

	generate_installer(name)

	green.Println("\nGenerated Installer!")
	blue.Println("\nPlease run following command to make installer usable:")
	blue.Println("\tsudo chmod +x ./" + INSTALLERS_PATH + "installer[" + name + "].sh")
	blue.Println("\nCheck out ./installer.sh")

}

func Status(c *ezcli.Command) {

	// check if dotman.json exists
	if !check_config_exist() {
		red.Println("Can't find the dotman.json file.")
		blue.Println("Run \"dotman init\" to initialize the configuration.")
		os.Exit(1)
	}

	// variables ._.
	var method string = "NULL"

	// Parse the positional arguments "<method>"
	for i, v := range c.CommandData.Arguments {
		switch i {
		case 0:
			method = v
		}
	}

	// Wrong usage output
	if method != "normal" && method != "dotfiles" && method != "installers" && method != "NULL" {
		red.Printf("\"%s\" is not a valid command.", method)
		white.Println("\n\nCommands:")
		blue.Println("\tnormal     | display general information about dot.")
		blue.Println("\tdotfiles   | display information about dotfiles in dot.")
		blue.Println("\tinstallers | display information about installers in dot")
		os.Exit(1)
	} else if method == "NULL" {
		show_status("normal")
	} else {
		show_status(method)
	}
}

/* ======================= */
