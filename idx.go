package indeks

import "time"

const (
	ResultUncheck = iota
	ResultOK
	ResultNG
	ResultChanged
)

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

func SumActualPoint(idx *Idx, start time.Time, end *time.Time) (result int) {
	var endt *time.Time
	if end == nil {
		now := time.Now()
		endt = &now
	} else {
		endt = end
	}
	for _, act := range idx.Actions {
		if act.Result == ResultOK && act.ActualTime.After(start) && act.ActualTime.Before(*endt) {
			result += act.Point
		}
	}
	return result
}

func SumTargetPoint(idx *Idx, start time.Time, end *time.Time) (result int) {
	var endt *time.Time
	if end == nil {
		now := time.Now()
		endt = &now
	} else {
		endt = end
	}
	for _, act := range idx.Actions {
		if act.Result != ResultChanged && act.TargetTime.After(start) && act.TargetTime.Before(*endt) {
			result += act.Point
		}
	}
	return result
}
