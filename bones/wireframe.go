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
		bone := Bone{id: config.Bones[i].Id, xpos: config.Bones[i].X, ypos: config.Bones[i].Y, width: config.Bones[i].Width}

		boneMap[i] = &bone
	}

	for i := 0; i < config.NumBones; i++ {
		for _, link := range config.Bones[i].Links {
			if link.Side == "left" {
				boneMap[i].leftAngles = append(boneMap[i].leftAngles, link.Angle)
				boneMap[i].lefts = append(boneMap[i].lefts, boneMap[link.Id])
			} else {
				boneMap[i].rightAngles = append(boneMap[i].rightAngles, link.Angle)
				boneMap[i].rights = append(boneMap[i].rights, boneMap[link.Id])
			}
		}
	}

}

//This will be a container for bones objects that are connected to eachother

// The algorithm will be, start at a bone, only go forward, figure oout the current frame, go to the next frame, and do a  matrix

// because the bones are double connected, we will need a visited map
func (s *WireFrame) Draw() {

}

func findFrames(currBone *Bone, frame *mat.Dense, frameMap map[int]*mat.Dense, visited map[int]bool) {

	//do a forward propagation for all the lefts

	visited[currBone.id] = true

	for i := 0; i < len(currBone.lefts); i++ {
		if _, ok := visited[currBone.lefts[i].id]; !ok {

			//figure out the translation vector first

			x_translation := -currBone.width/2 - (currBone.lefts[i].width/2)*math.Cos(currBone.leftAngles[i])
			y_translation := math.Sin(currBone.leftAngles[i])

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

			findFrames(currBone.lefts[i], newFrame, frameMap, visited)

		}

	}

	for i := 0; i < len(currBone.rights); i++ {
		if _, ok := visited[currBone.rights[i].id]; !ok {

			//figure out the translation vector first

			x_translation := currBone.width/2 + (currBone.rights[i].width/2)*math.Cos(currBone.leftAngles[i])
			y_translation := math.Sin(currBone.leftAngles[i])

			theta := currBone.leftAngles[i]

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

			findFrames(currBone.rights[i], newFrame, frameMap, visited)

		}

	}

}
