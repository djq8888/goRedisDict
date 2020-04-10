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
