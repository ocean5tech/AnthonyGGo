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


## ssh无法连接，报错“Socket error Event: 32 Error: 10053.”

restart 路由器

## listen tcp :8082: bind: An attempt was made to access a socket in a way forbidden by its access permissions.
1. check if the port8082 is excluded by tcp
    netsh interface ipv4 show excludedportrange protocol=tcp

2. Change Port OR change excludeportrange
    netsh int ipv4 set dynamicport tcp start=50000 num=500


## SCERET
secret -->> eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU

[]byte(secret)  --->> [101 121 74 104 98 71 99 105 79 105 74 73 85 122 73 49 78 105 73 115 73 110 82 53 99 67 73 54 73 107 112 88 86 67 74 57 46 101 121 74 109 98 50 56 105 79 105 74 105 89 88 73 105 76 67 74 117 89 109 89 105 79 106 69 48 78 68 81 48 78 122 103 48 77 68 66 57 46 117 49 114 105 97 68 49 114 87 57 55 111 112 67 111 65 117 82 67 84 121 52 119 53 56 66 114 45 90 107 45 98 104 55 118 76 105 82 73 115 114 112 85]


## JWT  : 使用jwt-go的时候遇到了一个报错
1. 注意tokenClaims.SignedString(jwtSecret)的参数jwtSecret为[]byte()而不是string
2. 注意jwt.NewWithClaims的加密方法用的是jwt.SigningMethodHS256而不是jwt.SigningMethodES256

## JWT TEST

&{dc:0xc000164090 releaseConn:0xacd300 rowsi:0xc000001680 cancel:<nil> closeStmt:<nil> closemu:{w:{state:0 sema:0} writerSem:0 readerSem:0 readerCount:0 readerWait:0} closed:false lasterr:<nil> lastcols:[]}
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo4NDk4MDgxLCJleHBpcmVzQXQiOjE1MDAwfQ.DkQ7O_10WnExoVTQt1AKiS1rECVQiZWnm8DTlU508U8
<nil>
JWT Token:  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo4NDk4MDgxLCJleHBpcmVzQXQiOjE1MDAwfQ.DkQ7O_10WnExoVTQt1AKiS1rECVQiZWnm8DTlU508U8

## Debug tip
panic("dd") 
panic(relfect.TypeOf(claims["accountNumber"]))
