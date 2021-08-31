package game

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Wait = iota
	Seed
	Grow
	Complete
)

type Action struct {
	Type          int
	TargetCellIdx int

	// SourceCellIdx only used for SEED right now
	SourceCellIdx int
}

var (
	ErrInvalidActionString = fmt.Errorf("invalid action string")
)

// NewAction parses an action string and returns an object
func NewAction(action string) (Action, error) {
	parts := strings.Split(action, " ")
	switch parts[0] {
	case "WAIT":
		return Action{Type: Wait}, nil
	case "SEED":
		source, err := strconv.Atoi(parts[1])
		if err != nil {
			return Action{}, ErrInvalidActionString
		}
		target, err := strconv.Atoi(parts[2])
		if err != nil {
			return Action{}, ErrInvalidActionString
		}
		return Action{Type: Seed, TargetCellIdx: target, SourceCellIdx: source}, nil
	case "GROW":
		target, err := strconv.Atoi(parts[1])
		if err != nil {
			return Action{}, ErrInvalidActionString
		}
		return Action{Type: Grow, TargetCellIdx: target}, nil
	case "COMPLETE":
		target, err := strconv.Atoi(parts[1])
		if err != nil {
			return Action{}, ErrInvalidActionString
		}
		return Action{Type: Complete, TargetCellIdx: target}, nil
	}
	return Action{}, ErrInvalidActionString
}

func (a Action) String() string {
	if a.Type == Wait {
		return "WAIT"
	}
	if a.Type == Seed {
		return fmt.Sprintf("SEED %d %d", a.SourceCellIdx, a.TargetCellIdx)
	}
	if a.Type == Complete {
		return fmt.Sprintf("COMPLETE %d", a.TargetCellIdx)
	}
	if a.Type == Grow {
		return fmt.Sprintf("GROW %d", a.TargetCellIdx)
	}
	// this shouldn't be possible
	return "ACTION_INVALID"
}
