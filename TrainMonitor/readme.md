## Usage
```go
func someHandler(r *http.Request, w *http.ResponseWriter) {
	// Create socket connection.
	conn, err := net.Dial(network, address)
	if err != nil {
		// error handling
	}
	
	// Binding epoch data from request body
	var epoch Epoch
	err = epoch.BindEpoch(r)
	if err != nil {
		// error handling
	}
	
	// Send to client through socket connection.
	err = epoch.PushToSocket(conn)
	if err != nil {
		// error handling
	}
}
```