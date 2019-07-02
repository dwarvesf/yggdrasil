package twilio

import (
	"testing"

	twilio "github.com/sfreiberg/gotwilio"
)

func TestClient_Send(t *testing.T) {
	type fields struct {
		t *twilio.Twilio
	}
	type args struct {
		from    string
		to      string
		content string
		appSid  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				t: twilio.NewTwilioClient("AC3e3cf85b4c6e9514f42cd04d2aad272c", "82bbff672390e6d3d808b64626a3573e"),
			},
			args: args{
				from:    "4703541628",
				to:      "+841212592491",
				content: "send from twilio",
				appSid:  "AC3e3cf85b4c6e9514f42cd04d2aad272a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := &Client{
				t: tt.fields.t,
			}
			if err := tw.Send(tt.args.from, tt.args.to, tt.args.content, tt.args.appSid); (err != nil) != tt.wantErr {
				t.Errorf("TwilioClient.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
