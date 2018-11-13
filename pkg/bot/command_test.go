package bot

import (
    "testing"
)

var commandProvider = []struct {
    in  string
    out string
}{
    {"/monthly_report some additional text", "/monthly_report"},
    {"/weekly_report some additional text", "/weekly_report"},
    {"/daily_report some additional text", "/daily_report"},
    {"/statistics some additional text", "/statistics"},
    {"buy products $1500 in market", "/budget"},
    {"/unknown_command", "/unknown_command"},
}

func TestDefineCommand(t *testing.T) {
    for _, data := range commandProvider {
        t.Run(data.in, func(t *testing.T) {
            command, err := defineCommand(data.in)
            if err != nil {
                t.Error(err.Error())
            }

            if data.out != command {
                t.Errorf("Command defined incorrectly. Expected: \"%s\", Received: \"%s\"", data.out, command)
            }
        })
    }
}
