package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rahulavasarala/supersmashnotes/bones"
	"github.com/rahulavasarala/supersmashnotes/graphics"
	"github.com/rahulavasarala/supersmashnotes/properties"
	"gopkg.in/yaml.v2"
)

var BASEX float64 = 200
var BASEY float64 = 200

type Game struct {
	debugMode        bool
	wireframe        *bones.WireFrame
	jointAngleConfig string
	writeDestination string
	frameList        []int
	jointMap         map[string][]float64
	baseBoneMap      map[string][]float64
	maxSaves         int
}

type JointAngles struct {
	J1    int     `yaml:"j1"`
	J2    int     `yaml:"j2"`
	Theta float64 `yaml:"theta"`
}

type JointAngleConfig struct {
	Frame          int           `yaml:"frame"`
	Save           bool          `yaml:"save"`
	XOffset        float64       `yaml:"xoffset"`
	YOffset        float64       `yaml:"yoffset"`
	BaseAngle      float64       `yaml:"baseangle"`
	JointAngleList []JointAngles `yaml:"jointangles"`
}

func (s *Game) InitAnimator(jointAngleConfig string, maxSaves int, debugMode bool, boneConfig string, animationConfigDestination string) {
	s.jointAngleConfig = jointAngleConfig

	wireFrame := bones.WireFrame{}
	wireFrame.InitWireFrame(boneConfig)
	s.wireframe = &wireFrame

	s.debugMode = debugMode
	s.frameList = []int{}
	s.jointMap = map[string][]float64{}
	s.baseBoneMap = map[string][]float64{}
	s.baseBoneMap["xoffset"] = []float64{}
	s.baseBoneMap["yoffset"] = []float64{}
	s.baseBoneMap["orientation"] = []float64{}

	s.maxSaves = maxSaves
	s.writeDestination = animationConfigDestination

}

// func (g *Game) Update() error {
// 	g.iteration++

// 	if g.iteration > 100 {
// 		g.iteration = 0
// 	}
// 	g.wireframe.SetOrientationOfBone(0, float64(g.iteration)/50)
// 	g.wireframe.ChangeAngleBetweenBones(0, 1, float64(g.iteration)/50)
// 	g.wireframe.ChangeAngleBetweenBones(1, 2, float64(g.iteration)/50)

// 	return nil
// }

func (g *Game) pollJointAngleConfig() *JointAngleConfig {
	data, err := os.ReadFile(g.jointAngleConfig)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Create an instance of Config
	var config JointAngleConfig

	// Unmarshal the YAML file into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling BoneConfig YAML: %v", err)
	}

	return &config

}

func (g *Game) Update() error {
	config := g.pollJointAngleConfig()

	g.wireframe.SetOrientationOfBone(0, config.BaseAngle)
	b0 := g.wireframe.GetBone(0)

	if b0 == nil {
		return nil
	}

	b0.SetPosition(config.XOffset+BASEX, config.YOffset+BASEY)

	for _, jointAngle := range config.JointAngleList {
		g.wireframe.ChangeAngleBetweenBones(jointAngle.J1, jointAngle.J2, jointAngle.Theta)
	}

	if config.Save {
		if len(g.frameList) != 0 && g.frameList[len(g.frameList)-1] > config.Frame {
			return fmt.Errorf("next frame %v is not greater than the previous frame %v", config.Frame, g.frameList[len(g.frameList)-1])
		} else if len(g.frameList) != 0 && g.frameList[len(g.frameList)-1] == config.Frame {
			return nil
		}

		g.frameList = append(g.frameList, config.Frame)
		fmt.Printf("Saved frame %v!\n", config.Frame)

		for _, jointAngle := range config.JointAngleList {
			jointKey := fmt.Sprintf("%v-%v", jointAngle.J1, jointAngle.J2)
			if _, ok := g.jointMap[jointKey]; !ok {
				g.jointMap[jointKey] = []float64{}
			}

			g.jointMap[jointKey] = append(g.jointMap[jointKey], jointAngle.Theta)
		}

		g.baseBoneMap["orientation"] = append(g.baseBoneMap["orientation"], config.BaseAngle)
		g.baseBoneMap["xoffset"] = append(g.baseBoneMap["xoffset"], config.XOffset)
		g.baseBoneMap["yoffset"] = append(g.baseBoneMap["yoffset"], config.YOffset)

		g.WriteAnimationConfigFile()

	}

	return nil
}

