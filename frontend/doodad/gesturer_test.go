package doodad

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	type testcase struct {
		name      string
		callbacks []func(*string) func() error
		expected  string
	}

	testcases := []testcase{
		{
			name: "Callback should be invoked",
			callbacks: []func(*string) func() error{
				func(s *string) func() error {
					return func() error {
						*s += "1"
						return nil
					}
				},
			},
			expected: "1",
		},
		{
			name: "Callbacks should be invoked in reverse order",
			callbacks: []func(*string) func() error{
				func(s *string) func() error {
					return func() error {
						*s += "1"
						return nil
					}
				},
				func(s *string) func() error {
					return func() error {
						*s += "2"
						return nil
					}
				},
			},
			expected: "21",
		},
		{
			name: "Stop propagation should stop invoking callbacks",
			callbacks: []func(*string) func() error{
				func(s *string) func() error {
					return func() error {
						*s += "1"
						return nil
					}
				},
				func(s *string) func() error {
					return func() error {
						*s += "2"
						return ErrStopPropagation
					}
				},
				func(s *string) func() error {
					return func() error {
						*s += "3"
						return nil
					}
				},
			},
			expected: "32",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cbRegistry := &CallbackRegistry[func() error]{}
			var stack = ""

			for _, cb := range tc.callbacks {
				cbRegistry.Add(cb(&stack))
			}

			cbRegistry.InvokeEndToStart(func(f func() error) error {
				return f()
			})

			assert.Equal(t, tc.expected, stack)
		})
	}

}
