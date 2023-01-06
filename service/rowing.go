package service

import (
	// "pm5-emulator/service/mux"
	"pm5-emulator/engine"
	"github.com/sirupsen/logrus"
	"github.com/bettercap/gatt"
	"time"
	// "strconv"
	// "math/big"
	"os"
	"fmt"
	"bufio"
	"strings"
	"encoding/hex"
)

/*
	C2 rowing primary service
*/

//C2 rowing primary service and characteristics UUIDs
var (
	attrRowingServiceUUID, _                                    = gatt.ParseUUID(getFullUUID("0030"))
	attrGeneralStatusCharacteristicsUUID, _                     = gatt.ParseUUID(getFullUUID("0031"))
	attrAdditionalStatus1CharacteristicsUUID, _                 = gatt.ParseUUID(getFullUUID("0032"))
	attrAdditionalStatus2CharacteristicsUUID, _                 = gatt.ParseUUID(getFullUUID("0033"))
	attrSampleRateCharacteristicsUUID, _                        = gatt.ParseUUID(getFullUUID("0034"))
	attrStrokeDataCharacteristicsUUID, _                        = gatt.ParseUUID(getFullUUID("0035"))
	attrAdditionalStrokeDataCharacteristicsUUID, _              = gatt.ParseUUID(getFullUUID("0036"))
	attrSplitIntervalDataCharacteristicsUUID, _                 = gatt.ParseUUID(getFullUUID("0037"))
	attrAdditionalSplitIntervalDataCharacteristicsUUID, _       = gatt.ParseUUID(getFullUUID("0038"))
	attrEndOfWorkoutSummaryDataCharacteristicsUUID, _           = gatt.ParseUUID(getFullUUID("0039"))
	attrAdditionalEndOfWorkoutSummaryDataCharacteristicsUUID, _ = gatt.ParseUUID(getFullUUID("003A"))
	attrHeartRateBeltInfoCharacteristicsUUID, _                 = gatt.ParseUUID(getFullUUID("003B"))
	attrForceCurveDataCharacteristicsUUID, _                    = gatt.ParseUUID(getFullUUID("003D"))
	attrMultiplexedInfoCharacteristicsUUID, _                   = gatt.ParseUUID(getFullUUID("0080"))
)

