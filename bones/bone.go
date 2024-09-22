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
	x     float64
	y     float64
	width float64

	orientation float64

	lefts      []*Bone
	leftAngles []float64

	rights      []*Bone //if the left and right bones are equal to nil, then ignore the angles
	rightAngles []float64
}

func (s *Bone) InitBone(id int, x float64, y float64, width float64) {
	s.id = id
	s.x = x
	s.y = y
	s.width = width
	s.lefts = []*Bone{}
	s.leftAngles = []float64{}
	s.rights = []*Bone{}
	s.rightAngles = []float64{}
}

func NewBone(id int, x float64, y float64, width float64) *Bone {
	newBone := Bone{}
	newBone.InitBone(id, x, y, width)

	return &newBone
}

func (s *Bone) GetWidth() float64 {
	return s.width
}

func (s *Bone) GetId() int {
	return s.id
}

func (s *Bone) SetOrientation(or float64) {
	s.orientation = or
}

// assumption is that the same bone cannot be connected to both left and right of the same bone
func (s *Bone) GetLink(id int) (*Bone, string) {
	if s.rights != nil {
		for i := 0; i < len(s.rights); i++ {
			if s.rights[i].id == id {
				return s.rights[i], "right"
			}
		}
	}

	if s.lefts != nil {
		for i := 0; i < len(s.lefts); i++ {
			if s.lefts[i].id == id {
				return s.lefts[i], "left"
			}
		}
	}

	return nil, ""

}

func (s *Bone) SetPosition(x float64, y float64) {
	s.x = x
	s.y = y
}

func (s *Bone) ChangeAngle(bone2 int, newAngle float64, side string) {

	if side == "left" {
		for i := 0; i < len(s.leftAngles); i++ {
			if s.lefts[i].id == bone2 {
				s.leftAngles[i] = newAngle
				break
			}
		}
	} else if side == "right" {
		for i := 0; i < len(s.rightAngles); i++ {
			if s.rights[i].id == bone2 {
				s.rightAngles[i] = newAngle
				break
			}
		}
	}
}
