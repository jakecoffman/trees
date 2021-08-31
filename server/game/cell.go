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

func check(err error) {
	if err != nil {
		panic(err)
	}
}
