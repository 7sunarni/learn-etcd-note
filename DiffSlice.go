package utils

type Comparable interface {
	NoMore(than Comparable) bool
}

func DiffSlice(s1, s2 []Comparable) ( /*onlyInS1*/ []Comparable /*onlyInS2*/, []Comparable) {
	if len(s1) == 0 {
		temp := make([]Comparable, len(s2))
		copy(temp, s2)
		return nil, temp
	}

	if len(s2) == 0 {
		temp := make([]Comparable, len(s1))
		copy(temp, s1)
		return temp, nil
	}

	// TODO: sort old new
	onlyInS1 := make([]Comparable, 0)
	onlyInS2 := make([]Comparable, 0)
	topIsS1 := true
	top := s1
	bottom := s2

	var maxLength int
	if len(s2) > len(s1) {
		maxLength = len(s2)
	} else {
		maxLength = len(s1)
	}
	for topIndex := 0; topIndex < maxLength; topIndex++ {

		// 最后一个元素
		if len(top)-1 == topIndex && len(bottom)-1 == topIndex {
			if bottom[topIndex] != top[topIndex] {
				if topIsS1 {
					onlyInS2 = append(onlyInS2, bottom[topIndex])
					onlyInS1 = append(onlyInS1, top[topIndex])
				} else {
					onlyInS2 = append(onlyInS2, top[topIndex])
					onlyInS1 = append(onlyInS1, bottom[topIndex])
				}
			}
			break
		}

		if len(top)-1 == topIndex {
			if topIsS1 {
				onlyInS2 = append(onlyInS2, bottom[topIndex+1:]...)
			} else {
				onlyInS1 = append(onlyInS1, bottom[topIndex+1:]...)
			}
			break
		}

		if len(bottom)-1 == topIndex {
			if topIsS1 {
				onlyInS1 = append(onlyInS1, top[topIndex+1:]...)
			} else {
				onlyInS2 = append(onlyInS2, top[topIndex+1:]...)
			}
			break
		}

		topItem := top[topIndex]

		bottomItem := bottom[topIndex]
		if bottomItem.NoMore(topItem) && topItem.NoMore(bottomItem) {
			continue
		}

		if bottomItem.NoMore(topItem) {
			// 交换数组，始终保证较大的数在上面
			temp := top
			top = bottom
			bottom = temp
			topIsS1 = !topIsS1
		}
		if topIsS1 {
			onlyInS1 = append(onlyInS1, top[topIndex])
		} else {
			onlyInS2 = append(onlyInS2, top[topIndex])
		}
		// 当前的值要小一点，移动一格
		top = append(top[0:topIndex], top[topIndex+1:]...)
		topIndex--
	}

	return onlyInS1, onlyInS2
}

func DiffApps(s1, s2 []Service) ( /*onlyInS1*/ []Service /*onlyInS2*/, []Service) {

	if len(s1) == 0 {
		temp := make([]Service, len(s2))
		copy(temp, s2)
		return nil, temp
	}

	if len(s2) == 0 {
		temp := make([]Service, len(s1))
		copy(temp, s1)
		return temp, nil
	}
	// TODO: sort old new
	onlyInS1 := make([]Service, 0)
	onlyInS2 := make([]Service, 0)
	topIsS1 := true
	top := s1
	bottom := s2

	var maxLength int
	if len(s2) > len(s1) {
		maxLength = len(s2)
	} else {
		maxLength = len(s1)
	}
	for topIndex := 0; topIndex < maxLength; topIndex++ {

		// 最后一个元素
		if len(top)-1 == topIndex && len(bottom)-1 == topIndex {
			if bottom[topIndex] != top[topIndex] {
				if topIsS1 {
					onlyInS2 = append(onlyInS2, bottom[topIndex])
					onlyInS1 = append(onlyInS1, top[topIndex])
				} else {
					onlyInS2 = append(onlyInS2, top[topIndex])
					onlyInS1 = append(onlyInS1, bottom[topIndex])
				}
			}
			break
		}

		if len(top)-1 == topIndex {
			if topIsS1 {
				onlyInS2 = append(onlyInS2, bottom[topIndex+1:]...)
			} else {
				onlyInS1 = append(onlyInS1, bottom[topIndex+1:]...)
			}
			break
		}

		if len(bottom)-1 == topIndex {
			if topIsS1 {
				onlyInS1 = append(onlyInS1, top[topIndex+1:]...)
			} else {
				onlyInS2 = append(onlyInS2, top[topIndex+1:]...)
			}
			break
		}

		topItem := top[topIndex]

		bottomItem := bottom[topIndex]
		if bottomItem.NoMore(topItem) && topItem.NoMore(bottomItem) {
			continue
		}

		if bottomItem.NoMore(topItem) {
			// 交换数组，始终保证较大的数在上面
			temp := top
			top = bottom
			bottom = temp
			topIsS1 = !topIsS1
		}
		if topIsS1 {
			onlyInS1 = append(onlyInS1, top[topIndex])
		} else {
			onlyInS2 = append(onlyInS2, top[topIndex])
		}
		// 当前的值要小一点，移动一格
		top = append(top[0:topIndex], top[topIndex+1:]...)
		topIndex--
	}

	return onlyInS1, onlyInS2
}
