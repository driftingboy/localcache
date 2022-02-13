# localcache
本地缓存管理

最后总结下来

避免绝大部分多余的内存分配行为，能复用绝不分配。
善用sync.Pool。
尽量避免[]byte与string之间转换带来的开销。
巧用[]byte相关的特性。


- 避免string与[]byte转换开销
这两种类型转换是带内存分配与拷贝开销的，但有一种办法(trick)能够避免开销。利用了string和slice在runtime里结构只差一个Cap字段实现的

- 不放过能复用内存的地方
有些地方需要kv型数据，一般使用map[string]string。但map不利于复用。所以fasthttp使用slice来实现了map

- 方法参数尽量用[]byte. 纯写场景可避免用bytes.Buffer
方法参数使用[]byte，这样做避免了[]byte到string转换时带来的内存分配和拷贝。毕竟本来从net.Conn读出来的数据也是[]byte类型。

某些地方确实想传string类型参数，fasthttp也提供XXString()方法。

String方法背后是利用了a = append(a, string…)。这样做不会造成string到[]byte的转换(该结论通过查看汇编得到，汇编里并没用到runtime.stringtoslicebyte方法)