package goRedisDict

//查找字典中key对应的value
func (d *Dict) FetchValue(key uint64) interface{} {
	if entry := d.find(key); entry != nil {
		return entry.getVal()
	} else {
		return nil
	}
}
