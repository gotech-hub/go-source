package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/sony/sonyflake"
	"math"
	"math/big"
	"sync"
	"time"
)

const (
	startCounter = 1
)

var (
	idGen *iDGenerator
)

type iDGenerator struct {
	current       int64
	counter       uint16
	processUnique uint16
	mu            *sync.Mutex

	node *sonyflake.Sonyflake
}

func init() {
	newIDGenerator()
}

func newIDGenerator() {
	node, _ := sonyflake.New(sonyflake.Settings{
		StartTime: time.Date(1999, 11, 9, 3, 33, 33, 333, time.Local),
		MachineID: func() (uint16, error) {
			return randomUint16(), nil
		},
	})

	idGen = &iDGenerator{
		current:       time.Now().Unix(),
		counter:       startCounter,
		processUnique: randomUint16(),
		mu:            new(sync.Mutex),
		node:          node,
	}
}

func GetIdGenerate() *iDGenerator {
	return idGen
}

func (g *iDGenerator) GetID() [8]byte {
	g.mu.Lock()
	defer g.mu.Unlock()

	var b [8]byte

	now := time.Now().Unix()
	binary.BigEndian.PutUint32(b[0:4], uint32(now))
	binary.BigEndian.PutUint16(b[4:6], g.processUnique)

	g.counter++

	binary.BigEndian.PutUint16(b[6:8], g.counter)

	return b
}

func (g *iDGenerator) GetIDFormatHex() string {
	b := g.GetID()
	return g.ParseHex(b)
}

func (g *iDGenerator) GetIDString() string {
	b := g.GetID()
	return g.toString(b)
}

func (g *iDGenerator) ParseHex(id [8]byte) string {
	var buf [20]byte
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}

func (g *iDGenerator) toString(id [8]byte) string {
	return fmt.Sprintf("%v%v%v", g.getTimeString(id), g.getProcess(id), g.getNumber(id))
}

func (g *iDGenerator) getTimeString(id [8]byte) string {
	unixSecs := binary.BigEndian.Uint32(id[0:4])
	t := time.Unix(int64(unixSecs), 0)
	s := t.Format("060102150405")
	return s
}

func (g *iDGenerator) getProcess(id [8]byte) string {
	n := binary.BigEndian.Uint16(id[4:6])
	return fmt.Sprintf("%05d", n)
}

func (g *iDGenerator) getNumber(id [8]byte) string {
	n := binary.BigEndian.Uint16(id[6:8])
	return fmt.Sprintf("%05d", n)
}

func (g *iDGenerator) GetIDStringV2() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	id, _ := g.node.NextID()
	return fmt.Sprintf("%d", id)
}

func randomUint16() uint16 {
	randInt, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %v", err))
	}
	return uint16(randInt.Int64())
}
