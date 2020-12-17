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

backend.readTx 的写入时机： 在调用 End 方法后，调用 batchTx 的 Unlock() 方法调用 commit 方法提交修改 batch_tx.go L202 ，会写入到backend.readTx中。 WriteBack 方法。 batch_tx.go L251 
是否会出现只复制了一半的情况，没有加锁，不是原子操作？ 


2. 事务同时更新


## backend Write 写事务的开启
batchTxBuffered 有两部分
一部分是boltdb 的 tx,另一部分是 buffer
写入数据的时候，buffer 和 tx 都写入了数据
buffer 的数据在最后提交的时候，写入到了 backend.readTx 中。


## store
currentRev: revision 只会在每次写事务提交只会才会增加
compactMainRev: revision 只会在每次 store 压缩后更新为当前的 currentRev。currentRev 和 compactMainRev 两个值相等的情况，则在上一次压缩之后没有新的写事务提交。

TODO: lease 的创建


## revision


## kvindex
带有版本信息的 kv 值


## txBuffer
txbuffer 是参考boltdb的临时数据库。每个 bucket 对应一个 bucketBuffer
bucketBuffer 记录了每次修改的记录，放入数组中。数组中是put记录的历史操作。每次放入最后一个。
bucketBuffer 会有扩容操作。

写事务每次 put 的时候会写入到数据库中 batch_tx.go L93 会创建 bucket。
写入 bucket 的数据的值是带有版本的。 kvstore_txn.go L203 d 是kv数据的值，来源于序列化的数据，其中包含了 revision 的信息。 revision 的 main 是事务开启时候的 revision, sub 是这个事务进过了多少的修改。

kvstore 当前的 revision 值为压实之后的值。 kvstore.go L397 L402

## txWriteBuffer 
txWriteBuffer 持有 txbuffer 对象。
TODO: batchTx.go L285 Rollback() 作用: 只是单纯的为了关闭事务。

## batchTxBuffered 对象
batchTxBuffered 有一个 batchTx 对象，有一个 txWriteBuffer 对象。

batchTx 持有真正的 BoltDB 的 tx 对象。txWriteBuffer 记录了每次修改的记录

在使用 batchTxBuffered 写入数据的时候，会直接写入到 BoldDB 的事务中和 buffer 中。batch_tx.go L300


## Store 开启 Write 事务
kvstore_txn.go L69 会调用 backend.BatchTx()，返回 backend 的 batchTxBuffered 对象，通过 batchTxBuffered 写入数据的时候，会同时写入数据到数据库中和 buffer 当中。


