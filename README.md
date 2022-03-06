# localcache

localcache 是一个分布式缓存系统，参考 groupcache，将在其基础上引入 节点增删、缓存失效、发布订阅等功能

服务端维护路由信息，简单，但有单点问题
![image](https://res.craft.do/user/full/c08a465b-93e5-c0ee-4310-637acb8215b4/doc/E9DDF3EA-D499-4564-AEB0-6111D0F760CA/4EE8B491-CE79-4C59-B609-4E5B90394387_2/R1g7Vq8JKEtwVKFVuHa9RvD7jsTC9vLXjPXH7y8zzq4z/Image.png)

客户端维护，复杂
![image](https://res.craft.do/user/full/c08a465b-93e5-c0ee-4310-637acb8215b4/doc/E9DDF3EA-D499-4564-AEB0-6111D0F760CA/61C2CA3E-3F5E-45F9-AE9C-D0C6E51588B0_2/19zsn1DqoKrQuJ8jx2bOlxJAoflHLBCUJrIBpCCFQLYz/Image.png)

当然两者可以结合使用

groupcache中每个节点都是相同的代码，即有分片路由能力，也保存缓存数据；这就导致路由信息保存在多个地方，不好维护，所以groupcache也没有提供增删节点的操作
![image](https://res.craft.do/user/full/c08a465b-93e5-c0ee-4310-637acb8215b4/doc/E9DDF3EA-D499-4564-AEB0-6111D0F760CA/CD5D8367-C985-4820-9866-E0DA9064D6F8_2/gLDfTg03Xf4cxVMVmFPO9ILRjENGPBHtX59rsIq6Lrwz/Image.png)

整个查询流程如下：
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶

