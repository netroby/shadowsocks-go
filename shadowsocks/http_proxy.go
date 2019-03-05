package shadowsocks

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
)

func HTTP_PROXY(network, addr string, auth *Auth, forward net.Dialer) (*http_proxy, error) {
	s := &http_proxy{
		network: network,
		addr:    addr,
		forward: forward,
	}
	if auth != nil {
		s.user = auth.User
		s.password = auth.Password
	}
	return s, nil
}

type http_proxy struct {
	user, password string
	network, addr  string
	forward        net.Dialer
}

func (s *http_proxy) Dial(network, addr string) (net.Conn, error) {
	switch network {
	case "tcp", "tcp6", "tcp4":
	default:
		return nil, errors.New("proxy: no support for HTTP proxy connections of type " + network)
	}

	conn, err := s.forward.Dial(s.network, s.addr)
	if err != nil {
		return nil, err
	}
	closeConn := &conn
	defer func() {
		if closeConn != nil {
			(*closeConn).Close()
		}
	}()

	fmt.Fprintf(conn, "CONNECT %s HTTP/1.1\r\n\r\n", addr)
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(status)
		return nil, err
	}

	closeConn = nil
	return conn, nil
}
