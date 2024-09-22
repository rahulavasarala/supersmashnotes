package bones

import (
	"fmt"
	"log"
	"math"
	"os"

	"gonum.org/v1/gonum/mat"
	"gopkg.in/yaml.v2"
)

type WireFrame struct {
	boneMap map[int]*Bone
}

type BoneConfig struct {
	NumBones int     `yaml:"numBones"`
	Bones    []BoneY `yaml:"bones"`
}

type BoneY struct {
	Id    int     `yaml:"id"`
	Width float64 `yaml:"width"`
	X     float64 `yaml:"x"`
	Y     float64 `yaml:"y"`
	Links []Link  `yaml:"links"`
}

type Link struct {
	Id    int     `yaml:"id"`
	Angle float64 `yaml:"angle"`
	Side  string  `yaml:"side"`
}

func (s *WireFrame) unfurl(boneConfig string) *BoneConfig {
	data, err := os.ReadFile(boneConfig)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Create an instance of Config
	var config BoneConfig

	// Unmarshal the YAML file into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling BoneConfig YAML: %v", err)
	}

	return &config
}

func (s *WireFrame) InitWireFrame(boneConfig string) error {

	config := s.unfurl(boneConfig)

	if config == nil {
		return fmt.Errorf("bone config file was not read properly")
	}

	boneMap := map[int]*Bone{}

	if config.NumBones == 0 {
		return fmt.Errorf("number of bones from config file was unspecified")
	}

	if len(config.Bones) != config.NumBones {
		return fmt.Errorf("amount of bones is not the same as amount specified")
	}

	for i := 0; i < config.NumBones; i++ {
		bone := NewBone(config.Bones[i].Id, config.Bones[i].X, config.Bones[i].Y, config.Bones[i].Width)
		boneMap[i] = bone
	}

	for i := 0; i < config.NumBones; i++ {
		for _, link := range config.Bones[i].Links {
			if link.Side == "left" {
				boneMap[i].leftAngles = append(boneMap[i].leftAngles, math.Pi*link.Angle)
				boneMap[i].lefts = append(boneMap[i].lefts, boneMap[link.Id])
				boneMap[link.Id].rights = append(boneMap[link.Id].rights, boneMap[i])
				boneMap[link.Id].rightAngles = append(boneMap[link.Id].rightAngles, math.Pi*link.Angle)
			} else {
				boneMap[i].rightAngles = append(boneMap[i].rightAngles, math.Pi*link.Angle)
				boneMap[i].rights = append(boneMap[i].rights, boneMap[link.Id])
				boneMap[link.Id].leftAngles = append(boneMap[link.Id].leftAngles, math.Pi*link.Angle)
				boneMap[link.Id].lefts = append(boneMap[link.Id].lefts, boneMap[i])
			}
		}
	}

	s.boneMap = boneMap

	return nil

}

//This will be a container for bones objects that are connected to eachother

// The algorithm will be, start at a bone, only go forward, figure oout the current frame, go to the next frame, and do a  matrix

// because the bones are double connected, we will need

func FindFrames(currBone *Bone, frame *mat.Dense, frameMap map[int]*mat.Dense, visited map[int]bool) {

	//do a forward propagation for all the lefts

	visited[currBone.id] = true

	for i := 0; i < len(currBone.lefts); i++ {
		if _, ok := visited[currBone.lefts[i].id]; !ok {

			//figure out the translation vector first

			x_translation := -currBone.width/2 - (currBone.lefts[i].width/2)*math.Cos(currBone.leftAngles[i])
			y_translation := (currBone.lefts[i].width / 2) * math.Sin(currBone.leftAngles[i])

			theta := -1 * currBone.leftAngles[i]

			fullRotationMatrix := []float64{
				math.Cos(theta), -1 * math.Sin(theta), x_translation,
				math.Sin(theta), math.Cos(theta), y_translation,
				0, 0, 1,
			}

			relative := mat.NewDense(3, 3, fullRotationMatrix)

			newFrame := new(mat.Dense)
			newFrame.Mul(frame, relative)

			frameMap[currBone.lefts[i].id] = newFrame

			visited[currBone.lefts[i].id] = true

			FindFrames(currBone.lefts[i], newFrame, frameMap, visited)

		}

	}

	for i := 0; i < len(currBone.rights); i++ {
		if _, ok := visited[currBone.rights[i].id]; !ok {

			//figure out the translation vector first

			x_translation := currBone.width/2 + (currBone.rights[i].width/2)*math.Cos(currBone.rightAngles[i])
			y_translation := (currBone.rights[i].width / 2) * math.Sin(currBone.rightAngles[i])

			theta := currBone.rightAngles[i]

			fullRotationMatrix := []float64{
				math.Cos(theta), -1 * math.Sin(theta), x_translation,
				math.Sin(theta), math.Cos(theta), y_translation,
				0, 0, 1,
			}

			relative := mat.NewDense(3, 3, fullRotationMatrix)

			newFrame := new(mat.Dense)
			newFrame.Mul(frame, relative)

			frameMap[currBone.rights[i].id] = newFrame

			visited[currBone.rights[i].id] = true

			FindFrames(currBone.rights[i], newFrame, frameMap, visited)

		}

	}

	//Now take the frames in the frame map and centralize them with the perspective of the map so that they can be drawn

}

