# SYSTEM MONITORING SERVICE

A CLI tool that can collect the system hardware data. It can collect: 

1. CPU data
2. Disk data
3. Connection data
4. Process data
5. System data

> Please note that this work the best on Linux platform. For Windows, some data cannot be collected, or missing. 

This tool can run as a standalone CLI, or it can send data to the server via HTTP request for online monitoring

For further information, you can run the following command:

```go   
go run main.go help
```

or if you prefer to run the executable:
```go
go build
./sys-check help
```

