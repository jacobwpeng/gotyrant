package gotyrant

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	MagicPut     = uint16(0xC810)
	MagicPutKeep = uint16(0xC811)
	MagicPutNR   = uint16(0xC818)
	MagicGet     = uint16(0xC830)
	//MagicMGet    = uint16(0xC831)
	MagicOut = uint16(0xC820)
)

type Client struct {
	conn   *net.TCPConn
	config Config
}

func establishConn(addr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("ResolveTCPAddr error: %s", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("DialTCP error: %s", err)
	}
	return conn, nil
}

func NewClient(config Config) (*Client, error) {
	client := &Client{
		config: config,
	}
	err := client.Reset()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (client *Client) Reset() error {
	client.Close()
	var err error
	client.conn, err = establishConn(client.config.Addr)
	return err
}

func (client *Client) Close() {
	if client.conn != nil {
		client.conn.Close()
		client.conn = nil
	}
}

func (client *Client) Put(key, value []byte) error {
	var req Request
	req.SetMagic(MagicPut)
	req.AddUInt32(uint32(len(key)))
	req.AddUInt32(uint32(len(value)))
	req.AddBinary(key)
	req.AddBinary(value)
	client.updateDeadline()
	_, err := req.WriteTo(client.conn)
	if err != nil {
		return err
	}
	return readCodeAsError(client.conn, CodeMisc)
}

func (client *Client) PutString(key, value string) error {
	return client.Put([]byte(key), []byte(value))
}

func (client *Client) PutNR(key, value []byte) error {
	var req Request
	req.SetMagic(MagicPutNR)
	req.AddUInt32(uint32(len(key)))
	req.AddUInt32(uint32(len(value)))
	req.AddBinary(key)
	req.AddBinary(value)
	return client.SendRequest(&req)
}

func (client *Client) PutKeep(key, value []byte) error {
	var req Request
	req.SetMagic(MagicPutKeep)
	req.AddUInt32(uint32(len(key)))
	req.AddUInt32(uint32(len(value)))
	req.AddBinary(key)
	req.AddBinary(value)
	err := client.SendRequest(&req)
	if err != nil {
		return err
	}
	return readCodeAsError(client.conn, CodeExist)
}

func (client *Client) Get(key []byte) ([]byte, error) {
	var req Request
	req.SetMagic(MagicGet)
	req.AddUInt32(uint32(len(key)))
	req.AddBinary(key)
	err := client.SendRequest(&req)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(client.conn)
	err = readCodeAsError(r, CodeNotExist)
	if err != nil {
		return nil, err
	}
	var valueLen uint32
	err = binary.Read(r, binary.BigEndian, &valueLen)
	if err != nil {
		return nil, err
	}
	value := make([]byte, valueLen)
	_, err = io.ReadFull(r, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

//func (client *Client) MGet(keys [][]byte) ([][]byte, error) {
//	var req Request
//	req.SetMagic(MagicMGet)
//	req.AddUInt32(uint32(len(keys)))
//	for _, key := range keys {
//		req.AddUInt32(uint32(len(key)))
//		req.AddBinary(key)
//	}
//	err := client.SendRequest(&req)
//	if err != nil {
//		return nil, err
//	}
//	r := bufio.NewReader(client.conn)
//	err = readCodeAsError(r, CodeNotExist)
//	if err != nil {
//		return nil, err
//	}
//	var numRecords uint32
//	err = binary.Read(r, binary.BigEndian, &numRecords)
//	if err != nil {
//		return nil, err
//	}
//	for i := 0; i < numRecords; i++ {
//		var keyLen, valueLen uint32
//		err = binary.Read(r, binary.BigEndian, &keyLen)
//		if err != nil {
//			return err
//		}
//	}
//}

func (client *Client) Out(key []byte) error {
	var req Request
	req.SetMagic(MagicOut)
	req.AddUInt32(uint32(len(key)))
	req.AddBinary(key)
	err := client.SendRequest(&req)
	if err != nil {
		return err
	}
	return readCodeAsError(client.conn, CodeNotExist)
}

func (client *Client) readCode() (error, int8) {
	return readCodeFrom(client.conn)
}

func (client *Client) updateDeadline() {
	var d time.Duration
	if client.config.Timeout == d {
		// No timeout
	} else {
		deadline := time.Now().Add(client.config.Timeout)
		client.conn.SetReadDeadline(deadline)
		client.conn.SetWriteDeadline(deadline)
	}
}

func (client *Client) SendRequest(req *Request) error {
	client.updateDeadline()
	_, err := req.WriteTo(client.conn)
	return err
}

func readCodeFrom(r io.Reader) (error, int8) {
	var buf [1]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return fmt.Errorf("Read code error: %s", err), 0
	}
	return nil, int8(buf[0])
}

func readCodeAsError(r io.Reader, errorCode int8) error {
	err, code := readCodeFrom(r)
	if err != nil {
		return err
	}
	if code != 0 {
		return &Error{Code: errorCode}
	}
	return nil
}
