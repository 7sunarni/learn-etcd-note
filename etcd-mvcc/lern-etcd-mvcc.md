# etcd mvcc 实现原理

etcd 可以看到没有提交的修改？

## backend 结构体
通过 `TestBackendWritebackForEach` 测试用例
在上一个事务 Unlock 之后，可以读取到没有上一个未提交的修改。

unsafeBegin() 方法是开启事务的方法。

## kvstore 的实现 store 结构体
ReadView 接口的实现，实际上是通过传入 backend ，通过 backend 的 ReadTx(), BatchTX() 函数实现。

## storeTxnRead storeTxnWrite 结构体
分别实现 Read Write 事务的功能。


## keyindex generations revision
keyindex 中的 revision 字段，指向的是当前修改的 revision
keyindex 中的 generations 字段，其中包含了历史的 revision?

TODO:
1. cindex
2. lease

每次提交之后都又新开启了一个事务

## Read事务的开启
1. store 结构体通过调用自己的 backend 对象的 ConcurrentReadTx() 开启一个只读事务。
End() 函数中最终会触发提交。

End() 的实现在 kvstore_txn.go 中。 storeTxnWrite 和 storeTxnRead。

backend.readTx 的写入时机： 在调用 End 方法后，调用 batchTx 的 Unlock() 方法调用 commit方法提交修改，会写入到backend.readTx中。 WriteBack 方法。 
是否会出现只复制了一半的情况，没有加锁，不是原子操作？


2. 事务同时更新


## backend Write 写事务的开启

## txBuffer
txbuffer 是参考boltdb的临时数据库。每个 bucket 对应一个 bucketBuffer
bucketBuffer 记录了每次修改的记录，放入数组中。数组中是put记录的历史操作。每次放入最后一个。
bucketBuffer 会有扩容操作。

写事务每次 put 的时候会写入到数据库中 batch_tx.go L93 会创建数据库

## txWriteBuffer 
TODO: batchTx.go L285 Rollback() 作用: 只是单纯的为了关闭事务。


## store
currentRev: revision 只会在每次写事务提交只会才会增加
compactMainRev: revision 只会在每次 store 压缩后更新为当前的 currentRev。currentRev 和 compactMainRev 两个值相等的情况，则在上一次压缩之后没有新的写事务提交。


TODO: lease 的创建



## revisio
