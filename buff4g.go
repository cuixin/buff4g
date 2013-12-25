package buff4g

// 内存缓冲器,此应用的目的就是为了减轻GC的负担
// 本应用的想法来源于达达的一段程序，Alloc的实现就是达达同学的想法。
// PAlloc是后来我加的一个多线程安全版本的Alloc,目的也是为了减少锁的征用,所以
// 上来先分配一个池的容量，然后每次CAS拿到一个数组的所以，然后在进行加锁,如果
// 使用全局锁必将增加负担,如果你有更好的实现方法请务必让我知道,我将把你加到贡献
// 者,谢谢.
// 在此特别感谢达达，这个简单、粗暴、有效的办法能解决很大问题，他的gihub:
// github.com/idada
import (
	"sync"
	"sync/atomic"
)

// 区块的定义,不要自己初始化,使用NewBlockBytes
type BlockBytes struct {
	blockSize int
	lock      sync.Mutex
	mutex     sync.Mutex
	curBytes  []byte
	newBytes  []byte
}

// buffer池的定义
type buffPool struct {
	poolMod    int32
	blockBytes []BlockBytes
	pos        int32
}

var pool buffPool

// 新申请一块内存块, size:区块大小
func NewBlockBytes(size int) *BlockBytes {
	return &BlockBytes{blockSize: size,
		lock:     sync.Mutex{},
		curBytes: make([]byte, size),
		newBytes: make([]byte, size)}
}

// 初始化Buffer, blockSize:每个区块的大小, poolSize:池的大小必须是2的N次方,例如1,2,4,8,16...
func InitBuffer(blockSize int, poolSize int32) {
	if poolSize <= 0 {
		panic("Cannot set poolSize <= 0")
	}
	if (poolSize & -poolSize) != poolSize {
		panic("You must be set poolSize is 1, 2, 4, 8, 16...")
	}
	pool = buffPool{
		poolMod:    poolSize - 1,
		blockBytes: make([]BlockBytes, poolSize),
		pos:        0}
	var i int32 = 0
	for ; i < poolSize; i++ {
		bb := BlockBytes{
			blockSize: blockSize,
			lock:      sync.Mutex{},
			mutex:     sync.Mutex{},
			curBytes:  make([]byte, blockSize),
			newBytes:  make([]byte, blockSize)}
		pool.blockBytes[i] = bb
	}
}

// 分配内存,非线程安全,只能在一个routine中
func (this *BlockBytes) Alloc(size int) []byte {
	if size > this.blockSize {
		return make([]byte, size)
	}
	if len(this.curBytes) < size {
		this.lock.Lock()
		this.curBytes = this.newBytes
		this.lock.Unlock()

		go func() {
			this.lock.Lock()
			this.newBytes = make([]byte, this.blockSize)
			this.lock.Unlock()
		}()
	}
	result := this.curBytes[0:size]
	this.curBytes = this.curBytes[size:]
	return result
}

// 多线程安全版本的,用之前必须先初始化: InitBuffer
func PAlloc(size int) []byte {
	num := atomic.AddInt32(&pool.pos, 1)
	mod := num & pool.poolMod
	blockBuf := pool.blockBytes[mod]
	blockBuf.mutex.Lock()
	defer blockBuf.mutex.Unlock()

	if size > blockBuf.blockSize {
		return make([]byte, size)
	}
	if len(blockBuf.curBytes) < size {
		blockBuf.lock.Lock()
		blockBuf.curBytes = blockBuf.newBytes
		blockBuf.lock.Unlock()

		go func() {
			blockBuf.lock.Lock()
			blockBuf.newBytes = make([]byte, blockBuf.blockSize)
			blockBuf.lock.Unlock()
			// fmt.Println("Resize")
		}()
	}
	result := blockBuf.curBytes[0:size]
	blockBuf.curBytes = blockBuf.curBytes[size:]
	return result
}
