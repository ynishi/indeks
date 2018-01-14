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
	Comment string
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