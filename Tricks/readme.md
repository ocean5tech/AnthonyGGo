# 迭代
for range 对字符串的迭代模拟实现：for range 迭代字符串时，每次解码一个 Unicode 字符，然后进入 for 循环体，遇到崩坏的编码并不会导致迭代停止。

    func forOnString(s string, forBody func(i int, r rune)) {
        for i := 0; len(s) > 0; {
            r, size := utf8.DecodeRuneInString(s)
            forBody(i, r)
            s = s[size:]
            i += size
        }
    }

# 解包
当可变参数是一个空接口类型时，调用者是否解包可变参数会导致不同的结果：

    func main() {
        var a = []interface{}{123, "abc"}

        Print(a...) // 123 abc
        Print(a)    // [123 abc]
    }

    func Print(a ...interface{}) {
        fmt.Println(a...)
    }   


# 单件模式 
sync/atomic 包 :sync/atomic 包对基本的数值类型及复杂对象的读写都提供了原子操作的支持。atomic.Value 原子对象提供了 Load 和 Store 两个原子方法，分别用于加载和保存数据，返回值和参数都是 interface{} 类型，因此可以用于任意的自定义复杂类型。

标准库中 sync.Once 的实现：

    type Once struct {
        m    Mutex
        done uint32
    }

    func (o *Once) Do(f func()) {
        if atomic.LoadUint32(&o.done) == 1 {
            return
        }

        o.m.Lock()
        defer o.m.Unlock()

        if o.done == 0 {
            defer atomic.StoreUint32(&o.done, 1)
            f()
        }
    }
基于 sync.Once 重新实现单件模式：

    var (
        instance *singleton
        once     sync.Once
    )

    func Instance() *singleton {
        once.Do(func() {
            instance = &singleton{}
        })
        return instance
    }

# 生产者消费者模型
：后台线程生成最新的配置信息；前台多个工作者线程获取最新的配置信息。所有线程共享配置信息资源。

    var config atomic.Value // 保存当前配置信息

    // 初始化配置信息
    config.Store(loadConfig())

    // 启动一个后台线程, 加载更新后的配置信息
    go func() {
        for {
            time.Sleep(time.Second)
            config.Store(loadConfig())
        }
    }()

    // 用于处理请求的工作者线程始终采用最新的配置信息
    for i := 0; i < 10; i++ {
        go func() {
            for r := range requests() {
                c := config.Load()
                // ...
            }
        }()
    }


# 阻塞式并发

在 main 函数所在线程中执行两次 mu.Lock()，当第二次加锁时会因为锁已经被占用（不是递归锁）而阻塞，main 函数的阻塞状态驱动后台线程继续向前执行。当后台线程执行到 mu.Unlock() 时解锁，此时打印工作已经完成了，解锁会导致 main 函数中的第二个 mu.Lock() 阻塞状态取消，此时后台线程和主线程再没有其它的同步事件参考，它们退出的事件将是并发的：在 main 函数退出导致程序退出时，后台线程可能已经退出了，也可能没有退出。虽然无法确定两个线程退出的时间，但是打印工作是可以正确完成的。

    func main() {
        var mu sync.Mutex

        mu.Lock()
        go func(){
            fmt.Println("你好, 世界")
            mu.Unlock()
        }()

        mu.Lock()
    }

# 无缓存的管道

Go 语言内存模型规范，对于从无缓冲 Channel 进行的接收，发生在对该 Channel 进行的发送完成之前。因此，后台线程 <-done 接收操作完成之后，main 线程的 done <- 1 发送操作才可能完成（从而退出 main、退出程序），而此时打印工作已经完成了。

    func main() {
        done := make(chan int)

        go func(){
            fmt.Println("你好, 世界")
            <-done
        }()

        done <- 1
    }

# 有缓存的管道
对于带缓冲的 Channel，对于 Channel 的第 K 个接收完成操作发生在第 K+C 个发送操作完成之前，其中 C 是 Channel 的缓存大小。虽然管道是带缓存的，main 线程接收完成是在后台线程发送开始但还未完成的时刻，此时打印工作也是已经完成的。

    func main() {
        done := make(chan int, 1) // 带缓存的管道

        go func(){
            fmt.Println("你好, 世界")
            done <- 1
        }()

        <-done
    }

开启 10 个后台线程分别打印：

    func main() {
        done := make(chan int, 10) // 带 10 个缓存

        // 开 N 个后台打印线程
        for i := 0; i < cap(done); i++ {
            go func(){
                fmt.Println("你好, 世界")
                done <- 1
            }()
        }

        // 等待 N 个后台线程完成
        for i := 0; i < cap(done); i++ {
            <-done
        }
    }

