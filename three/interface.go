package three

type IMap[Key, Value any] interface {
	Load(key Key) (value Value, ok bool)
	Store(key Key, value Value)
	Delete(Key Key)
	Range(f func(key Key, value Value) bool)
}
