package main

import (
	"math/rand"
	"log"
)

//读写锁加锁  解锁的问题
//解锁需要client relase请求
func MapAddNewClient(id int64) {
	_, exits := ClientState[id];if exits {
		//FIXME 每次新的acquire client_id已经存在 不会继续添加
		//FIXME 这里曲解了分布式中at most once的语义
		return
	} else {
		MapRWMutex.Lock()
		ClientState[id] = false //默认无lock
		MapRWMutex.Unlock()
		log.Print("add new client:", id)
	}
}

func IsLockAlready() bool{
	flag := false
	MapRWMutex.RLock()
	for _,v := range ClientState {
		if v == true {
			flag = true
			break
		}
	}
	MapRWMutex.RUnlock()
	return flag
}

//前提 所有client均未被加锁  随机给一个client加锁
//如果已经存在client锁 返回-1
func RandomLock() int64{
	var changeID int64 = -1
	if IsLockAlready(){
		return changeID
	}
	i := rand.Intn(len(ClientState))
	MapRWMutex.RLock()
	for k := range ClientState {
		if i == 0 {
			changeID = k
			break
		}
		i--
	}
	MapRWMutex.RUnlock()
	if changeID != -1{
		MapRWMutex.Lock()
		ClientState[changeID] = true
		MapRWMutex.Unlock()
	}
	return changeID
}

func GetCurrentLock() int64{
	var res int64 = -1
	MapRWMutex.RLock()
	for key,val := range ClientState{
		if val == true{
			res = key
		}
	}
	MapRWMutex.RUnlock()
	return res
}

func UnlockClient(id int64){
	_,exits := ClientState[id]; if exits{
		ClientState[id] = false
	}
}

func LockClear(){
	var changeID int64 = -1
	MapRWMutex.RLock()
	for key,val := range ClientState{
		if val == true{
			changeID = key
			break
		}
	}
	MapRWMutex.RUnlock()
	if changeID != -1{
		MapRWMutex.Lock()
		ClientState[changeID] = false
		MapRWMutex.Unlock()
	}
}