func (s *WireFrame) FindGlobalBoneFrames(originBone *Bone) map[int]*mat.Dense {
	frameMap := map[int]*mat.Dense{}

	if originBone == nil {
		return frameMap
	}

	theta := originBone.orientation * math.Pi //value between 0 and 2 pi
	x_translation := originBone.x
	y_translation := originBone.y

	firstFrameRotationMatrix := []float64{
		math.Cos(theta), -1 * math.Sin(theta), x_translation,
		math.Sin(theta), math.Cos(theta), y_translation,
		0, 0, 1,
	}

	firstRotation := mat.NewDense(3, 3, firstFrameRotationMatrix)

	frame := mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	})

	frameMap[originBone.id] = frame

	visited := map[int]bool{}

	FindFrames(originBone, frame, frameMap, visited)

	for key := range frameMap {
		transformed := new(mat.Dense)
		transformed.Mul(firstRotation, frameMap[key])

		frameMap[key] = transformed
	}

	return frameMap

}

func (s *WireFrame) GetBone(id int) *Bone {
	val, ok := s.boneMap[id]

	if !ok {
		return nil
	}

	return val
}

func (s *WireFrame) SetOrientationOfBone(bone int, orientation float64) {
	val, ok := s.boneMap[bone]

	if !ok {
		return
	}

	val.orientation = orientation
}

func (s *WireFrame) ChangeAngleBetweenBones(bone1 int, bone2 int, newAngle float64) {
	b1, ok := s.boneMap[bone1]

	if !ok {
		return
	}

	b2, ok2 := s.boneMap[bone2]

	if !ok2 {
		return
	}

	//check if the bones are connected

	val, side := b1.GetLink(bone2)

	if val == nil {
		return
	}

	if side == "left" {
		b2.ChangeAngle(bone1, newAngle*math.Pi, "right")
		b1.ChangeAngle(bone2, newAngle*math.Pi, "left")
	} else if side == "right" {
		b2.ChangeAngle(bone1, newAngle*math.Pi, "left")
		b1.ChangeAngle(bone2, newAngle*math.Pi, "right")
	}

}

func (s *WireFrame) ApplyAnimation(anim Animation, frame int, x float64, y float64) {

	orientationInterface := anim.baseOrientation.Read(frame)

	orientation, ok := orientationInterface.(float64)

	if !ok {
		return
	}

	xOffsetInterface := anim.xOffset.Read(frame)

	xOffset, ok := xOffsetInterface.(float64)

	if !ok {
		return
	}

	yOffsetInterface := anim.yOffset.Read(frame)

	yOffset, ok := yOffsetInterface.(float64)

	if !ok {
		return
	}

	s.boneMap[0].orientation = orientation
	s.boneMap[0].x = x + xOffset
	s.boneMap[0].y = y + yOffset

	for _, jointProperty := range anim.jointPropertyList {
		b1 := jointProperty.j1
		b2 := jointProperty.j2

		thetaInterface := jointProperty.thetaProperty.Read(frame)

		theta, ok := thetaInterface.(float64)

		if !ok {
			return
		}

		s.ChangeAngleBetweenBones(b1, b2, theta)
	}

}
