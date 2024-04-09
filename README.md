# Ouro

Ouro is a command line utility for planning and executing changes within a software development project using LLMs.

## Installation
**MacOS**
```shell
git clone github.com/pbjer/ouro
cd ouro
make install
```
Running `make install` will install ouro to `/usr/local/bin/`

## Usage
To get started, initialize ouro at the root of your project.
```bash
> cd path/to/project
> ouro new
```
To plan a change to your project, start by loading directories or files relevant to the change
```bash
> ouro load cli/editor.go cmd
Loaded cli/editor.go
Loaded cmd/main.go
```
You can list what's currently loaded and an estimate of the token count
```bash
> ouro list
cli/editor.go - 177
cmd/main.go - 63
```
You can unload files or directories of files that you don't want to use
```bash
> ouro unload cli
Unloaded cli/editor.go
```
You can create changes to your project by running a plan, which will call the llm and use the loaded files as implementation references 
```bash
> export OPENAI_API_KEY="YOUR-API-KEY"
> ouro plan "update main so that all errors are handled"
Updated cmd/main.go
```
