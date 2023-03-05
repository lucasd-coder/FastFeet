package model_test

import (
	"testing"
	"time"

	model "github.com/lucasd-coder/user-manger-service/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUser_Validate(t *testing.T) {
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
				Name:  "",
				Email: "test validate email",
				CPF:   "",
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
