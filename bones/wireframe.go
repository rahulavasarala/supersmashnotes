package bones

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/rahulavasarala/supersmashnotes/collisions"
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
	Id        int     `yaml:"id"`
	Width     float64 `yaml:"width"`
	Thickness float64 `yaml:"thickness"`
	X         float64 `yaml:"x"`
	Y         float64 `yaml:"y"`
	Links     []Link  `yaml:"links"`
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
		bone := NewBone(config.Bones[i].Id, config.Bones[i].X, config.Bones[i].Y, config.Bones[i].Width, config.Bones[i].Thickness)
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

// This function writes the hurtbox points as a string
func (s *WireFrame) SampleHurtboxPoints(dpi int, frameMap map[int]*mat.Dense, frame int) [][]string {

	samplePoints := [][]string{}

	for key := range frameMap {

		samples := []string{}
		samples = append(samples, fmt.Sprintf("%v", frame))
		samples = append(samples, fmt.Sprintf("%v", key))

		bone := s.GetBone(key)
		frame := frameMap[key]

		if dpi == 1 {
			samples = append(samples, fmt.Sprintf("%v", frame.At(0, 2)), fmt.Sprintf("%v", frame.At(1, 2)))
			samplePoints = append(samplePoints, samples)
			continue
		}

		bw := bone.hurtbox.GetWidth()
		bh := bone.hurtbox.GetHeight()

		posx := frame.At(0, 2) - bw/2
		posy := frame.At(1, 2) - bh/2

		incr_x := bw / (float64(dpi) - 1)
		incr_y := bh / (float64(dpi) - 1)

		for i := 0; i < dpi; i++ {
			for j := 0; j < dpi; j++ {
				px := fmt.Sprintf("%v", posx+float64(i)*incr_x)
				py := fmt.Sprintf("%v", posy+float64(j)*incr_y)

				samples = append(samples, px, py)

			}
		}

		samplePoints = append(samplePoints, samples)

	}

	return samplePoints

}

func (s *WireFrame) SampleBoneHurtBoxPoints(dpi int, frameMap map[int]*mat.Dense, key int, frame_num int) [][]float64 {
	samplePoints := [][]float64{}

	bone := s.GetBone(key)

	frame := frameMap[key]

	if dpi == 1 {
		samplePoints = append(samplePoints, []float64{frame.At(0, 2), frame.At(1, 2)})
		return samplePoints
	}

	posx := frame.At(0, 2) - bone.hurtbox.GetWidth()/2
	posy := frame.At(1, 2) - bone.hurtbox.GetHeight()/2
	incr_x := bone.hurtbox.GetWidth() / (float64(dpi) - 1)
	incr_y := bone.hurtbox.GetHeight() / (float64(dpi) - 1)

	for i := 0; i < dpi; i++ {
		for j := 0; j < dpi; j++ {
			samplePoints = append(samplePoints, []float64{posx + float64(i)*incr_x, posy + float64(j)*incr_y})
		}
	}

	return samplePoints

}

func (s *WireFrame) GetBoneOrientationGlobal(frameMap map[int]*mat.Dense, key int) float64 {
	boneFrame := frameMap[key]

	y_orientation_i := boneFrame.At(0, 1)
	y_orientation_j := boneFrame.At(1, 1)

	orient := math.Atan2(y_orientation_j, y_orientation_i)

	return orient
}

