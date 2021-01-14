package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	ps "github.com/mitchellh/go-ps"
)

type StateProcessStruct struct {
	AccID     int    `json:"accid"`
	Msisdn    string `json:"msisdn"`
	Pid       int    `json:"pid"`
	PMsg      int    `json:"pmsg"`
	Timestamp int64  `json:"timestamp"`
}

func saveProcessFile(accID int, msisdn string, pid int, cMsg int) {
	timesTamp := time.Now().Unix()

	filename := "state_process.json"

	fmt.Printf("Pending msg: %v", cMsg)
	fmt.Printf("do saveProcessFile timestamp: %v", timesTamp)

	// creating json file section
	//check if file exist or not
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		//if file doesn't exist, so i'll create one
		_, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Cant create file, error: %v", err)
		} else {
			// first time file is created
			stateStructArr := []StateProcessStruct{}
			stateStruct := StateProcessStruct{AccID: accID, Msisdn: msisdn, Pid: pid, PMsg: cMsg, Timestamp: timesTamp}

			stateStructArr = append(stateStructArr, stateStruct)

			jsonByte, err := json.Marshal(stateStructArr)
			if err != nil {
				fmt.Printf("Fail to parse struct into json : %v\n", err.Error())
			} else {
				err = ioutil.WriteFile(filename, jsonByte, 0644)
				if err != nil {
					fmt.Printf("Fail to write json : %v\n", err.Error())
				}
			}

		}
	} else {
		fmt.Printf("Update content jsonFile\n")
		// if file exists , insert or update json file
		filename := "state_process.json"

		jsonObjArr := []StateProcessStruct{} // temp json content

		file, _ := ioutil.ReadFile(filename)
		json.Unmarshal(file, &jsonObjArr)

		if len(jsonObjArr) > 0 {
			for index := range jsonObjArr {
				fmt.Println("======jsonObjArr========")
				fmt.Println(jsonObjArr)
				fmt.Println("==============")

				if os.Getpid() == jsonObjArr[index].Pid {
					fmt.Println("Preparing to replace string: ", jsonObjArr[index])
					tempNewJsonObjArr := StateProcessStruct{
						AccID:     jsonObjArr[index].AccID,
						Msisdn:    jsonObjArr[index].Msisdn,
						Pid:       jsonObjArr[index].Pid,
						PMsg:      jsonObjArr[index].PMsg,
						Timestamp: time.Now().Unix(),
					}

					fmt.Println("tempNewJsonObjArr: ", tempNewJsonObjArr)
					jsonObjArr[index] = tempNewJsonObjArr

				} else {
					// TODO: add new element
					fmt.Println("os.Getpid() != jsonObjArr[index].Pid: ")

				}
			}
		} else {
			// TODO: do something

		}

		jsonArrByte, err := json.Marshal(jsonObjArr)
		if err != nil {
			fmt.Printf("Fail to parse struct into json : %v\n", err.Error())
		} else {
			err = ioutil.WriteFile(filename, jsonArrByte, 0644)
			if err != nil {
				fmt.Printf("Fail to write json : %v\n", err.Error())
			}
		}
	}

}

func listAllProcess() {
	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
		return
	}

	// map ages
	for x := range processList {
		var process ps.Process
		process = processList[x]
		log.Printf("%d\t%s\n", process.Pid(), process.Executable())

		// check process by pids
		if process.Executable() == "wasapbro" {
			// insert to json
			findProcessByPID(process.Pid())
			return
		}
	}
}

func findProcessByPID(pids int) bool {
	findProcess, err := ps.FindProcess(pids)
	if err != nil {
		log.Println("ps.findProcess() Failed, are you using windows?")
		return false
	}

	if findProcess == nil && err == nil {
		log.Println("prosesnya kosong bambang")
		return false
	} else {
		log.Println("prosesnya ada euyyy")
	}

	log.Printf("findProcessByPID %v\n", findProcess.Pid())
	return true
}

func main() {
	msisdn := "6287800000000" // unique
	accID := 1                // unique
	pids := os.Getpid()       // get self pid
	pendingMsg := 150         // count pending message
	for {
		saveProcessFile(accID, msisdn, pids, pendingMsg)
		time.Sleep(time.Second * 1)
	}
}
