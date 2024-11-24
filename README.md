# XM-Company

A simple project to demonstrate microservice development in Go with PostgreSQL.

## Pre-requisites:
Before you begin, make sure you have the following installed:

1. **make** - A build automation tool.
2. **docker** - A platform for developing, shipping, and running applications in containers.
3. **docker-compose** - A tool for defining and running multi-container Docker applications.
4. **Internet connection** - Required to pull Docker images from the internet.

## A. Steps to run the service:
To get the service up and running, follow these steps:

1. Clone the repository:

    ```bash
    $ git clone <repository_url>
    $ cd <xm-company>  # Navigate to the repo directory
    ```

2. Run the following command to start the services:

    ```bash
    $ make up
    ```

    This command will:
    - Build the Docker containers.
    - Pull any necessary Docker images.
    - Start the services in the containers.

    **Note**: It might take some time to build and start the containers for the first time.

## B. Steps to run integration tests:
Once the service is running, you can run the integration tests using the following command:

```bash
$ go test -tags=integration ./...
```


## B. Steps to run unit tests:
Unit tests won't need any db service running:

```bash
$ go test -tags=init ./...
```