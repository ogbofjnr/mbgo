// Copyright (c) 2018 Senseye Ltd. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in the LICENSE file.

package rest_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/ogbofjnr/mbgo/internal/assert"
	"github.com/ogbofjnr/mbgo/internal/rest"
)

func TestClient_NewRequest(t *testing.T) {
	cases := []struct {
		// general
		Description string

		// inputs
		Root   *url.URL
		Method string
		Path   string
		Body   io.Reader
		Query  url.Values

		// output expectations
		AssertFunc func(*testing.T, *http.Request, error)
		Request    *http.Request
		Err        error
	}{
		{
			Description: "should return an error if the provided request method is invalid",
			Root:        &url.URL{},
			Method:      "bad method",
			AssertFunc: func(t *testing.T, _ *http.Request, err error) {
				assert.Equals(t, errors.New(`net/http: invalid method "bad method"`), err)
			},
		},
		{
			Description: "should construct the URL based on provided root URL, path and query parameters",
			Root: &url.URL{
				Scheme: "http",
				Host:   net.JoinHostPort("localhost", "2525"),
			},
			Method: http.MethodGet,
			Path:   "foo",
			Query: url.Values{
				"replayable": []string{"true"},
			},
			AssertFunc: func(t *testing.T, actual *http.Request, err error) {
				assert.Ok(t, err)
				expected := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme:   "http",
						Host:     net.JoinHostPort("localhost", "2525"),
						Path:     "/foo",
						RawQuery: "replayable=true",
					},
					Host:       net.JoinHostPort("localhost", "2525"),
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     http.Header{"Accept": []string{"application/json"}},
				}
				assert.Equals(t, expected.WithContext(context.Background()), actual)
			},
		},
		{
			Description: "should only set the 'Accept' header if method is GET",
			Root:        &url.URL{},
			Method:      http.MethodGet,
			AssertFunc: func(t *testing.T, actual *http.Request, err error) {
				assert.Ok(t, err)
				expected := &http.Request{
					Method:     http.MethodGet,
					URL:        &url.URL{},
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     http.Header{"Accept": []string{"application/json"}},
				}
				assert.Equals(t, expected.WithContext(context.Background()), actual)
			},
		},
		{
			Description: "should only set the 'Accept' header if method is DELETE",
			Root:        &url.URL{},
			Method:      http.MethodDelete,
			AssertFunc: func(t *testing.T, actual *http.Request, err error) {
				assert.Ok(t, err)
				expected := &http.Request{
					Method:     http.MethodDelete,
					URL:        &url.URL{},
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     http.Header{"Accept": []string{"application/json"}},
				}
				assert.Equals(t, expected.WithContext(context.Background()), actual)
			},
		},
		{
			Description: "should set both the 'Accept' and 'Content-Type' headers if method is POST",
			Root:        &url.URL{},
			Method:      http.MethodPost,
			AssertFunc: func(t *testing.T, actual *http.Request, err error) {
				assert.Ok(t, err)
				expected := &http.Request{
					Method:     http.MethodPost,
					URL:        &url.URL{},
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header: http.Header{
						"Accept":       []string{"application/json"},
						"Content-Type": []string{"application/json"},
					},
				}
				assert.Equals(t, expected.WithContext(context.Background()), actual)
			},
		},
		{
			Description: "should set both the 'Accept' and 'Content-Type' headers if method is PUT",
			Root:        &url.URL{},
			Method:      http.MethodPut,
			AssertFunc: func(t *testing.T, actual *http.Request, err error) {
				assert.Ok(t, err)
				expected := &http.Request{
					Method:     http.MethodPut,
					URL:        &url.URL{},
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header: http.Header{
						"Accept":       []string{"application/json"},
						"Content-Type": []string{"application/json"},
					},
				}
				assert.Equals(t, expected.WithContext(context.Background()), actual)
			},
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.Description, func(t *testing.T) {
			t.Parallel()

			cli := rest.NewClient(nil, c.Root)
			req, err := cli.NewRequest(context.Background(), c.Method, c.Path, c.Body, c.Query)
			c.AssertFunc(t, req, err)
		})
	}
}

type testDTO struct {
	Test bool   `json:"test"`
	Foo  string `json:"foo"`
}

func TestClient_DecodeResponseBody(t *testing.T) {
	cases := []struct {
		// general
		Description string

		// inputs
		Body  io.ReadCloser
		Value interface{}

		// output expectations
		Expected interface{}
		Err      error
	}{
		{
			Description: "should return an error if the JSON cannot be decoded into the value pointer",
			Body:        ioutil.NopCloser(strings.NewReader(`"foo"`)),
			Value:       &testDTO{},
			Expected:    &testDTO{},
			Err: &json.UnmarshalTypeError{
				Offset: 5, // 5 bytes read before first full JSON value
				Value:  "string",
				Type:   reflect.TypeOf(testDTO{}),
			},
		},
		{
			Description: "should unmarshal the expected JSON into value pointer when valid",
			Body:        ioutil.NopCloser(strings.NewReader(`{"test":true,"foo":"bar"}`)),
			Value:       &testDTO{},
			Expected: &testDTO{
				Test: true,
				Foo:  "bar",
			},
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.Description, func(t *testing.T) {
			t.Parallel()

			cli := rest.NewClient(nil, nil)
			err := cli.DecodeResponseBody(c.Body, c.Value)
			if c.Err != nil {
				assert.Equals(t, c.Err, err)
			} else {
				assert.Ok(t, err)
			}
			assert.Equals(t, c.Expected, c.Value)
		})
	}
}
