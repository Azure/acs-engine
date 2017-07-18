package statsd

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"
)

const (
	testAddr = ":0"
	testKey  = "test_key"
)

var testDate = time.Date(2015, 10, 22, 16, 53, 0, 0, time.UTC)

func TestCount(t *testing.T) {
	testOutput(t, "test_key:5|c", func(c *Client) {
		c.Count(testKey, 5)
	})
}

func TestIncrement(t *testing.T) {
	testOutput(t, "test_key:1|c", func(c *Client) {
		c.Increment(testKey)
	})
}

func TestGauge(t *testing.T) {
	testOutput(t, "test_key:5|g\ntest_key:0|g\ntest_key:-10|g", func(c *Client) {
		c.Gauge(testKey, 5)
		c.Gauge(testKey, -10)
	})
}

func TestTiming(t *testing.T) {
	testOutput(t, "test_key:6|ms", func(c *Client) {
		c.Timing(testKey, 6)
	})
}

func TestHistogram(t *testing.T) {
	testOutput(t, "test_key:17|h", func(c *Client) {
		c.Histogram(testKey, 17)
	})
}

func TestNumbers(t *testing.T) {
	testOutput(t,
		"test_key:1|g\n"+
			"test_key:1|g\n"+
			"test_key:2|g\n"+
			"test_key:2|g\n"+
			"test_key:15|g\n"+
			"test_key:15|g\n"+
			"test_key:137|g\n"+
			"test_key:137|g\n"+
			"test_key:1|g\n"+
			"test_key:1|g\n"+
			"test_key:17.6|g\n"+
			"test_key:0|g\n"+
			"test_key:-42.5|g\n"+
			"test_key:|g",
		func(c *Client) {
			c.Gauge(testKey, 1)
			c.Gauge(testKey, uint(1))
			c.Gauge(testKey, int64(2))
			c.Gauge(testKey, uint64(2))
			c.Gauge(testKey, int32(15))
			c.Gauge(testKey, uint32(15))
			c.Gauge(testKey, int16(137))
			c.Gauge(testKey, uint16(137))
			c.Gauge(testKey, int8(1))
			c.Gauge(testKey, uint8(1))
			c.Gauge(testKey, float64(17.6))
			c.Gauge(testKey, float32(-42.5))
			c.Gauge(testKey, "invalid")
		})
}

func TestNewTiming(t *testing.T) {
	i := 0
	now = func() time.Time {
		i++
		switch i {
		default:
			return testDate
		case 2:
			return testDate.Add(10 * time.Millisecond)
		case 3:
			return testDate.Add(100 * time.Millisecond)
		case 4:
			return testDate.Add(time.Second)
		}
	}
	defer func() { now = time.Now }()

	testOutput(t, "test_key:10|ms\ntest_key:1000|ms", func(c *Client) {
		timing := c.NewTiming()
		timing.Send(testKey)

		got := timing.Duration().Nanoseconds()
		want := int64(100 * time.Millisecond)
		if got != want {
			t.Errorf("Duration() = %v, want %v", got, want)
		}

		timing.Send(testKey)
	})
}

func TestUnique(t *testing.T) {
	testOutput(t, "test_key:foo|s", func(c *Client) {
		c.Unique(testKey, "foo")
	})
}

func TestMute(t *testing.T) {
	dialTimeout = func(string, string, time.Duration) (net.Conn, error) {
		t.Fatal("net.Dial should not be called")
		return nil, nil
	}
	defer func() { dialTimeout = net.DialTimeout }()

	c, err := New(Mute(true))
	if err != nil {
		t.Errorf("New() = %v", err)
	}
	c.Increment(testKey)
	c.Gauge(testKey, 1)
	c.Timing(testKey, 1)
	c.Histogram(testKey, 1)
	c.Unique(testKey, "1")
	c.Flush()
	c.Close()
}

func TestSamplingRateOK(t *testing.T) {
	testOutput(t, "test_key:3|c|@0.6\ntest_key:4|ms|@0.6", func(c *Client) {
		randFloat = func() float32 { return 0.5 }
		c.Count(testKey, 3)
		c.Timing(testKey, 4)
	}, SampleRate(0.6))
}

func TestSamplingRateKO(t *testing.T) {
	testOutput(t, "", func(c *Client) {
		randFloat = func() float32 { return 0.5 }
		c.Count(testKey, 1)
		c.Timing(testKey, 2)
	}, SampleRate(0.3))
}

func TestPrefix(t *testing.T) {
	testOutput(t, "foo.test_key:1|c", func(c *Client) {
		c.Increment(testKey)
	}, Prefix("foo"))
}

func TestNilTags(t *testing.T) {
	testOutput(t, "test_key:1|c", func(c *Client) {
		c.Increment(testKey)
	}, TagsFormat(InfluxDB), Tags())
}

func TestInfluxDBTags(t *testing.T) {
	testOutput(t, "test_key,tag1=value1,tag2=value2:1|c", func(c *Client) {
		c.Increment(testKey)
	}, TagsFormat(InfluxDB), Tags("tag1", "value1", "tag2", "value2"))
}

