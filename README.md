# Basic POS

Basic POS with UI to "sell" products.

## Getting Started

Populate the database with the following command:

```bash
go run ./scripts/import.go -path items-sale.csv
```

The CSV file must have the following format:

```csv
description,qty,price
item,1,200
```

## Run

Run the server with the following command:
```bash
    go run main.go -port 8080
```
 
## Run with Docker
```bash
  docker build --tag pos .
  docker run --net=host --rm --name pos pos
```

## Screenshots

![img.png](img.png)
