package resp

import (
  "net"
)

type Conn struct {
  net.Conn
  scanner *Scanner
}

func Dial(network string, ip string) (*Conn, error) {
  conn, err := net.Dial(network, ip)
  if err != nil {
    return nil, err
  }

  return &Conn{conn, NewScanner(conn)}, nil
}

func (self *Conn) Send(request ...interface{}) (interface{}, error) {
  req, err := Marshall(request)

  if err != nil {
    return nil, err
  }

  self.Write([]byte(req))

  if self.scanner.Scan() {
    return self.scanner.Obj(), nil
  } else {
    return nil, self.scanner.Err()
  }
}
