#### v0.13.4

- package circularqueue:
  - [feat] 新增package，底层基于切片实现的固定容量大小的FIFO的环形队列
- package nazanet:
  - [feat] 这是一个新增的package，其中struct AvailUDPConnPool，可以从指定的UDP端口范围内，寻找可绑定监听的端口，绑定监听并返回
  - [feat] 新增struct UDPConnection，对UDP连接对象的简易封装
- package bele:
  - [feat] 新增func BEPutUint16, BEPutUint64

#### v0.13.3

- package nazahttp:
  - [fix] 当request line和status line存在多个空格时，解析错误
- package bininfo:
  - [feat] 增加git tag信息

#### v0.13.2

- package nazareflect:
  - [feat] 新增package，提供三个函数IsNil，Equal，EqualInteger

#### v0.13.1

- package connection:
  - [feat] 增加连接关闭标志，使用channel发送数据时，如果连接已关闭，可以向调用方返回错误
  - [feat] 增加Option.WriteChanFullBehavior，使用channel发送数据时，如果channel满了，可以配置是阻塞还是返回错误
  - [fix] 设置wChan大小时，应该使用WriteChanSize而不是WriteBufSize
  - [refactor] 不同错误返回不同的错误值
  - [refactor] 去除一些debug日志

#### v0.13.0

- package nazabits:
    - [feat] BitReader的所有函数增加读取越界检查
    - [feat] 增加BitReader::ReadGolomb函数，读取0阶指数哥伦布编码
- package nazahttp:
    - [feat] nazahttp: 增加函数ReadHTTPHeader，ParseHTTPRequestLine，ParseHTTPStatusLine，读取HTTP头部信息
    - [refactor] 函数GetHttpFile，DownloadHTttpFile重命名为GetHTTPFile，DownloadHTTPFile

#### v0.12.3

- package lru:
    - [feat] 新增package，一个基础的LRU缓存
- package nazahttp:
    - [feat] 新增函数GetHttpFile，用于下载HTTP文件

#### v0.12.2

- package nazabits:
    - [fix] BitWriter::WriteBit如果原数据不为零值时，会错误覆盖非写入的位

#### v0.12.1

- package nazabits:
    - [feat] 新增 BitReader::ReadBits16，ReadBits32，ReadBytes函数

#### v0.12.0

- package nazabits:
    - [refactor] BitReader::ReadBits重命名为ReadBits8，BitWriter::WriteBits重命名为WriteBits8，WriteBits16
    - [fix] BitWriter::WriteBit传入的值不为0和1时，只取最低位

#### v0.11.0

- package nazabits:
    - [feat] 增加BitReader用于按位读取，增加BitWriter用于按位写入
    - [refactor] GetBit8等函数修改整型类型

#### v0.10.0

- package nazalog:
    - [feat] 新增WithPrefix函数，用于支持设置前缀，并且前缀可叠加，使得可以按repo ，package，对象等维度添加不同的前缀
    - [feat] 新增Println等函数，方便替换标准库日志
    - [perf] 减小锁粒度
    - [test] 测试覆盖率增加至100%
    - [refactor] 删除FatalIfErrorNotNil, PanicIfErrorNotNil, Outputf, writeShortFile四个函数函数
- package nazabits:
    - [feat] 新增GetBit16，GetBits16函数
- package fake:
    - [feat] 新增WithRecover函数
    - [feat] 新增time.Now hook相关的接口
    - [refactor] 重新命名os.Exit hook相关的接口

#### v0.9.1

- package nazabits:
    - [feat] 新增package，提供一些位运算函数

#### v0.9.0

- package taskpool:
    - [feat] 增加Option.MaxWorkerNum，可配置最大协程并发数量
    - [feat] task任务函数可携带参数
- package nazalog:
    - [feat] 新增Assert函数，并可配置断言失败后的行为
- package bele:
    - [feat] 增加一些ReadXXX函数，从io.Reader中读取数据
    - [feat] 新增BEUint64函数

#### v0.8.0

- package ratelimit:
    - [feat] 新增漏桶LeakyBucket和令牌桶TokenBucket，把以前的RateLimit删了
- package nazalog:
    - [feat] 新增Sync函数，将日志刷盘
    - [feat] panic和fatal前调用Sync
    - [fix] 全局的Panic和Panicf忘记调用标准库中的panic
    - [fix] 使用IsRotateDaily控制日志是否按天翻转，之前没有判断这个标志，导致总是按天翻转
- package nazahttp:
    - [feat] 新增DownloadHttpFile函数，下载http保持为本地文件

#### v0.7.1

- package snowflake
    - [feat] 新增的包，分布式唯一性64位ID生成器

#### v0.7.0

