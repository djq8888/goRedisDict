package goRedisDict

/*Redis字典（使用哈希表实现）*/

const(
	//哈希表初始大小
	DICT_HT_INITIAL_SIZE = 4
	//哈希表最大容量（uint64最大值）
	LONG_MAX = ^uint32(0)
	//哈希表是否可以扩容
	dict_can_resize = true
	//强制扩容比率
	dict_force_resize_ratio = 2
)

type Dict struct {
	//两个哈希表
	Ht [2]Dictht
	//记录rehash进度的标志，值为-1表示rehash未进行
	rehash int
}

//创建字典
func Create() *Dict {
	d := new(Dict)
	d.init()
	return d
}

//初始化字典
func (d *Dict) init() {
	d.rehash = -1
}

//判断字典是否正在进行rehash
func (d *Dict) isRehashing() bool {
	return d.rehash != -1
}
