# Pension Data API

The Pension Data API lets you easily access information and quotes about pension funds with tax benefits available in Belgium and Luxembourg.

This project is the REST API behind [api.pensiondata.eu](https://api.pensiondata.eu) built using Go and PostgreSQL.

For the documentation on how to use the API checkout [the Pension Data website](https://www.pensiondata.eu).

# Development
## Docker

Build the image 
`docker build -t pensiondata-api .`

Run the container
`docker run -it -p 8080:8080 --env-file .env pensiondata-api`

If the database is on localhost the following environment variable need to be modified in the .env file: `DATABASE_HOST=host.docker.internal`

## Build on Windows for Linux

`set GOOS=linux`
`set GOARCH=amd64`
`go build -ldflags "-s -w" -o pensiondata-api`
