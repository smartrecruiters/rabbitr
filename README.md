<p align="center">
	<h1 align="center">rabbitr</h1>
	<p align="center">
		<a href="https://travis-ci.org/smartrecruiters/rabbitr"><img alt="Build" src="https://travis-ci.org/smartrecruiters/rabbitr.svg?branch=master"></a>	
		<a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>	
		<a href="https://goreportcard.com/report/github.com/smartrecruiters/rabbitr"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/smartrecruiters/rabbitr?style=flat-square"></a>
		<a href="http://godoc.org/github.com/smartrecruiters/rabbitr"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-brightgreen.svg?style=flat-square"></a>
	</p>
</p>
<!-- MarkdownTOC autolink="true" bracket="round" autoanchor="true" -->

- [Rabbitr](#rabbitr)
  - [Implemented features](#implemented-features)
  - [Installation](#installation)
    - [Tapping the repo](#tapping-the-repo)
    - [Installing](#installing)
  - [Configuration](#configuration)
  - [Usage](#usage)
  - [Advanced filtering](#advanced-filtering)
  - [TODO Ideas](#todo-ideas)

<!-- /MarkdownTOC -->

<a id="rabbitr"></a>
# Rabbitr

Small CLI application written in GoLang for easier management of RabbitMQ related tasks. 

<a id="implemented-features"></a>
## Implemented features

 - connections
     - [x] list
     - [x] close
 - messages
     - [x] move
 - queue  
     - [x] list
     - [x] purge  
     - [x] sync 
     - [x] delete
 - exchange  
     - [x] list
     - [x] delete
 - server
     - [x] add
     - [x] delete
     - [x] list
 - shovels     
     - [x] delete
     - [x] list
 - policies     
     - [x] list
 
<a id="installation"></a>
## Installation

<a id="tapping-the-repo"></a>
#### Tapping the repo
`brew tap smartrecruiters/public-homebrew-taps git@github.com:smartrecruiters/public-homebrew-taps.git`

<a id="installing"></a>
#### Installing
`brew update && brew install rabbitr`
 
<a id="configuration"></a>
## Configuration
After downloading or building the source code `rabbitr` application needs to be configured with coordinates to 
RabbitMQ servers that it will operate on. It uses REST API when communicating with RabbitMQ server. 
To add server configuration invoke:

`rabbitr server add -s my-server-name -api-url http://localhost:15672 -u user -p pass`

After the server has been configured it can be used in context of other commands such as `queues`

<a id="usage"></a>
## Usage
Each command comes with a description and examples. Start with `rabbitr -h` to check all the commands. 
Lower level commands provide their own usage, for example `rabbitr queues -h` or `rabbitr queues list -h`

!["Example flow"](rabbitr-demo.gif)

Example commands:

```
rabbitr server add -s my-server-name -api-url http://localhost:15672 -u user -p pass
rabbitr queues list -s my-server-name
rabbitr queues sync -s my-server-name
rabbitr queues list -s my-server-name --filter="queue.Consumers==0"
rabbitr queues list -s my-server-name --filter="queue.Consumers==0 && queue.Messages>=200"
rabbitr queues purge -s my-server-name --filter="queue.Consumers==0 && queue.Messages>=200"
rabbitr messages move -s my-server-name --src-vhost vhost1 --src-queue test-queue --duplicate"
```

<a id="advanced-filtering"></a>
## Advanced filtering
Rabbitr uses [goevaluate](https://github.com/Knetic/govaluate#govaluate) library for dynamic filtering of the resources.
It can be useful to determine list of subjects that match given criteria.
It allows for creating flexible conditions considering for example only queues with particular name, vhost, defined number of consumers or messages.
Check command descriptions for a list of properties available for use on given resource type.     

<a id="todo-ideas"></a>
## TODO Ideas
- move messages from one queue to the other and strip some headers on the way
- dump messages from a queue to a file
