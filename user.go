package laforge

import (
	"errors"
	"os/user"

	"github.com/google/uuid"

	"gopkg.in/AlecAivazis/survey.v1"
)

// User defines a laforge command line user and their properties
type User struct {
	Name      string `hcl:"name,attr" cty:"name" json:"name,omitempty"`
	ID        string `hcl:"id,attr" cty:"id" json:"id,omitempty"`
	Email     string `hcl:"email,attr" cty:"email" json:"email,omitempty"`
	Confirmed bool   `json:"-"`
}

// UserWizard runs an interactive prompt to get the user's information.
func UserWizard() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	user := User{
		ID: uuid.New().String(),
	}
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
			Validate: survey.Required,
		},
		{
			Name: "Confirmed",
			Prompt: &survey.Confirm{
				Message: "Write configuration to ~/.laforge/global.laforge?",
			},
		},
	}
	err = survey.Ask(qs, &user)
	if err != nil {
		return err
	}
	if !user.Confirmed {
		return errors.New("write authorization not granted")
	}
	err = CreateGlobalConfig(user)
	if err != nil {
		return err
	}
	Logger.Warnf("Global configuration written to ~/.laforge/global.laforge")
	return nil
}
