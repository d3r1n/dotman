<div align="center">

<img src="./assets/dotman_logo.png" height="250" width="250" alt="dotman's logo">

</div>

# Dotman
The dotfile manager you are searching for

## Version History
- v0.4 [Next]
	* Better Documentation
	* Website?
	* Multiple Installer Scripts - [Main Feature]
	* Publishing Dotman to Package Managers - [Main Feature]
		* Arch 
		* Ubuntu
		* Gentoo
		* Void
		* And Many others
- v0.3 [Now]
	* Automatic Generated installer script - [Main Feature]
	* Status Commands
	* Dotman Binary installing script
	* Bug fixes
- v0.2
	* Automatic git support added - [Main Feature]
	* Lots of bug fixes
- v0.1
	* Initial version of Dotman

## Installation

Just Execute this command to install the dotman from source

**Requirements**:
- Go **1.17** or *higher*.

```bash
curl -L https://git.io/J68cu | bash
```
**or**

```bash
curl https://raw.githubusercontent.com/d3r1n/dotman/master/linux_installer.sh | bash
```

## Usage:

### Commands

```python
# Init
dotman init <name> [--description=<text>, -d=<text>] [--git, -g] [--remote_url=<git url> --branch_name=<text>] [--install_path=<location>, -i=<location>]

# Show information about the dot
dotman [status, ss] <[normal, dotfiles, commands]>

# Add
dotman add <name> <location> [--description=<text>, -d=<text>] [--priority=<int>, -p=<int>]

# Remove
dotman remove <name>

# Update
dotman update <name>

# Update All
dotman update @a

# Add a Command
dotman [command, cmd] add [--sudo, -s] <name> <command>

# Remove a Command
dotman [command, cmd] remove <name>

# Generate Installer Script
dotman [generate, gen]
```
---

<div align="center">

If you liked this project please leave a star, it helps a lot :3

</div>