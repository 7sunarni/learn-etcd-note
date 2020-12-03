package utils

import (
	"fmt"
	"log"
	"testing"
	"time"

	"git.forchange.cn/launcher/launcher-api/v1beta1"
)

func TestPutAndRead(t *testing.T) {
	cacher := NewIndexCacher()
	testCases := []struct {
		Name      string
		GroupName string
		Reg       v1beta1.Service
		Method    string
		WantReg   v1beta1.Service
		IsError   bool
	}{
		{
			Name:      "Put first",
			GroupName: "001",
			Reg: v1beta1.Service{
				Name:      "name-1",
				GroupName: "001",
				Namespace: "namesapce-1",
				Endpoint:  v1beta1.Endpoint{Comment: "hello"},
			},
			Method:  "PUT",
			WantReg: v1beta1.Service{},
			IsError: false,
		},

		{
			Name:      "Re Put it",
			GroupName: "001",
			Reg: v1beta1.Service{
				Name:      "name-1",
				GroupName: "001",
				Namespace: "namesapce-1",
				Endpoint:  v1beta1.Endpoint{Comment: "hello"},
			},
			Method:  "PUT",
			WantReg: v1beta1.Service{},
			IsError: true,
		},

		{
			Name:      "Get it",
			GroupName: "001",
			Reg: v1beta1.Service{
				Name:      "name-1",
				GroupName: "001",
				Namespace: "namesapce-1",
			},
			Method: "GET",
			WantReg: v1beta1.Service{
				Name:      "name-1",
				Namespace: "namesapce-1",
				GroupName: "001",
				Endpoint:  v1beta1.Endpoint{Comment: "hello"},
			},
			IsError: false,
		},

		{
			Name:      "Delete it",
			GroupName: "001",
			Reg: v1beta1.Service{
				Name:      "name-1",
				GroupName: "001",
				Namespace: "namesapce-1",
			},
			Method: "DELETE",
			WantReg: v1beta1.Service{
				Name:      "name-1",
				GroupName: "001",
				Namespace: "namesapce-1",
				Endpoint:  v1beta1.Endpoint{Comment: "hello"},
			},
			IsError: false,
		},

		{
			Name:      "Re Get it",
			GroupName: "001",
			Reg: v1beta1.Service{
				Name:      "name-1",
				GroupName: "001",
				Namespace: "namesapce-1",
			},
			Method: "DELETE",
			WantReg: v1beta1.Service{
				Name:      "name-1",
				GroupName: "001",
				Namespace: "namesapce-1",
				Endpoint:  v1beta1.Endpoint{Comment: "hello"},
			},
			IsError: true,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.Name, func(t *testing.T) {
			if tCase.Method == "PUT" {
				err := cacher.Put(tCase.Reg)
				isErr := err != nil
				if isErr != tCase.IsError {
					t.Fatalf("expect success, but got error %s", err)
				}
			}

			if tCase.Method == "GET" {
				got := cacher.GetByGroup(tCase.Reg)
				if tCase.WantReg != *got {
					t.Fatalf("expect %+v, but got error %+v", tCase.WantReg, *got)
				}
			}

			if tCase.Method == "DELETE" {
				err := cacher.Remove(tCase.Reg)
				isErr := err != nil
				if isErr != tCase.IsError {
					t.Fatalf("expect success, but got error %s", err)
				}
			}
		})
	}

}

func TestRange(t *testing.T) {
	cacher := NewIndexCacher()
	for i := 0; i < 30000; i++ {
		if err := cacher.Put(v1beta1.Service{
			GroupName: "group-001",
			Name:      fmt.Sprintf("name-%d", i),
			Namespace: "namesapce-1",
			Endpoint:  v1beta1.Endpoint{Comment: fmt.Sprintf("hello, there are some comment=>%d", i)},
		}); err != nil {
			t.Fatalf("put error")
		}
	}

	{
		start := time.Now()
		if err := cacher.Put(v1beta1.Service{
			GroupName: "group-001",
			Name:      fmt.Sprintf("name-%d", -1),
			Namespace: "namesapce-1",
			Endpoint:  v1beta1.Endpoint{Comment: fmt.Sprintf("hello, there are some comment=>%d", -1)},
		}); err != nil {
			t.Fatalf("put error")
		}
		fmt.Println(time.Since(start).Microseconds())
	}

	{
		// start := time.Now()
		// launcher := cacher.GetByGroup("group-001", "namesapce-1", "name-01")
		// fmt.Println(time.Since(start).Microseconds(), launcher)
	}
}

func TestRemoveMiddle(t *testing.T) {
	cacher := NewIndexCacher()
	for i := 0; i < 1000; i++ {
		if err := cacher.Put(v1beta1.Service{
			GroupName: "group-001",
			Name:      fmt.Sprintf("name-%d", i),
			Namespace: "namesapce-1",
			Endpoint:  v1beta1.Endpoint{Comment: fmt.Sprintf("hello, there are some comment=>%d", i)},
		}); err != nil {
			t.Fatalf("put error")
		}
	}

	log.Println(cacher.findPositionIndexByGroup("group-001", "namesapce-1", fmt.Sprintf("name-%d", 500)))

	cacher.Remove(v1beta1.Service{
		GroupName: "group-001",
		Name:      fmt.Sprintf("name-%d", 500),
		Namespace: "namesapce-1",
		Endpoint:  v1beta1.Endpoint{Comment: fmt.Sprintf("hello, there are some comment=>%d", 500)},
	})

	if err := cacher.Put(v1beta1.Service{
		GroupName: "group-001",
		Name:      fmt.Sprintf("name-%d", 2000),
		Namespace: "namesapce-1",
		Endpoint:  v1beta1.Endpoint{Comment: fmt.Sprintf("hello, there are some comment=>%d", 2000)},
	}); err != nil {
		t.Fatalf("put error")
	}

	// 位置应该和之前的相同
	log.Println(cacher.findPositionIndexByGroup("group-001", "namesapce-1", fmt.Sprintf("name-%d", 2000)))
}
