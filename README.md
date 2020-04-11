# goRedisDict
redis字典的go语言实现<br>
* 基于redis2.6
* 参考《redis设计与实现》
## 已实现方法
* Create
* Add
* Replace
* FetchValue
* Delete
* Next(迭代器)
## rehash
* 支持渐进式rehash（增加/删除/查找时进行）
## 与源码的差异
* key只支持uint64类型
* 哈希函数只支持求余法
* 未实现内存释放相关函数（内存释放通过go语言垃圾回收实现）