# sync.WaitGroup
其中 wg.Add(1) 用于增加等待事件的个数，必须确保在后台线程启动之前执行（如果放到后台线程之中执行则不能保证被正常执行到）。当后台线程完成打印工作之后，调用 wg.Done() 表示完成一个事件。main 函数的 wg.Wait() 是等待全部的事件完成。

    func main() {
        var wg sync.WaitGroup

        // 开 N 个后台打印线程
        for i := 0; i < 10; i++ {
            wg.Add(1)

            go func() {
                fmt.Println("你好, 世界")
                wg.Done()
            }()
        }

        // 等待 N 个后台线程完成
        wg.Wait()
    }

# 生产者消费者模型

    // 生产者: 生成 factor 整数倍的序列
    func Producer(factor int, out chan<- int) {
        for i := 0; ; i++ {
            out <- i*factor
        }
    }

    // 消费者
    func Consumer(in <-chan int) {
        for v := range in {
            fmt.Println(v)
        }
    }
    func main() {
        ch := make(chan int, 64) // 成果队列

        go Producer(3, ch) // 生成 3 的倍数的序列
        go Producer(5, ch) // 生成 5 的倍数的序列
        go Consumer(ch)    // 消费生成的队列

        // 运行一定时间后退出
        time.Sleep(5 * time.Second)
    }

# 发布订阅模型
发布订阅（publish-and-subscribe）模型通常被简写为 pub/sub 模型。在这个模型中，消息生产者成为发布者（publisher），而消息消费者则成为订阅者（subscriber），生产者和消费者是 M:N 的关系。在传统生产者和消费者模型中，是将消息发送到一个队列中，而发布订阅模型则是将消息发布给一个主题。

    // Package pubsub implements a simple multi-topic pub-sub library.
    package pubsub

    import (
        "sync"
        "time"
    )

    type (
        subscriber chan interface{}         // 订阅者为一个管道
        topicFunc  func(v interface{}) bool // 主题为一个过滤器
    )

    // 发布者对象
    type Publisher struct {
        m           sync.RWMutex             // 读写锁
        buffer      int                      // 订阅队列的缓存大小
        timeout     time.Duration            // 发布超时时间
        subscribers map[subscriber]topicFunc // 订阅者信息
    }

    // 构建一个发布者对象, 可以设置发布超时时间和缓存队列的长度
    func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
        return &Publisher{
            buffer:      buffer,
            timeout:     publishTimeout,
            subscribers: make(map[subscriber]topicFunc),
        }
    }

    // 添加一个新的订阅者，订阅全部主题
    func (p *Publisher) Subscribe() chan interface{} {
        return p.SubscribeTopic(nil)
    }

    // 添加一个新的订阅者，订阅过滤器筛选后的主题
    func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
        ch := make(chan interface{}, p.buffer)
        p.m.Lock()
        p.subscribers[ch] = topic
        p.m.Unlock()
        return ch
    }

    // 退出订阅
    func (p *Publisher) Evict(sub chan interface{}) {
        p.m.Lock()
        defer p.m.Unlock()

        delete(p.subscribers, sub)
        close(sub)
    }

    // 发布一个主题
    func (p *Publisher) Publish(v interface{}) {
        p.m.RLock()
        defer p.m.RUnlock()

        var wg sync.WaitGroup
        for sub, topic := range p.subscribers {
            wg.Add(1)
            go p.sendTopic(sub, topic, v, &wg)
        }
        wg.Wait()
    }

    // 关闭发布者对象，同时关闭所有的订阅者管道。
    func (p *Publisher) Close() {
        p.m.Lock()
        defer p.m.Unlock()

        for sub := range p.subscribers {
            delete(p.subscribers, sub)
            close(sub)
        }
    }

    // 发送主题，可以容忍一定的超时
    func (p *Publisher) sendTopic(
        sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup,
    ) {
        defer wg.Done()
        if topic != nil && !topic(v) {
            return
        }

        select {
        case sub <- v:
        case <-time.After(p.timeout):
        }
    }

# 线程安全退出或超时的控制

在 Go1.7 发布时，标准库增加了一个 context 包，用来简化对于处理单个请求的多个 Goroutine 之间与请求域的数据、超时和退出等操作，官方有博文对此做了专门介绍。我们可以用 context 包来重新实现前面的线程安全退出或超时的控制:

    func worker(ctx context.Context, wg *sync.WaitGroup) error {
        defer wg.Done()

        for {
            select {
            default:
                fmt.Println("hello")
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }

    func main() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

        var wg sync.WaitGroup
        for i := 0; i < 10; i++ {
            wg.Add(1)
            go worker(ctx, &wg)
        }

        time.Sleep(time.Second)
        cancel()

        wg.Wait()
    }
