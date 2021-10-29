package main

const FILES_PATH string = "./files/"

type Config struct {
	Name        string // 		name of the Dotfile Project				(required)
	Description string // 		description of the Dotfile Project		(optional) (default: "No Description Provided")
	InstallPath string // 		path to the installation 				(required)

	// if true, a git repo will be initialized and automatically updated with every action (optional) (default: false)
	Git bool

	Repository GitConfig
	Dotfiles   []Dotfile
	Commands   []Command
}

type GitConfig struct {
	RemoteName string // name of the git repository 			(required)
	RemoteUrl  string // name of the git repository 			(required)
	Branch     string // branch name of the git repository 	(optional) (default: "master")
}

type Dotfile struct {
	Name        string // the name of the dotfile				(required)
	Description string // the description of the dotfile 		(optional) (default: "No Description Provided")
	Location    string // the location of the dotfile			(required)
	Type        string // the type of the dotfile				(auto) 	   (file, directory)
	Priority    int64  // the priority of the dotfile			(optional) (default: 1) (1-3)
	LastUpdate  string // the last update date of the dotfile 	(auto)
}

type Command struct {
	Name    string // the name of the command 								(required)
	Execute string // actual command										(required)
	Sudo    bool   // whether the command should be executed in sudo mode 	(optional) (default: false)
}

/* Template For The Installer Script */

// I hate go templates >:[

const installer_template = `

# -----------------------------------
# |									|
# |		   Dotman Installer			|
# |	https://github.com/d3r1n/dotman |
# |									|
# -----------------------------------

# Colored Output

# Reset
Color_Off='\033[0m'       # Text Reset

# Regular Colors
Black='\033[0;30m'        # Black
Red='\033[0;31m'          # Red
Green='\033[0;32m'        # Green
Yellow='\033[0;33m'       # Yellow
Blue='\033[0;34m'         # Blue
Purple='\033[0;35m'       # Purple
Cyan='\033[0;36m'         # Cyan
White='\033[0;37m'        # White

# Bold
BBlack='\033[1;30m'       # Black
BRed='\033[1;31m'         # Red
BGreen='\033[1;32m'       # Green
BYellow='\033[1;33m'      # Yellow
BBlue='\033[1;34m'        # Blue
BPurple='\033[1;35m'      # Purple
BCyan='\033[1;36m'        # Cyan
BWhite='\033[1;37m'       # White

# Installation Location
install_path={{ .InstallPath }}

# Setup Function
function setup_point() {
	printf "${BWhite}"
	printf "\n --- Dotman Installer --- \n"
	printf "https://github.com/d3r1n/dotman\n"
	printf "${Color_Off}"
	printf "\n"
	printf "Copying Files from ${BBlue}./files/${Color_Off} to ${BBlue}${install_path}${Color_Off}\n"
	cp -r ./files/* $install_path

	printf "${BWhite}"
	printf "\n --- User Defined Commands --- \n"
	printf "${Color_Off}"
	printf "\n"
}

# User defined commands
function run_point() {
	{{ range .Commands }}
	{{ if .Sudo }}# Command: {{ .Name}}
	sudo {{ .Execute}}{{ else }}
	# Command {{ .Name }}
	{{ .Execute }}{{ end }}
	{{ end }}
}

# Entry points
setup_point
run_point
`
