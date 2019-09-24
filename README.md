# Painless Script Testing and Execution

Elasticsearch Painless Script lacks a bit of tooling. This project allows you developing Painless scripts "painlessly".
It is a  command line tool. `plsd`
You can use `plsd` in your build tools, continuous deployment pipelines(jenkins etc). `plsd` would help writing scripts and with its `--watch` feature whenever you change your script/parameters file you can see its output from Elasticsearch. 

# plsd

``` bash
Painless Script Development Toolkit

Usage:
  plsd [command]

Available Commands:
  exec        Executes Painless Script and returns output
  help        Help about any command
  perf        This command can be used to help understanding impacts of sorting with painless script. This feature is experimental. Output/methodology might change
  test        Runs test files woth the script

Flags:
      --es-endpoint string   Elasticsearch Painless Execution API Endpoint (default "http://localhost:9200/_scripts/painless/_execute")
  -h, --help                 help for plsd

Use "plsd [command] --help" for more information about a command.

```


# Install

## Linux
 You can either download the .deb or .rpm from the [releases page](https://github.com/csenol/plsd/releases) and install with dpkg -i and rpm -i respectively or use homebrew same as macOS
## macOS
You can install it by using homebrew taps  
`brew tap csenol/plsd`  
`brew install plsd`
## Build from source
If you have go 1.12>= installed  
`go install ./...`

