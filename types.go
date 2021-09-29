package main

const FILES_PATH string = "./files/"

type Config struct {
	Name 		string // 		name of the Dotfile Project				(required)
	Description string // 		description of the Dotfile Project		(optional) (default: "No Description Provided")
	InstallPath string // 		path to the installation 				(required)

	// if true, a git repo will be initialized and automatically updated with every action (optional) (default: false)
	Git 		bool 

	Repository GitConfig 
	Dotfiles []Dotfile
}

type GitConfig struct {
	Name 			string // 		name of the git repository 			(required)
	Description 	string // 		description of the git repository 	(required)
	Branch 			string // 		branch name of the git repository 	(optional) (default: "master")
	Origin 			string // 		name of the git repository 			(required)
}

type Dotfile struct {
	Name        string // 		the name of the dotfile					(required)
	Description string // 		the description of the dotfile 			(optional) (default: "No Description Provided")
	Location    string // 		the location of the dotfile				(required)
	Type        string // 		the type of the dotfile					(auto) 	   (file, directory)
	Priority    int64  // 		the priority of the dotfile				(optional) (default: 1) (1-3)
	LastUpdate  string // 		the last update date of the dotfile 	(auto)
}