package driver

import ()

// 非キャッシュ用。
type MemoryEventRegistry struct {
	trees map[string]*eventTree
}

func NewMemoryEventRegistry() *MemoryEventRegistry {
	return &MemoryEventRegistry{map[string]*eventTree{}}
}

func (reg *MemoryEventRegistry) Handler(usrUuid, event string) (Handler, error) {
	tree := reg.trees[usrUuid]
	if tree == nil {
		return nil, nil
	}
	return tree.handler(event), nil
}
func (reg *MemoryEventRegistry) AddHandler(usrUuid, event string, hndl Handler) error {
	tree := reg.trees[usrUuid]
	if tree == nil {
		tree = newEventTree()
		reg.trees[usrUuid] = tree
	}
	tree.add(event, hndl)
	return nil
}
func (reg *MemoryEventRegistry) RemoveHandler(usrUuid, event string) error {
	tree := reg.trees[usrUuid]
	if tree == nil {
		return nil
	}
	tree.remove(event)
	return nil
}
