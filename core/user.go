package core

import (
	"context"
	"errors"
	"net/mail"
	"os/user"
	"reflect"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/google/uuid"

	"gopkg.in/AlecAivazis/survey.v1"
)

// User defines a laforge command line user and their properties
//easyjson:json
type User struct {
	ID    string `hcl:"id,label" cty:"id" json:"id,omitempty"`
	Name  string `hcl:"name,attr" cty:"name" json:"name,omitempty"`
	UUID  string `hcl:"uuid,optional" cty:"uuid" json:"uuid,omitempty"`
	Email string `hcl:"email,attr" cty:"email" json:"email,omitempty"`
}

// UserWizard runs an interactive prompt to get the user's information.
func UserWizard() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	user := User{
		ID:   u.Username,
		UUID: uuid.New().String(),
	}
	confirmed := false
	qs := []*survey.Question{
		{
			Name: "Name",
			Prompt: &survey.Input{
				Message: "Enter your name:",
				Default: u.Username,
			},
			Validate: survey.Required,
		},
		{
			Name: "Email",
			Prompt: &survey.Input{
				Message: "Enter your email address:",
			},
			Validate: func(val interface{}) error {
				value := reflect.ValueOf(val)
				_, err := mail.ParseAddress(value.String())
				return err
			},
		},
	}
	err = survey.Ask(qs, &user)
	if err != nil {
		return err
	}
	qs2 := []*survey.Question{
		{
			Name: "confirmed",
			Prompt: &survey.Confirm{
				Message: "Write configuration to ~/.laforge/global.laforge?",
			},
		},
	}
	err = survey.Ask(qs2, &confirmed)
	if !confirmed {
		return errors.New("write authorization not granted")
	}
	if err != nil {
		return err
	}

	err = CreateGlobalConfig(user)
	if err != nil {
		return err
	}
	cli.Logger.Warnf("Global configuration written to ~/.laforge/global.laforge")
	return nil
}

// CreateUserEntry ...
func (u *User) CreateUserEntry(ctx context.Context, client *ent.Client) (*ent.User, error) {
	user, err := client.User.
		Create().
		SetName(u.Name).
		SetUUID(u.UUID).
		SetEmail(u.Email).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating user: %v", err)
		return nil, err
	}

	cli.Logger.Debugf("user was created: ", user)
	return user, nil
}
