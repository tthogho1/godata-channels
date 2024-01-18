// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

// data-channels is a Pion WebRTC application that shows how you can send/recv DataChannel messages from a web browser
package main

import (
	"browser/signal"
	"fmt"
	"syscall/js"
	"time"

	"github.com/pion/webrtc/v4"
)

func main() {
	// Everything below is the Pion WebRTC API! Thanks for using it

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	setimage()

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Peer Connection has gone to failed exiting")
			//os.Exit(0)
		}

		if s == webrtc.PeerConnectionStateClosed {
			// PeerConnection was explicitly closed. This usually happens from a DTLS CloseNotify
			fmt.Println("Peer Connection has gone to closed exiting")
			//os.Exit(0)
		}
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(5 * time.Second).C {
				message := signal.RandSeq(15)
				fmt.Printf("Sending '%s'\n", message)

				// Send the message as text
				sendErr := d.SendText(message)
				if sendErr != nil {
					panic(sendErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

	// Wait for the offer to be pasted
	js.Global().Set("startSession", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		go func() {
			el := getElementByID("remoteSessionDescription")
			sd := el.Get("value").String()
			if sd == "" {
				js.Global().Call("alert", "Session Description must not be empty")
				return
			}

			offer := webrtc.SessionDescription{}
			signal.Decode(sd, &offer)

			if err := peerConnection.SetRemoteDescription(offer); err != nil {
				handleError(err)
			}

			// Create an answer
			answer, err := peerConnection.CreateAnswer(nil)
			fmt.Println("CreateAnswer")
			fmt.Println(answer)

			if err != nil {
				handleError(err)
			}

			// Create channel that is blocked until ICE Gathering is complete
			gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
			fmt.Println("GatheringCompletePromise")

			// Sets the LocalDescription, and starts our UDP listeners
			err = peerConnection.SetLocalDescription(answer)
			if err != nil {
				handleError(err)
			}
			fmt.Println("SetLocalDescription")

			// Block until ICE Gathering is complete, disabling trickle ICE
			// we do this because we only can exchange one signaling message
			// in a production application you should exchange ICE Candidates via OnICECandidate
			<-gatherComplete

			// Output the answer in base64 so we can paste it in browser
			//log(signal.Encode(*peerConnection.LocalDescription()))
			localDescription := signal.Encode(*peerConnection.LocalDescription())
			fmt.Println(localDescription)
			local_el := getElementByID("localSessionDescription")
			local_el.Set("value", localDescription)

		}()
		return js.Undefined()
	}))

	js.Global().Set("copySDP", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					switch e := e.(type) {
					case error:
						handleError(e)
					default:
						handleError(fmt.Errorf("recovered with non-error value: (%T) %s", e, e))
					}
				}
			}()

			browserSDP := getElementByID("localSessionDescription")

			browserSDP.Call("focus")
			browserSDP.Call("select")

			copyStatus := js.Global().Get("document").Call("execCommand", "copy")
			if copyStatus.Bool() {
				log("Copying SDP was successful")
			} else {
				log("Copying SDP was unsuccessful")
			}
		}()
		return js.Undefined()
	}))
	// Block forever
	select {}
}

func setimage() {
	imgbase64 := "data:image/png;base64,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"

	el := getElementByID("qrcode")
	el.Set("src", imgbase64)
	//img := js.Global().Get("document").Call("createElement", "img")
	// src属性に文字列を設定する
	//img.Set("src", imgbase64)
	// imgタグをDOMに追加する
	//js.Global().Get("document").Get("body").Call("appendChild", img)
}

func getElementByID(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

func log(msg string) {
	el := getElementByID("logs")
	el.Set("innerHTML", el.Get("innerHTML").String()+msg+"<br>")
}

func handleError(err error) {
	fmt.Println("Unexpected error. Check console.")
	panic(err)
}
