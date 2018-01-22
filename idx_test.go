package indeks

import (
	"context"
	"reflect"
	"testing"
	"time"

	firebase "firebase.google.com/go"
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

func TestChangeTargetTimeAction(t *testing.T) {
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

func TestRemoveAction(t *testing.T) {
	raIdx := &Idx{Name: "idxra", Desc: "idxra", DefaultPoint: 1, Actions: nil}
	action := CreateAction(raIdx, time.Date(2018, 1, 13, 0, 0, 0, 0, time.UTC))
	action2 := CreateAction(raIdx, time.Date(2018, 1, 14, 0, 0, 0, 0, time.UTC))

	raIdx = RemoveAction(raIdx, action)
	if len(raIdx.Actions) != 1 {
		t.Fatalf("Removed raIdx not 1: %v", len(raIdx.Actions))
	}
	if !reflect.DeepEqual(action2, raIdx.Actions[0]) {
		t.Fatalf("Failed match remained action\nhave: %v\nwant: %v", raIdx.Actions[0], action2)
	}
}

func TestSumActualPointIdx(t *testing.T) {
	spIdx := &Idx{Name: "spidx", Desc: "spidx", DefaultPoint: 1, Actions: nil}
	action1 := CreateAction(spIdx, time.Date(2018, 1, 13, 0, 0, 0, 0, time.UTC))
	action1.Point = 2
	action2 := CreateAction(spIdx, time.Date(2018, 1, 14, 0, 0, 0, 0, time.UTC))
	CreateAction(spIdx, time.Date(2018, 1, 15, 0, 0, 0, 0, time.UTC))

	t1 := time.Date(2018, 1, 12, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2018, 1, 13, 0, 0, 0, 0, time.UTC)
	Do(action1, &t1)
	Do(action2, &t2)

	start := time.Date(2018, 1, 10, 0, 0, 0, 0, time.UTC)

	expected := 3
	actual := SumActualPoint(spIdx, start, nil)
	if expected != actual {
		t.Fatal("Sum result not match\nhave %v\nwant %v", actual, expected)
	}
}

func TestSumTargetPointIdx(t *testing.T) {
	spIdx := &Idx{Name: "spidx", Desc: "spidx", DefaultPoint: 1, Actions: nil}
	action1 := CreateAction(spIdx, time.Date(2018, 1, 13, 0, 0, 0, 0, time.UTC))
	action1.Point = 2
	action2 := CreateAction(spIdx, time.Date(2018, 1, 14, 0, 0, 0, 0, time.UTC))
	action3 := CreateAction(spIdx, time.Date(2018, 1, 15, 0, 0, 0, 0, time.UTC))

	t1 := time.Date(2018, 1, 12, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2018, 1, 13, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2018, 1, 16, 0, 0, 0, 0, time.UTC)

	Do(action1, &t1)
	Do(action2, &t2)
	ChangeTargetTime(action3, &t3, "test")

	start := time.Date(2018, 1, 13, 1, 0, 0, 0, time.UTC)
	end := time.Date(2018, 1, 17, 0, 0, 0, 0, time.UTC)

	expected := 2
	actual := SumTargetPoint(spIdx, start, &end)
	if expected != actual {
		t.Fatalf("Sum result not match\nhave %v\nwant %v", actual, expected)
	}
}

func TestSummaryPointIdx(t *testing.T) {
	spIdx := &Idx{Name: "spidx", Desc: "spidx", DefaultPoint: 1, Actions: nil}
	action1 := CreateAction(spIdx, time.Date(2018, 1, 13, 0, 0, 0, 0, time.UTC))
	action1.Point = 2
	action2 := CreateAction(spIdx, time.Date(2018, 1, 15, 0, 0, 0, 0, time.UTC))
	action2.Point = 3
	action3 := CreateAction(spIdx, time.Date(2018, 1, 16, 0, 0, 0, 0, time.UTC))

	t1 := time.Date(2018, 1, 12, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2018, 1, 14, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2018, 1, 17, 0, 0, 0, 0, time.UTC)

	Do(action1, &t1)
	Do(action2, &t2)
	ChangeTargetTime(action3, &t3, "test")

	start := time.Date(2018, 1, 13, 1, 0, 0, 0, time.UTC)
	end := time.Date(2018, 1, 18, 0, 0, 0, 0, time.UTC)

	expected := Summary{Target: 4, Actual: 3, Ratio: 0.75, TargetC: 2, ActualC: 1, RatioC: 0.5}
	actual := SummaryPoint(spIdx, start, &end)
	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("Sum result not match\nhave %v\nwant %v", actual, expected)
	}
}

func TestInitFirebase(t *testing.T) {
	// FireStore
	projectID := ""
	conf := &firebase.Config{ProjectID: projectID}
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		t.Fatalf("Firebase app init failed: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		t.Fatalf("Firestore client init failed: %v", err)
	}
	defer client.Close()
}
