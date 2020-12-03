package utils

import (
	"errors"
	"log"
	"sync"

	launcher "git.forchange.cn/launcher/launcher-api/v1beta1"

	"github.com/google/btree"
)

const (
	defaultTreeDepth = 10
	defaultArraySize = 1 << 20
)

var (
	ErrGroupNotExist  = errors.New("Group not exist")
	ErrItemNotExist   = errors.New("Item not exist")
	ErrItemHasExisted = errors.New("Item has existed")
	ErrSetItemIndex   = errors.New("Miss free position")
)

type BPlusTreeCacher struct {
	indexs []*btree.BTree
	// 数组中空闲的位置
	freePosition []int
	// 存放数据的数组
	freeArray []container
	lock      *sync.Mutex
}

type container struct {
	hasV bool
	v    interface{}
}

var _ btree.Item = &indexNode{}

// 索引节点，节点的值(value)存放的是数据所在的下标
type indexNode struct {
	key string
	// positions 的值指向的是 freeArray 的下标
	positions []int
}

func (c *BPlusTreeCacher) getGroup(grpName string) *indexNode {
	var grpIndexNode *indexNode
	item := c.indexs[0].Get(&indexNode{key: grpName})
	if item == nil {
		grpIndexNode = &indexNode{key: grpName, positions: make([]int, 0)}
	} else {
		grpIndexNode = item.(*indexNode)
	}
	return grpIndexNode
}

func (c *BPlusTreeCacher) getNamespace(namespace string) *indexNode {
	var nsIndexNode *indexNode
	item := c.indexs[1].Get(&indexNode{key: namespace})
	if item == nil {
		nsIndexNode = &indexNode{key: namespace, positions: make([]int, 0)}
	} else {
		nsIndexNode = item.(*indexNode)
	}
	return nsIndexNode
}

func (in *indexNode) Less(i btree.Item) bool {
	than := i.(*indexNode)
	return in.key < than.key
}

func (c *BPlusTreeCacher) putArrayItem(r launcher.Service) int {
	if len(c.freePosition) == 0 {
		// TODO: 扩容动作
		log.Printf("Error: array full, need expand")
		return -1
	}

	// 取最后一个元素
	freePos := c.freePosition[len(c.freePosition)-1]
	c.freePosition = c.freePosition[0 : len(c.freePosition)-2]
	c.freeArray[freePos] = container{hasV: true, v: r}
	return freePos
}

