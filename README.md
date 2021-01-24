
`jtopthreads` prints Java stack traces ordered by CPU usage. It can be run against a running JVM or previously captured `jstack` output.

## Installation

From a directory outside of `$GOPATH` and without a `go.mod` file run,

```
GO111MODULE=on go get github.com/c2nes/jtopthreads@latest
```

## Usage

``` shellsession
$ jtopthreads -h
usage: jtopthreads [options] <stack-file> [stack-file]
   or: jtopthreads [options] [-sample <duration>] <pid | main-class>

  -n N
        limit output to the top N threads
  -sample duration
        sample process for duration
  -summary
        omit stacks
```

### Examples

Sample the process with the main-class `net.qrono.server.Main` for 5 seconds and print a summary (omitting stack traces) for the 10 busiest threads:

``` shellsession
$ jtopthreads -n 10 -summary -sample 5s net.qrono.server.Main
[ 77.97%] "epollEventLoopGroup-5-3" #27 prio=10 os_prio=0 cpu=7653.72ms elapsed=10.29s tid=0x00007fdf6000b000 nid=0x1eb8 runnable  [0x00007fdf4e8f4000]
[ 76.99%] "epollEventLoopGroup-5-4" #28 prio=10 os_prio=0 cpu=7659.95ms elapsed=10.29s tid=0x00007fdf60007000 nid=0x1eb9 runnable  [0x00007fdf4e7f3000]
[ 43.22%] "Thread-1" #25 daemon prio=5 os_prio=0 cpu=15317.64ms elapsed=60.64s tid=0x00007fdf48001000 nid=0x1d07 runnable  [0x00007fdf4fffd000]
[  0.79%] "GC Thread#0" os_prio=0 cpu=900.33ms elapsed=105.53s tid=0x00007fe01c042000 nid=0x1b7c runnable
[  0.78%] "GC Thread#3" os_prio=0 cpu=919.06ms elapsed=59.57s tid=0x00007fdfe0114000 nid=0x1d0b runnable
[  0.64%] "GC Thread#1" os_prio=0 cpu=1081.23ms elapsed=59.57s tid=0x00007fdfe0111000 nid=0x1d09 runnable
[  0.60%] "GC Thread#2" os_prio=0 cpu=1009.09ms elapsed=59.57s tid=0x00007fdfe0112800 nid=0x1d0a runnable
[  0.19%] "VM Thread" os_prio=0 cpu=119.19ms elapsed=105.48s tid=0x00007fe01d168000 nid=0x1b81 runnable
[  0.08%] "G1 Young RemSet Sampling" os_prio=0 cpu=74.55ms elapsed=105.50s tid=0x00007fe01d106000 nid=0x1b80 runnable
[  0.08%] "grpc-default-worker-ELG-3-1" #19 daemon prio=5 os_prio=0 cpu=299.90ms elapsed=104.77s tid=0x00007fdf58003000 nid=0x1b9c runnable  [0x00007fdfba8b5000]
[201.25%] Total (elapsed 5.05s)
```

Same as above, but using previously captured `jstack` output (handy if you are grabbing data from a machine without `jtopthreads` installed):

