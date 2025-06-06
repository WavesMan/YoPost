package protocol_test

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/protocol"
	"github.com/stretchr/testify/assert"
)

type testConn struct {
	buf bytes.Buffer
}

func (c *testConn) Read(b []byte) (int, error)         { return 0, nil }
func (c *testConn) Write(b []byte) (int, error)        { return c.buf.Write(b) }
func (c *testConn) Close() error                       { return nil }
func (c *testConn) LocalAddr() net.Addr                { return nil }
func (c *testConn) RemoteAddr() net.Addr               { return nil }
func (c *testConn) SetDeadline(t time.Time) error      { return nil }
func (c *testConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *testConn) SetWriteDeadline(t time.Time) error { return nil }

func newTestConn() *testConn {
	return &testConn{buf: bytes.Buffer{}}
}

var testConfig = &config.Config{
	SMTP: config.SMTPConfig{Port: 2525},
}

type mockMailCore struct {
	cfg *config.Config
}

func (m mockMailCore) ValidateUser(email string) bool {
	return email == "valid@test.com"
}

func (m mockMailCore) GetConfig() *config.Config {
	return m.cfg
}

func TestSMTPServer_HandleEHLO(t *testing.T) {
	mockCore := mockMailCore{cfg: testConfig}
	srv := protocol.NewSMTPServer(testConfig, mockCore)

	conn := newTestConn()
	defer conn.Close()

	err := srv.HandleCommand(conn, "EHLO example.com")
	assert.NoError(t, err)
	assert.Contains(t, conn.buf.String(), "250-HELO")
}

func TestSMTPServer_InvalidCommand(t *testing.T) {
	mockCore := mockMailCore{cfg: testConfig}
	srv := protocol.NewSMTPServer(testConfig, mockCore)

	conn := newTestConn()
	defer conn.Close()

	err := srv.HandleCommand(conn, "INVALID")
	assert.NoError(t, err)
	assert.Contains(t, conn.buf.String(), "500 Unknown command")
}
