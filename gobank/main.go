package main

func main() {

	server := NewAPIServer(":8082")
	server.Run()

}
