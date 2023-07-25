


one 中，用Golang实现了一个高性能的kv，但是由于锁的关系，即使有多个go协程并发对kv进行操作，但还是会因为全局的锁导致一go干活，多go看戏的场面，今天就来优化一下这个kv，让它实现真正的并发干活
使用 hash 来讲请求 key 进行 转发,引入了一个新的机制：将dirty数据分散存储到若干个map里，这样锁就只会锁住单个map，其余的map还是可以进行读写操作的，通过这样来让多go程真正地并发工作

---------------------------------
作者: SunistC
本文来自于: https://www.sunist.cn/post/KeyValueStore-GolangImplement-2
博客内容遵循 署名-非商业性使用-相同方式共享 4.0 国际 (CC BY-NC-SA 4.0) 协议