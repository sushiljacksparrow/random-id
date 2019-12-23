package main

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"math"
	"net"
	"time"

	"github.com/gin-gonic/gin"
)

// generates 64 bit Ids out of which first 42 bits are timestamp, 10 bits are node
// and 12 bits are sequence number
const totalBits int = 64
const epochBits int = 48
const nodeBits int = 10
const sequenceBits int = 12
const customEpoch int = 1420070400000

var lastTimestamp int64 = -1
var sequence int64
var nodeID = generateNodeID()

func maxNodeID() int64 {
	return int64(math.Pow(2, float64(nodeBits)) - float64(1))
}

func maxSequence() int64 {
	return int64(math.Pow(2, float64(sequenceBits)) - float64(1))
}

func timestamp() int64 {
	return time.Now().Unix()*1000 - int64(customEpoch)
}

func waitTillNextTimestamp(currentTimestamp int64) int64 {
	for {
		currentTimestamp = timestamp()
		if currentTimestamp != lastTimestamp {
			break
		}
	}
	return currentTimestamp
}

func getMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func generateNodeID() int64 {
	var macAddress = getMacAddr()
	var hashCode = hash(macAddress)
	return int64(hashCode) & maxNodeID()
}

func nextID() int64 {
	currentTimestamp := timestamp()

	if currentTimestamp == lastTimestamp {
		sequence = (sequence + 1) & maxSequence()

		if sequence == 0 {
			currentTimestamp = waitTillNextTimestamp(currentTimestamp)

		}
	} else {
		sequence = 0
	}

	lastTimestamp = currentTimestamp

	id := currentTimestamp << (totalBits - epochBits)
	id |= (0 << (totalBits - epochBits - nodeBits))
	id |= sequence
	return id
}

func main() {
	fmt.Println(time.Now().Unix())
	r := gin.Default()

	r.GET("/next", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": nextID(),
		})
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
