package shadowsocks

import (
	"errors"
	"net"
	"net/url"
)

type ParentProxy interface {
	Dial(network string, addr string) (net.Conn, error)
}

type NoProxy struct {
}

func (p *NoProxy) Dial(network string, addr string) (net.Conn, error) {
	return net.Dial(network, addr)
}

type Socks5Proxy struct {
	socks5 *socks5
}

func (p *Socks5Proxy) Dial(network string, addr string) (net.Conn, error) {
	return p.socks5.Dial(network, addr)
}

type HttpProxy struct {
	http_proxy *http_proxy
}

func (p *HttpProxy) Dial(network string, addr string) (net.Conn, error) {
	return p.http_proxy.Dial(network, addr)
}

type ShadowSocksProxy struct {
	cipher *Cipher
	server string
}

func (p *ShadowSocksProxy) Dial(network string, addr string) (net.Conn, error) {
	return Dial(addr, p.server, p.cipher.Copy())
}

func CreateParentProxy(proxy_url string) (ParentProxy, error) {
	if proxy_url != "" {
		var d net.Dialer
		url1, _ := url.Parse(proxy_url)
		var auth *Auth
		if url1.User != nil {
			auth = &Auth{}
			auth.User = url1.User.Username()
			auth.Password, _ = url1.User.Password()
		}
		switch url1.Scheme {
		case "socks5":
			s, _ := SOCKS5("tcp", url1.Host, auth, d)
			return &Socks5Proxy{
				socks5: s,
			}, nil
		case "http":
			s, _ := HTTP_PROXY("tcp", url1.Host, auth, d)
			return &HttpProxy{
				http_proxy: s,
			}, nil
		case "shadowsocks":
			cipher, _ := NewCipher(auth.User, auth.Password)
			return &ShadowSocksProxy{
				cipher: cipher,
				server: url1.Host,
			}, nil
		default:
			return nil, errors.New("unknown parent proxy type:" + url1.Scheme)
		}
	} else {
		return &NoProxy{}, nil
	}
}
