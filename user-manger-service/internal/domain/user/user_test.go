package model_test

import (
	"testing"
	"time"

	model "github.com/lucasd-coder/user-manger-service/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUser_Validate(t *testing.T) {
	type fields struct {
		UserID     string
		Name       string
		Email      string
		CPF        string
		Attributes map[string]string
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should validate model",
			fields: fields{
				Name:  "",
				Email: "test validate email",
				CPF:   "",
			},
			wantErr: true,
		},
		{
			name: "should validate field email",
			fields: fields{
				Name:   "maria",
				Email:  "test validate email",
				CPF:    "901.940.000-28",
				UserID: "b6390035-5728-4d51-8d1c-2d9e049b8b77",
			},
			wantErr: true,
		},
		{
			name: "should validate field userID",
			fields: fields{
				Name:   "maria",
				Email:  "maria2@gmail.com",
				CPF:    "532.895.180-86",
				UserID: "test validate userID",
			},
			wantErr: true,
		},
		{
			name: "should validate with success",
			fields: fields{
				Name:   "maria",
				Email:  "maria3@gmail.com",
				CPF:    "901.940.000-28",
				UserID: "b6390035-5728-4d51-8d1c-2d9e049b8b77",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &model.User{
				UserID:     tt.fields.UserID,
				Name:       tt.fields.Name,
				Email:      tt.fields.Email,
				CPF:        tt.fields.CPF,
				Attributes: tt.fields.Attributes,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
			}
			if err := user.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_ValidatePattern(t *testing.T) {
	type fields struct {
		ID         primitive.ObjectID
		Name       string
		Email      string
		CPF        string
		Attributes map[string]string
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should validate model",
			fields: fields{
				Name:  "Maria#$$%%()",
				Email: "maria@gmail.com",
				CPF:   "12345678",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &model.User{
				ID:         tt.fields.ID,
				Name:       tt.fields.Name,
				Email:      tt.fields.Email,
				CPF:        tt.fields.CPF,
				Attributes: tt.fields.Attributes,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
			}
			if err := user.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_GetCreatedAt(t *testing.T) {
	type fields struct {
		ID         primitive.ObjectID
		Name       string
		Email      string
		CPF        string
		Attributes map[string]string
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get field createAt",
			fields: fields{
				Name:      "Maria",
				Email:     "maria@gmail.com",
				CPF:       "1234567",
				CreatedAt: time.Date(2023, time.March, 5, 10, 22, 30, 0, time.UTC),
			},
			want: "2023-03-05T10:22:30Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &model.User{
				ID:         tt.fields.ID,
				Name:       tt.fields.Name,
				Email:      tt.fields.Email,
				CPF:        tt.fields.CPF,
				Attributes: tt.fields.Attributes,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
			}
			if got := user.GetCreatedAt(); got != tt.want {
				t.Errorf("User.GetCreatedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
