//go:build task_template

package main

type Message map[string]string

func (m *Message) UnmarshalJSON(bytes []byte) error {
	return nil
}
