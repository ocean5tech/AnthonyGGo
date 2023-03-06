# RPC 版 “Hello, World”

 Hello 方法必须满足 Go 语言的 RPC 规则：方法只能有两个可序列化的参数，其中第二个参数是指针类型，并且返回一个 error 类型，同时必须是公开的方法。

    type HelloService struct {}

    func (p *HelloService) Hello(request string, reply *string) error {
        *reply = "hello:" + request
        return nil
    }

然后就可以将 HelloService 类型的对象注册为一个 RPC 服务：
rpc.Register 函数调用会将对象类型中所有满足 RPC 规则的对象方法注册为 RPC 函数，所有注册的方法会放在 “HelloService” 服务空间之下。然后我们建立一个唯一的 TCP 连接，并且通过 rpc.ServeConn 函数在该 TCP 连接上为对方提供 RPC 服务。

    func main() {
        rpc.RegisterName("HelloService", new(HelloService))

        listener, err := net.Listen("tcp", ":1234")
        if err != nil {
            log.Fatal("ListenTCP error:", err)
        }

        conn, err := listener.Accept()
        if err != nil {
            log.Fatal("Accept error:", err)
        }

        rpc.ServeConn(conn)
    }

下面是客户端请求 HelloService 服务的代码：

    func main() {
        client, err := rpc.Dial("tcp", "localhost:1234")
        if err != nil {
            log.Fatal("dialing:", err)
        }

        var reply string
        err = client.Call("HelloService.Hello", "hello", &reply)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(reply)
    }

# 更安全的 RPC 接口

重构 HelloService 服务，第一步需要明确服务的名字和接口：
 RPC 服务的接口规范分为三个部分：首先是服务的名字，然后是服务要实现的详细方法列表，最后是注册该类型服务的函数。为了避免名字冲突，我们在 RPC 服务的名字中增加了包路径前缀（这个是 RPC 服务抽象的包路径，并非完全等价 Go 语言的包路径）。RegisterHelloService 注册服务时，编译器会要求传入的对象满足 HelloServiceInterface 接口。

    const HelloServiceName = "path/to/pkg.HelloService"

    type HelloServiceInterface interface {
        Hello(request string, reply *string) error
    }

    func RegisterHelloService(svc HelloServiceInterface) error {
        return rpc.RegisterName(HelloServiceName, svc)
    }

在定义了 RPC 服务接口规范之后，客户端就可以根据规范编写 RPC 调用的代码了：

    func main() {
        client, err := rpc.Dial("tcp", "localhost:1234")
        if err != nil {
            log.Fatal("dialing:", err)
        }

        var reply string
        err = client.Call(HelloServiceName+".Hello", "hello", &reply)
        if err != nil {
            log.Fatal(err)
        }
    }

通过 client.Call 函数调用 RPC 方法依然比较繁琐，同时参数的类型依然无法得到编译器提供的安全保障。
为了简化客户端用户调用 RPC 函数，我们在可以在接口规范部分增加对客户端的简单包装：
    type HelloServiceClient struct {
        *rpc.Client
    }

    var _ HelloServiceInterface = (*HelloServiceClient)(nil)

    func DialHelloService(network, address string) (*HelloServiceClient, error) {
        c, err := rpc.Dial(network, address)
        if err != nil {
            return nil, err
        }
        return &HelloServiceClient{Client: c}, nil
    }

    func (p *HelloServiceClient) Hello(request string, reply *string) error {
        return p.Client.Call(HelloServiceName+".Hello", request, reply)
    }

接口规范中针对客户端新增加了 HelloServiceClient 类型，该类型也必须满足 HelloServiceInterface 接口，这样客户端用户就可以直接通过接口对应的方法调用 RPC 函数。同时提供了一个 DialHelloService 方法，直接拨号 HelloService 服务。

基于新的客户端接口，我们可以简化客户端用户的代码：

    func main() {
        client, err := DialHelloService("tcp", "localhost:1234")
        if err != nil {
            log.Fatal("dialing:", err)
        }

        var reply string
        err = client.Hello("hello", &reply)
        if err != nil {
            log.Fatal(err)
        }
    }