- package consistenthash:
    - [feat] Nodes 接口返回 node 对应的 point 范围，供调用方判断 hash 是否均匀
    - [feat] hash 函数可由外部配置
    - [refactor] 增加 struct Option
    - [refactor] 内部 point 类型 int -> uint32
- package bitrate:
    - [feat]可配置 bitrate 返回时的单位
    - [feat] unix 时间戳可选择由外部传入
    - [refactor] struct Bitrate -> interface Bitrate
    - [fix] 遍历切片时删除了元素导致崩溃
- package fake:
    - [feat] 添加 func Exit，它是对 os.Exit 的封装，便于其他代码做单元测试
- package log:
    - [test] 使用 fake.Exit

#### v0.6.0

- 新增 package ratelimit：限速器，令牌桶
- 新增 package bitrate：计算带宽
- 新增 package fake
- 删除 package mockwriter
- 删除 package mockserver
- `demo/add_blog_license`：行尾增加两个空格，便于部分 markdown 解析器解析为为换行

#### v0.5.1

- package ic:
  - 新增的包，将整型切片压缩序列化成二进制字符切片
- package bininfo:
  - 增加注释
- package assert:
  - 增加注释

#### v0.5.0

- package filebatch:
    - 遍历读取文件发生错误时，不退出遍历，而是将错误在回调中返回给上层
- package connection:
    - bugfix，初始化 write chan 相关的信息是通过 write chan 的配置，而不是 write buf 的配置
- package slicebytepool:
    - 新增的包，一个 []byte 内存池
- package nazamd5:
    - 新增的包，md5 操作
- package consistenthash:
    - 新增的包，一致性hash
- package bufferpool:
    - 删除 bufferpool 包
- demo/myapp:
    - 用于演示 package bininfo 的使用
- `demo/add_blog_license`:
    - 修改 license 内容
- 其他：
    - 统一error变量的命名方式及内容格式，涉及到的 package：filebatch, connection, taskpool, nazalog

#### v0.4.3

- package bufferpool
    - 新增的包，bytes.Buffer 池
- package nazaatomic
    - 新增的包，对 sync.atomic 的再封装
- package taskpool
    - 新增的包，协程池
- test.sh
    - 做更多的 go tools 检查
- `demo/add_blog_license`
    - 修改许可证样式
    - 检查许可证是否存在时，只检查声明两个字
- `demo/add_go_license`
    - 用户名和邮箱由命令行参数传入

#### v0.4.2

- package filebatch:
    - 新增的包，用于文件批处理操作
- 新增 demo/add_go_license：给 Go 仓库的所有go源码文件添加MIT许可证
- 新增 demo/add_blog_license：给我自己博客的所有文章尾部添加声明

#### v0.4.1

- package nazastring:
    - 新增 func SliceByteToStringTmp 和 func StringToSliceByteTmp，用于无拷贝的做string和[]byte的转换

#### v0.4.0

- package log:
    - rename -> package nazalog
    - mkdir 0777 and create file 0666, append if file exist
    - 配置使用 Option
    - 配置默认值修改： 打印至控制台开关默认打开，打印源码文件行号开关默认打开
- package connection:
    - erase func Printf
    - 配置使用 Option
- package nazajson:
    - 新增包，作为系统包 json 的补充
- 其它:
    - repo name nezha -> naza

#### v0.3.0

- package connection:
    - 可配置使用 channel 进行异步发送：Config 中增加 WChanSize。增加 Flush, Done, ModWriteChanSize 三个方法
- package log:
    - 增加 panic 相关的方法
- 其它:
    - test.sh 中添加 gofmt 检查

#### v0.2.0

- package log:
    - 去除了对标准库中log的依赖
    - 日志支持按天翻转
    - 增加 ShortFileFlag 可配置是否打印源码文件及行号的信息
    - 添加一个fatal日志级别，打印完后exit程序
    - 当同时打印至控制台和文件时，打印至文件中的level字段也带颜色属性
    - 增加 Out 接口
    - 日志不再支持按固定大小翻转 [不兼容]
    - 日志级别从0 -> 1开始 [不兼容]

#### v0.1.0

- 删除 /pkg/errors [不兼容]
- package log:
    - 增加 FatalIfErrorNotNil 接口函数，打印错误并退出程序
    - 日志内容中的级别字段右对齐
    - 日志内容中的源码文件名和行号放在整行日志的末尾
    - 增加一些 benchmark
- package assert: 打印正确的源码文件名和行号信息 [bugfix]
- package bele: 增加一些 benchmark
- package unique:
    - 不同的 key 使用不同的自增计数
    - 增加一些 benchmark
- package mockserver: 模拟一些服务端，用于快速测试其它代码
- package mockwriter: 模拟 Writer 接口，用于快速测试其它代码
- 删除 /demo/connstat
- test.sh 脚本只测试 /pkg 目录下的源码

#### v0.0.1

第一个版本
