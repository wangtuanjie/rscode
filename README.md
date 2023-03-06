# RSCODE
Rscode provides a command-line tool that uses the OpenAI ChatGPT API to suggest improvements on function names,
variable names, comments, and log content of a given code while keeping its structure.

Installation:
```
1. go get github.com/wangtuanjie/rscode
2. export OPENAI_API_KEY=<your_key_here>
```


Example usage:

```shell
cat main.go | rscode

rscode -f main.go -w
```
