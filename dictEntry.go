package goRedisDict

/*哈希表节点*/

type DictEntry struct {
	key uint64
	val interface{}
	next *DictEntry
}

//设置哈希表节点的key
func (de *DictEntry) setKey(key uint64)  {
	de.key = key
}

//设置哈希表节点的value
func (de *DictEntry) setVal(key uint64, val interface{})  {
	de.val = val
}

//获取哈希表节点的value
func (de *DictEntry) getVal() interface{} {
	return de.val
}

//获取哈希表节点的key，value（源码中不包含本函数，仅用于测试迭代器）
func (de *DictEntry) Get() (uint64, interface{}) {
	return de.key, de.val
}
