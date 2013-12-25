Buff4g
======

设计目的:
-------------
内存缓冲器,此应用的目的就是为了减轻GC的负担  

本应用的想法来源于达达的一段程序，Alloc的实现就是达达同学的想法。  
PAlloc是后来我加的一个多线程安全版本的Alloc,目的也是为了减少锁的征用,所以  
上来先分配一个池的容量，然后每次CAS拿到一个数组的所以，然后在进行加锁,如果  
使用全局锁必将增加负担,如果你有更好的实现方法请务必让我知道,我将把你加到贡献  
者,谢谢.  

在此特别感谢达达，这个简单、粗暴、有效的办法能解决很大问题，他的gihub:
[达达](http://github.com/idada)

安装说明
------------
go get github.com/cuixin/proxy

使用说明
------------
在不用关心多线程问题的时候，可以直接使用:
```
	bb := NewBlockBytes(1024 * 200)
	buf := bb.Alloc(64)
```

多线程安全的版本，必须先在线程安全的程序中初始化:  
InitBuffer(1024 * 200, 4)  
1024*200代表每个内存区块的大小(200K)，4代表池大小  
池大小必须为2的次方,1,2,4,8,16,32...  
	buf := bb.PAlloc(64)

License
-------
Buff4g is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).