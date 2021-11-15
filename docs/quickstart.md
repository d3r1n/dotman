# Welcome to the Quick Start Guide!

In this guide we will learn how to use dotman to create dotfiles. \
This guide is intended to be a quick introduction to dotman and its features. \
We will cover the following topics:
- Installation
- Dotfiles
- Updating
- Installers
- Commands
- Generating
- Next Steps

Lets get started!

## Installation
Installing dotman is as simple as running the following command:

```bash
$ curl -L https://git.io/J68cu | bash
```
or
```bash
$ curl https://raw.githubusercontent.com/d3r1n/dotman/master/linux_installer.sh | bash
```

What this command does is download the latest version of dotman's auto installer and run it. \

## Dotfiles
Dotfiles are files that are used to store configuration information. \
Dotfiles are usually stored in the `.config` directory in your home directory but they can be anywhere. \
With dotfiles you can configure your system in a way you can constumize anthing and maintain. \

But maintaing dotfiles is not easy. \
You have to know what files you want to store and where. \
That's why dotman is here to help you. \

### Initializing
we will be using a github repo to store our dotfiles. \
To initialize dotman you need to run the following command:

```bash
$ mkdir MyDotfiles
$ dotman init MyDotfiles --git --remote="https://github.com/YOUR_USERNAME/YOUR_REPO" --branch="master"
```

this command will initialize dotman in your current working directory. \

### Adding Dotfiles
To add dotfiles you need to have a file to add... :D
So lets create a file called `.bashrc` and add it to the dotfiles repo:

```bash
$ touch .bashrc
$ echo "echo Hello World!" > .bashrc
```

Now we can add this file to the dotfiles repo:

```bash
$ dotman add my_bashrc .bashrc --description="My bashrc file"
```

You can see 2 things are generated:
- A new folder and file, `files/.bashrc`
- A new dotfile object in `dotman.json`

these files are used to store the dotfile. \
Yes! dotman does not uses symlinks to store dotfiles. \
It makes it easier to manage dotfiles and features. \

## Updating
Let's change something on our .bashrc file:
``bash
$ echo "echo Hello From dotman!" > .bashrc
```

You can see that the file has changed. \
but our dotfile has not. \

We can update our dotfile using the following command:
```bash
$ dotman update my_bashrc
```

Now everything is up-to-date.

## Installers
With dotman you can create installers to install dotfiles. \
You can create installers to install dotfiles from a git repo, \
You can do anything you want with installers. \
They are basically automaticly generated bash scripts that will install dotfiles and run your custom commands.

We will create an basic installer that installs our dotfiles:

```bash
$ dotman installer add my_installer --description="My installer"
```

## Commands

This command will create a new installer called `my_installer` and add it to `dotman.json`. \
But it does not have any custom commands. \
So why not add some?

```bash
$ dotman command add "echo I'm Pickle Riiiick!!"
```

## Generating

We added a dotfile, \
an installer, \
and a command. \
Now we need to generate the installer. \

In dotman installers will generate in the `installers` directory. \
with their own unique name: `installer[my_installer].sh` \

So let's generate the installer:
```bash
$ dotman generate my_installer
```

Tada...! You created your own dot with dotman! \
and you can use it to install dotfiles.

## Next Steps
You can learn more about dotman with reading the commands documentation. \
You can also learn more about dotman by reading the [documentation](commands)
