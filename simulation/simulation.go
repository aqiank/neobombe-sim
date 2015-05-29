package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jackyb/enigma"
	"github.com/jackyb/enigma/stringutil"
	"github.com/jackyb/neobombe-sim/twitter"
)

const (
	NumUnits        = 36
	NumCombinations = 26 * 26 * 26

	AngleStep = math.Pi / enigma.NumAlphabets
)

type Rotor struct {
	Angle float64 `json:"angle"`
}

func (r *Rotor) offset() int {
	return int(math.Mod(r.Angle, math.Pi) / AngleStep)
}

type Unit struct {
	Enigma enigma.Enigma `json:"enigma"`
	Rotors [3]Rotor      `json:"rotors"`
}

func (u *Unit) init(idx int) {
	Enigma := enigma.CreateStandardEnigma()
	Enigma.Step(idx * (NumCombinations / NumUnits))
	for i := range u.Rotors {
		u.Rotors[i].Angle = float64(Enigma.Offset(i)) * AngleStep
	}
	u.Enigma = Enigma
}

func (u *Unit) run() {
	for i := range u.Rotors {
		u.Rotors[i].Angle = float64(u.Enigma.Offset(i)) * AngleStep
	}
	u.Enigma.Step(1)
}

type Bombe struct {
	Units [NumUnits]Unit
	Spinning bool
}

var bombe Bombe

func init() {
	for i := range bombe.Units {
		bombe.Units[i].init(i)
	}
}

func update(stateChan chan<- Bombe) {
	for i := range bombe.Units {
		bombe.Units[i].run()
	}
	stateChan <- bombe
}

func Run(cs map[string]interface{}) {
	msgChan := cs["message"].(chan string)
	decChan := cs["decrypted"].(chan string)
	encChan := cs["encrypted"].(chan string)
	stateChan := cs["state"].(chan Bombe)
	oscChan := cs["osc"].(chan Bombe)
	sigChan := cs["signal"].(chan os.Signal)

	fmt.Println("simulation starting..")
	go func() {
		for {
			msg := <-msgChan

			println("simulation: Received message, encrypting:", msg)
			msg = stringutil.Sanitize(msg)
			enc := encrypt(msg)

			// start spinning
			bombe.Spinning = true

			println("simulation: Decrypting:", enc)
			key := stringutil.Sanitize(twitter.Track())
			decrypt(enc, key, encChan, stateChan)

			println("simulation: Decrypted:", msg)
			decChan <- msg

			// stop spinning
			bombe.Spinning = false
		}
	}()

	go func() {
		for {
			oscChan <- bombe
			time.Sleep(16 * time.Millisecond)
		}
	}()

	sigChan <- <-sigChan
	fmt.Println("simulation exiting..")
}

func encrypt(msg string) string {
	idx := rand.Intn(NumUnits)
	m := bombe.Units[idx].Enigma
	return m.Encrypt(msg)
}

func decrypt(msg, orig string, encChan chan string, stateChan chan Bombe) {
	for {
		for i := range bombe.Units {
			e := bombe.Units[i].Enigma.Clone()
			s := e.Encrypt(msg)
			encChan <- s
			if strings.Contains(s, orig) {
				return
			}
		}
		update(stateChan)
		time.Sleep(16 * time.Millisecond)
	}
}
