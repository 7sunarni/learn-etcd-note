package utils

func TestDiff(t *testing.T) {
	testCases := []struct {
		Name        string
		Old         []Comparable
		New         []Comparable
		WantOnlyOld []Comparable
		WantOnlyNew []Comparable
	}{
		{
			Name:        "",
			Old:         []Comparable{ComparableInt(3), ComparableInt(5), ComparableInt(6), ComparableInt(7), ComparableInt(8), ComparableInt(9), ComparableInt(10), ComparableInt(12)},
			New:         []Comparable{ComparableInt(3), ComparableInt(4), ComparableInt(6), ComparableInt(7), ComparableInt(9), ComparableInt(10), ComparableInt(11)},
			WantOnlyOld: []Comparable{ComparableInt(5), ComparableInt(8)},
			WantOnlyNew: []Comparable{ComparableInt(4), ComparableInt(11)},
		},
		{
			Name:        "",
			Old:         []Comparable{ComparableInt(1)},
			New:         []Comparable{ComparableInt(2)},
			WantOnlyOld: []Comparable{ComparableInt(1)},
			WantOnlyNew: []Comparable{ComparableInt(2)},
		},
		{
			Name:        "",
			Old:         []Comparable{ComparableInt(1), ComparableInt(2), ComparableInt(3)},
			New:         []Comparable{ComparableInt(2), ComparableInt(3), ComparableInt(4)},
			WantOnlyOld: []Comparable{ComparableInt(1)},
			WantOnlyNew: []Comparable{ComparableInt(4)},
		},
		{
			Name:        "",
			Old:         []Comparable{ComparableInt(1), ComparableInt(2), ComparableInt(3), ComparableInt(4)},
			New:         []Comparable{ComparableInt(2), ComparableInt(3), ComparableInt(4)},
			WantOnlyOld: []Comparable{ComparableInt(1)},
			WantOnlyNew: []Comparable{},
		},

		{
			Name:        "",
			Old:         []Comparable{ComparableInt(2), ComparableInt(3), ComparableInt(4)},
			New:         []Comparable{ComparableInt(2), ComparableInt(3), ComparableInt(4), ComparableInt(5)},
			WantOnlyOld: []Comparable{},
			WantOnlyNew: []Comparable{ComparableInt(5)},
		},

		{
			Name:        "",
			Old:         []Comparable{},
			New:         []Comparable{ComparableInt(2)},
			WantOnlyOld: []Comparable{},
			WantOnlyNew: []Comparable{ComparableInt(2)},
		},

		{
			Name:        "",
			Old:         []Comparable{ComparableInt(2)},
			New:         []Comparable{},
			WantOnlyOld: []Comparable{ComparableInt(2)},
			WantOnlyNew: []Comparable{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			gotOnlyOld, gotOnlyNew := DiffSlice(testCase.Old, testCase.New)
			if !reflect.DeepEqual(gotOnlyOld, testCase.WantOnlyOld) && !reflect.DeepEqual(gotOnlyNew, testCase.WantOnlyNew) {
				t.Fatalf("want only old => %+v, got only old => %+v \r\n want only new => %+v, got only new => %+v",
					testCase.WantOnlyOld,
					gotOnlyOld,
					testCase.WantOnlyNew,
					gotOnlyNew)
			}
		})
	}
}

var _ Comparable = ComparableInt(1)

type ComparableInt int

// TODO: 使用 sort.Interface 接口来做比较
func (i ComparableInt) NoMore(than Comparable) bool {
	t := than.(ComparableInt)
	return i <= t
}
