# BoltDB 源码学习笔记
boltdb 是一个用 Go 写的内嵌 Key-Value 数据库，支持读写事务。Etcd 的存储全部依赖 boltdb 来实现，其中包括了 StableLog, UnstableLog 和具体的 Key-Value 存储。
Etcd 的 MVCC 实现也依赖于 BoltDB 的事务。因此在学习 Etcd 的 MVCC 之前，了解 boltdb 的实现是很有必要的。

boltdb 最开始的仓库已经没有维护了，Etcd 使用的 boltdb 是 fork 出来的。具体在 [https://github.com/etcd-io/bbolt](https://github.com/etcd-io/bbolt).

boltdb 源码定义分析
- db
- bucket
- node
- page
- freelist
- cursor
- tx

## db
db 对象可以理解为数据库对象，其实也就是一个文件，创建 db 其实也就是创建了一个文件，然后从这个数据库中开始读、写事务，对文件中的数据进行读写。
StrictMode: 严格模式，只在调试的时候使用。每次提交之后会调用 Check 函数。

NoSync 使用字后，会使每次提交之后会跳过 fsync() 方法。 使用 fdatasync() 调用。

## bucket
bucket 可以理解为一张表，可以新建，删除一个 bucket，然后具体的 k-v 数据是在 bucket 中操作的。

## inode
inode 是数据存储的最小单元，数据库中所有的数据，无论类型都存放在 inode 中。

## page 
page 对应的是磁盘的数据库文件中的某一个部分，将磁盘文件按页的方式读到内存中。page 本身可以理解为一块部分映射磁盘文件到内存的对象。page 有大量的 inode 数组对象组成。page 可以有不同的类型：branchPage 和 leafPage，对应的是b树中的叶节点和主干节点。如果 page 是主干节点，那么里面的 inode 对象 key 存放的是索引值，inode 对象的 pageid 对应的是索引对应的数据存放的 pageid，可以通过 pageid 快速找到数据所在的 page。如果是 leafPage，那么 inode 的 key，value 对应的是真正存放的 key,value 数据，其中的 pageid 是没有意义的。可以通过 inode 的 flag 来区分 inode 存放的数据类型。
branchPage 是由叶节点分裂出来的，branchPage 的 key 是节点的第一个 inode 的key [node.go L381]()

## node
node 是一棵树。
在写事务写入数据的时候，全部都在一个 node 的 inode 上新增。很有可能这个 node 的存放数量已经超过了最大限制。
事务提交的时候，对node 的树结构从下至上进行递归分裂。如果一个 node 的数量超过了最大限制，就进行分裂，分裂的大小对应一个 page 的大小，然后为这个 node 申请一个 page，将 node 的pageid 设置为申请的pageid。然后向该 node 的父节点插入一个 inode，这个 inode 可以看做是一个索引节点，通过这个 inode 的pageid 就可以很快找到存放数据节点的 page 了。


## tx
在创建事务的时候，会给数据库对象加锁 [db.go L593]()，所以写事物只能有一个。
在事务关闭的时候，会释放锁。[tx.go L303]()。

事务开始的时候会创建一个空的 Bucket 对象。事务在创建或者获取打开一个 Bucket 的时候，会用这个 Bucket 对象去打开。

## bucket
事务的读写都是通过 bucket 对象来实现的。

## freelist
> 用于空闲内存管理 

绑定事务和内存页的分配
allocs: 字段存放了 `起始页ID` 和 `事务ID` 的映射。为事务分配虚拟内存页的时候，会寻找连续的页码，然后记录起始页ID和事务ID的映射。

pending: 存放了即将释放的事务页码。绑定 `事务ID` 和即将释放的页码ID关系。

ids: 存放了所有的可以用的页码 id

cache: 存放 页码ID 是否是空闲或者即将空闲。 TRUE 表示即将空闲。

txp.alloctx 作用？

free() 函数：移除 allocs 的页码和事务ID的绑定关系。将事务ID绑定的页码 ID 放入 pending 对象中。知识取消了绑定，并没有释放页码ID。

release() 函数：将事务ID 在 pending 中的页码ID释放出来。

rollback(): 将 pending 中指定事务ID对应的页码ID释放出来，其他事务的页码ID标记正在被事务ID使用。

只有 db 对象持有 freelist 对象。
freelist 只是持有id，并不持有真正的 page 对象

在 db 初始化的时候，db.freepages() 函数的作用是最开始填充 freelist？

## cursor
search() 递归查询
最开始传入的 pgid 是 root 节点的pgid
inode 是存放 key, value 的结构体。
seek() 函数开始 search 会将游标移动到指定的 page 上。
然后调用最近的 keyValue 读取 page 上指定下标的 inode 的数据。
node branchPage 中的 key 只是索引？

写入数据都是存放在node中，在写入完成之后提交的时候，会检查 node 的长度之后会分裂node

inode 的 pageid 在节点分裂的时候回发生改变，node.go L384

key,value 节点的 pgid 是没有意义的？
leaf-page 是存放数据的， barnch-page 是存放数据索引的。

## Questiones
- 为什么 node 的 rabalance 和 spill 两个函数是分开的？
- deletebucket 的操作?
- freelist 
- 可能有多个根节点，因为最开始的根节点分裂之后，不会创建新的根节点? node.go L378
- db.meta 的数据？ tx.root.bucket 是多少
- inline

