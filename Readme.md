# Golang App for some kind of smart home

## Build
- build container
  `docker build --tag {{NAS_ip:registry_port}}/home_app .`
- push it to NAS docker registry
  `docker push {{NAS_ip:registry_port}}/home_app:latest`

## Run
- ssh to raspberry ssh `pi@{{NAS_ip}}`
- clone this repo with `https`
- pull docker container: `docker pull {{NAS_ip:registry_port}}/home_app:latest`
- start it: `docker-compose up -d`
- to stop: `docker-compose stop`