func WriteHurtBoxPointsToCSV(samplePoints [][]string, destFile string) error {

	file, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Comma = ','

	writer.WriteAll(samplePoints)

	fmt.Printf("Created %s or appended samplePoints to %s successfully!", destFile, destFile)
	return nil
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

//This is where the area score is calculated
//I think I should update the calculate average area function to have a

func (s *WireFrame) CalculateAverageArea(anim Animation, groups [][]int, dpi int) float64 {

	area := 0.

	for frame_num := 0; frame_num < anim.length; frame_num++ {
		frame_area := s.FindAreaScore(anim, groups, dpi, frame_num)

		sum_area := 0.
		for _, a := range frame_area {
			sum_area += a
		}

		area += sum_area

	}

	return area / float64(anim.length)

}

// the utility of this function is to calculate the area score of a grouping at any frame
func (s *WireFrame) FindAreaScore(anim Animation, groups [][]int, dpi int, frame int) []float64 {
	area_list := []float64{}

	s.ApplyAnimation(anim, frame, 200, 200)

	frameMap := s.FindGlobalBoneFrames(s.GetBone(0))

	for _, group := range groups {

		samplePoints := [][]float64{}

		for _, boneId := range group {
			sP := s.SampleBoneHurtBoxPoints(dpi, frameMap, boneId, frame)
			samplePoints = append(samplePoints, sP...)
		}

		_, _, w, h := collisions.FindEncompassingRectangle(samplePoints)

		area_list = append(area_list, w*h)
	}

	return area_list
}

func (s *WireFrame) FindAreaScoreDebug(anim Animation, groups [][]int, dpi int, frame int) []float64 {
	area_list := []float64{}

	s.ApplyAnimation(anim, frame, 200, 200)

	frameMap := s.FindGlobalBoneFrames(s.GetBone(0))

	for i, group := range groups {
		samplePoints := [][]float64{}

		for _, boneId := range group {
			sP := s.SampleBoneHurtBoxPoints(dpi, frameMap, boneId, frame)
			samplePoints = append(samplePoints, sP...)
		}

		_, _, w, h := collisions.FindEncompassingRectangle(samplePoints)

		area_list = append(area_list, w*h)

		fmt.Printf(" group: %v, sample points: %v, area: %v, w,h : %v,%v", i, samplePoints, w*h, w, h)
	}

	return area_list
}

// return the center points for each of the 4 groups for a given frame
func (s *WireFrame) FindCenters(anim Animation, groups [][]int, dpi int, frame int) [][]float64 {
	centers := [][]float64{}
	s.ApplyAnimation(anim, frame, 200, 200)
	frameMap := s.FindGlobalBoneFrames(s.GetBone(0))

	for _, group := range groups {

		samplePoints := [][]float64{}

		for _, boneId := range group {
			sP := s.SampleBoneHurtBoxPoints(dpi, frameMap, boneId, frame)
			samplePoints = append(samplePoints, sP...)
		}

		//find the average point, and find the bone that is the farthest away from that point

		sum_x := 0.
		sum_y := 0.

		for _, point := range samplePoints {
			sum_x += point[0]
			sum_y += point[1]
		}

		sum_x /= float64(len(samplePoints))
		sum_y /= float64(len(samplePoints))

		centers = append(centers, []float64{sum_x, sum_y})

	}

	return centers

}

// The way that the outlier bone positions are encoded are through 3 values, pos(x,y) and orientation(-\pi to pi)
func (s *WireFrame) FindOutlierBonePos(anim Animation, groups [][]int, dpi int, frame int) ([][]float64, []float64, []int) {

	ob_pos_list := [][]float64{}
	ob_orient_list := []float64{}
	outlier_bones := []int{}

	s.ApplyAnimation(anim, frame, 200, 200)
	frameMap := s.FindGlobalBoneFrames(s.GetBone(0))

	centers := s.FindCenters(anim, groups, dpi, frame)

	for i, group := range groups {
		center := centers[i]

		outlier_bone := -1
		max_dist := -1.
		ob_pos := []float64{}
		ob_orient := -1.

		for _, boneId := range group {
			sP := s.SampleBoneHurtBoxPoints(2, frameMap, boneId, frame)

			dist := 0.

			for _, p := range sP {
				dist += (p[0]-center[0])*(p[0]-center[0]) + (p[1]-center[1])*(p[1]-center[1])
			}

			if dist > max_dist {
				outlier_bone = boneId
				ob_pos = s.SampleBoneHurtBoxPoints(1, frameMap, boneId, frame)[0]
				ob_orient = s.GetBoneOrientationGlobal(frameMap, outlier_bone)
			}
		}

		ob_pos_list = append(ob_pos_list, ob_pos)
		ob_orient_list = append(ob_orient_list, ob_orient)
		outlier_bones = append(outlier_bones, outlier_bone)

	}

	return ob_pos_list, ob_orient_list, outlier_bones

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
