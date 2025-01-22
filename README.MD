# Twitter-Reddit Clone

A social media platform combining features from Twitter and Reddit, built with a Go backend and Vue.js frontend.

## Getting Started

Follow these steps to set up the project on your local machine for development and testing.

### Prerequisites
- Go 1.22 or higher
- Node.js and npm
- SQLite3

### Backend Setup
1. Start the backend server:
```bash
go run main.go
```

### Frontend Setup

1. Navigate to the frontend directory:

```bash
   cd frontend
```

2. Install the required dependencies:

```bash
   npm install
```

3. Start the development server:

```bash
   npm run serve
```

### Test Performance

```bash
wrk -t8 -c200 -d30s http://localhost:8080 
```


### Test Parameters
- **Duration**: 30 seconds  
- **Target URL**: `http://localhost:8080`  
- **Threads**: 8  
- **Connections**: 200  

### Test Output
| Metric           | Average  | Standard Deviation | Maximum  | Percentage (+/- Stdev) |
|------------------|----------|-------------------|----------|------------------------|
| **Latency**      | 60.77 ms | 8.70 ms           | 124.19 ms | 89.07%                |
| **Requests/Sec** | 412.83   | 79.62             | 505.00   | 57.29%                |

### Summary Statistics
- **Total Requests**: 98,678  
- **Total Data Read**: 84.88 MB  
- **Requests per Second**: 3,286.71  
- **Data Transfer per Second**: 2.83 MB  

These results indicate that the server can handle a significant number of concurrent requests efficiently.
