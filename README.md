# BayarInd Book

A simple RESTful API for managing books and authors with user authentication and authorization using JWT token.

## Overview
there's schema of the apps
<img width="1245" alt="Screenshot 2024-09-24 at 19 19 51" src="https://github.com/user-attachments/assets/04ee7539-1a94-44e3-9f81-83a83d429f43">

## Build instructions
### Prerequisites
Clone the repository 
```
git clone https://github.com/storyofhis/bayarind-book.git
```
### Run 
to run applications independently 
```
go run cmd/main.go
```
to run with docker 
```
task compose
```
### Tools

Install the required tools by running the following command:

```shell
task tools
```

### Generate Mocks

To generate mock files that are mainly used for testing, run the following command:
```shell
task mocks
```

### Linting

To check the code for linting errors, run the following command:
```shell
task lint
```

### Unit Tests

To run the unit tests, run the following command:
```shell
task test:unit
```

### Coverage

To run the tests and generate the coverage report, run the following command:
```shell
task coverage
```
