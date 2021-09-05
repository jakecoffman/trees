package bot

import (
	"fmt"
)

const (
	SizeSeed = iota
	SizeSmall
	SizeMedium
	SizeLarge
)

type Tree struct {
	CellIndex int8
	Size      int8
	IsMine    bool
	IsDormant bool
	Exists    bool
}

func (t *Tree) String() string {
	//return fmt.Sprintf("%#v", *t)
	if !t.Exists {
		return "No Tree!"
	}
	if t.IsMine {
		if t.IsDormant {
			return fmt.Sprintf("My:%v@%v (zzz)", Size(t.Size), t.CellIndex)
		}
		return fmt.Sprintf("My:%v@%v", Size(t.Size), t.CellIndex)
	} else {
		if t.IsDormant {
			return fmt.Sprintf("Op:%v@%v (zzz)", Size(t.Size), t.CellIndex)
		}
		return fmt.Sprintf("Op:%v@%v", Size(t.Size), t.CellIndex)
	}
}

func Size(size int8) string {
	if size == SizeLarge {
		return "L"
	} else if size == SizeMedium {
		return "M"
	} else if size == SizeSmall {
		return "S"
	} else {
		return "s"
	}
}
