package bot

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"reflect"
)

type Settings struct {
	//InDev          bool
	MinimalEconomy int
	MaximumEconomy int
	InDev          bool
}

var settings = Settings{
	MinimalEconomy: 16,
	MaximumEconomy: 25,
}

func (s Settings) String() string {
	return fmt.Sprintf("%#v", s)
}

func (s Settings) Randomize() Settings {

	return s
}

func (s Settings) Tweak() Settings {
	value := reflect.ValueOf(&s)
	i := rand.Intn(value.Elem().NumField())
	field := value.Elem().Field(i)
	switch v := field.Interface().(type) {
	case float64:
		field.SetFloat(field.Float() + rand.Float64()*20 - 10)
	case int8, int, int64:
		field.SetInt(field.Int() + int64(rand.Intn(20)-5))
	default:
		panic(v)
	}
	return s
}

// Priority allows us to compare states to see which is better.
func (g *State) Priority(s *Settings) int {
	if g == nil {
		return -1
	}
	if g.Day >= 24-chokudaiMaxTurns {
		end := g.MyScore + int(math.Floor(float64(g.MySun)))/3
		// tested: remove pts for wasted energy
		end -= bits.OnesCount64(g.MyTrees) * 10
		return end
	}

	var score int

	score += bits.OnesCount64(g.MyTrees&g.Seeds) * 10
	score += bits.OnesCount64(g.MyTrees&g.Smalls) * 20
	score += bits.OnesCount64(g.MyTrees&g.Mediums) * 30
	score += bits.OnesCount64(g.MyTrees&g.Larges) * 45
	//score /= int(g.Day)

	//score += int(g.MySun)

	// bonus points for never being shaded
	//score += bits.OnesCount64(g.MyTrees & g.Larges & ^g.MaxShadows) * 5
	//score += bits.OnesCount64(g.MyTrees & g.Mediums & ^g.MaxShadows) * 3
	//score += bits.OnesCount64(g.MyTrees & g.Smalls & ^g.MaxShadows)

	// bonus points for shading opponent
	score += bits.OnesCount64(g.EnemyTrees&g.Larges&g.Shadows[0][SizeLarge]) * 5

	const (
		maxLargeTrees  = 4
		maxMediumTrees = 4
		maxSmallTrees  = 2
	)
	if g.Num[SizeLarge] > maxLargeTrees {
		score -= (int(g.Num[SizeLarge]) - maxLargeTrees) * 10
	}
	if g.Num[SizeMedium] > maxMediumTrees {
		score -= (int(g.Num[SizeMedium]) - maxMediumTrees) * 10
	}
	if g.Num[SizeSmall] > maxSmallTrees {
		score -= (int(g.Num[SizeSmall]) - maxSmallTrees) * 5
	}

	score += g.MyScore * 100

	if g.Num[SizeSeed] == 0 {
		score -= 100
	}

	return score / int(g.Day)
}

func (g *State) SunNextRound() (earning int) {
	for i := 0; i < 37; i++ {
		t := g.Trees[i]
		if t.Exists && t.IsMine && !IsSet(g.Shadows[1][t.Size], t.CellIndex) {
			earning += int(t.Size) * 2 // multiplier?
		}
	}
	return
}

func (g *State) Economy() int {
	return int(g.Num[SizeLarge]*3 + g.Num[SizeMedium]*2 + g.Num[SizeSmall])
}
