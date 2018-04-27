package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca"
	"github.com/rodolfoag/gow32"
)

type FrequencyFlag struct {
	Value float32
	IsSet bool
}

type VolumeFlag struct {
	Value float32
	IsSet bool
}

func (f *FrequencyFlag) Set(value string) (err error) {
	var v float64
	if v, err = strconv.ParseFloat(value, 64); err != nil {
		return
	}
	if v > 3600.0 || v < 0.001 {
		err = fmt.Errorf("Invalid frequency")
		return
	}
	f.Value = float32(v)
	f.IsSet = true
	return
}

func (f *VolumeFlag) Set(value string) (err error) {
	var v float64
	if v, err = strconv.ParseFloat(value, 64); err != nil {
		return
	}
	if v > 0.99 || v < 0.0 {
		err = fmt.Errorf("Minimum volume range is 0 to 0.99")
		return
	}
	f.Value = float32(v)
	f.IsSet = true
	return
}

func (f *VolumeFlag) String() string {
	return fmt.Sprintf("%v", f.Value)
}

func (f *FrequencyFlag) String() string {
	return fmt.Sprintf("%v", f.Value)
}
func onReady() {
	systray.SetIcon(Data)
	systray.SetTitle("Ear Saver")
	systray.SetTooltip("Ear Saver")

	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()
	go func() {
		var first = true

		var volumeFlag VolumeFlag

		var frequencyFlag FrequencyFlag

		volumeFlag.Value = 0.15
		frequencyFlag.Value = 10

		f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		f.Var(&volumeFlag, "minvol", "Specify minimum volume as a scalar value")
		f.Var(&frequencyFlag, "freq", "Specify frequency in seconds when to lower volume by one step")
		f.Parse(os.Args[1:])

		for {
			var err error
			if err = run(volumeFlag, first); err != nil {
				log.Fatal(err)
				systray.Quit()
				return
			}
			first = false

			time.Sleep(time.Second * time.Duration(frequencyFlag.Value))
		}

	}()
}

func onExit() {
	// clean up here
}

func main() {

	_, err := gow32.CreateMutex("winearsave")
	if err != nil {
		fmt.Printf("Error: %d - %s\n", int(err.(syscall.Errno)), err.Error())
		return
	}

	systray.Run(onReady, onExit)

}

func run(volumeFlag VolumeFlag, first bool) (err error) {

	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator
	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err = mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		return
	}
	defer mmd.Release()

	var ps *wca.IPropertyStore
	if err = mmd.OpenPropertyStore(wca.STGM_READ, &ps); err != nil {
		return
	}
	defer ps.Release()

	var pv wca.PROPVARIANT
	if err = ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv); err != nil {
		return
	}
	if first {
		fmt.Printf("Device: %s\n", pv.String())
	}
	var aev *wca.IAudioEndpointVolume
	if err = mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
		return
	}
	defer aev.Release()

	var mute bool
	if err = aev.GetMute(&mute); err != nil {
		return
	}
	if mute {
		return
	}
	var masterVolumeLevelScalar float32
	if err = aev.GetMasterVolumeLevelScalar(&masterVolumeLevelScalar); err != nil {
		return
	}
	var newvol float32
	newvol = masterVolumeLevelScalar - 0.01
	if newvol >= 0 && newvol > volumeFlag.Value {
		if err = aev.SetMasterVolumeLevelScalar(masterVolumeLevelScalar-0.01, nil); err != nil {
			return
		}
	}

	return
}
