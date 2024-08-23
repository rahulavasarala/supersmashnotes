package collisions

type CollisionHandler interface {
	HandleCollisions(sg SpatialGrid, collisionMap map[string]map[string]bool) error
}

type SimpleCollisionHandler struct {
}

// Implement a logic by sunday where collisions are properly handled, addition and summation of collisions are handled
// conservation of momentum
func (s *SimpleCollisionHandler) HandleCollisions(collisionMap map[string][]Thing) {

}
