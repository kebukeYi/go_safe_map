package three

// RWMap 使用分段 map+读写锁
type RWMap[Key, Value any] struct {
	tables map[uint]*Table[Key, Value]
	length uint
}

func (r *RWMap[Key, Value]) Load(key Key) (value Value, ok bool) {
	index := hash(key) % r.length
	return r.tables[index].get(key)
}
func (r *RWMap[Key, Value]) Store(key Key, value Value) {
	index := hash(key) % r.length
	r.tables[index].set(key, value)
}
func (r *RWMap[Key, Value]) Delete(key Key) {
	index := hash(key) % r.length
	r.tables[index].delete(key)
}
func (r *RWMap[Key, Value]) Range(f func(key Key, value Value) bool) {
	for _, table := range r.tables {
		if !table.rangeDB(f) {
			break
		}
	}
}

func newRWMap[Key, Value any](length uint) IMap[Key, Value] {
	if length == 0 {
		length = 1007
	}
	rwMap := &RWMap[Key, Value]{
		tables: make(map[uint]*Table[Key, Value], length),
		length: length,
	}
	for i := uint(0); i < length; i++ {
		rwMap.tables[i] = newTable[Key, Value]()
	}
	return rwMap
}