func (c *BPlusTreeCacher) Put(r launcher.Service) error {
	pos := c.findPositionIndexByGroup(r.GroupName, r.Namespace, r.Name)
	if pos > -1 {
		return ErrItemHasExisted
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	emptyPos := c.putArrayItem(r)
	if emptyPos < 0 {
		return ErrSetItemIndex
	}
	c.setGroupIndex(r.GroupName, emptyPos)
	c.setNamespaceIndex(r.Namespace, emptyPos)
	return nil
}

func (c *BPlusTreeCacher) setNamespaceIndex(namespace string, pos int) {
	nsIndexNode := c.getNamespace(namespace)
	nsIndexNode.positions = append(nsIndexNode.positions, pos)
	c.indexs[1].ReplaceOrInsert(nsIndexNode)
}

func (c *BPlusTreeCacher) setGroupIndex(groupName string, pos int) {
	grpIndexNode := c.getGroup(groupName)
	grpIndexNode.positions = append(grpIndexNode.positions, pos)
	c.indexs[0].ReplaceOrInsert(grpIndexNode)
}
func removeItem(slice []int, value int) {
	for i, _value := range slice {
		if _value == value {
			slice = append(slice[0:i], slice[i+1:]...)
			return
		}
	}
}

func (c *BPlusTreeCacher) clearGroupIndex(groupName string, pos int) {
	grpIndexNode := c.getGroup(groupName)
	removeItem(grpIndexNode.positions, pos)
	c.indexs[0].ReplaceOrInsert(grpIndexNode)
}

func (c *BPlusTreeCacher) clearNamespaceIndex(namesapce string, pos int) {
	nsIndexNode := c.getNamespace(namesapce)
	removeItem(nsIndexNode.positions, pos)
	c.indexs[1].ReplaceOrInsert(nsIndexNode)
}

func NewIndexCacher() *BPlusTreeCacher {
	ic := &BPlusTreeCacher{
		freePosition: make([]int, 0),
		freeArray:    make([]container, defaultArraySize),
		indexs:       []*btree.BTree{btree.New(defaultTreeDepth), btree.New(defaultTreeDepth)},
		lock:         &sync.Mutex{},
	}

	for i := 0; i < defaultArraySize; i++ {
		ic.freePosition = append(ic.freePosition, i)
	}
	return ic
}

func (c *BPlusTreeCacher) Remove(r launcher.Service) error {
	pos := c.findPositionIndexByGroup(r.GroupName, r.Namespace, r.Name)
	if pos < 0 {
		return ErrItemNotExist
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.freePosition = append(c.freePosition, pos)
	c.freeArray[pos].hasV = false
	c.clearGroupIndex(r.GroupName, pos)
	c.clearNamespaceIndex(r.Namespace, pos)
	return nil
}

func (c *BPlusTreeCacher) findPositionIndexByGroup(groupName, namespace, name string) int {
	grpIndexNode := c.getGroup(groupName)
	for _, pos := range grpIndexNode.positions {
		if len(c.freeArray)-1 < pos {
			log.Printf("GetByGroup index length error")
			return -1
		}
		container := c.freeArray[pos]
		if !container.hasV {
			log.Printf("GetByGroup get empty value")
			continue
		}
		launcher := container.v.(launcher.Service)
		if launcher.Name == name && launcher.Namespace == namespace {
			return pos
		}
	}
	return -2
}

func (c *BPlusTreeCacher) findPositionIndexByNamesapce(groupName, namespace, name string) int {
	nsIndexNode := c.getNamespace(namespace)
	for _, pos := range nsIndexNode.positions {
		if len(c.freeArray)-1 < pos {
			log.Printf("Error: index length error")
			return -1
		}
		container := c.freeArray[pos]
		if !container.hasV {
			log.Printf("Error: got empty value")
			continue
		}
		launcher := container.v.(launcher.Service)
		if launcher.Name == name && launcher.Namespace == namespace {
			return pos
		}
	}
	return -2
}

func (c *BPlusTreeCacher) GetByGroup(r launcher.Service) *launcher.Service {
	/*
		用两个索引中数量较少的来查询
		grpIndexNode := c.getGroup(groupName)
		nsIndexNode := c.getNamespace(namespace)
		var positions []int
		if len(grpIndexNode.positions) < len(nsIndexNode.positions) {
			positions = grpIndexNode.positions
		} else {
			positions = nsIndexNode.positions
		}

		for _, pos := range positions {
			if len(c.freeArray)-1 < pos {
				log.Printf("Error: index length error")
				return nil
			}
			container := c.freeArray[pos]
			if !container.hasV {
				log.Printf("Error: got empty value")
				continue
			}
			launcher := container.v.(launcher.Service)
			if launcher.Name == name && launcher.Namespace == namespace {
				return &launcher
			}
		}
	*/
	pos := c.findPositionIndexByGroup(r.GroupName, r.Namespace, r.Name)
	if pos < 0 {
		return nil
	}
	ret := c.freeArray[pos].v.(launcher.Service)
	return &ret
}

func (c *BPlusTreeCacher) GetGroupRegistrations(groupName string) []launcher.Service {
	var services []launcher.Service
	grpIndexNode := c.getGroup(groupName)
	for _, pos := range grpIndexNode.positions {
		container := c.freeArray[pos]
		service := container.v.(launcher.Service)
		if container.hasV {
			services = append(services, service)
		}
	}
	return services
}
