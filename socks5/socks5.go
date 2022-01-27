package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const MaxWorkerCount = 40

func main() {
	port := os.Getenv("SOCKS_PORT")
	whiteAddr := os.Getenv("WHITE_ADDR")

	arrWhiteAddr := strings.Split(whiteAddr, ",")
	if len(arrWhiteAddr) == 0 {
		log.Println("[WARN] none of your services are available to the public")
	}

	var goWorkerCount int32 = 0

	if len(port) == 0 {
		port = "13030"
	}
	addr := fmt.Sprintf("%s:%s", "0.0.0.0", port)

	req, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("start server failed : ", err.Error())
	}

	for {
		var cli net.Conn
		cli, err = req.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		if c := atomic.LoadInt32(&goWorkerCount); c >= MaxWorkerCount {
			log.Println("too many connections are loaded simultaneously")
			cli.Close()
			continue
		}

		atomic.AddInt32(&goWorkerCount, 1)

		rAddr := cli.RemoteAddr().String()
		rAddr = strings.Split(rAddr, ":")[0]
		if !inArr(arrWhiteAddr, rAddr) {
			log.Printf("%s without authorization\n", rAddr)
			cli.Close()
			continue
		}

		log.Println(rAddr ,"Connected")
		go worker(cli, &goWorkerCount)
	}
}

func worker(cli net.Conn, count *int32) {
	defer func() {
		atomic.AddInt32(count, -1)
	}()

	err := socks5ConnAuth(cli)
	if err != nil {
		log.Println(err)
	}

	target, err := socks5ConnectMethod(cli)
	if err != nil {
		log.Println(err)
		cli.Close()
		return
	}
	err = relay(cli, target)
	if err != nil {
		log.Printf("relay error: %v\n", err)
	}
}

func socks5ConnAuth(cli net.Conn) error {
	buf := make([]byte, 256)
	n, err := io.ReadFull(cli, buf[:2])
	if n != 2 || err != nil {
		return errors.New("reading header: " + err.Error())
	}

	ver, methods := int(buf[0]), int(buf[1])

	if ver != 5 {
		return errors.New("invalid socks version")
	}

	n, err = io.ReadFull(cli, buf[:methods])
	if n != methods || err != nil {
		return errors.New("reading methods : " + err.Error())
	}

	// the socks5 protocol is end

	n, err = cli.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return errors.New("write response error : " + err.Error())
	}

	return nil
}

func socks5ConnectMethod(cli net.Conn) (net.Conn, error) {
	buf := make([]byte, 256)

	_, err := io.ReadFull(cli, buf[:4])
	if err != nil {
		return nil, errors.New("read socks header : " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, errors.New("invalid ver/cmd")
	}

	addr := ""

	// VER, CMD, RSV
	switch atyp {
	case 1:
		_, err = io.ReadFull(cli, buf[:4])
		if err != nil {
			return nil, err
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])

	case 3:
		_, err = io.ReadFull(cli, buf[:1])
		if err != nil {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		_, err = io.ReadFull(cli, buf[:addrLen])
		if err != nil {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])
	case 4:
		return nil, errors.New("IPv6: no supported yet")

	default:
		return nil, errors.New("invalid atyp")
	}

	_, err = io.ReadFull(cli, buf[:2])
	if err != nil {
		return nil, err
	}

	port := binary.BigEndian.Uint16(buf[:2])

	dstAddrPort := fmt.Sprintf("%s:%d", addr, port)
	dest, err := net.Dial("tcp", dstAddrPort)
	if err != nil {
		return nil, errors.New("dial dst : " + err.Error())
	}

	_, err = cli.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		err = dest.Close()
		if err != nil {
			return nil, err
		}
		return nil, errors.New("write rsp : " + err.Error())
	}

	return dest, nil
}

func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)

	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait))
	}()

	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait))
	wg.Wait()

	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}

	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) {
		return err1
	}

	return nil
}

func inArr(all []string, object string) bool {
	for _, v := range all {
		if v == object {
			return true
		}
	}

	return false
}
