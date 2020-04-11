package goRedisDict

//渐进式rehash
func (d *Dict) rehashStep() {
	//只在安全迭代器为0的情况下进行
	//避免重复遍历（在哈希表0遍历过key之后，key被rehash到哈希表1，又被遍历一次）
	if d.iterators == 0 {
		d.dictRehash(1)
	}
}

//rehash（rehash n个哈希表节点）
//如果rehash完成，返回0，否则返回1
func (d *Dict) dictRehash(n int) int {
	if !d.isRehashing() {
		return 0
	}

	for ; n > 0; n-- {
		//如果哈希表0已经为空（说明已经rehash完毕），则用哈希表1替换哈希表0
		if d.Ht[0].used == 0 {
			d.Ht[0] = d.Ht[1]
			d.Ht[1] = nil
			d.rehash = -1
			return 0
		}

		//找到下一个不为空的索引
		for d.Ht[0].table[d.rehash] == nil {
			d.rehash++
		}

		//将链表中所有节点迁移到哈希表1
		de := d.Ht[0].table[d.rehash]
		for de != nil {
			nextde := de.next
			//计算节点在哈希表1中的索引
			h := d.hashKey(de.key) & d.Ht[1].sizemask
			de.next = d.Ht[1].table[h]
			d.Ht[1].table[h] = de
			de = nextde
			//更新计数器
			d.Ht[0].used--
			d.Ht[1].used++
		}

		//迁移后将原数组节点置为nil
		d.Ht[0].table[d.rehash] = nil
		//前进至下一索引
		d.rehash++
	}

	return 1
}