package rpc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

const constHeaderLength = 4

type rpcData struct {
	Name string
	Arg  interface{}
	Err  string
}

func packet(message []byte) []byte {
	return append(intToBytes(len(message)), message...)
}

func unpack(reader *bufio.Reader) ([]byte, error) {
	//消息头长度
	peek, err := reader.Peek(constHeaderLength)
	if err != nil {
		return nil, err
	}
	//消息长度
	bodylen := bytesToInt(peek)
	_, err = reader.Peek(bodylen + constHeaderLength)
	if err != nil {
		return nil, err
	}
	data := make([]byte, bodylen+constHeaderLength)
	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}
	return data[constHeaderLength:], nil
}

//整形转换成字节
func intToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func bytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func encode(data rpcData) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&data)
	return buffer.Bytes(), err
}

func decode(data []byte) (rpcData, error) {
	var message rpcData
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&message)
	return message, err
}
