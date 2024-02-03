# GO TODO [![build-project](https://github.com/GurYN/go-todo/actions/workflows/build-project.yml/badge.svg)](https://github.com/GurYN/go-todo/actions/workflows/build-project.yml)
A quick example of a TODO list app using Go language.

![Screenshot](/doc/medias/screenshot.png)

# Benchmark
Amazing results with 50 000 requests executed in less than 1 second!

Note: The benchmark was executed on a MacBook Pro M1 Pro and 16 GB of RAM using Apache HTTP server benchmarking tool.

![Benchmark](/doc/medias/benchmark.png)

# Run it
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
