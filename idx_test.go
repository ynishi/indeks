package indeks

import (
	"reflect"
	"testing"
	"time"
)

var (
	testIdx    *Idx
	testAction *Action
)

func init() {
	testIdx = &Idx{Name: "idx1", Desc: "One of idx", DefaultDuration: time.Duration(24) * time.Hour, DefaultPoint: 1, Actions: nil}
	testAction = &Action{Idx: testIdx}
}

func TestInit(t *testing.T) {
	idx := &Idx{Name: "idx1", Desc: "One of idx", DefaultDuration: time.Duration(24) * time.Hour, DefaultPoint: 1, Actions: nil}

	if !reflect.DeepEqual(testIdx, idx) {
		t.Fatalf("Idx Init Failed\nhave: %v\nwant: %v", idx, testIdx)
	}
}

func TestCreateAction(t *testing.T) {
	idx := &Idx{Name: "idx1", Desc: "One of idx", DefaultPoint: 1, Actions: nil}

	action := CreateAction(idx, time.Time{})
	action.Idx.Name = "changed"

	if !reflect.DeepEqual(idx, action.Idx) {
		t.Fatalf("Failed action links idx\nhave: %v\nwant: %v", action.Idx, idx)
	}
	if !reflect.DeepEqual(idx.Actions[0], action) {
		t.Fatalf("Failed add action to idx\nhave: %v\nwant: %v", idx.Actions[0], action)
	}
}

func TestCheckAction(t *testing.T) {

}
