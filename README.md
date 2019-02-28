# mini-redis

Simple Go implementation of a subset of Redis commands.

## How to run the application

To use the application, you must first build a Docker image. From the shell, go to the project directory and type:

```
docker build -t miniredis .
```

Then, you can run the application in a container with the following instruction:

```
docker run --rm -ti miniredis
```

The container exposes the port 8080 to handle http requests too. To be able to interact with the container, run the container with the *-p* option:

```
docker run --rm -ti -p 8080:8080 miniredis
```

After that, you can send http requests like:

```
curl http://localhost:8080/?cmd=INCR%20xyz
> 1
```

## How to run tests

You can use Docker to run the tests. From the shell, just change directory to the project and run:

```
docker run --rm -ti -v "$PWD":"/go/src/miniredis" -w "/go/src/miniredis" golang:1.12 go test
```
