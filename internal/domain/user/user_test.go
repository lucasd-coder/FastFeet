package model_test

import (
	"reflect"
	"testing"
	"time"

	model "github.com/lucasd-coder/business-service/internal/domain/user"
)

func TestPayloadValidate(t *testing.T) {
	type fields struct {
		Name       string
		Email      string
		CPF        string
		Attributes map[string]string
		Password   string
		Authority  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "should validate model",
			fields:  fields{"", "test validate email", "", map[string]string{}, "", ""},
			wantErr: true,
		},
		{
			name:    "should validate field email",
			fields:  fields{"maria", "test validate email", "901.940.000-28", map[string]string{}, "USER", "12345678"},
			wantErr: true,
		},
		{
			name:    "should validate field cpf",
			fields:  fields{"maria", "maria@gmail.com", "test validate cpf", map[string]string{}, "USER", "12345678"},
			wantErr: true,
		},
		{
			name:    "should validate field password",
			fields:  fields{"maria", "maria2@gmail.com", "995.563.460-07", map[string]string{}, "USER", ""},
			wantErr: true,
		},
		{
			name:    "should validate field authority",
			fields:  fields{"maria", "maria3@gmail.com", "495.211.400-70", map[string]string{}, "test validate authority", "123456"},
			wantErr: true,
		},
		{
			name:    "should validate with success",
			fields:  fields{"maria", "maria4@gmail.com", "999.388.560-63", map[string]string{}, "USER", ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := &model.Payload{
				Data: model.Data{
					Name:       tt.fields.Name,
					Email:      tt.fields.Email,
					CPF:        tt.fields.CPF,
					Attributes: tt.fields.Attributes,
					Password:   tt.fields.Password,
					Authority:  tt.fields.Authority,
				},
				EventDate: time.Now(),
			}
			if err := payload.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Payload.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPayloadToRegister(t *testing.T) {
	type fields struct {
		Name       string
		Email      string
		CPF        string
		Attributes map[string]string
		Password   string
		Authority  string
	}
	tests := []struct {
		name   string
		fields fields
		want   *model.Register
	}{
		{
			name: "should model register",
			fields: fields{
				Name:       "maria",
				Email:      "manoel@gmail.com",
				CPF:        "858.416.310-71",
				Password:   "1234567",
				Authority:  "USER",
				Attributes: map[string]string{},
			},
			want: &model.Register{
				Name:      "maria",
				Username:  "manoel@gmail.com",
				Password:  "1234567",
				Authority: "USER",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := &model.Payload{
				Data: model.Data{
					Name:       tt.fields.Name,
					Email:      tt.fields.Email,
					CPF:        tt.fields.CPF,
					Attributes: tt.fields.Attributes,
					Password:   tt.fields.Password,
					Authority:  tt.fields.Authority,
				},
				EventDate: time.Now(),
			}
			if got := payload.ToRegister(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Payload.ToRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPayloadToUser(t *testing.T) {
	type fields struct {
		Name       string
		Email      string
		CPF        string
		Attributes map[string]string
		Password   string
		Authority  string
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *model.User
	}{
		{
			name: "should model register",
			fields: fields{
				Name:       "maria",
				Email:      "jose@gmail.com",
				CPF:        "415.746.480-04",
				Password:   "1234567",
				Authority:  "USER",
				Attributes: map[string]string{},
			},
			args: args{
				userID: "0c05421e-eb2e-43fd-add1-98690e0432ba",
			},
			want: &model.User{
				Name:       "maria",
				Email:      "jose@gmail.com",
				UserID:     "0c05421e-eb2e-43fd-add1-98690e0432ba",
				CPF:        "415.746.480-04",
				Attributes: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := &model.Payload{
				Data: model.Data{
					Name:       tt.fields.Name,
					Email:      tt.fields.Email,
					CPF:        tt.fields.CPF,
					Attributes: tt.fields.Attributes,
					Password:   tt.fields.Password,
					Authority:  tt.fields.Authority,
				},
				EventDate: time.Now(),
			}
			if got := payload.ToUser(tt.args.userID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Payload.ToUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
