# Simple Go HTTP Load Balancer

This project implements a basic HTTP load balancer in Go, demonstrating core concepts such as round-robin traffic distribution, health checking of backend servers, and request logging. It is built using only Go's standard library, making it lightweight and easy to understand.

## Features

-   **HTTP Request Distribution**: Accepts incoming HTTP requests on port `8080` and forwards them to a configurable list of backend servers.
-   **Round-Robin Algorithm**: Distributes requests to backend servers sequentially, ensuring an even load distribution.
-   **Health Checks**: Periodically checks the health of registered backend servers. If a server is detected as down, it is temporarily removed from the rotation until it becomes healthy again.
-   **Request Logging**: Logs each forwarded request, including the backend server that handled the request.
-   **Standard Library Only**: Built exclusively with Go's standard packages, avoiding external dependencies for simplicity and transparency.

## Project Structure

-   `main.go`: Contains the core logic for the load balancer, including the `Backend` and `ServerPool` structs, health check mechanisms, and the main HTTP handler.
-   `backend1.go`, `backend2.go`, `backend3.go`: Simple HTTP servers acting as backend services. Each backend listens on a different port (`9001`, `9002`, `9003` respectively) and provides a basic `/` endpoint and a `/health` endpoint for health checks.

## How to Run

To compile and run the load balancer and its backend servers, follow these steps:

### Prerequisites

-   Go (version 1.18 or higher) installed on your system.

### 1. Clone the Repository (if applicable)

If you received these files as part of a repository, clone it first:

```bash
git clone <repository_url>
cd <repository_directory>
```

### 2. Compile the Backend Servers

Open your terminal or command prompt, navigate to the project directory, and compile each backend server:

```bash
go build -o backend1 backend1.go
go build -o backend2 backend2.go
go build -o backend3 backend3.go
```

This will create executable files named `backend1`, `backend2`, and `backend3` in your project directory.

### 3. Compile the Load Balancer

In the same directory, compile the load balancer:

```bash
go build -o loadbalancer main.go
```

This will create an executable file named `loadbalancer`.

### 4. Run the Backend Servers

It's recommended to run each backend server in a separate terminal window or in the background, as they will continuously listen for requests:

```bash
./backend1 &
./backend2 &
./backend3 &
```

### 5. Run the Load Balancer

Finally, run the load balancer. It will start listening on port `8080`:

```bash
./loadbalancer &
```

## How to Test

Once all servers are running, you can test the load balancer's functionality using `curl` or by accessing `http://localhost:8080` in your web browser.

### Basic Request Distribution

Send multiple requests to the load balancer. You should observe that requests are distributed in a round-robin fashion among the active backend servers:

```bash
for i in $(seq 1 10); do curl http://localhost:8080; echo; done
```

Expected output will show responses from `Backend 1`, `Backend 2`, and `Backend 3` alternating.

### Health Check Functionality

To test the health check mechanism, stop one of the backend servers (e.g., `backend1`). The load balancer should detect that `backend1` is down and stop forwarding requests to it. Requests will then only be distributed between the remaining healthy servers.

After a few seconds, restart the stopped backend server. The load balancer's health check will detect that it's back online and include it in the rotation again.

To stop a backend server, you can find its process ID (PID) and kill it. For example, to stop `backend1`:

```bash
pgrep backend1 # Find the PID of backend1
kill <PID>     # Replace <PID> with the actual process ID
```

Alternatively, if you ran them in separate terminal windows, you can simply close the terminal window for the desired backend.