``` shellsession
$ JVMID=$(jps -l | awk '$2 == "net.qrono.server.Main" {print $1}'); jstack "$JVMID" > stack.0 && sleep 5 && jstack "$JVMID" > stack.1

$ jtopthreads -n 10 -summary stack.0 stack.1
[ 81.82%] "epollEventLoopGroup-5-5" #29 prio=10 os_prio=0 cpu=7477.64ms elapsed=9.43s tid=0x00007fdf6000a000 nid=0x25eb runnable  [0x00007fdf4e6f1000]
[ 80.03%] "epollEventLoopGroup-5-6" #30 prio=10 os_prio=0 cpu=7324.02ms elapsed=9.43s tid=0x00007fdf60008800 nid=0x25ec waiting for monitor entry  [0x00007fdf4e5f1000]
[ 42.61%] "Thread-1" #25 daemon prio=5 os_prio=0 cpu=25618.41ms elapsed=431.62s tid=0x00007fdf48001000 nid=0x1d07 waiting on condition  [0x00007fdf4fffe000]
[  0.17%] "G1 Young RemSet Sampling" os_prio=0 cpu=245.53ms elapsed=476.47s tid=0x00007fe01d106000 nid=0x1b80 runnable
[  0.08%] "VM Thread" os_prio=0 cpu=264.03ms elapsed=476.46s tid=0x00007fe01d168000 nid=0x1b81 runnable
[  0.06%] "C2 CompilerThread0" #6 daemon prio=9 os_prio=0 cpu=8711.94ms elapsed=476.44s tid=0x00007fe01d18d800 nid=0x1b86 waiting on condition  [0x0000000000000000]
[  0.05%] "C1 CompilerThread0" #8 daemon prio=9 os_prio=0 cpu=1879.74ms elapsed=476.44s tid=0x00007fe01d18f800 nid=0x1b87 waiting on condition  [0x0000000000000000]
[  0.03%] "grpc-default-worker-ELG-3-1" #19 daemon prio=5 os_prio=0 cpu=692.85ms elapsed=475.74s tid=0x00007fdf58003000 nid=0x1b9c runnable  [0x00007fdfba8b5000]
[  0.03%] "VM Periodic Task Thread" os_prio=0 cpu=426.87ms elapsed=476.40s tid=0x00007fe01d240000 nid=0x1b89 waiting on condition
[  0.02%] "grpc-default-executor-0" #20 daemon prio=5 os_prio=0 cpu=287.92ms elapsed=475.16s tid=0x00007fdf5c0b3800 nid=0x1ba2 waiting on condition  [0x00007fdfbafb8000]
[204.76%] Total (elapsed 5.54s)
```

Show 3 busiest threads since application start up (not sampled):

