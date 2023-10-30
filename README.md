# proxy-go

## Project Description

This project is a proxy server written in the Go programming language. It leverages the power of Docker Compose for effortless deployment and management. The server collects metrics using Prometheus, while Grafana is used for visually displaying these metrics.

![example](https://github.com/MaksimPozharskiy/proxy-go/assets/47855527/46ade724-da27-4123-aba0-b0016a1efa49)


### Usage:

Fill `./.env` file with environment variables as in the example: `./example.env` 

Run in docker:
```shell
cd ./deployments && sudo docker compose up --build
```

### Metrics Available in Grafana

- **HTTP Requests**: Track the total number of incoming HTTP requests to the proxy server. 
- **HTTP Response Codes**: Monitor different HTTP response codes returned by the server. Detect and troubleshoot potential issues occurring during the request-response cycle.
- **Average HTTP Request Time**: Calculate the average duration of HTTP requests processed by the server. Understand the server's performance and responsiveness.
- **Average Throughput**: Measure the average amount of data (in bytes) transferred per unit of time. Evaluate the server's capacity to handle data efficiently.
- **CPU Load**: Keep track of the CPU utilization of the proxy server. Observe the percentage of CPU resources being utilized, enabling performance monitoring and optimization.
- **Memory Usage**: Monitor the amount of memory consumed by the proxy server. Ensure resource allocation efficiency and identify any potential memory leaks.

By utilizing Prometheus and Grafana, this project provides comprehensive insights and aesthetically pleasing visualizations for monitoring the proxy server's performance. It enables identifying bottlenecks, optimizing resource allocation, and ensuring the seamless operation of the application.

Experience the power of this proxy server project and empower your application with efficient and customizable monitoring capabilities!
