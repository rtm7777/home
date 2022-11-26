# Golang App for some kind of smart home

## Install
- install golang https://go.dev/doc/install
- install packages
  `go install`

## Build
- build container
  `docker build --tag {{NAS_ip:registry_port}}/home_app .`
- push it to NAS docker registry
  `docker push {{NAS_ip:registry_port}}/home_app:latest`

## Run
- ssh to raspberry ssh `pi@{{RPI_ip}}`
- clone this repo with `https`
- pull docker container: `docker pull {{NAS_ip:registry_port}}/home_app:latest`
- start it: `docker-compose up -d`
- to stop: `docker-compose stop`