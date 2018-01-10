package indeks

import "time"

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
