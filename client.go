// Package sib is a library for constructing a SendInBlue API2.0 client.
package sib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// The Client type is the primary type in the package.
type Client struct {
	apiKey string
	Client *http.Client
}

// NewClient takes a private SendInBlue API key
// and constructs a Client Object that can be used
// to talk to the SendInBlue API via the Client methods.
func NewClient(apiKey string) (*Client, error) {

	emptyClient := &Client{}

	if apiKey == "" {
		err := fmt.Errorf("Error: Please provide a SendInBlue API Key.")
		return emptyClient, err
	}

	return &Client{
		apiKey: apiKey,
		Client: &http.Client{ // could consider using fasthttp client -- but would introduce vendor dep
			Timeout: time.Second * 10,
		},
	}, nil
}

// AggregateReport is a Client Method for the SMTP API.
// Developers can access information about aggregate / date-wise report of the SendinBlue SMTP account using this API.
// https://apidocs.sendinblue.com/statistics/
func (c *Client) AggregateReport(a *AggregateReport) (AggregateResponse, error) {

	emptyResp := AggregateResponse{}

	body, err := json.Marshal(a)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return emptyResp, err
	}
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "https://api.sendinblue.com/v2.0/statistics", r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return emptyResp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return emptyResp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Could not recognize API response format: %+v", err)
		return emptyResp, err
	}

	var response AggregateResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		err := fmt.Errorf("Could not decode response format: %+v", err)
		return emptyResp, err
	}

	return response, nil
}

// CreateSMSCampaign ...
func (c *Client) CreateSMSCampaign(s *SMSCampaign) (SMSCampaignResponse, error) {

	emptyResp := SMSCampaignResponse{}

	body, err := json.Marshal(s)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return emptyResp, err
	}
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "https://api.sendinblue.com/v2.0/sms", r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return emptyResp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return emptyResp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Could not recognize API response format: %+v", err)
		return emptyResp, err
	}

	var response SMSCampaignResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		err := fmt.Errorf("Could not decode response format: %+v", err)
		return emptyResp, err
	}

	return response, nil
}

// CreateTemplate ...
func (c *Client) CreateTemplate(t *Template) (TemplateResponse, error) {

	emptyResp := TemplateResponse{}

	body, err := json.Marshal(t)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return emptyResp, err
	}
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "https://api.sendinblue.com/v2.0/template", r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return emptyResp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return emptyResp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Could not recognize API response format: %+v", err)
		return emptyResp, err
	}

	var response TemplateResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		err := fmt.Errorf("Could not decode response format: %+v", err)
		return emptyResp, err
	}

	return response, nil
}

// DeleteBouncedEmails ...
// Start and End dates must be in YYYY-MM-DD format
// Start date must be before end date, and end date must be after start date
func (c *Client) DeleteBouncedEmails(start, end, email string) error {

	request := DeleteBouncesRequest{
		Start_date: start,
		End_date:   end,
		Email:      email,
	}

	body, err := json.Marshal(request)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return err
	}
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "https://api.sendinblue.com/v2.0/bounces", r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("Request error: ", resp.Status)
		return err
	}

	return nil
}

// GetTemplate ...
func (c *Client) GetTemplate(template_id int) (CampaignResponse, error) {

	emptyResp := CampaignResponse{}

	url := fmt.Sprintf("https://api.sendinblue.com/v2.0/campaign/%v/detailsv2", template_id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return emptyResp, err
	}
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return emptyResp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Could not recognize API response format: %+v", err)
		return emptyResp, err
	}

	var response CampaignResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		err := fmt.Errorf("Could not decode response format: %+v", err)
		return emptyResp, err
	}

	return response, nil
}

// ListTemplates ...
func (c *Client) ListTemplates(t *TemplateList) (TemplateListResponse, error) {

	emptyResp := TemplateListResponse{}

	body, err := json.Marshal(t)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return emptyResp, err
	}
	r := bytes.NewReader(body)

	url := fmt.Sprintf("https://api.sendinblue.com/v2.0/campaign/detailsv2")
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return emptyResp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return emptyResp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Could not recognize API response format: %+v", err)
		return emptyResp, err
	}

	var response TemplateListResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		err := fmt.Errorf("Could not decode response format: %+v", err)
		return emptyResp, err
	}

	return response, nil
}

