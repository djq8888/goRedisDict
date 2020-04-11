package goRedisDict

import "errors"

//如果key不存在则新建哈希表节点
//如果key存在则更新值
func (d *Dict) Replace(key uint64, val interface{}) error {
	//首先尝试Add，如果key存在会返回error
	if err := d.Add(key, val); err == nil {
		return nil
	}

	//key存在
	//首先找到key对应的哈希表节点
	entry := d.find(key)
	if entry == nil {
		return errors.New("key exists but not found.")
	}
	//如果找到key对应的哈希表节点，则更新该节点的value
	entry.setVal(key, val)
	//源码这里释放了旧值，go语言通过垃圾回收机制自动释放指针，这里无需手动释放

	return nil
}

//查找key对应的哈希表节点
func (d *Dict) find(key uint64) *DictEntry {
	//首先判断哈希表是否为空
	if d.Ht[0].size == 0 {
		return nil
	}

	//渐进式rehash
	if d.isRehashing() {
		d.rehashStep()
	}

	//计算哈希值
	h := d.hashKey(key)

	//在两个哈希表中查找
	for table := 0; table <= 1; table++ {
		idx := h & d.Ht[table].sizemask
		he := d.Ht[table].table[idx]
		for he != nil {
			if compare(key, he.key) {
				return he
			}
			he = he.next
		}

		//如果当前没有进行rehash，则无需遍历哈希表1
		if !d.isRehashing() {
			break
		}
	}

	return nil
}

func compare(key1, key2 uint64) bool {
	return key1 == key2
}