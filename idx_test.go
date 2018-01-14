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
	action := &Action{Idx: testIdx}
	action.Idx = nil
	action.TargetTime = time.Date(2018, 1, 13, 0, 0, 0, 0, time.UTC)
	if action.Result != ResultUncheck {
		t.Fatalf("Failed check unknown action\nhave: %v\nwant: %v", action.Result, ResultUncheck)
	}
	d := time.Date(2018, 1, 12, 0, 0, 0, 0, time.UTC)
	action1 := Do(action, &d)
	if action1.Result != ResultOK {
		t.Fatalf("Failed check ok action\nhave: %v\nwant: %v", action1.Result, ResultOK)
	}
	d = time.Date(2018, 1, 14, 0, 0, 0, 0, time.UTC)
	action2 := Do(action, &d)
	if action2.Result != ResultNG {
		t.Fatalf("Failed check ok action\nhave: %v\nwant: %v", action2.Result, ResultOK)
	}
	if !reflect.DeepEqual(action2.ActualTime, d) {
		t.Fatalf("ActualTime is not Matched\nhave: %v\nwant: %v", action2.ActualTime, d)
	}
}

func TestChangeTargetDateAction(t *testing.T) {
	action := &Action{Idx: testIdx}
	action.TargetTime = time.Date(2018, 1, 13, 0, 0, 0, 0, time.UTC)

	d := time.Date(2018, 1, 12, 0, 0, 0, 0, time.UTC)
	changed := ChangeTargetTime(action, &d, "reschedule")

	if reflect.DeepEqual(d, changed.TargetTime) {
		t.Fatalf("Failed Changed TargetTime\nhave: %v\nwant: %v", d, changed.TargetTime)
	}
	if changed.Result != ResultUncheck {
		t.Fatalf("changed ResultStatus is not Uncheck: %v", changed.Result)
	}
	if action.Result != ResultChanged {
		t.Fatalf("original ResultStatus is not ResultChanged: %v", action.Result)
	}

}