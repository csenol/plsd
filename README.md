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


# Examples


## Script Execution and Development
Let's create an index as in Elasticsearch [docs](https://www.elastic.co/guide/en/elasticsearch/painless/7.3/painless-execute-api.html)

``` bash
curl -X PUT "localhost:9200/my-index?pretty" -H 'Content-Type: application/json' -d'
{
  "mappings": {
      "properties": {
        "field": {
          "type": "keyword"
        },
        "rank": {
          "type": "long"
        }
      }
  }
}
'
```


With a context file example-context-file.json

``` json
{
    "index": "my-index",
    "context": "score",
    "document": {
        "rank": 4
    },
    "params" : {
        "max_rank": 5
    }

}
```


and with a painless script as example-script.painless  

``` groovy
(double)doc['rank'].value / params.max_rank
```

Running script with *plsd*  
``` bash
plsd exec --context-file example-context-file.json --script-file example-script.painless
0.8
```


You can also *watch* files and run script with every single change. This can be useful with development

``` bash
plsd exec --context-file example-context-file.json --script-file example-script.painless --watch
```
## Testing 

plsd can be also to write tests for painless scripts
An example test file is as follows  

``` json
[
    {
        "description":"This Test Should pass",
        "index": "my-index",
        "params": {
            "max_rank":5
        },
        "document": {
            "rank":4
        },
        "expected_result": 0.8
    },
    {
        "description":"This Test Should FAIL",
        "index": "my-index",
        "params": {
            "max_rank":0
        },
        "document": {
            "rank":4
        },
        "expected_result": 5
    }
]

```

And you can run the test with plsd

``` bash
plsd test --test-file example-test.json < example-script.painless
Test Passed: This Test Should pass
This Test Should FAIL 
 TestCase Failed at example-test.json
 Expected 5.000000 Got Infinity
```

## REPL
I did not bother creating a REPL for painless but it can be achieved simply by  
`alias plsd-repl='while true; do printf ">" ; plsd exec --index my-index ; echo "" ; done'`  
Where you can have an alias with your favorite index and context-params etc.
