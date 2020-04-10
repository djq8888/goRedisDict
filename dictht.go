package goRedisDict

/*哈希表*/

type Dictht struct {
	//哈希表节点指针数组
	table []*DictEntry
	//指针数组大小
	size uint64
	//指针数组的长度掩码，用于计算索引值
	sizemask uint64
	//哈希表现有的节点数量
	used uint64
}
