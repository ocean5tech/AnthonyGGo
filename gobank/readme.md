## Run Make&Makefile on Windows

<https://earthly.dev/blog/makefiles-on-windows/>

1. Start Powershell as Admin

2. Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

3. choco install make

## ListenAndServe exit immediately without blocking

1. use Println to get error from ListenAndServe

Println(http.ListenAndServe(":3000", nil))

## go run VS go build

1. `go run main.go` only run main.go , `go build` check all go file in same folder

## Tag : should not include space between  [:"]

   Right:  `json:"firstName"`
   Wrong:  `json: "firstName"`

## Use Add() to setup content-type, not set

        func WriteJSON(w http.ResponseWriter, status int, v any) error {
            w.Header().Add("Content-Type", "application/json")
            w.WriteHeader(status)
            //w.Header().Set("Content-Type", "application/json")
            return json.NewEncoder(w).Encode(v)
        }

## Install docker

1. sudo  docker run --name some-postgres -e POSTGRES_PASSWORD=gobank -p 5432:5432 -d postgres

2. sudo docker ps
 CONTAINER ID   IMAGE      COMMAND                  CREATED              STATUS              PORTS                                       NAMES
b41611b671a3   postgres   "docker-entrypoint.s…"   About a minute ago   Up About a minute   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp   some-postgres

3. go get github.com/lib/pq
https://pkg.go.dev/github.com/lib/pq

## dial tcp 54.250.166.42:5432: connectex: No connection could be made because the target machine actively refused it.

1. restart server
2. restart docker
    `sudo docker start some-postgres`

## JWT

https://pkg.go.dev/github.com/golang-jwt/jwt/v4

1. go get -u github.com/golang-jwt/jwt/v4
2. import "github.com/golang-jwt/jwt/v4"

## JWT Secret
1. windows : set environment variables
2. Linux: export JWT_SECRET=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU