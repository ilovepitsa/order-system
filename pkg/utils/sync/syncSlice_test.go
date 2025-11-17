package utils_test

import (
	"testing"
	"testing/synctest"

	utils "github.com/ilovepitsa/orders/pkg/utils/sync"
)

func Test_SyncSlice(t *testing.T) {
	tests := []struct {
		testName        string
		AmountGorutines int
		initSlice       []int
		gorutinesAct    []func(*utils.SyncSlice[int])
		want            []int
	}{
		{
			testName:        "Test append",
			AmountGorutines: 1,
			initSlice:       []int{},
			gorutinesAct: []func(*utils.SyncSlice[int]){
				func(ss *utils.SyncSlice[int]) { ss.Push(1) },
			},
			want: []int{1},
		},
		{
			testName:        "Test pop",
			AmountGorutines: 1,
			initSlice:       []int{1},
			gorutinesAct: []func(*utils.SyncSlice[int]){
				func(ss *utils.SyncSlice[int]) { ss.Pop() },
			},
			want: []int{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				testSlice := utils.NewSyncSlice[int](len(tc.initSlice), cap(tc.initSlice))
				for i, v := range tc.initSlice {
					testSlice.Set(i, v)
				}
				for i := range tc.AmountGorutines {
					go tc.gorutinesAct[i](testSlice)
				}
				synctest.Wait()
				if len(tc.want) != testSlice.Len() {
					t.Errorf("slices not equal. Want:%v . Got: %v", tc.want, testSlice)
					return
				}

				for i, v := range tc.want {
					testVal := testSlice.Get(i)
					if v != testVal {
						t.Errorf("slices not equalt. Index %v. Want val: %v. Got val: %v", i, v, testVal)
					}
				}
			})
		})
	}
}
