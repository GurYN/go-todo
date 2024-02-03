# GO TODO [![build-project](https://github.com/GurYN/go-todo/actions/workflows/build-project.yml/badge.svg)](https://github.com/GurYN/go-todo/actions/workflows/build-project.yml) [![Publish Docker image](https://github.com/GurYN/go-todo/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/GurYN/go-todo/actions/workflows/docker-publish.yml)

A quick example of a TODO list app using Go language. The app include a REST API and a dynamic web interface using websockets.

![Screenshot](/doc/medias/demo.gif)

# Benchmark
Amazing results with 50 000 requests executed in less than 1 second!

Note: The benchmark was executed on a MacBook Pro M1 Pro and 16 GB of RAM using Apache HTTP server benchmarking tool.

![Benchmark](/doc/medias/benchmark.png)

# Run it
## Local
Get the dependencies:
```bash
go mod download
```

Copy the `.env.example` file to `.env` and set the environment variables:
```bash
cp .env.example .env
```

And execute the makefile:
```bash
make
```
## Docker
Image available on Docker Hub: [vcibelli/go-todo](https://hub.docker.com/r/vcibelli/go-todo)
```bash
docker run -p 3000:3000 -e SERVER_PORT=3000 -e API_URL=http://localhost:3000 --name go-todo -d vcibelli/go-todo:latest
```