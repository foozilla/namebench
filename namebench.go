package main

import (
	"log"
	
	"github.com/andlabs/ui"

	"github.com/foozilla/namebench/dnschecks"
	"github.com/foozilla/namebench/dnsqueue"
	"github.com/foozilla/namebench/io"
	"fmt"
)

const (
	// QUEUELENGTH How many requests/responses can be queued at once
	QUEUELENGTH = 65535

	// WORKERS count of workers (same as Chrome's DNS prefetch queue)
	WORKERS = 8

	// COUNT of tests to run
	COUNT = 100
)

func main() {
	err := ui.Main(func() {			
		outText := ui.NewEntry()
		outText.SetReadOnly(true)
		
		outBox := ui.NewVerticalBox()
		outBox.Append(outText, true)

		title := ui.NewLabel("Find the fastest DNS server, \ntuned just for you.")
		//title.SetText
		button := ui.NewButton("Go")
		button.OnClicked(func(*ui.Button) {
			//TODO: call dns checks, update new scrolling box
			startDNSTest(outText)
		})
		dnsBox := ui.NewVerticalBox()
		dnsBox.Append(title, true)
		dnsBox.Append(button, true)
		
		title2 := ui.NewLabel("DNS SEC")
		button2 := ui.NewButton("Go")
		button2.OnClicked(func(*ui.Button) {
			//TODO: call dns checks, update new scrolling box
			dnsSec(outText)
		})
		secBox := ui.NewVerticalBox()
		secBox.Append(title2, true)
		secBox.Append(button2, true)		

		outerBox := ui.NewVerticalBox()
		outerBox.Append(outBox, true)
		outerBox.Append(ui.NewHorizontalSeparator(), false)
		outerBox.Append(dnsBox, true)
		outerBox.Append(ui.NewHorizontalSeparator(), false)
		outerBox.Append(secBox, true)

		window := ui.NewWindow("namebench", 200, 600, false)
		window.SetChild(outerBox)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func dnsSec(textBox *ui.Entry) {
	servers := []string{
		"8.8.8.8:53",
		"75.75.75.75:53",
		"4.2.2.1:53",
		"208.67.222.222:53",
	}
	for _, ip := range servers {
		result, err := dnschecks.DnsSec(ip)
		resultString := fmt.Sprintf("%s DNSSEC: %s (%s)", ip, result, err)
		log.Println(resultString)
		
		outText := fmt.Sprintf("%s\n%s", textBox.Text(), resultString)
		textBox.SetText(outText)
	}
}

func startDNSTest(textBox *ui.Entry) {
	records, err := io.GetSites()
	if err != nil {
		panic(err)
	}

	q := dnsqueue.StartQueue(QUEUELENGTH, WORKERS)

	for _, record := range records {
		q.Add("8.8.8.8:53", "A", record+".")
		log.Printf("Added %s", record)
	}
	q.SendCompletionSignal()
	answered := 0
	for {
		if answered == len(records) {
			break
		}
		result := <-q.Results
		answered++
		log.Printf("%s", result)
	}
	return
}
