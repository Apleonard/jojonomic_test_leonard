# jojonomic_test_leonard

## Requirements
  - GO 1.19
  - Go Module
  - Postgresql Database
  - kafka sarama (github.com/Shopify/sarama)

## Installing
  - Install Dependecies

    console in each microservices
     ```
     $ go mod tidy
     ```
  - Fill your local .env configuration in each microservices

## Run Program
  - Run Project

    console in each microservices
     ```
     $ go run cmd/main.go
     ```
  - Run Docker
     ```
    docker-compose -f ./misc/docker-compose.yaml up -d
     ```

## Postman Collection
    https://www.getpostman.com/collections/29c97b71c76c0dc58558


## Test Run Video
    https://www.loom.com/share/5a1644a93ec1446aade318a11b2ce173