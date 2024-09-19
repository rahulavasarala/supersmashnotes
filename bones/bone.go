package bones

//hurtboxes and hitboxes will be included in the bones package
//bones will support task control for animations and integrate the task control into the statemachine config files through timelines

//individual bones will be modelled with centerpoints and angles, and width

//There will be an algorithm that is similar to dfs where you iterate trhough all the joints, figure out the positions of the joints
// in terms of the start joint, given their connection angles
// and draw all of the joints

//I need to define the way to denote a connection angle, ok, connection angles will have |theta| < pi, so negative pi to positive pi
//actuall angles will have values from 0 - 2pi

//in a wireframe , joints will be hashed by their id(simple number), joint 1,2,3.... and skeletons will support the draw method
//cool

//figure out whether we will need master and slave bones, lol

//I will make a standard that bending outwards is -Pi and bending inwards is pi, this way the left and right joints will be synchronized
//and you can have multi directional traversing

//the problem will turn into representing the centerpoints of the frame in terms of the frames of other centerpoints,
//and then drawing all of the joints

//bones can be mapped based on a config file

type Bone struct {
	id    int
	xpos  float64
	ypos  float64
	width float64

	angle float64 //this is the angle in radians, from 0-2pi that will determine the orientation

	lefts      []*Bone
	leftAngles []float64

	rights      []*Bone //if the left and right bones are equal to nil, then ignore the angles
	rightAngles []float64
}
