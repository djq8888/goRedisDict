package goRedisDict

import "errors"

//添加key-value对到字典
func (d *Dict) Add(key uint64, val interface{}) error {
	if entry, err := d.addRaw(key); err != nil {
		return err
	} else {
		entry.setVal(key, val)
		return nil
	}
}

//添加一个key节点
func (d *Dict) addRaw(key uint64) (*DictEntry, error) {
	//渐进式rehash
	if d.isRehashing() {
		d.rehashStep()
	}

	//计算新节点的索引，如果key已经存在，返回nil
	index, err := d.keyIndex(key)
	if err != nil {
		return nil, err
	}

	//判断新节点存放在哪个哈希表
	table := 0
	if d.isRehashing() {
		table = 1
	}

	//创建哈希表节点
	entry := new(DictEntry)

	//在哈希表中使用头插法插入新的哈希节点
	entry.next = d.Ht[table].table[index]
	d.Ht[table].table[index] = entry

	d.Ht[table].used++

	//设置哈希表节点的key
	entry.setKey(key)

	return entry, nil
}

func (d *Dict) keyIndex(key uint64) (uint64, error) {
	//首先判断是否需要扩容哈希表
	if err := d.expandIfNeeded(); err != nil {
		return 0, err
	}

	//计算key的哈希
	h := d.hashKey(key)

	//查找key是否已经存在
	var idx uint64
	for table := 0; table <= 1; table++  {
		//计算key在哈希表中的索引
		idx = h & d.Ht[table].sizemask
		//查找该索引对应的哈希表节点
		he := d.Ht[table].table[idx]
		for he != nil {
			//如果key已经存在，则插入失败
			if he.key == key {
				return 0, errors.New("key exists.")
			}
			//hash冲突
			he = he.next
		}
		//如果当前没有进行rehash，则无需查找哈希表1
		if !d.isRehashing() {
			break
		}
	}
	return idx, nil
}

//哈希函数
func (d *Dict) hashKey(key uint64) uint64 {
	//TODO:源码是根据数据类型采用不同的哈希函数，这部分有些复杂，目前先假设key均为int类型，使用求余法计算哈希值
	return key % 5381
}

//判断是否需要扩展（新建或扩容）哈希表
func (d *Dict) expandIfNeeded() error {
	//如果正在rehash，直接返回
	if d.isRehashing() {
		return nil
	}

	//如果哈希表未创建（说明是第一次添加key），则新建一个默认大小的哈希表
	if d.Ht[0] == nil {
		return d.expand(DICT_HT_INITIAL_SIZE)
	}

	//如果哈希表中节点个数已经到达哈希表大小，或者哈希表中节点个数占比大于dict_force_resize_ratio，则进行扩容
	if (dict_can_resize && d.Ht[0].used >= d.Ht[0].size) || (float32(d.Ht[0].used) / float32(d.Ht[0].size) > float32(dict_force_resize_ratio)) {
		return d.expand(d.Ht[0].used * 2)
	}

	return nil
}

//扩展哈希表
func (d *Dict) expand(size uint64) error {
	//新建哈希表
	ht := Dictht{}

	//计算扩展后哈希表的大小
	realsize := d.nextPower(size)

	//size不应该小于哈希表中元素个数，此时也不应正在rehash
	if d.isRehashing() || (d.Ht[0] != nil && size <= d.Ht[0].used) {
		return errors.New("error while expanding.")
	}

	//初始化新建的哈希表
	ht.used = 0
	ht.size = realsize
	ht.sizemask = realsize - 1
	ht.table = make([]*DictEntry, realsize)

	//如果字典中的哈希表0此时为空，说明是第一次添加key，应将新建的哈希表赋给哈希表0
	if d.Ht[0] == nil {
		d.Ht[0] = &ht
		return nil
	}

	//如果字典中的哈希表0此时不为空，说明是哈希表需要扩容，应将新建的哈希表赋给哈希表1，并将rehash置为0
	d.Ht[1] = &ht
	d.rehash = 0
	return nil
}

//计算拓展后哈希表大小
func (d *Dict) nextPower(size uint64) uint64 {
	var i uint64 = DICT_HT_INITIAL_SIZE

	//确保扩容后大小不会溢出
	if size >= uint64(LONG_MAX) {
		return uint64(LONG_MAX)
	}

	//如果小于DICT_HT_INITIAL_SIZE（说明是新建哈希表），则返回DICT_HT_INITIAL_SIZE
	//否则返回第一个>=size的二次幂
	for {
		if i >= size {
			return i
		}
		i *= 2
	}
}