package bot

const (
	RichnessUnusable = iota
	RichnessLow
	RichnessMedium
	RichnessHigh
)

type Cell struct {
	Index     int8
	Richness  int8
	Neighbors []int8
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
