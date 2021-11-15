# Documentation for commands

## **dotman init**
### *Description:*
Initialize a dotman project. \
Automatically creates a `dotman.json` file in the current directory. \
If you already have a `dotman.json` file, it will be overwritten. 

### *Usage:*
```bash
dotman init <name> [options]
```

### *Options:*
```bash
help:
    description: Show help information
    alias: none
    type: boolean
    default: false
```

```bash
name:
    description: The name of the project
    type: string
    required: true
```

```bash
--description:, -d:
    description: A description: of the project
    alias: -d
    type: string
    default: "No description: provided."
```

```bash
--install_path, -i:
    description: The location to install the project
    alias: -i
    type: string
    default: none
```

```bash
--git, -g:
    description: Initialize a git repository
    alias: -g
    type: boolean
    default: false
```

```bash
--remote, -r:
    description: The remote url of the git repository
    alias: -r
    type: string
    default: none
```

```bash
--branch, -b:
    description:*The branch name of the git repository
    alias: -b
    type: string
    default: none
```
---

## **dotman add**
### *Description:*
Add a dotfile to the project.


### *Usage:*
```bash
dotman add <name> <path> [options]
```

### *Options:*

```bash
name:
    description The name of the dotfile
    type: string
    required: true
```

```bash
path:
    description: The path to the dotfile
    type: string
    required: true
```

```bash
--description, -d:
    description: A description of the project
    alias: -d
    type: string
    default: "No description provided."
```

```bash
--priority, -p:
    description The priority of the dotfile
    alias: -p
    type: integer
    default: 0
```
---

## **dotman remove**
### *Description:*
Remove a dotfile from the project.


### *Usage:*
```bash
dotman remove <name>
```

### *Options:*

```bash
name:
    description: The name of the dotfile to be removed
    type: string
    required: true
```
---

## **dotman update**
### *Description:*
Update a dotfile from the project.


### *Usage:*
```bash
dotman update <name>
```

### *Options:*

```bash
name:
    description: The name of the dotfile to be updated
    type: string
    required: true
```

### *Special Usage:*  
```bash
dotman update @a
```
This usage will update all dotfiles.

---

## **dotman installer**
### *Description:*
Installer Utilities.

### *Usage:*
```bash
dotman installer <method> <name> [options]
```

### *Method:*
```bash
add:
    description: Add an installer to the project

    name:
        description: The name of the installer to be added
        type: string
        required: true
```
```bash
remove:
    description: Remove an installer from the project

    name:
        description: The name of the installer to be removed
        type: string
        required: true
```

### *Options:*

```bash
--description, -d:
    description: A description of the project
    alias: -d
    type: string
    default: "No description provided."
```
---

## **dotman command**
### *Description:*
Command Utilities.


### *Usage:*
```bash
dotman command <method> <name> <command> [options]
```

### *Method:*
```bash
add:
    description: Add a command to the project

    name:
        description: The name of the command to be added
        type: string
        required: true
    
    command:
        description: actual command to be run
        type: string
        required: true
```
```bash
remove:
    description: Remove a command from the project

    name:
        description: The name of the command to be removed
        type: string
        required: true
```

### *Options:*

```bash
--sudo, -s:
    description: Whether the command should be run with sudo
    alias: -s
    type: boolean
    default: false
```
---

## **dotman generate**
### *Description:*
Generate installer scripts.


### *Usage:*
```bash
dotman generate <installer name>
```

### *Installer Name:*
```bash
name of an installer defined in the project.
```
---

## **dotman status**
### *Description:*
Show the status of the project.


### *Usage:*
```bash
dotman status <type>
```

### *Type:*

```bash
normal:
    description: show basic indformation about the project
```

```bash
dotfiles:
    description: show extended indformation about the projects dotfiles
```

```bash
installers:
    description: show basic indformation about the projects installers
```
---