package osc

import (
	"fmt"
	"math"
	"net"
	"os"
	"strconv"

	"bitbucket.org/liamstask/gosc"
	"github.com/jackyb/neobombe-sim/simulation"
)

var addr *net.UDPAddr
var conn *net.UDPConn

func Run(cs map[string]interface{}) {
	var err error

	if addr, err = net.ResolveUDPAddr("udp", "localhost:7770"); err != nil {
		fmt.Println("osc:", err)
		os.Exit(1)
	}
	if conn, err = net.DialUDP("udp", nil, addr); err != nil {
		fmt.Println("osc:", err)
		os.Exit(1)
	}

	fmt.Println("osc starting..")

	oscChan := cs["osc"].(chan simulation.Bombe)
	sigChan := cs["signal"].(chan os.Signal)
	go func() {
		for {
			send(<-oscChan)
		}
	}()

	sigChan <- <-sigChan
	fmt.Println("osc exiting..")
}

func send(bombe simulation.Bombe) {
	idx := 0
	for _, u := range bombe.Units {
		for _, r := range u.Rotors {
			m := buildMsg(r, idx, bombe.Spinning)
			m.WriteTo(conn)
			idx++
		}
	}
}

func buildMsg(rotor simulation.Rotor, idx int, spinning bool) osc.Message {
	angle := rotor.Angle

	speed := 0.0
	if spinning {
		speed = simulation.AngleStep / math.Pow(26.0, float64(idx%3))
	}

	m := osc.Message{Address: "/rotor/" + strconv.Itoa(idx)}
	m.Args = append(m.Args, float32(angle))
	m.Args = append(m.Args, float32(speed))
	return m
}
