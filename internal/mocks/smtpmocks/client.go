// Code generated by mockery v1.0.0. DO NOT EDIT.

package smtpmocks

import io "io"
import mock "github.com/stretchr/testify/mock"
import smtp "net/smtp"

import tls "crypto/tls"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Auth provides a mock function with given fields: a
func (_m *Client) Auth(a smtp.Auth) error {
	ret := _m.Called(a)

	var r0 error
	if rf, ok := ret.Get(0).(func(smtp.Auth) error); ok {
		r0 = rf(a)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *Client) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Data provides a mock function with given fields:
func (_m *Client) Data() (io.WriteCloser, error) {
	ret := _m.Called()

	var r0 io.WriteCloser
	if rf, ok := ret.Get(0).(func() io.WriteCloser); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.WriteCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Extension provides a mock function with given fields: ext
func (_m *Client) Extension(ext string) (bool, string) {
	ret := _m.Called(ext)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(ext)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(ext)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// Hello provides a mock function with given fields: localName
func (_m *Client) Hello(localName string) error {
	ret := _m.Called(localName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(localName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Mail provides a mock function with given fields: from
func (_m *Client) Mail(from string) error {
	ret := _m.Called(from)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(from)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Noop provides a mock function with given fields:
func (_m *Client) Noop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Quit provides a mock function with given fields:
func (_m *Client) Quit() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Rcpt provides a mock function with given fields: to
func (_m *Client) Rcpt(to string) error {
	ret := _m.Called(to)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(to)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Reset provides a mock function with given fields:
func (_m *Client) Reset() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StartTLS provides a mock function with given fields: config
func (_m *Client) StartTLS(config *tls.Config) error {
	ret := _m.Called(config)

	var r0 error
	if rf, ok := ret.Get(0).(func(*tls.Config) error); ok {
		r0 = rf(config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TLSConnectionState provides a mock function with given fields:
func (_m *Client) TLSConnectionState() (tls.ConnectionState, bool) {
	ret := _m.Called()

	var r0 tls.ConnectionState
	if rf, ok := ret.Get(0).(func() tls.ConnectionState); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(tls.ConnectionState)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Verify provides a mock function with given fields: addr
func (_m *Client) Verify(addr string) error {
	ret := _m.Called(addr)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}