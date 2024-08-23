package collisions

type CollisionFinder interface {
	FindPotentialCollisions(sg SpatialGrid) map[string]map[string]bool
}

//the idea is you have an id to id map, which models the potential collisions that are going on.

//then you can use the thingMap in the spatial grid, for quick look ups
// I think the long term decision of getting rid of the extra characters list and immovables list, and just looking things up by id would
//work well

//This is pretty true
//certain questions that may be asked are like, why do you need an immovables list