``` shellsession
$ jtopthreads -n 3 net.qrono.server.Main
[ 79.03%] "epollEventLoopGroup-5-7" #31 prio=10 os_prio=0 cpu=4235.79ms elapsed=5.36s tid=0x00007fdf60012000 nid=0x2f1d runnable  [0x00007fdf31ffb000]
   java.lang.Thread.State: RUNNABLE
        at io.netty.buffer.DefaultByteBufHolder.content(DefaultByteBufHolder.java:36)
        at io.netty.handler.codec.MessageAggregator.decode(MessageAggregator.java:284)
        at io.netty.handler.codec.MessageToMessageDecoder.channelRead(MessageToMessageDecoder.java:88)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:379)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:365)
        at io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:357)
        at io.netty.handler.codec.ByteToMessageDecoder.fireChannelRead(ByteToMessageDecoder.java:324)
        at io.netty.handler.codec.ByteToMessageDecoder.channelRead(ByteToMessageDecoder.java:296)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:379)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:365)
        at io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:357)
        at io.netty.channel.DefaultChannelPipeline$HeadContext.channelRead(DefaultChannelPipeline.java:1410)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:379)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:365)
        at io.netty.channel.DefaultChannelPipeline.fireChannelRead(DefaultChannelPipeline.java:919)
        at io.netty.channel.epoll.AbstractEpollStreamChannel$EpollStreamUnsafe.epollInReady(AbstractEpollStreamChannel.java:792)
        at io.netty.channel.epoll.EpollEventLoop.processReady(EpollEventLoop.java:475)
        at io.netty.channel.epoll.EpollEventLoop.run(EpollEventLoop.java:378)
        at io.netty.util.concurrent.SingleThreadEventExecutor$4.run(SingleThreadEventExecutor.java:989)
        at io.netty.util.internal.ThreadExecutorMap$2.run(ThreadExecutorMap.java:74)
        at io.netty.util.concurrent.FastThreadLocalRunnable.run(FastThreadLocalRunnable.java:30)
        at java.lang.Thread.run(java.base@11.0.10/Thread.java:834)

[ 78.34%] "epollEventLoopGroup-5-8" #32 prio=10 os_prio=0 cpu=4191.43ms elapsed=5.35s tid=0x00007fdf60014000 nid=0x2f1e runnable  [0x00007fdf31efa000]
   java.lang.Thread.State: RUNNABLE
        at com.lmax.disruptor.BlockingWaitStrategy.signalAllWhenBlocking(BlockingWaitStrategy.java:68)
        at com.lmax.disruptor.SingleProducerSequencer.publish(SingleProducerSequencer.java:207)
        at com.lmax.disruptor.RingBuffer.translateAndPublish(RingBuffer.java:1004)
        at com.lmax.disruptor.RingBuffer.publishEvent(RingBuffer.java:555)
        at net.qrono.server.Queue.enqueueAsync(Queue.java:139)
        at net.qrono.server.redis.RedisChannelInitializer$RequestHandler.lambda$handleEnqueue$0(RedisChannelInitializer.java:123)
        at net.qrono.server.redis.RedisChannelInitializer$RequestHandler$$Lambda$70/0x00000008001c4040.apply(Unknown Source)
        at net.qrono.server.QueueManager$QueueWrapper.apply(QueueManager.java:160)
        - locked <0x0000000683652818> (a net.qrono.server.QueueManager$QueueWrapper)
        at net.qrono.server.QueueManager.withQueueAsync(QueueManager.java:115)
        at net.qrono.server.redis.RedisChannelInitializer$RequestHandler.handleEnqueue(RedisChannelInitializer.java:122)
        at net.qrono.server.redis.RedisChannelInitializer$RequestHandler.handleMessage(RedisChannelInitializer.java:337)
        at net.qrono.server.redis.RedisChannelInitializer$RequestHandler.handleMessageSafe(RedisChannelInitializer.java:405)
        at net.qrono.server.redis.RedisChannelInitializer$RequestHandler.channelRead0(RedisChannelInitializer.java:415)
        at net.qrono.server.redis.RedisChannelInitializer$RequestHandler.channelRead0(RedisChannelInitializer.java:65)
        at io.netty.channel.SimpleChannelInboundHandler.channelRead(SimpleChannelInboundHandler.java:99)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:379)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:365)
        at io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:357)
        at io.netty.handler.codec.MessageToMessageDecoder.channelRead(MessageToMessageDecoder.java:103)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:379)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:365)
        at io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:357)
        at io.netty.handler.codec.MessageToMessageDecoder.channelRead(MessageToMessageDecoder.java:103)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:379)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:365)
        at io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:357)
        at io.netty.handler.codec.ByteToMessageDecoder.fireChannelRead(ByteToMessageDecoder.java:324)
        at io.netty.handler.codec.ByteToMessageDecoder.channelRead(ByteToMessageDecoder.java:296)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:379)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:365)
        at io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:357)
        at io.netty.channel.DefaultChannelPipeline$HeadContext.channelRead(DefaultChannelPipeline.java:1410)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:379)
        at io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:365)
        at io.netty.channel.DefaultChannelPipeline.fireChannelRead(DefaultChannelPipeline.java:919)
        at io.netty.channel.epoll.AbstractEpollStreamChannel$EpollStreamUnsafe.epollInReady(AbstractEpollStreamChannel.java:792)
        at io.netty.channel.epoll.EpollEventLoop.processReady(EpollEventLoop.java:475)
        at io.netty.channel.epoll.EpollEventLoop.run(EpollEventLoop.java:378)
        at io.netty.util.concurrent.SingleThreadEventExecutor$4.run(SingleThreadEventExecutor.java:989)
        at io.netty.util.internal.ThreadExecutorMap$2.run(ThreadExecutorMap.java:74)
        at io.netty.util.concurrent.FastThreadLocalRunnable.run(FastThreadLocalRunnable.java:30)
        at java.lang.Thread.run(java.base@11.0.10/Thread.java:834)

[  4.35%] "epollEventLoopGroup-5-5" #29 prio=10 os_prio=0 cpu=18547.01ms elapsed=426.43s tid=0x00007fdf6000a000 nid=0x25eb runnable  [0x00007fdf4e6f2000]
   java.lang.Thread.State: RUNNABLE
        at io.netty.channel.epoll.Native.epollWait(Native Method)
        at io.netty.channel.epoll.Native.epollWait(Native.java:148)
        at io.netty.channel.epoll.Native.epollWait(Native.java:141)
        at io.netty.channel.epoll.EpollEventLoop.epollWaitNoTimerChange(EpollEventLoop.java:290)
        at io.netty.channel.epoll.EpollEventLoop.run(EpollEventLoop.java:347)
        at io.netty.util.concurrent.SingleThreadEventExecutor$4.run(SingleThreadEventExecutor.java:989)
        at io.netty.util.internal.ThreadExecutorMap$2.run(ThreadExecutorMap.java:74)
        at io.netty.util.concurrent.FastThreadLocalRunnable.run(FastThreadLocalRunnable.java:30)
        at java.lang.Thread.run(java.base@11.0.10/Thread.java:834)

[ 20.26%] Total (elapsed 14m53.53s)
```

## Supported Platforms

`jtopthreads` has only been tested with HotSpot. It supports Java 8 on Linux and Java 11+ on other platforms.

