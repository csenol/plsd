# Painless Script Testing and Execution

Elasticsearch Painless Script lacks a bit of tooling. This project allows you developing Painless scripts "painlessly".
It has to command line programs. `plst` for testing and `plte` for executing scripts continuously. 
You can use `plst` in your build tools, continuous deployment pipelines(jenkins etc)
`plte` would help writing scripts and with its `-watch` feature whenever you change your script/parameters file you can see its output from Elasticsearch. 

# plst
An example use-case for plst is as follow.


