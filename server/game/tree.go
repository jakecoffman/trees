package game

const (
	SizeSeed = iota
	SizeSmall
	SizeMedium
	SizeLarge
)

type Tree struct {
	CellIndex int
	Size      int
	Player    int
	IsDormant bool
}
