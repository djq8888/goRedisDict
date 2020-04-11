package goRedisDict

/*字典迭代器*/

type DictIterator struct {
	//正在迭代的字典
	d *Dict
	//正在迭代的哈希表下标（0或者1）
	table int
	//正在迭代的哈希表节点数组下标
	index int64
	//是否为安全迭代器（安全迭代器允许迭代过程中对字典进行Add，Find等操作）
	safe bool
	//当前哈希表节点
	entry *DictEntry
	//当前哈希表节点的后继节点
	nextEntry *DictEntry
}

//创建一个不安全迭代器
func GetIterator(d *Dict) *DictIterator {
	iter := new(DictIterator)
	iter.d = d
	iter.index = -1
	iter.entry = nil
	iter.nextEntry = nil
	return iter
}

//创建一个安全迭代器
func GetSafeIterator(d *Dict) *DictIterator {
	iter := GetIterator(d)
	iter.safe = true
	return iter
}

//返回迭代器指向的后继节点，遍历完返回nil
func (di *DictIterator) Next() *DictEntry {
	//找到下一个非空节点
	for {
		if di.entry == nil {
			//第一次遍历或者链表已经遍历到尾节点

			ht := di.d.Ht[di.table]
			//如果是安全迭代器，则增加字典的安全迭代器数量
			if di.safe && di.index == -1 && di.table == 0 {
				di.d.iterators++
			}
			//增加索引
			di.index++
			//判断是否已经遍历完当前哈希表
			if di.index >= int64(ht.size) {
				//判断是否需要遍历下一个哈希表
				if di.d.isRehashing() && di.table == 0 {
					di.table = 1
					di.index = 0
					ht = di.d.Ht[1]
				} else {
					break
				}
			}
			//指向哈希表下一个链表头节点
			di.entry = ht.table[di.index]
		} else {
			di.entry = di.nextEntry
		}

		//保存后继节点（因为当前节点可能被修改）
		if di.entry != nil {
			di.nextEntry = di.entry.next
			return di.entry
		}
	}
	return nil
}