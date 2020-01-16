# Rabbitr

Small CLI application written in GoLang for easier management of RabbitMQ related tasks. 

## Implemented features

 - connections
     - [x] list
     - [x] close
 - queue  
     - [x] list
     - [x] purge  
     - [x] sync 
     - [x] delete
     - [x] move
 - server
     - [x] add
     - [x] delete
     - [x] list
 - policies     
     - [x] list
 
## Installation

#### Making sure the environment is set
`export HOMEBREW_GITHUB_API_TOKEN=<MY_GITHUB_TOKEN>`

#### Tapping the repo
`brew tap smartrecruiters/homebrew-taps`

#### Installing!
`brew update && brew install rabbitr`
 
## Configuration
After downloading or building the source code `rabbitr` application needs to be configured with coordinates to 
RabbitMQ servers that it will operate on. It uses REST API when communicating with RabbitMQ server. 
To add server configuration invoke:

`rabbitr server add -s my-server-name -api-url http://localhost:15672 -u user -p pass`

After server has been configured it can be used in context of other commands such as `queues`

## Usage
Each command comes with a description and examples. Start with `rabbitr -h` to check all the commands. 
Lower level commands provide their own usage, for example `rabbitr queues -h` or `rabbitr queues list -h`

Example commands:

```
rabbitr server add -s my-server-name -api-url http://localhost:15672 -u user -p pass

rabbitr queues list -s my-server-name

rabbitr queues sync -s my-server-name

rabbitr queues list -s my-server-name --filter="queue.Consumers==0"

rabbitr queues list -s my-server-name --filter="queue.Consumers==0 && queue.Messages>=200"

rabbitr queues purge -s my-server-name --filter="queue.Consumers==0 && queue.Messages>=200"

rabbitr queue move -s my-server-name --src-vhost vh1 --src-queue q1 --dst-vhost vh2 --dst-queue q2

```

## Advanced filtering
Rabbitr uses [goevaluate](https://github.com/Knetic/govaluate#govaluate) library for dynamic filtering of the resources.
It can be useful to determine list of subjects that match given criteria.
It allows for creating flexible conditions considering for example only queues with particular name, vhost, defined number of consumers or messages.
Check command descriptions for a list of properties available for use on given resource type.     
