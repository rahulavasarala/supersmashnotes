package bones

import (
	"log"
	"os"

	"github.com/rahulavasarala/supersmashnotes/properties"
	"gopkg.in/yaml.v2"
)

type Animation struct {
	xOffset         *properties.DoubleProperty
	yOffset         *properties.DoubleProperty
	baseOrientation *properties.DoubleProperty

	jointPropertyList []JointProperty
}

type JointProperty struct {
	j1            int
	j2            int
	thetaProperty *properties.DoubleProperty
}

type AnimationConfig struct {
	XOffset         []properties.RangeValue `yaml:"xoffset"`
	YOffset         []properties.RangeValue `yaml:"yoffset"`
	BaseOrientation []properties.RangeValue `yaml:"baseorientation"`
	JointConfigList []JointConfig           `yaml:"joints"`
}

func WriteAnimationConfigToYAML(filename string, data *AnimationConfig) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	return encoder.Encode(data)

}

type JointConfig struct {
	J1        int                     `yaml:"j1"`
	J2        int                     `yaml:"j2"`
	ThetaList []properties.RangeValue `yaml:"thetas"`
}

func (s *Animation) unfurl(animationConfig string) *AnimationConfig {
	data, err := os.ReadFile(animationConfig)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Create an instance of Config
	var config AnimationConfig

	// Unmarshal the YAML file into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling BoneConfig YAML: %v", err)
	}

	return &config
}

func (s *Animation) InitAnimation(animationConfig string) {
	config := s.unfurl(animationConfig)

	xoffset := properties.BuildDoubleProperty("xoffset", config.XOffset)
	yoffset := properties.BuildDoubleProperty("yoffset", config.YOffset)
	orientation := properties.BuildDoubleProperty("baseorientation", config.BaseOrientation)

	s.xOffset = xoffset
	s.yOffset = yoffset
	s.baseOrientation = orientation

	s.jointPropertyList = []JointProperty{}

	for _, jointConfig := range config.JointConfigList {

		jP := JointProperty{}
		jP.j1 = jointConfig.J1
		jP.j2 = jointConfig.J2
		jP.thetaProperty = properties.BuildDoubleProperty("theta", jointConfig.ThetaList)

		s.jointPropertyList = append(s.jointPropertyList, jP)
	}

}

//animation is a all encompassing data structure that allows you to read animations joint values
