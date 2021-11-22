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

# install Function
function install() {
	printf "${BWhite}"
	printf "\n --- Dotman Installer --- \n"
	printf "https://github.com/d3r1n/dotman\n"
	printf "${Color_Off}"
	printf "\n"

	printf "${BPurple}"
	printf "Cloning Repository..."
	printf "${Color_Off}"

	git clone https://github.com/d3r1n/dotman.git

	printf "${BPurple}"
	printf "\nBuilding Dotman...\n"
	printf "${Color_Off}"

	cd ./dotman
	go get github.com/5elenay/ezcli
	go get github.com/fatih/color
	go get github.com/tidwall/pretty
	go build

	printf "${BPurple}"
	printf "\nMoving Executable to /bin/ ..\n."
	printf "${Color_Off}"

	sudo mv ./dotman /usr/bin/

	printf "${BPurple}"
	printf "\nCleaning...\n"
	printf "${Color_Off}"
	
	cd ..
	rm -rf ./dotman
}

install
