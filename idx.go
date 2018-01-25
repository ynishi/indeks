package indeks

import (
	"time"
)

const (
	ResultUncheck = iota
	ResultOK
	ResultNG
	ResultChanged
)

var Idxs []*Idx

type Idx struct {
	Name            string
	Desc            string
	DefaultPoint    int
	DefaultDuration time.Duration
	Actions         []*Action
}

type Action struct {
	Idx        *Idx
	TargetTime time.Time
	ActualTime time.Time
	Point      int
	Result     int
	Comment    string
}

type Summary struct {
	Target  int
	Actual  int
	Ratio   float32
	TargetC int
	ActualC int
	RatioC  float32
}

func CreateAction(idx *Idx, now time.Time) (action *Action) {
	d := now.Add(idx.DefaultDuration)
	action = &Action{
		Idx:        idx,
		TargetTime: d,
		ActualTime: time.Time{},
		Point:      idx.DefaultPoint,
		Result:     0,
	}
	idx.Actions = append(idx.Actions, action)
	return action
}

func Do(action *Action, t *time.Time) *Action {
	if t == nil {
		action.ActualTime = time.Now()
	} else {
		action.ActualTime = *t
	}
	if action.ActualTime.Before(action.TargetTime) {
		action.Result = ResultOK
	} else {
		action.Result = ResultNG
	}
	return action
}

func ChangeTargetTime(action *Action, t *time.Time, comment string) (changed *Action) {
	changed = CreateAction(action.Idx, *t)
	action.Result = ResultChanged
	action.ActualTime = time.Now()
	action.Comment = comment
	return changed
}

func RemoveAction(idx *Idx, action *Action) (removed *Idx) {
	result := []*Action{}
	for _, act := range idx.Actions {
		if act != action {
			result = append(result, act)
		}
	}
	idx.Actions = result
	return idx
}

func CheckOK(act *Action, start *time.Time, end *time.Time) (ok bool) {
	ok = act.Result == ResultOK && act.ActualTime.After(*start) && act.ActualTime.Before(*end)
	return ok
}

func CheckTarget(act *Action, start *time.Time, end *time.Time) (isT bool) {
	isT = act.Result != ResultChanged && act.TargetTime.After(*start) && act.TargetTime.Before(*end)
	return isT
}

func SumActualPoint(idx *Idx, start time.Time, end *time.Time) (result int) {
	endt := fillEndtime(end)
	for _, act := range idx.Actions {
		if CheckOK(act, &start, endt) {
			result += act.Point
		}
	}
	return result
}

func SumTargetPoint(idx *Idx, start time.Time, end *time.Time) (result int) {
	endt := fillEndtime(end)
	for _, act := range idx.Actions {
		if CheckTarget(act, &start, endt) {
			result += act.Point
		}
	}
	return result
}

func SummaryPoint(idx *Idx, start time.Time, end *time.Time) (summary *Summary) {
	endt := fillEndtime(end)

	stp := SumTargetPoint(idx, start, endt)
	sap := SumActualPoint(idx, start, endt)

	stc := 0
	for _, act := range idx.Actions {
		if CheckTarget(act, &start, endt) {
			stc++
		}
	}

	sac := 0
	for _, act := range idx.Actions {
		if CheckOK(act, &start, endt) {
			sac++
		}
	}

	summary = &Summary{
		Target:  stp,
		Actual:  sap,
		Ratio:   float32(sap) / float32(stp),
		TargetC: stc,
		ActualC: sac,
		RatioC:  float32(sac) / float32(stc),
	}
	return summary
}

func fillEndtime(end *time.Time) *time.Time {
	if end == nil {
		now := time.Now()
		return &now
	} else {
		return end
	}
}
