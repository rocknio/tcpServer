package main

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"
)

// TMsg 消息结构体
type TMsg struct {
	head   uint8
	addr   uint8
	len    uint16
	sn     uint16
	serial uint32
	msisdn string
	iccid  string
	cmd    uint8
	param  uint8
	data   []byte
	crc    uint8
}

func readPacket(conn net.Conn) ([]byte, uint16, error) {
	var msgHeader [4]byte

	if _, err := io.ReadFull(conn, msgHeader[:]); err != nil {
		return nil, 0, err
	}

	size := binary.BigEndian.Uint16(msgHeader[2:])
	packet := make([]byte, size)

	if _, err := io.ReadFull(conn, packet); err != nil {
		return nil, 0, err
	}

	return packet, size, nil
}

// Unmarshal TMsg的解析方法
func (m *TMsg) Unmarshal(b []byte, packetLen uint16) {
	var idx uint16
	idx = 0

	//m.head = uint8(b[idx])
	//idx++
	//
	//m.addr = uint8(b[idx])
	//idx++
	//
	//m.len = binary.BigEndian.Uint16(b[idx:])
	//idx += 2

	m.sn = binary.BigEndian.Uint16(b[idx:])
	idx += 2

	m.serial = binary.BigEndian.Uint32(b[idx:])
	idx += 4

	m.msisdn = string(b[idx : idx+11])
	idx += 11

	m.iccid = string(b[idx : idx+20])
	idx += 20

	m.cmd = uint8(b[idx])
	idx++

	m.param = uint8(b[idx])
	idx++

	m.data = b[idx : idx+(packetLen-2-4-11-20-1-1-1)]
	idx += packetLen - 2 - 4 - 11 - 20 - 1 - 1 - 1

	m.crc = b[idx]
}

// MsgHandler 消息处理函数
func MsgHandler(conn net.Conn) {
	defer conn.Close()

	for {
		msg, packetLen, err := readPacket(conn)
		if err != nil {
			MainLogger.Error("MsgHandler err: " + err.Error())
			return
		}

		SendStatus("R")

		MainLogger.Debug("Recv Msg, LEN = " + strconv.Itoa(int(packetLen)))
		var tMsg TMsg
		tMsg.Unmarshal(msg, packetLen)
		conn.Write([]byte("7a000049000f00050781343434343535353533333333343434343535353536363636373737373838380800023039303630393235353900026928000730393036303932353539000156940007810d0a"))
		MainLogger.Debug("Write Resp SUCC!")

		SendStatus("P")

		// 去掉每条消息后的0d 0a
		var tmp [2]byte
		if _, err := io.ReadFull(conn, tmp[:]); err != nil {
			MainLogger.Error("MsgHandler err: " + err.Error())
			return
		}
	}
}
