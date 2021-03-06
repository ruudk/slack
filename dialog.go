package slack

import (
	"context"
	"errors"
	"encoding/json"
)

type DialogTrigger struct {
	TriggerId  string  `json:"trigger_id"`  //Required. Must respond within 3 seconds.
	Dialog     Dialog  `json:"dialog"`      //Required.
}

type Dialog struct {
	Title           string           `json:"title"`                         //Required.
	CallbackId      string           `json:"callback_id"`                   //Required.
	Elements        []DialogElement  `json:"elements"`                      //Required.
	SubmitLabel     string           `json:"submit_label,omitempty"`        //Optional. Default value is 'Submit'
	NotifyOnCancel  bool             `json:"notify_on_cancel,omitempty"`    //Optional. Default value is false
}

type DialogElement interface {}

type DialogTextElement struct {
	Label        string     `json:"label"`                  //Required.
	Name         string     `json:"name"`                   //Required.
	Type         string     `json:"type"`                   //Required. Allowed values: "text", "textarea", "select".
	Placeholder  string     `json:"placeholder,omitempty"`  //Optional.
	Optional     bool       `json:"optional,omitempty"`     //Optional. Default value is false
	Value        string     `json:"value,omitempty"`        //Optional.
	MaxLength    int        `json:"max_length,omitempty"`   //Optional.
	MinLength    int        `json:"min_length,omitempty"`   //Optional,. Default value is 0
	Hint         string     `json:"hint,omitempty"`         //Optional.
	Subtype      string     `json:"subtype,omitempty"`      //Optional. Allowed values: "email", "number", "tel", "url".
}

type DialogSelectElement struct {
	Label            string                 `json:"label"`                          //Required.
	Name             string                 `json:"name"`                           //Required.
	Type             string                 `json:"type"`                           //Required. Allowed values: "text", "textarea", "select".
	Placeholder      string                 `json:"placeholder,omitempty"`          //Optional.
	Optional         bool                   `json:"optional,omitempty"`             //Optional. Default value is false
	Value            string                 `json:"value,omitempty"`                //Optional.
	DataSource       string                 `json:"data_source,omitempty"`          //Optional. Allowed values: "users", "channels", "conversations", "external".
	SelectedOptions  string                 `json:"selected_options,omitempty"`     //Optional. Default value for "external" only
	Options          []DialogElementOption  `json:"options,omitempty"`              //One of options or option_groups is required.
	OptionGroups     []DialogElementOption  `json:"option_groups,omitempty"`        //Provide up to 100 options.
}

type DialogElementOption struct {
	Label  string   `json:"label"`	//Required.
	Value  string   `json:"value"`	//Required.
}

// DialogCallback is sent from Slack when a user submits a form from within a dialog
type DialogCallback struct {
	Type         string             `json:"dialog_submission"`
	CallbackID   string             `json:"callback_id"`
	Team         Team               `json:"team"`
	Channel      Channel            `json:"channel"`
	User         User               `json:"user"`
	ActionTs     string             `json:"action_ts"`
	Token        string             `json:"token"`
	ResponseURL  string             `json:"response_url"`
	Submission   map[string]string  `json:"submission"`
}

// DialogSuggestionCallback is sent from Slack when a user types in a select field with an external data source
type DialogSuggestionCallback struct {
	Type         string             `json:"dialog_suggestion"`
	Token        string             `json:"token"`
	ActionTs     string             `json:"action_ts"`
	Team         Team               `json:"team"`
	User         User               `json:"user"`
	Channel      Channel            `json:"channel"`
	ElementName  string             `json:"name"`
	Value        string             `json:"value"`
	CallbackID   string             `json:"callback_id"`
}

// OpenDialog opens a dialog window where the triggerId originated from
func (api *Client) OpenDialog(triggerId string, dialog Dialog) (err error) {
	return api.OpenDialogContext(context.Background(), triggerId, dialog)
}

// OpenDialogContext opens a dialog window where the triggerId originated from with a custom context
func (api *Client) OpenDialogContext(ctx context.Context, triggerId string, dialog Dialog) (err error) {
	if triggerId == "" {
		return errors.New("received empty parameters")
	}

	resp := DialogTrigger{
		TriggerId:  triggerId,
		Dialog:     dialog,
	}
	jsonResp, _ := json.Marshal(resp)
	response := &SlackResponse{}
	if err := postJson(ctx, api.httpclient, "dialog.open", api.token, jsonResp, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}