//NewRowingService advertises rowing service defined by PM5 device
func NewRowingService() *gatt.Service {
	s := gatt.NewService(attrRowingServiceUUID)
	
	rowingEngine := engine.NewRowingEngine()

	replayFileName := "replaylog.erg"

	
	/*
		C2 rowing general status characteristic
	*/
	rowingGenStatusChar := s.AddCharacteristic(attrGeneralStatusCharacteristicsUUID)

	rowingGenStatusChar.HandleNotifyFunc(
		func(r gatt.Request, n gatt.Notifier) {
			logrus.Info("General Status Char Notify Request")
			go func() {
				for true {
					n.Write(rowingEngine.GenerateGeneralStatusChar())										
					time.Sleep(500 * time.Millisecond)
				}
			}()
		})	

	/*
		C2 rowing additional status 1 characteristic
	*/
	additionalStatus1Char := s.AddCharacteristic(attrAdditionalStatus1CharacteristicsUUID)
	additionalStatus1Char.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("Additional Status 1 Char Notify Request - launching goroutine")
		go func() {
			for true {
				logrus.Info("Sending Additional Status 1 Notification from goroutine")				
				n.Write(rowingEngine.GenerateAdditionalStatus1Char())										
				time.Sleep(500 * time.Millisecond)
			}
		}()
	})

	/*
		C2 rowing additional status 2 characteristic
	*/
	additionalStatus2Char := s.AddCharacteristic(attrAdditionalStatus2CharacteristicsUUID)
	additionalStatus2Char.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("Additional Status 2 Char Notify Request - launching goroutine")
		go func() {
			for true {
				n.Write(rowingEngine.GenerateAdditionalStatus2Char())
				time.Sleep(500 * time.Millisecond)
			}
		}()	
	})

	/*
		C2 rowing general status and additional status sample rate characteristic 0x0034
	*/
	sampleRateChar := s.AddCharacteristic(attrSampleRateCharacteristicsUUID)
	sampleRateChar.HandleReadFunc(func(rsp gatt.ResponseWriter, req *gatt.ReadRequest) {
		logrus.Info("Sample Rate Char Read Request")
		data := make([]byte, 1)
		rsp.Write(data)
	})

	sampleRateChar.HandleWriteFunc(func(req gatt.Request, data []byte) (status byte) {
		logrus.Info("Sample Rate Char Write Request: ", string(data))
		if (len(data) > 1){
			logrus.Error("Sample Rate Char Write Request received more than one byte")
		}
		return gatt.StatusSuccess
	})

	/*
		C2 rowing stroke data  characteristic 0x0035
	*/
	strokeDataChar := s.AddCharacteristic(attrStrokeDataCharacteristicsUUID)
	strokeDataChar.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("Stroke Data Char Notify Request - launching goroutine")
		go func() {
			for true {
				logrus.Info("Stroke Data Notification from goroutine")
				n.Write(rowingEngine.GenerateStrokeDataChar())
				time.Sleep(1000 * time.Millisecond)
			}
		}()	
	})

	/*
		C2 rowing additional stroke data characteristic 0x0036
	*/
	additionalStrokeDataChar := s.AddCharacteristic(attrAdditionalStrokeDataCharacteristicsUUID)
	additionalStrokeDataChar.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("Additional Stroke Data Char Notify Request - launching goroutine")
		go func() {
			for true {
				logrus.Info("Additional Stroke Data Notification from goroutine")
				n.Write(rowingEngine.GenerateStrokeData2Char())
				time.Sleep(1000 * time.Millisecond)
			}
		}()	
	})

	/*
		C2 rowing split/interval data characteristic
	*/
	splitIntervalDataChar := s.AddCharacteristic(attrSplitIntervalDataCharacteristicsUUID)
	splitIntervalDataChar.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("Split/Interval Data Char Notify Request - launching goroutine")
		go func() {
			for true {
				logrus.Info("Split/Interval Data Notification from goroutine")
				n.Write(rowingEngine.GenerateSplitIntervalChar())
				time.Sleep(50000 * time.Millisecond)
			}
		}()	
	})


	/*
		C2 rowing additional split/interval data characteristic
	*/
	additionalSplitIntervalDataChar := s.AddCharacteristic(attrAdditionalSplitIntervalDataCharacteristicsUUID)
	additionalSplitIntervalDataChar.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("Additional Split/Interval Data Char Notify Request - launching goroutine")
		go func() {
			for true {
				logrus.Info("Additional Split/Interval Data Notification from goroutine")
				n.Write(rowingEngine.GenerateSplitInterval2Char())
				time.Sleep(50000 * time.Millisecond)
			}
		}()	
	})

	/*
		C2 rowing end of workout summary data characteristic
	*/
	endOfWorkoutSummaryDataChar := s.AddCharacteristic(attrEndOfWorkoutSummaryDataCharacteristicsUUID)
	endOfWorkoutSummaryDataChar.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("End of workout summary Data Char Notify Request - launching goroutine")
		go func() {
			for true {
				time.Sleep(200000 * time.Millisecond)
				logrus.Info("End of workout summary Data Notification from goroutine")
				n.Write(rowingEngine.GenerateWorkoutSummaryChar())
			}
		}()	
	})

	/*
		C2 rowing end of workout additional summary data characteristic
	*/
	additionalEndOfWorkoutSummaryDataChar := s.AddCharacteristic(attrAdditionalEndOfWorkoutSummaryDataCharacteristicsUUID)
	additionalEndOfWorkoutSummaryDataChar.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("End of workout Additional summary Data Char Notify Request - launching goroutine")
		go func() {
			for true {
				time.Sleep(200000 * time.Millisecond)
				logrus.Info("End of workout Additional summary Data Notification from goroutine")
				n.Write(rowingEngine.GenerateWorkoutSummary2Char())
			}
		}()	
	})


	/*
		C2 rowing heart rate belt information characteristic
	*/
	heartRateBeltInfoChar := s.AddCharacteristic(attrHeartRateBeltInfoCharacteristicsUUID)
	heartRateBeltInfoChar.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("Heart Rate Belt Info Char Notify Request - launching goroutine")
		go func() {
			for true {
				logrus.Info("Heart Rate Belt Data Notification from goroutine")
				n.Write(make([]byte, 6))
				time.Sleep(100000 * time.Millisecond)

			}
		}()	
	})

	/*
		C2 force curve data characteristic
	*/
	forceCurveDataChar := s.AddCharacteristic(attrForceCurveDataCharacteristicsUUID)
	forceCurveDataChar.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("Force Curve Data Char Notify Request - launching goroutine")
		go func() {
			for true {
				logrus.Info("Force Curve Data Notification from goroutine")
				n.Write(rowingEngine.GenerateForceCurveChar())
				time.Sleep(1000 * time.Millisecond)
			}
		}()	
	})

	/*
		C2 multiplexed information 	characteristic

		0x0080 | Up to 20 bytes | NOTIFY Permission
	*/
	multiplexedInfoChar := s.AddCharacteristic(attrMultiplexedInfoCharacteristicsUUID)

	multiplexedInfoChar.HandleNotifyFunc(func(r gatt.Request, n gatt.Notifier) {
		logrus.Info("Multiplex Info Char Notify Func")
		//generate a rowing general status payload here
		// m:=mux.Multiplexer{}
		// var count = 0

		replayFile, err := os.OpenFile(replayFileName, os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("file error")
			fmt.Println(err)
		}


		logrus.Info("Multiplexed Data Char Notify Request - launching goroutine")
		go func() {
			scanner := bufio.NewScanner(replayFile)
			for scanner.Scan() {
				text := scanner.Text()
				textsplit := strings.Split(text, ":") 
				delta := textsplit[0]
				data := textsplit[1] 
				fmt.Println(data)
				
				logrus.Info("Multiplexed Data Notification from goroutine")
				// bytes := 
				// copy(bytes[0:], count)
				dur,err := time.ParseDuration(delta+"ms")
				if err != nil {
					fmt.Println("error parsing duration from string")
					fmt.Println(err)
				}

				bytes,decodeerr := hex.DecodeString(data)
				if decodeerr != nil {
					fmt.Println("error decoding bytes from hex")
					fmt.Println(decodeerr)
				}

				if dur.Milliseconds() > 50 {
					fmt.Println("sleeping for " + delta + " ms")
					time.Sleep(dur)
				} else {
					time.Sleep(50 * time.Millisecond)
				}
				n.Write(bytes)
				
			}

			if err := scanner.Err(); err != nil {
				fmt.Println(err)
			}

			defer replayFile.Close() 
		}()	
	})


	return s
}
