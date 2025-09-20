// Filename: cmd/examples/cors/preflight/main.go
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
</head>
   <h1>Appletree Preflight CORS</h1>
    <div id="output"></div>
    <script>
		document.addEventListener('DOMContentLoaded', function() {
		fetch("http://localhost:4000/v1/quotes") {
			method: "POST",
			headers: {'Content-Type': 'application/json'}}
			body: JSON.stringify({"type":"funny", "quote":"I am fond of pigs. Dogs look up to us. Cats look down on us. Pigs treat us as equals.", "author":"Winston S. Churchill"
			}).then(function (response) {
				response.json().then(function (data) {
					 document.getElementById("output").innerhtml = JSON.stringify(data, null, 2);
				});
				}, function(err) {
					document.getElementById("output").innerhtml = err;
				}
			);
		});

  </script>
</body>
</html>`

// A very simple HTTP server
func main() {

	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()

	log.Printf("starting server on %s", *addr)

	err := http.ListenAndServe(*addr,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(html))
		}))
	log.Fatal(err)
}
