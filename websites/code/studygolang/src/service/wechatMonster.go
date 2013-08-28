package service

import (
	//"logger"
	"math/rand"
	"model"
)

var Map_MonsterData map[int]*model.WechatMonster

//id int, name string, hp int, attack int, defense int, agility int, resistance int, exp int, props []int

func init() {

	Map_MonsterData = make(map[int]*model.WechatMonster)

	Map_MonsterData[0] = model.NewWechatMonster(0, "风铃怪", 15, 2, 2, 4, 1, 1, []int{1, 2}, []int{30, 30})
	Map_MonsterData[1] = model.NewWechatMonster(1, "泥巴怪", 15, 3, 4, 1, 4, 2, []int{1, 3}, []int{20, 40})
	Map_MonsterData[2] = model.NewWechatMonster(2, "浪鱼妖", 50, 10, 5, 5, 4, 10, []int{1}, []int{50})
}

//一个玩家与一个怪物对战情况
func Player_VS_Moster(player *model.WechatPlayer, mosterIndex int) (b_Win bool, HPLoss int) {
	//初始化
	HPLoss = 0
	b_Win = false

	//玩家的攻击倍率
	rate := player.Agility / Map_MonsterData[mosterIndex].Agility

	//玩家的DPS
	playerHurt := (player.Attack - Map_MonsterData[mosterIndex].Defense) * rate
	if playerHurt <= 0 {
		playerHurt = 1
	}

	//怪物的DPS
	mosterHurt := Map_MonsterData[mosterIndex].Attack - player.Defense
	if mosterHurt <= 0 {
		mosterHurt = 1
	}

	playerDPSTime := Map_MonsterData[mosterIndex].HP / playerHurt
	mosterDPSTime := player.Cur_HP / mosterHurt

	//logger.Debugln(playerHurt)
	//logger.Debugln(mosterHurt)
	//logger.Debugln(playerDPSTime)
	//logger.Debugln(mosterDPSTime)

	if playerDPSTime <= mosterDPSTime {
		HPLoss = playerDPSTime * mosterHurt
		player.Cur_HP -= HPLoss
		b_Win = true
	} else {
		HPLoss = player.Cur_HP
		player.Cur_HP = 0
		b_Win = false
	}

	return b_Win, HPLoss
}

//获取怪物将会掉落的物品，假设返回-1代表没有获取到东西
func GetMosterProp(mosterIndex int) (propIndex int) {
	propIndex = -1

	number := rand.Intn(100)

	for i := 0; i < len(Map_MonsterData[mosterIndex].PropsGetRate); i++ {
		if number < Map_MonsterData[mosterIndex].PropsGetRate[i] {
			propIndex = Map_MonsterData[mosterIndex].Props[i]
			break
		}
		number -= Map_MonsterData[mosterIndex].PropsGetRate[i]
	}

	return propIndex
}