func TestDatadogTags(t *testing.T) {
	testOutput(t, "test_key:1|c|#tag1:value1,tag2:value2", func(c *Client) {
		c.Increment(testKey)
	}, TagsFormat(Datadog), Tags("tag1", "value1", "tag2", "value2"))
}

func TestNoTagFormat(t *testing.T) {
	testOutput(t, "test_key:1|c", func(c *Client) {
		c.Increment(testKey)
	}, Tags("tag1", "value1", "tag2", "value2"))
}

func TestOddTagsArgs(t *testing.T) {
	dialTimeout = mockDial
	defer func() { dialTimeout = net.DialTimeout }()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Tags should panic when only one argument is provided")
		}
	}()

	_, _ = New(TagsFormat(InfluxDB), Tags("tag1"))
	t.Fatal("A panic should occur")
}

func TestErrorHandler(t *testing.T) {
	errorCount := 0
	testClient(t, func(c *Client) {
		getBuffer(c).err = errors.New("test error")
		c.Increment(testKey)
		c.Close()
		if errorCount != 2 {
			t.Errorf("Wrong error count, got %d, want 2", errorCount)
		}
	}, ErrorHandler(func(err error) {
		if err == nil {
			t.Error("Error should not be nil")
		}
		errorCount++
	}))
}

func TestFlush(t *testing.T) {
	testClient(t, func(c *Client) {
		c.Increment(testKey)
		c.Flush()
		got := getOutput(c)
		want := "test_key:1|c"
		if got != want {
			t.Errorf("Invalid output, got %q, want %q", got, want)
		}
		c.Close()
	})
}

func TestFlushPeriod(t *testing.T) {
	testClient(t, func(c *Client) {
		c.Increment(testKey)
		time.Sleep(time.Millisecond)
		c.conn.mu.Lock()
		got := getOutput(c)
		want := "test_key:1|c"
		if got != want {
			t.Errorf("Invalid output, got %q, want %q", got, want)
		}
		c.conn.mu.Unlock()
		c.Close()
	}, FlushPeriod(time.Nanosecond))
}

func TestMaxPacketSize(t *testing.T) {
	testClient(t, func(c *Client) {
		c.Increment(testKey)
		conn := getBuffer(c)
		got := conn.buf.String()
		if got != "" {
			t.Errorf("Output should be empty, got %q", got)
		}

		c.Increment(testKey)
		got = conn.buf.String()
		want := "test_key:1|c"
		if got != want {
			t.Errorf("Invalid output, got %q, want %q", got, want)
		}
		conn.buf.Reset()
		c.Close()

		got = conn.buf.String()
		if got != want {
			t.Errorf("Invalid output, got %q, want %q", got, want)
		}
	}, MaxPacketSize(15))
}

func TestClone(t *testing.T) {
	testOutput(t, "test_key:5|c", func(c *Client) {
		c.Clone().Count(testKey, 5)
	})
}

func TestCloneInherits(t *testing.T) {
	testOutput(t, "app.test_key:5|c|@0.5|#tag1:value1", func(c *Client) {
		clone := c.Clone()
		randFloat = func() float32 { return 0.3 }
		clone.Count(testKey, 5)
		randFloat = func() float32 { return 0.8 }
		clone.Count(testKey, 5)
	},
		TagsFormat(Datadog),
		Prefix("app"),
		SampleRate(0.5),
		Tags("tag1", "value1"),
	)
}

func TestCloneFromMuted(t *testing.T) {
	testOutput(t, "", func(c *Client) {
		c.Clone().Count(testKey, 5)
	}, Mute(true))
}

func TestMuteClone(t *testing.T) {
	testOutput(t, "", func(c *Client) {
		c.Clone(Mute(true)).Count(testKey, 5)
	})
}

func TestClonePrefix(t *testing.T) {
	testOutput(t, "app.http.test_key:5|c", func(c *Client) {
		c.Clone(Prefix("http")).Count(testKey, 5)
	}, Prefix("app"))
}

func TestCloneRate(t *testing.T) {
	testOutput(t, "", func(c *Client) {
		randFloat = func() float32 { return 0.8 }
		c.Clone(SampleRate(0.5)).Count(testKey, 5)
	})
}

func TestCloneInfluxDBTags(t *testing.T) {
	testOutput(t, "test_key,tag1=value1,tag2=value2:5|c", func(c *Client) {
		clone := c.Clone(Tags("tag1", "value3", "tag2", "value2"))
		clone.Count(testKey, 5)
	}, TagsFormat(InfluxDB), Tags("tag1", "value1"))
}

func TestCloneDatadogTags(t *testing.T) {
	testOutput(t, "test_key:5|c|#tag1:value1,tag2:value2", func(c *Client) {
		clone := c.Clone(Tags("tag1", "value3", "tag2", "value2"))
		clone.Count(testKey, 5)
	}, TagsFormat(Datadog), Tags("tag1", "value1"))
}

