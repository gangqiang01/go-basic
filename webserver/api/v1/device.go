package v1

import (
	"bytes"
	"encoding/hex"
	"net"

	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"

	"strings"
)

func AwakeDevice(c *gin.Context) {
	mac := c.Param("mac")

	if err := sendMagicPacket(mac); err != nil {
		responce.FailWithMessage("send broadcast error", c)
		return
	}
	responce.Ok(c)
}

func buildMagicPacket(macAddress string) ([]byte, error) {
	header := bytes.Repeat([]byte{0xFF}, 6)
	mac, err := hex.DecodeString(strings.ReplaceAll(macAddress, ":", ""))
	if err != nil {
		return nil, err
	}
	magicPacket := append(header, bytes.Repeat(mac, 16)...)
	return magicPacket, nil
}

// sendMagicPacket 发送WOL（Wake-on-LAN）数据包
func sendMagicPacket(macAddress string) error {
	magicPacket, err := buildMagicPacket(macAddress)
	if err != nil {
		return err
	}

	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(magicPacket)
	if err != nil {
		return err
	}

	return nil
}
