const treeBaseCost = [0, 1, 3, 7]
const CostLifecycleEnd = 4
const RichnessBonusMedium = 2
const RichnessBonusHigh = 4
const SizeSeed = 0
const SizeSmall = 1
const SizeMedium = 2
const SizeLarge = 3

// growthCost returns the cost of growing a tree, including the cost of selling it
export function growthCost(player, state, tree) {
  const targetSize = tree.Size + 1
  if (targetSize > SizeLarge) {
    return CostLifecycleEnd
  } else {
    return cost(player, targetSize, state)
  }
}

// seedCost returns the cost of sowing a seed
export function seedCost(player, state) {
  return cost(player, 0, state)
}

function cost(player, size, state) {
  let c = treeBaseCost[size]
  for (const index in state.Trees) {
    if (state.Trees[index].Size === size && player === state.Trees[index].Owner) {
      c++
    }
  }
  return c
}
