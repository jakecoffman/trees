package bot

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
	Type          int8
	TargetCellIdx int8

	// SourceCellIdx only used for SEED right now
	SourceCellIdx int8
}

func NewAction(action string) Action {
	parts := strings.Split(action, " ")
	switch parts[0] {
	case "WAIT":
		return Action{Type: Wait}
	case "SEED":
		source, err := strconv.Atoi(parts[1])
		check(err)
		target, err := strconv.Atoi(parts[2])
		check(err)
		return Action{Type: Seed, TargetCellIdx: int8(target), SourceCellIdx: int8(source)}
	case "GROW":
		target, err := strconv.Atoi(parts[1])
		check(err)
		return Action{Type: Grow, TargetCellIdx: int8(target)}
	case "COMPLETE":
		target, err := strconv.Atoi(parts[1])
		check(err)
		return Action{Type: Complete, TargetCellIdx: int8(target)}
	}
	panic("Invalid action: " + parts[0])
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
	return "ASDF_QWERT"
}
