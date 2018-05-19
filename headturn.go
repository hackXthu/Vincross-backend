package headturn

import (
	"os"
	"fmt"
	"time"
	"math"
	"mind/core/framework/skill"
	"mind/core/framework/drivers/distance"
  	"mind/core/framework/drivers/hexabody"
  	"mind/core/framework/log"
	"mind/core/framework/drivers/media"
	"net/http"
)

const (
	STAND_TIME = 2.0
	TURN_TIME  = 5.0
)

type headturn struct {
	skill.Base
}

func NewSkill() skill.Interface {
	// Use this method to create a new skill.

	return &headturn{}
}

func (d *headturn) circlehead() {
	log.Info.Println("Turning head")
	hexabody.Start()
	/* var angle float64
	for {
		hexabody.MoveHead(angle, TURN_TIME)
		angle += 60.0
		if angle>=360{
			angle -= 360.0
		time.Sleep(STAND_TIME)
		}
	} */

	hexabody.RotateHeadContinuously(1, 60)
	hexabody.Close()
}
func (d *headturn) OnStart() {
	// Use this method to do something when this skill is starting.
	log.Info.Println("OnStart, Program started.")
	err := hexabody.Start()
	if err != nil {
		log.Error.Println("Hexabody start err:", err)
		return
	}
	hexabody.Stand()
	if !media.Available() {
		log.Error.Println("Media driver not available")
		return
	}
	if err = media.Start(); err != nil {
		log.Error.Println("Media driver could not start")
	}
  	distance.Start()
}

func (d *headturn) OnClose() {
	// Use this method to do something when this skill is closing.
	log.Info.Println("OnClose: System Closed.")
	hexabody.Close()
  	distance.Close()
}

func (d *headturn) OnConnect() {
	// Use this method to do something when the remote connected.
	// go d.circlehead()
	var angle float64
	var img_count int
	for {
		hexabody.MoveHead(angle, TURN_TIME)
		angle += 1
		if angle>=360{
			angle -= 360.0
		}
		time.Sleep(STAND_TIME)
		robo, _ := skill.RobotInfo() // cs
		if math.Mod(angle, 45) == 0{
			filepath, _ := skill.SkillDataPath()
			log.Info.Println(filepath)
			log.Info.Println(fmt.Sprintf("Saving Img %d.jpeg",img_count))
			media.SnapshotJPEG(fmt.Sprintf("%s/%d.jpeg",filepath,img_count),70)
			img_count ++
			
			var s float64 = float64(robo.DiskSpace)/1024/1024/1024 // cs
			log.Info.Println("Available diskspace is ", s, "G") // cs

			err := media.Start()
			if err != nil {
			log.Error.Println("Media start err:", err)
			return
			}
			for {
			log.Info.Println("Connected")
			buf := new(bytes.Buffer)
			log.Info.Println("JPEG")
			jpeg.Encode(buf, media.SnapshotYCbCr(), nil)
			log.Info.Println("BASE64")
			str := base64.StdEncoding.EncodeToString(buf.Bytes())
			log.Info.Println("SENDING")
			framework.SendString(str)
			log.Info.Println("Sent:", str[:20], len(str))
			}
			resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)

			// 上传
		}

		}

	

	}



func (d *headturn) OnDisconnect() {
	// Use this method to do something when the remote disconnected.
	log.Info.Println("OnDisconnect")
	os.Exit(0) // Closes the process when remote disconnects
}

func (d *headturn) OnRecvJSON(data []byte) {
	// Use this method to do something when skill receive json data from remote client.
}

func (d *headturn) OnRecvString(data string) {
	// Use this method to do something when skill receive string from remote client.
	go d.circlehead()
}
