package game

const (
	RichnessUnusable = iota
	RichnessLow
	RichnessMedium
	RichnessHigh
)

type Cell struct {
	Index     int
	Richness  int
	Neighbors []int
}
