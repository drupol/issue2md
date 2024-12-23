# issue2md

A cli and web tool to convert GitHub issue into Markdown.

## Command-line mode

### Install issue2md cli

```
$go install github.com/bigwhite/issue2md/cmd/issue2md@latest
```

### Convert issue to markdown

```
Usage: issue2md issue-url [markdown-file]
Arguments:
  issue-url      The URL of the github issue to convert.
  markdown-file  (optional) The output markdown file.
```

## Web mode

### Install and run issue2md web

```
$git clone https://github.com/bigwhite/issue2md.git
$make web
$./issue2mdweb   
Server is running on http://0.0.0.0:8080
```

### Convert issue to markdown

Open localhost:8080 with the browser: 

![](./screen-snapshot.png)

Input the issue url you want to convert and click "Convert" button!
