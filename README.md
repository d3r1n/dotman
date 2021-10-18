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
	* Dotman Binary installing script
	* Bug fixes
- v0.2
	* Automatic git support added - [Main Feature]
	* Lots of bug fixes
- v0.1
	* Initial version of Dotman

## Usage:

### Commands

```python
# Init
dotman init <name> [--description=<text>, -d=<text>] [--git, -g] [--remote_url=<git url> --branch_name=<text>] [--install_path=<location>, -i=<location>]

# Add
dotman add <name> <location> [--description=<text>, -d=<text>] [--priority=<int>, -p=<int>]

# Remove
dotman remove <name>

# Update
dotman update <name>

# Update All
dotman update @a

# Add a Command
dotman [command, cmd] add <command>

# Remove a Command
dotman [command, cmd] remove <command>

# Generate Installer Script
dotman [generate, gen]
```
---

<div align="center">

If you liked this project please leave a star, it helps a lot :3

</div>