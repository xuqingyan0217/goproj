package bitmap

type Bitmap struct {
	bits []byte
	size int
}

func NewBitmap(size int) *Bitmap {
	if size == 0 {
		// 250个字节，共2000位，表示一个2000人的群
		size = 250
	}
	return &Bitmap{
		// 字节数组初始化时是指定多少个字节，因为go里面最小就是字节
		bits: make([]byte, size),
		// 字节转化为bit，就是实际bitmap的大小，用int接收，仅表数值
		size: size * 8,
	}
}

// Set 设置值
func (b *Bitmap) Set(id string)  {
	// id在哪个bit
	idx := hash(id) % b.size
	// 计算在哪个byte
	byteIdx := idx / 8
	// 在该byte中的哪个bit
	bitIdx := idx % 8
	// 1表示读了，将1左移bitIdx后，在与字节数组相应的字节做或，就可以把那一位变为1
	// 以bitIdx为4为例
	// 在bitmap中：[0,0,0,0,0,0,0,0]
	// 1 << 4   : [      1,0,0,0,0]
	// 做或之后，在bitmap中：[0,0,0,1,0,0,0,0]
	// 如此，就实现了某一用户在一bit的实现
	b.bits[byteIdx] |= 1 << bitIdx
}

// IsSet 检查给定id的位是否在Bitmap中被设置为1。
// 参数:
//   id - 要检查的位的唯一标识符。
// 返回值:
//   如果位被设置为1，则返回true；否则返回false。
func (b *Bitmap) IsSet(id string) bool {
    // 计算id对应的位索引。
    idx := hash(id) % b.size
    // 计算字节索引和位索引。
    byteIdx := idx / 8
    bitIdx := idx % 8
    // 检查位是否被设置为1。
    return (b.bits[byteIdx] & (1 << bitIdx)) != 0
}


// Export 导出bitmap
func (b *Bitmap) Export() []byte {
	return b.bits
}

// Load 导入bitmap
func Load(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return NewBitmap(0)
	}
	return &Bitmap{
		bits: bits,
		size: len(bits) * 8,
	}
}

func hash(id string) int {
	// 使用 BKDR 哈希算法
	seed := 131313 // 31 131 1313 13131 131313 ...
	hash := 0
	for _, c := range id {
		hash = hash * seed + int(c)
	}
	return hash & 0x7FFFFFFF
}

