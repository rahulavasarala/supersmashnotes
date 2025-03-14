package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rahulavasarala/supersmashnotes/bones"
)

//this will be the animation end point where a python script can call api's to essentially get the information from the
//environment given a certain amount of groupings

var AnimationMap map[string]*bones.Animation
var WireFrame bones.WireFrame
var BoneConfig = "../bones/boneconfig1.yaml"
var AnimationConfigs = []string{"../animationtool/jab.yaml", "../animationtool/knee.yaml", "../animationtool/flipkick.yaml", "../animationtool/bounce.yaml"}
var AnimationNames = []string{"jab", "knee", "flipkick", "bounce"}

type AnimationRequestData struct {
	Groups     [][]int `json:"groups"`
	Name       string  `json:"name"`
	PollFrames []int   `json:"pollframe"`
}

type AnimationResponseData struct {
	DataList []FrameData `json:"datalist"`
}

type FrameData struct {
	OutlierBones   []int       `json:"outlierbones"`
	OutlierPos     [][]float64 `json:"outlierpos"`
	OutlierOrients []float64   `json:"outlierorients"`
	CenterPos      [][]float64 `json:"centerpos"`
	AreaScore      []float64   `json:"areascore"`
}

func SendAnimationData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var animRequestData AnimationRequestData

	if err := json.NewDecoder(r.Body).Decode(&animRequestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	animationData := retrieve_animation_data(animRequestData.Groups, animRequestData.Name, animRequestData.PollFrames)

	if err := json.NewEncoder(w).Encode(animationData); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

}

func main() {

	Init()
	groups := [][]int{{1, 2, 3, 4}, {5, 6, 7}, {0, 8, 9}}
	anim := "jab"
	poll_frames := []int{1, 2, 3}

	fmt.Println(retrieve_animation_data(groups, anim, poll_frames))

	http.HandleFunc("/sendanimdata", SendAnimationData)

	// Start the HTTP server on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

}

// poll will be given as a list of frames in which data should be polled from
func retrieve_animation_data(groups [][]int, anim_name string, poll_frames []int) AnimationResponseData {

	animRespData := AnimationResponseData{}

	animRespData.DataList = []FrameData{}

	anim := AnimationMap[anim_name]

	for _, frame := range poll_frames {

		frameData := FrameData{}

		ob_pos_list, ob_orient_list, outlier_bones := WireFrame.FindOutlierBonePos(*anim, groups, 2, frame)
		centers := WireFrame.FindCenters(*anim, groups, 2, frame)

		area_score := WireFrame.FindAreaScore(*anim, groups, 2, frame)

		frameData.AreaScore = area_score
		frameData.CenterPos = centers
		frameData.OutlierPos = ob_pos_list
		frameData.OutlierOrients = ob_orient_list
		frameData.OutlierBones = outlier_bones

		animRespData.DataList = append(animRespData.DataList, frameData)
	}

	return animRespData

}

func Init() {

	WireFrame = bones.WireFrame{}
	WireFrame.InitWireFrame(BoneConfig)
	AnimationMap = map[string]*bones.Animation{}

	for i, acfg := range AnimationConfigs {
		anim_name := AnimationNames[i]

		anim := bones.Animation{}
		anim.InitAnimation(acfg)

		AnimationMap[anim_name] = &anim
	}

}
