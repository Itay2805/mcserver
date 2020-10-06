package socket

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"github.com/itay2805/mcserver/common"
	"github.com/itay2805/mcserver/minecraft"
	"io"
	"log"
	"net"
)

type Packet interface {
	Encode(writer *minecraft.Writer)
}

const (
	Disconnected = -1
	Handshaking = 0
	Status = 1
	Login = 2
	Play = 3
)


type Socket struct {
	net.Conn

	// the protocol state
	State 			int

	// the reader and writer
	// for this socket
	reader 			*bufio.Reader
	writer 			io.Writer

	// the send queues
	// TODO: maybe don't use this...
	sendQueueIn		chan<- interface{}
	sendQueueOut	<-chan interface{}
	sendDone		chan bool

	// The buffer used for recving data
	recvBuffer		[65565]byte
	decompBuffer	[65565]byte

	// Compression threshold
	compression		bool
}

func NewSocket(conn net.Conn) Socket {
	in, out := common.MakeInfinite()
	return Socket{
		Conn:       	conn,
		State:      	Handshaking,
		reader:     	bufio.NewReader(conn),
		writer:     	conn,
		sendQueueIn: 	in,
		sendQueueOut: 	out,
		sendDone:		make(chan bool),
		recvBuffer: 	[65565]byte{},
		decompBuffer: 	[65565]byte{},
		compression:	false,
	}
}

func (s*Socket) EnableCompression() {
	s.compression = true
}

func (s *Socket) readVarint() (int, error) {
	numRead := 0
	result := 0
	for {
		read, err := s.reader.ReadByte()
		if err != nil {
			return 0, err
		}
		value := read & 0b01111111
		result |= int(value) << (7 * numRead)

		numRead++
		if numRead > 5 {
			return 0, errors.New("varint is too big")
		}

		if (read & 0b10000000) == 0 {
			return result, nil
		}
	}
}

func (s *Socket) Recv() ([]byte, error) {
	// read the length
	packetLength, err := s.readVarint()
	if err != nil {
		return nil, err
	}

	// read the full data of the packet
	_, err = io.ReadFull(s.reader, s.recvBuffer[:packetLength])
	if err != nil {
		return nil, err
	}

	if s.compression {
		reader := minecraft.Reader{Data: s.recvBuffer[:packetLength]}
		dataLength := reader.ReadVarint()

		if dataLength != 0 {
			// compressed
			r, err := zlib.NewReader(&reader)
			if err != nil {
				return nil, err
			}

			// recompress it fully
			_, err = io.ReadFull(r, s.decompBuffer[:dataLength])
			if err != nil {
				return nil, err
			}

			// return the uncompressed data
			return s.decompBuffer[:dataLength], nil
		} else {
			// uncompressed, skip the length byte
			return s.recvBuffer[1:packetLength], nil
		}
	} else {
		// return the packet
		return s.recvBuffer[:packetLength], nil
	}
}

func (s *Socket) writeVarint(val int) error {
	data := [5]byte{}
	offset := 0

	temp := byte(0)
	for {
		temp = byte(val & 0b01111111)
		val >>= 7
		if val != 0 {
			temp |= 0b10000000
		}
		data[offset] = temp
		offset++
		if val == 0 {
			break
		}
	}

	_, err := s.Write(data[:offset])
	return err
}

type sendRequest struct {
	data	[]byte
	done 	chan bool
}

func (s *Socket) Send(packet Packet) {
	writer := minecraft.Writer{}
	packet.Encode(&writer)
	s.SendRaw(writer.Bytes())
}

func (s *Socket) SendRaw(data []byte) {
	s.SendRawChan(data, nil)
}

func (s *Socket) SendSync(packet Packet) {
	writer := minecraft.Writer{}
	packet.Encode(&writer)
	s.SendRawSync(writer.Bytes())
}

func (s *Socket) SendRawSync(data []byte) {
	done := make(chan bool)
	s.SendRawChan(data, done)
	<-done
}

func (s *Socket) SendChan(packet Packet, done chan bool) {
	writer := minecraft.Writer{}
	packet.Encode(&writer)
	s.SendRawChan(writer.Bytes(), done)
}

func (s *Socket) SendRawChan(data []byte, done chan bool) {
	s.sendQueueIn <- sendRequest{
		data: data,
		done: done,
	}
}

func (s *Socket) StartSend() {
	buf := bytes.Buffer{}

	for req := range s.sendQueueOut {
		data := req.(sendRequest)

		if s.compression {
			if len(data.data) > 128 {
				// attempt to compress
				buf.Reset()
				w := zlib.NewWriter(&buf)
				_, _ = w.Write(data.data)
				_ = w.Close()

				compressed := buf.Bytes()[:buf.Len()]

				// write the size of compressed + data size
				err := s.writeVarint(
					common.VarintSize(int32(len(data.data))) + len(compressed),
				)
				if err != nil {
					log.Println("go error", err, "on", s.RemoteAddr())
					goto gotError
				}

				// write the data size
				err = s.writeVarint(len(data.data))
				if err != nil {
					log.Println("go error", err, "on", s.RemoteAddr())
					goto gotError
				}

				// write the data
				_, err = s.Write(compressed)
				if err != nil {
					log.Println("go error", err, "on", s.RemoteAddr())
					goto gotError
				}

				// skip the rest of the code, we are done
			} else {
				// don't compress
				// write the packet length
				err := s.writeVarint(len(data.data) + 1)
				if err != nil {
					log.Println("go error", err, "on", s.RemoteAddr())
					goto gotError
				}

				// set no compression
				err = s.writeVarint(0)
				if err != nil {
					log.Println("go error", err, "on", s.RemoteAddr())
					goto gotError
				}

				// write the data
				_, err = s.Write(data.data)
				if err != nil {
					log.Println("go error", err, "on", s.RemoteAddr())
					goto gotError
				}
			}
		} else {
			err := s.writeVarint(len(data.data))
			if err != nil {
				log.Println("go error", err, "on", s.RemoteAddr())
				goto gotError
			}

			_, err = s.Write(data.data)
			if err != nil {
				log.Println("go error", err, "on", s.RemoteAddr())
				goto gotError
			}
		}

		if data.done != nil {
			data.done <- true
		}
		continue

	gotError:
		data.done <- false
		break
	}

	s.sendDone <- true
}


func (s *Socket) Close() {
	// close the queue in, which will close
	// the send queue
	close(s.sendQueueIn)

	// wait for all the packets that need to
	// be sent to be sent
	<-s.sendDone

	// now close the connection itself
	_ = s.Conn.Close()
}
