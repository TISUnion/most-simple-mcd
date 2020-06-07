package utils

import "encoding/json"

// 链表
type List struct {
	_list  []*Item
	length int
}

// 链表元素
type Item struct {
	value    interface{}
	typeName string
}

func (l *List) LPush(value interface{}) {
	l.LPushWithTypeName(value, "")
}

func (l *List) RPush(value interface{}) {
	l.RPushWithTypeName(value, "")
}

func (l *List) LPushWithTypeName(value interface{}, typeName string) {
	temp := append(make([]*Item, 0), &Item{
		value:    value,
		typeName: typeName,
	})
	temp = append(temp, l._list...)
	l._list = temp
	l.length++
}

func (l *List) RPushWithTypeName(value interface{}, typeName string) {
	l._list = append(l._list, &Item{
		value:    value,
		typeName: typeName,
	})
	l.length++
}

func (l *List) InsertWithIndex(i int, value interface{}) {
	l.InsertWithIndexAndTypeName(i, value, "")
}

func (l *List) InsertWithIndexAndTypeName(i int, value interface{}, typeName string) {
	if i > l.length {
		i = l.length
	}
	it := &Item{
		value:    value,
		typeName: typeName,
	}
	temp := make([]*Item, 0)
	temp = append(temp, l._list[:i]...)
	temp = append(temp, it)
	temp = append(temp, l._list[i:]...)
	l._list = temp
	l.length++
}

func (l *List) DelWithIndex(i int) {
	if i > l.length {
		i = l.length
	}
	l._list = append(l._list[:i], l._list[i+1:]...)
	l.length--
}

func (l *List) DelByCallback(cb func(int, interface{}) bool) {
	for i, v := range l._list {
		if cb(i, v) {
			l.DelWithIndex(i)
			l.length--
		}
	}
}

func (l *List) GetLen() int {
	return l.length
}

func (l *List) ToArray() []interface{} {
	result := make([]interface{}, 0)
	list := l._list
	for _, v := range list {
		result = append(result, v)
	}

	return result
}

func (l *List) ToJson() (string, error) {
	objArr := l.ToArray()
	jsonByte, e := json.Marshal(objArr)
	if e != nil {
		return "", e
	}

	return string(jsonByte), nil
}