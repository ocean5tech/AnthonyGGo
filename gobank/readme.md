## Run Make&Makefile on Windows

https://earthly.dev/blog/makefiles-on-windows/

1. Start Powershell as Admin

2. Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

3. choco install make


## ListenAndServe exit immediately without blocking

1. use Println to get error from ListenAndServe

Println(http.ListenAndServe(":3000", nil))

## go run VS go build

1. `go run main.go` only run main.go , `go build` check all go file in same folder 