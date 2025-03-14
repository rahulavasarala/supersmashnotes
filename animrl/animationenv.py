import requests
import random
import torch
import math

go_endpoint = "http://localhost:8080/sendanimdata"
ACTION_SPACE = 7

class AnimationEnvironment:
    def __init__(self, anim, init_groupings, anim_len, total_frames):
        self.curr_frame = 0
        self.groupings = init_groupings
        self.start_groupings = init_groupings
        self.anim_name = anim
        self.anim_len = anim_len
        self.total_frames = total_frames  
        self.action_space = ACTION_SPACE 
        self.poll_frames = [0, 2, 5]
        self.light_action_penalty = False

        device = torch.device(
            "cuda" if torch.cuda.is_available() else
            "mps" if torch.backends.mps.is_available() else
            "cpu"
        )

        self.device = device

        _, self.area_score, self.outlierbones = self.get_state()

    def SpecifyAnimation(self, anim, total_frames):
        self.anim_name = anim
        self.total_frames = total_frames

    def get_state(self):#data is now normalized

        poll_frames = [(self.curr_frame + offset)%self.total_frames for offset in self.poll_frames]

        response = GetAnimationData(self.groupings , self.anim_name, poll_frames)
        data_dict = response.json()

        area_score = sum(data_dict["datalist"][0]["areascore"])
        outlier_bones = data_dict["datalist"][0]["outlierbones"]

        pos = torch.tensor([])
        areas = torch.tensor([])
        orients = torch.tensor([])

        for data in data_dict["datalist"]:

            pos = torch.cat((pos, torch.tensor(data["centerpos"]).flatten()))
            pos = torch.cat((pos, torch.tensor(data["outlierpos"]).flatten()))
            areas = torch.cat((areas, torch.tensor(data["areascore"]).flatten()))
            orients = torch.cat((orients, torch.tensor(data["outlierorients"]).flatten()))

        next_state = torch.tensor([])
        norm_pos = (pos - 200)/10
        norm_areas = areas/1000
        norm_orients = orients/math.pi

        next_state = torch.cat((next_state, norm_pos))
        next_state = torch.cat((next_state, norm_areas))
        next_state = torch.cat((next_state, norm_orients))
        next_state = torch.cat((next_state, torch.tensor([self.curr_frame/self.total_frames]).flatten()))

        return next_state.flatten().unsqueeze(0).to(self.device), area_score/1000, outlier_bones


    def step(self, action):

        self.curr_frame += 1
        self.switch_groupings(action)
        terminated = False
        if self.curr_frame >= self.total_frames:
            terminated = True

        next_state, next_area_score, outlierbones = self.get_state()

        immediate_reward = 0
        self.area_score = next_area_score
        self.outlierbones = outlierbones
        area_reward, action_penalty = self.calc_reward(action)

        immediate_reward = area_reward + action_penalty

        return next_state, immediate_reward , terminated
    
    def calc_reward(self, action):

        area_reward = 0
        if self.area_score <= 2:
            area_reward = 1
        elif self.area_score > 2 and self.area_score < 7:
            area_reward = (-11 * self.area_score)/5 + 27/5
        else:
            area_reward = -12

        action_penalty = 0
        PENALTY_FACTOR = 1.2
        PERIOD_THRESHOLD = 2.8
        period = self.curr_frame/(self.anim_len)

        if self.light_action_penalty:
            if period > PERIOD_THRESHOLD and action != 0:
                action_penalty = -2
        else:
            if action != 0:
                x = self.curr_frame//(self.anim_len)
                action_penalty = -PENALTY_FACTOR*x

        print("Area reward: ", area_reward, "Action pen: ", action_penalty)

        return area_reward, action_penalty

    def switch_groupings(self, action):

        g1 = -1
        g2 = -1

        if action == 1:
            g1 = 0
            g2 = 1
        elif action == 2:
            g1 = 0
            g2 = 2
        elif action == 3:
            g1 = 0
            g2 = 3
        elif action == 4:
            g1 = 1
            g2 = 2
        elif action == 5:
            g1 = 1
            g2 = 3
        elif action == 6:
            g1 = 2
            g2 = 3
        else:
            return
        
        b1 = self.outlierbones[g1]
        b2 = self.outlierbones[g2]

        i1 = -1

        for i in range(len(self.groupings[g1])):
            if self.groupings[g1][i] == b1:
                i1 = i
                break

        i2 = -1

        for i in range(len(self.groupings[g2])):
            if self.groupings[g2][i] == b2:
                i2 = i
                break

        #Then just swap the two positions
        temp = self.groupings[g1][i1]
        self.groupings[g1][i1] = self.groupings[g2][i2]
        self.groupings[g2][i2] = temp

    def reset(self):
        self.curr_frame = 0
        self.groupings = self.start_groupings

    def sampleActionSpace(self):
        return random.randint(0,ACTION_SPACE-1)


def GetAnimationData(groups: list[list[int]], anim: str, poll_frames: list[int]):
    payload = {"groups" : groups, "name": anim, "pollframe" : poll_frames}
    response = requests.post(go_endpoint, json=payload)

    if response.status_code != 200:
        print("NO/Invalid data received")
        return None

    return response

def main():

    #This will just be a test to see if the get animation data function is working
    start_groupings = [[0,1,2,3],[4,5],[6,7],[8,9]]
    anim = "jab"
    pollframe = [0,15]

    response = GetAnimationData(start_groupings, anim, pollframe)

    print(response.json())

    env = AnimationEnvironment(anim, start_groupings, 26, 100)

    state, _, _ = env.get_state()

    print("outlier bones: ", env.outlierbones)
    print("area score", env.area_score)

    # print(state)
    # print(state.shape)

    _, reward, _ = env.step(5)

    print("groupings: ", env.groupings)
    print("reward: ", reward)
    print("new area score: ", env.area_score)

    #curse of dimensionality, need to figure out how to make the data more thin
    return

if __name__ == '__main__':
    main()