# issue2md

A CLI and web tool to convert GitHub issues into Markdown format.

[中文文档](./README-zh.md)

## Command-line Mode

### Install issue2md CLI

```bash
$ go install github.com/bigwhite/issue2md/cmd/issue2md@latest
```

### Convert an Issue to Markdown

```bash
Usage: issue2md issue-url [markdown-file]
Arguments:
  issue-url      The URL of the GitHub issue to convert.
  markdown-file  (optional) The output markdown file.
```

## Web Mode

### Install and Run issue2md Web

```bash
$ git clone https://github.com/bigwhite/issue2md.git
$ make web
$ ./issue2mdweb
Server is running on http://0.0.0.0:8080
```

### Convert an Issue to Markdown

Open `localhost:8080` in your browser:

![Screenshot](./screen-snapshot.png)

Input the issue URL you wish to convert and click the "Convert" button!