// SendEmail ...
func (c *Client) SendEmail(e *Email) (EmailResponse, error) {

	emptyResp := EmailResponse{}

	body, err := json.Marshal(e)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return emptyResp, err
	}
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "https://api.sendinblue.com/v2.0/email", r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return emptyResp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return emptyResp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Could not recognize API response format: %+v", err)
		return emptyResp, err
	}

	var response EmailResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		err := fmt.Errorf("Could not decode response format: %+v", err)
		return emptyResp, err
	}

	return response, nil
}

// SendSMS ...
func (c *Client) SendSMS(s *SMSRequest) (SMSResponse, error) {

	emptyResp := SMSResponse{}

	body, err := json.Marshal(s)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return emptyResp, err
	}
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "https://api.sendinblue.com/v2.0/sms", r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return emptyResp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return emptyResp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Could not recognize API response format: %+v", err)
		return emptyResp, err
	}

	var response SMSResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		err := fmt.Errorf("Could not decode response format: %+v", err)
		return emptyResp, err
	}

	return response, nil
}

// SendTemplateEmail ...
func (c *Client) SendTemplateEmail(id int, to []string, e *EmailOptions) (EmailResponse, error) {

	toString := strings.Join(to, "|")

	email := TemplateEmail{}
	email.To = toString

	if e != nil {
		email.Cc = e.Cc
		email.Bcc = e.Bcc
		email.ReplyTo = e.ReplyTo
		email.Attr = e.Attr
		email.Attachment_url = e.Attachment_url
		email.Attachment = e.Attachment
		email.Headers = e.Headers
	}

	emptyResp := EmailResponse{}

	body, err := json.Marshal(email)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return emptyResp, err
	}
	r := bytes.NewReader(body)

	url := fmt.Sprintf("https://api.sendinblue.com/v2.0/template/%v", id)
	req, err := http.NewRequest("PUT", url, r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return emptyResp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return emptyResp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Could not recognize API response format: %+v", err)
		return emptyResp, err
	}

	var response EmailResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		err := fmt.Errorf("Could not decode response format: %+v", err)
		return emptyResp, err
	}

	return response, nil
}

// SMSCampaignTest ...
func (c *Client) SMSCampaignTest(id int, to string) (SMSResponse, error) {

	request := SMSTest{
		To: to,
	}

	emptyResp := SMSResponse{}

	body, err := json.Marshal(request)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return emptyResp, err
	}
	r := bytes.NewReader(body)

	url := fmt.Sprintf("https://api.sendinblue.com/v2.0/sms/%v", id)
	req, err := http.NewRequest("GET", url, r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return emptyResp, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return emptyResp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Could not recognize API response format: %+v", err)
		return emptyResp, err
	}

	var response SMSResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		err := fmt.Errorf("Could not decode response format: %+v", err)
		return emptyResp, err
	}

	return response, nil
}

// UpdateSMSCampaign ...
func (c *Client) UpdateSMSCampaign(id int, s *SMSCampaign) error {

	body, err := json.Marshal(s)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return err
	}
	r := bytes.NewReader(body)

	url := fmt.Sprintf("https://api.sendinblue.com/v2.0/sms/%v", id)
	req, err := http.NewRequest("PUT", url, r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("Request error: ", resp.Status)
		return err
	}

	return nil
}

// UpdateTemplate ...
func (c *Client) UpdateTemplate(id int, t *Template) error {

	body, err := json.Marshal(t)
	if err != nil {
		err = fmt.Errorf("Could not marshal JSON: %+v", err)
		return err
	}
	r := bytes.NewReader(body)

	url := fmt.Sprintf("https://api.sendinblue.com/v2.0/template/%v", id)
	req, err := http.NewRequest("PUT", url, r)
	if err != nil {
		err := fmt.Errorf("Could not create http request: %+v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", c.apiKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		err := fmt.Errorf("Could not send http request: %+v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("Request error: ", resp.Status)
		return err
	}

	return nil
}
