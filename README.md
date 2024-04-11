# Ouro
![stability-unstable](https://img.shields.io/badge/stability-unstable-red.svg)

Ouro is a command line utility for planning and executing changes within a software development project using LLMs.

## Installation
**MacOS**
```shell
$ git clone https://github.com/pbjer/ouro.git
$ cd ouro
$ make install
```
Running `make install` will install ouro to `/usr/local/bin/`

## Usage
To get started, initialize ouro at the root of your project. An `.ouro/` directory will be created at the root of the project.
```bash
$ cd path/to/project
$ ouro new
Initialized ouro
```
To plan a change to your project, start by loading directories or files relevant to the change
```bash
$ ouro load cli/editor.go cmd
cli/editor.go - 177
cmd/main.go - 63
TOTAL: 240
```
You can list what's currently loaded and an estimate of the token count
```bash
$ ouro list
cli/editor.go - 177
cmd/main.go - 63
TOTAL: 240
```
You can unload files or directories of files that you don't want to use
```bash
$ ouro unload cli
Unloaded cli/editor.go
```
By default, ouro uses `openai` with `gpt-4-0125-preview` for generation. There is also the option to run `groq` with `mixtral-8x7b-32768`. Support for more providers and models coming soon.
```bash
$ export OPENAI_AI_KEY="YOUR_OPENAI_API_KEY"
$ export GROQ_AI_KEY="YOUR_GROQ_API_KEY"
$ ouro use openai
Using openai
$ ouro use groq
Using groq
```
You can create changes to your project by running a plan, which will call the llm and use the loaded files as implementation references 
```bash
$ ouro plan "update main so that all errors are handled"
Request tokens: 2069
Updated cmd/main.go
```
You can reload context to include the latest changes since the last load
```bash
$ ouro reload
cli/editor.go - 226
cmd/main.go - 68
LOADED: 294
```