package goRedisDict

import "errors"

func (d *Dict) Delete(key uint64) error {
	return d.genericDelete(key)
}

//源码可以选择释放key或者不释放key，go语言变量没有引用之后会被垃圾回收自动释放，因此不考虑key的释放问题
func (d *Dict) genericDelete(key uint64) error {
	//判断是否为空表
	if d.Ht[0].size == 0 {
		return errors.New("delete failed, hash table is nil.")
	}

	//TODO:渐进式rehash
	//if d.isRehashing() {}

	//计算哈希值
	h := d.hashKey(key)

	//在两个哈希表中查找
	for table := 0; table <= 1; table++ {
		idx := h & d.Ht[table].sizemask
		he := d.Ht[table].table[idx]
		//单链表删除
		prev := new(DictEntry)
		for he != nil {
			if compare(key, he.key) {
				if prev != nil {
					//删除的是头节点
					d.Ht[table].table[idx] = he.next
				} else {
					//删除的不是头节点
					prev.next = he.next
				}
				return nil
			}
			prev = he
			he = he.next
		}

		//如果当前没有进行rehash，则无需遍历哈希表1
		if !d.isRehashing() {
			break
		}
	}

	return nil
}