func TestDialError(t *testing.T) {
	dialTimeout = func(string, string, time.Duration) (net.Conn, error) {
		return nil, errors.New("")
	}
	defer func() { dialTimeout = net.DialTimeout }()

	c, err := New()
	if c == nil || !c.muted {
		t.Error("New() did not return a muted client")
	}
	if err == nil {
		t.Error("New() did not return an error")
	}
}

func TestConcurrency(t *testing.T) {
	testOutput(t, "test_key:1|c\ntest_key:1|c\ntest_key:1|c", func(c *Client) {
		var wg sync.WaitGroup
		wg.Add(1)
		c.Increment(testKey)
		go func() {
			c.Increment(testKey)
			wg.Done()
		}()
		c.Increment(testKey)
		wg.Wait()
	})
}

func TestUDPNotListening(t *testing.T) {
	dialTimeout = mockUDPClosed
	defer func() { dialTimeout = net.DialTimeout }()

	c, err := New()
	if c == nil || !c.muted {
		t.Error("New() did not return a muted client")
	}
	if err == nil {
		t.Error("New should return an error")
	}
}

type mockClosedUDPConn struct {
	i int
	net.Conn
}

func (c *mockClosedUDPConn) Write(p []byte) (int, error) {
	c.i++
	if c.i == 2 {
		return 0, errors.New("test error")
	}
	return 0, nil
}

func (c *mockClosedUDPConn) Close() error {
	return nil
}

func mockUDPClosed(string, string, time.Duration) (net.Conn, error) {
	return &mockClosedUDPConn{}, nil
}

func testClient(t *testing.T, f func(*Client), options ...Option) {
	dialTimeout = mockDial
	defer func() { dialTimeout = net.DialTimeout }()

	options = append([]Option{
		FlushPeriod(0),
		ErrorHandler(expectNoError(t)),
	}, options...)
	c, err := New(options...)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	f(c)
}

func testOutput(t *testing.T, want string, f func(*Client), options ...Option) {
	testClient(t, func(c *Client) {
		f(c)
		c.Close()

		got := getOutput(c)
		if got != want {
			t.Errorf("Invalid output, got:\n%q\nwant:\n%q", got, want)
		}
	}, options...)
}

func expectNoError(t *testing.T) func(error) {
	return func(err error) {
		t.Errorf("ErrorHandler should not receive an error: %v", err)
	}
}

type testBuffer struct {
	buf bytes.Buffer
	err error
	net.Conn
}

func (c *testBuffer) Write(p []byte) (int, error) {
	if c.err != nil {
		return 0, c.err
	}
	return c.buf.Write(p)
}

func (c *testBuffer) Close() error {
	return c.err
}

func getBuffer(c *Client) *testBuffer {
	if mock, ok := c.conn.w.(*testBuffer); ok {
		return mock
	}
	return nil
}

func getOutput(c *Client) string {
	if c.conn.w == nil {
		return ""
	}
	return getBuffer(c).buf.String()
}

func mockDial(string, string, time.Duration) (net.Conn, error) {
	return &testBuffer{}, nil
}

func TestUDP(t *testing.T) {
	testNetwork(t, "udp")
}

func TestTCP(t *testing.T) {
	testNetwork(t, "tcp")
}

func testNetwork(t *testing.T, network string) {
	received := make(chan bool)
	server := newServer(t, network, testAddr, func(p []byte) {
		s := string(p)
		if s != "test_key:1|c" {
			t.Errorf("invalid output: %q", s)
		}
		received <- true
	})
	defer server.Close()

	c, err := New(
		Address(server.addr),
		Network(network),
		ErrorHandler(expectNoError(t)),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	c.Increment(testKey)
	c.Close()
	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("server received nothing after 100ms")
	case <-received:
	}
}

type server struct {
	t      testing.TB
	addr   string
	closer io.Closer
	closed chan bool
}

func newServer(t testing.TB, network, addr string, f func([]byte)) *server {
	s := &server{t: t, closed: make(chan bool)}
	switch network {
	case "udp":
		laddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
			t.Fatal(err)
		}
		s.closer = conn
		s.addr = conn.LocalAddr().String()
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					s.closed <- true
					return
				}
				if n > 0 {
					f(buf[:n])
				}
			}
		}()
	case "tcp":
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		s.closer = ln
		s.addr = ln.Addr().String()
		go func() {
			for {
				conn, err := ln.Accept()
				if err != nil {
					s.closed <- true
					return
				}
				p, err := ioutil.ReadAll(conn)
				if err != nil {
					t.Fatal(err)
				}
				if err := conn.Close(); err != nil {
					t.Fatal(err)
				}
				f(p)
			}
		}()
	default:
		t.Fatalf("Invalid network: %q", network)
	}

	return s
}

func (s *server) Close() {
	if err := s.closer.Close(); err != nil {
		s.t.Error(err)
	}
	<-s.closed
}

func Benchmark(b *testing.B) {
	serv := newServer(b, "udp", testAddr, func([]byte) {})
	c, err := New(Address(serv.addr), FlushPeriod(0))
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		c.Increment(testKey)
		c.Count(testKey, i)
		c.Gauge(testKey, i)
		c.Timing(testKey, i)
		c.NewTiming().Send(testKey)
	}
	c.Close()
	serv.Close()
}
