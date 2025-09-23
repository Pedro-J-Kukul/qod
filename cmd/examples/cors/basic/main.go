package main

import (
	"flag"
	"log"
	"net/http"
)

const html = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Health Check Web App</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
        }
        #output {
            background-color: #f5f5f5;
            padding: 15px;
            border-radius: 8px;
            border: 1px solid #ddd;
            min-height: 100px;
        }
        button {
            background-color: #4CAF50;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            margin: 10px 0;
        }
        button:hover {
            background-color: #45a049;
        }
        .loading {
            color: #666;
            font-style: italic;
        }
        .error {
            color: red;
        }
        .success {
            color: green;
        }
    </style>
</head>
<body>
    <h1>Lab 5 - Health Check Web App</h1>
    <button onclick="checkHealth()">Check Health</button>
    <pre id="output">Press "Check Health" to get API status</pre>
    
    <script>
        function checkHealth() {
            const output = document.getElementById("output");
            output.textContent = "Loading...";
            output.className = "loading";
            
            // Make sure this matches your API server port
            fetch("http://localhost:4000/v5/healthcheck")
                .then(function(response) {
                    console.log('Response status:', response.status);
                    
                    if (!response.ok) {
                        throw new Error('HTTP ' + response.status + ': ' + response.statusText);
                    }
                    
                    return response.json();
                })
                .then(function(data) {
                    console.log('Response data:', data);
                    output.textContent = JSON.stringify(data, null, 2);
                    output.className = "success";
                })
                .catch(function(err) {
                    console.error('Fetch error:', err);
                    output.textContent = "Error: " + err.message;
                    output.className = "error";
                });
        }
    </script>
</body>
</html>`

// A very simple HTTP server
func main() {
	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()

	log.Printf("starting web server on %s", *addr)
	log.Printf("open browser to: http://localhost:%s", (*addr)[1:])

	err := http.ListenAndServe(*addr,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(html))
		}))
	log.Fatal(err)
}
