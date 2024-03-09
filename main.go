package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/konradit/gowpd"
)

var (
	Devices []Device
)

type Device struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Manufacture string `json:"manufacture"`
	Size        string
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "devices":
			devices()
		case "send":
			if len(os.Args) > 4 {
				copyTo(os.Args[2], os.Args[3], os.Args[4])
			} else {
				copyTo(os.Args[2], os.Args[3])
			}
		default:
			fmt.Println("ERROR")
		}
	}
}

func devices() {
	err := gowpd.Init()
	if err != nil {
		panic(err)
	}
	defer gowpd.Destroy()

	n := gowpd.GetDeviceCount()
	for i := 0; i < n; i++ {
		Devices = append(Devices, Device{
			ID:          i,
			Name:        gowpd.GetDeviceName(i),
			Description: gowpd.GetDeviceDescription(i),
			Manufacture: gowpd.GetDeviceManufacturer(i),
		})
	}

	json, err := json.Marshal(Devices)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))
}

func copyTo(file, to string, dir ...string) {
	err := gowpd.Init()
	if err != nil {
		panic(err)
	}
	defer gowpd.Destroy()

	deviceid, err := strconv.Atoi(to)
	if err != nil {
		panic(err)
	}

	device, err := gowpd.ChooseDevice(deviceid)
	if err != nil {
		panic(err)
	}
	defer device.Release()
	var path string
	o := device.FindObject(gowpd.PathSeparator)
	if o == nil {
		panic("no device object id")
	}
	path = o.Id

	if len(dir) > 0 {
		d := device.FindObject(gowpd.PathSeparator + dir[0])
		if d == nil { // create filder if does not exists
			n, err := device.CreateFolder(o.Id, dir[0])
			if err != nil {
				panic("cant create folder")
			}
			path = n
		} else { // if folder already exists
			path = d.Id
		}
	}

	// check if file exists
	pathArr := strings.Split(file, "\\")
	filename := pathArr[len(pathArr)-1]
	filePathInDevice := gowpd.PathSeparator
	if len(dir) > 0 {
		filePathInDevice += dir[0] + gowpd.PathSeparator
	}
	filePathInDevice += filename
	isFileExists := device.FindObject(filePathInDevice)
	if isFileExists != nil {
		fmt.Println("done")
		return
	}
	_, err = device.CopyToDevice(path, file)
	if err != nil {
		panic(err)
	}
	fmt.Println("done")
}
