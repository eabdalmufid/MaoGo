package commands

import (
	"encoding/json"
	"mao/src/libs"

	"github.com/robertkrimen/otto"
)

func init() {
	libs.NewCommands(&libs.ICommand{
		Name:     `;`,
		As:       []string{";"},
		Tags:     "owner",
		IsPrefix: false,
		Exec: func(client *libs.NewClientImpl, m *libs.IMessage) {
			vm := otto.New()
			vm.Set("M", m)

			h, err := vm.Run(m.Querry)
			if err != nil {
				m.Reply(err.Error())
				return
			}

			if h.IsObject() {
				var data interface{}
				if r, err := vm.Run("JSON.stringify(" + m.Querry + ")"); err != nil {
					json.Unmarshal([]byte(r.String()), &data)
				}
				pe, _ := json.MarshalIndent(data, "", "  ")
				m.Reply(string(pe))
			} else {
				m.Reply(h.String())
			}
		},
	})
}
