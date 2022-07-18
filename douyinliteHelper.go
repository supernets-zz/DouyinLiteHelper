package main

import (
	"fmt"
	"log"
	"time"

	. "github.com/electricbubble/gwda"
)

var (
	driver           WebDriver
	chestBtn         WebElement
	errChestBtn      error
	getBtns          []WebElement
	errGetBtns       error
	timeLimitTips    WebElement
	errTimeLimitTips error
)

func watchAD(driver WebDriver) error {
	log.Println("watchAD begin")
	defer log.Println("watchAD end")
	startTick := time.Now()
	for {
		_, err := driver.FindElement(BySelector{Name: "广告"})
		if err == nil {
			break
		}

		time.Sleep(time.Second)
		log.Printf("等待广告开始 %ds\n", time.Now().Unix()-startTick.Unix())
		if time.Now().Unix()-startTick.Unix() > 15 {
			log.Println("AD loading timeout")
			return fmt.Errorf("广告加载失败")
		}
	}

	startTick = time.Now()
	for {
		succTips, errSucc := driver.FindElement(BySelector{Name: "领取成功"})
		backBtn, errBackBtn := driver.FindElement(BySelector{Name: "返回"})
		closeBtn, errCloseBtn := driver.FindElement(BySelector{Name: "关闭，按钮"})
		if errSucc == nil {
			log.Println("点击 领取成功")
			err := succTips.Click()
			if err != nil {
				log.Println("点击 领取成功: " + err.Error())
				return err
			}
			return nil
		}

		if errBackBtn == nil {
			log.Println("点击 返回")
			err := backBtn.Click()
			if err != nil {
				log.Println("点击 领取成功: " + err.Error())
				return err
			}
		}

		if errCloseBtn == nil {
			log.Println("点击 关闭，按钮")
			err := closeBtn.Click()
			if err != nil {
				log.Println("点击 关闭，按钮: " + err.Error())
				return err
			}
			return nil
		}

		time.Sleep(time.Second)
		log.Printf("watchAD %ds, 领取成功: %v, 返回: %v, 关闭，按钮: %v", time.Now().Unix()-startTick.Unix(), errSucc == nil, errBackBtn == nil, errCloseBtn == nil)
		if time.Now().Unix()-startTick.Unix() > 60 {
			log.Println("watchAD timeout")
			break
		}
	}
	return fmt.Errorf("watchAD timeout")
}

func watchChestAD(driver WebDriver, chestBtn WebElement) error {
	log.Println("观看宝箱广告")
	defer log.Println("宝箱广告 结束")
	log.Println("点击 开宝箱得音符")
	err := chestBtn.Click()
	if err != nil {
		log.Println("点击 开宝箱得音符: " + err.Error())
		return err
	}

	time.Sleep(3 * time.Second)
	confirmBtn, err := driver.FindElement(BySelector{Name: "看广告视频再赚"})
	if err == nil {
		log.Println("点击 看广告视频再赚")
		err = confirmBtn.Click()
		if err != nil {
			log.Println("点击 看广告视频再赚: " + err.Error())
			return err
		}
		time.Sleep(3 * time.Second)
	}

	err = watchAD(driver)
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)
	nextADBtn, err := driver.FindElement(BySelector{Name: "领取奖励"})
	if err == nil {
		log.Println("点击 领取奖励")
		err = nextADBtn.Click()
		if err != nil {
			log.Println("点击 领取奖励: " + err.Error())
			return err
		}
	}

	err = watchAD(driver)
	return err
}

func watchTimeLimitAD(driver WebDriver, timeLimitTips WebElement) error {
	log.Println("观看限时广告")
	defer log.Println("限时广告 结束")
	rect, err := timeLimitTips.Rect()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("点击 每5分钟完成一次广告任务，单日最高可赚20000音符 去领取")
	err = driver.Tap(rect.X+rect.Width+48, rect.Y)
	if err != nil {
		log.Println("点击 去领取: " + err.Error())
		return err
	}

	err = watchAD(driver)
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)
	nextADBtn, err := driver.FindElement(BySelector{Name: "领取奖励"})
	if err == nil {
		log.Println("点击 领取奖励")
		err = nextADBtn.Click()
		if err != nil {
			log.Println("点击 领取奖励: " + err.Error())
			return err
		}
	} else {
		log.Println("无 领取奖励")
	}

	err = watchAD(driver)
	return err
}

func douyinADAutoClose(driver WebDriver) {
	for {
		log.Println("寻找 去领取 或 开宝箱得音符")
		chestBtn, errChestBtn = driver.FindElement(BySelector{Name: "开宝箱得音符"})
		getBtns, errGetBtns = driver.FindElements(BySelector{Name: "去领取"})
		if len(getBtns) == 2 {
			timeLimitTips, errTimeLimitTips = driver.FindElement(BySelector{Name: "每5分钟完成一次广告任务，单日最高可赚20000音符"})
		}
		log.Printf("开宝箱得音符: %v, 去领取: %v, 限时广告: %v", errChestBtn == nil, len(getBtns), len(getBtns) == 2 && errTimeLimitTips == nil)
		if errChestBtn == nil {
			time.Sleep(3 * time.Second)
			err := watchChestAD(driver, chestBtn)
			if err != nil {
				log.Printf("watchChestAD return " + err.Error())
			}
			continue
		}

		if len(getBtns) == 2 && errTimeLimitTips == nil {
			err := watchTimeLimitAD(driver, timeLimitTips)
			if err != nil {
				log.Printf("watchTimeLimitAD return " + err.Error())
			}
			continue
		}

		log.Println("等待30s")
		time.Sleep(30 * time.Second)
	}
}

func main() {
	log.Println("start")
	driver, err := NewUSBDriver(nil)
	checkErr(err)
	log.Println("connected")

	windowSize, _ := driver.WindowSize()
	log.Println(windowSize)

	douyinADAutoClose(driver)
}

func checkErr(err error, msg ...string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