func (g *Game) WriteAnimationConfigFile() {

	//first let us generate a time line, this same time line will be applicable to all the joints

	config := bones.AnimationConfig{}
	config.JointConfigList = []bones.JointConfig{}

	//First step is to generate the range list
	rangeList := []string{}

	for i := 0; i < (len(g.frameList) - 1); i++ {
		rangeString := fmt.Sprintf("%v-%v", g.frameList[i], g.frameList[i+1])

		rangeList = append(rangeList, rangeString)
	}

	for key := range g.jointMap {
		thetaList := g.jointMap[key]

		boneNums := properties.ParseIntRange(key)

		jC := bones.JointConfig{}

		jC.J1 = boneNums.First
		jC.J2 = boneNums.Second

		values := []string{}

		for i := 0; i < len(thetaList)-1; i++ {
			values = append(values, fmt.Sprintf("%v-%v", thetaList[i], thetaList[i+1]))
		}

		rangeVals := properties.BuildRangeVals(rangeList, values)

		jC.ThetaList = rangeVals

		config.JointConfigList = append(config.JointConfigList, jC)
	}

	//let us fill the animationConfig with rangeVals forxoffset

	xOffsetVals := []string{}

	for i := 0; i < len(g.baseBoneMap["xoffset"])-1; i++ {
		xOffsetVals = append(xOffsetVals, fmt.Sprintf("%v-%v", g.baseBoneMap["xoffset"][i], g.baseBoneMap["xoffset"][i+1]))
	}

	fmt.Printf("%v", g.baseBoneMap["xoffset"])

	XORangeVals := properties.BuildRangeVals(rangeList, xOffsetVals)

	config.XOffset = XORangeVals

	//let us fill the animationConfig with rangeVals for yoffset

	yOffsetVals := []string{}

	for i := 0; i < len(g.baseBoneMap["yoffset"])-1; i++ {
		yOffsetVals = append(yOffsetVals, fmt.Sprintf("%v-%v", g.baseBoneMap["yoffset"][i], g.baseBoneMap["yoffset"][i+1]))
	}

	YORangeVals := properties.BuildRangeVals(rangeList, yOffsetVals)

	config.YOffset = YORangeVals

	//let us fill the animationConfig with rangeVals forxoffset

	bOVals := []string{}

	for i := 0; i < len(g.baseBoneMap["orientation"])-1; i++ {
		bOVals = append(bOVals, fmt.Sprintf("%v-%v", g.baseBoneMap["orientation"][i], g.baseBoneMap["orientation"][i+1]))
	}

	bORangeVals := properties.BuildRangeVals(rangeList, bOVals)

	config.BaseOrientation = bORangeVals

	bones.WriteAnimationConfigToYAML(g.writeDestination, &config)

}

// func (g *Game) Draw(screen *ebiten.Image) {
// 	screen.Fill(color.White)

// 	frameMap := g.wireframe.FindGlobalBoneFrames(g.wireframe.GetBone(0))

// 	for key := range frameMap {
// 		bone := g.wireframe.GetBone(key)

// 		graphics.DrawBone(screen, color.Black, 2, 500, false, true, bone, frameMap[key])
// 	}

// }

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	frameMap := g.wireframe.FindGlobalBoneFrames(g.wireframe.GetBone(0))

	for key := range frameMap {
		bone := g.wireframe.GetBone(key)

		graphics.DrawBone(screen, color.Black, 2, 500, false, true, bone, frameMap[key])
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
}

func main() {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Super Smash Notes")

	game := Game{}

	game.InitAnimator("./thetaconfig.yaml", 10, true, "../bones/boneconfig1.yaml", "./testanimationconfig.yaml")

	ebiten.SetTPS(1)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

//I will make a simple animation rig that reads from a config file, constantly
//polling it for changes in angles, then I will render the object

//this will happen for a while

//change the value of one thing to save the angles in a timeline

//you have something called an animation map
//it will be like the state map, each animation will have certain bone positions

//position of the character will be the torso mid point

//First what is going to happen is, basically, you have a frameList, which has the saves of all the frames

//then you have a map that is keyed by a string(bone_bone), which corresponds to a list of values per save frame

//then you basically create a new datastructure that holds bone, bone, double property,  for each joint

//then you populate this datastructure, adn then you can implement a read method for the wire frame, and then you are don
