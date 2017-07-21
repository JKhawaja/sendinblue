package sib

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

/* Request Types */

type AggregateReport struct {
	Aggregate  int    `json:"aggregate"` // 0 or 1, 0 means no aggregation (stats per day)
	Start_date string `json:"start_date"`
	End_date   string `json:"end_date"`
	Days       int    `json:"days"`
	Tag        string `json:"tag"`
}

// API Docs: https://apidocs.sendinblue.com/tutorial-sending-transactional-email/
type Email struct {
	To           map[string]string `json:"to"`
	Subject      string            `json:"subject"`
	From         [2]string         `json:"from"`
	HTML         string            `json:"html"`
	Text         string            `json:"text"`
	CC           map[string]string `json:"cc"`
	Bcc          map[string]string `json:"bcc"`
	ReplyTo      map[string]string `json:"replyto"`
	Attachment   map[string]string `json:"attachment"` // must be URL
	Headers      map[string]string `json:"headers"`
	Inline_image map[string]string `json:"inline_image"`
}

type EmailOptions struct {
	Cc             string // multiple addresses, delimiter = pipe
	Bcc            string // multiple addresses, delimiter = pipe
	ReplyTo        string
	Attr           map[string]string
	Attachment_url string
	Attachment     map[string]string
	Headers        map[string]string
}

type DeleteBouncesRequest struct {
	Start_date string `json:"start_date"`
	End_date   string `json:"end_date"`
	Email      string `json:"email"`
}

// API Docs: https://apidocs.sendinblue.com/template/
type Template struct {
	From_name      string `json:"from_name"`
	Template_name  string `json:"template_name"` // Mandatory
	Bat            string `json:"bat"`
	Html_content   string `json:"html_content"` // Mandatory (if no html_url)
	Html_url       string `json:"html_url"`     // Mandatory (if no html_content)
	Subject        string `json:"subject"`      // Mandatory
	From_email     string `json:"from_email"`   // Mandatory
	Reply_to       string `json:"reply_to"`
	To_field       string `json:"to_field"`
	Status         int    `json:"status"` // 0 (inactive -- default) or 1 (active)
	Attachment_url string `json:"attachment_url"`
}

type TemplateEmail struct {
	To             string            `json:"to"`  // multiple addresses, delimiter = pipe
	Cc             string            `json:"cc"`  // multiple addresses, delimiter = pipe
	Bcc            string            `json:"bcc"` // multiple addresses, delimiter = pipe
	ReplyTo        string            `json:"replyto"`
	Attr           map[string]string `json:"attr"`
	Attachment_url string            `json:"attachment_url"`
	Attachment     map[string]string `json:"attachment"`
	Headers        map[string]string `json:"headers"`
}

type TemplateList struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	Page       int    `json:"page"`
	Page_limit int    `json:"page_limit"`
}

/* Response Types*/

type AggregateData struct {
	Date          string `json:"date"`
	Tag           string `json:"tag"`
	Requests      int    `json:"requests"`
	Delivered     int    `json:"delivered"`
	Bounces       int    `json:"bounces"`
	Clicks        int    `json:"clicks"`
	Unique_clicks int    `json:"unique_clicks"`
	Opens         int    `json:"opens"`
	Unique_opens  int    `json:"unique_opens"`
	SpamReports   int    `json:"spamreports"`
	Blocked       int    `json:"blocked"`
	Invalid       int    `json:"invalid"`
}

type AggregateResponse struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    []AggregateData `json:"data"`
}

type CampaignData struct {
	ID            int    `json:"id"`
	Campaign_name string `json:"campaign_name"`
	Subject       string `json:"subject"`
	Bat_sent      string `json:"bat_sent"`
	Type          string `json:"type"`
	Html_content  string `json:"html_content"`
	Entered       string `json:"entered"`
	Modified      string `json:"modified"`
	Templ_status  string `json:"templ_status"`
	From_name     string `json:"from_name"`
	From_email    string `json:"from_email"`
	Reply_to      string `json:"reply_to"`
	To_field      string `json:"to_field"`
}

type CampaignResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Data    []CampaignData `json:"data"`
}

type EmailData struct {
	Message_id string `json:"message-id"`
}

type EmailResponse struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    EmailData `json:"data"`
}

type EmailResponseError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    []EmailData `json:"data"`
}

type TemplateData struct {
	ID int `json:"id"`
}

type TemplateResponse struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Data    TemplateData `json:"data"`
}

type TemplateListData struct {
	Campaign_records       []CampaignData `json:"campaign_records"`
	Page                   int            `json:"page"`
	Page_limit             int            `json:"page_limit"`
	Total_campaign_records int            `json:"total_campaign_records"`
}

type TemplateListResponse struct {
	Code    string           `json:"code"`
	Message string           `json:"message"`
	Data    TemplateListData `json:"data"`
}

func NewEmail() *Email {

	var from [2]string
	attachment := make(map[string]string)
	bcc := make(map[string]string)
	cc := make(map[string]string)
	headers := make(map[string]string)
	to := make(map[string]string)
	inline_image := make(map[string]string)

	return &Email{
		Attachment:   attachment,
		Bcc:          bcc,
		CC:           cc,
		From:         from,
		Headers:      headers,
		To:           to,
		Inline_image: inline_image,
	}
}

func NewEmailOptions(replyto, attachment_url string, cc, bcc []string) *EmailOptions {

	attr := make(map[string]string)
	attachment := make(map[string]string)
	headers := make(map[string]string)
	ccString := strings.Join(cc, "|")
	bccString := strings.Join(bcc, "|")

	return &EmailOptions{
		Cc:             ccString,
		Bcc:            bccString,
		ReplyTo:        replyto,
		Attr:           attr,
		Attachment_url: attachment_url,
		Attachment:     attachment,
		Headers:        headers,
	}
}

// AddImage() Method ...
// returns the filename, which can (and should) be used as a variable in HTML
// < img src="{{{filename}}}" alt="image" border="0" >
func (e *Email) AddImage(f *os.File) string {
	filename := f.Name()

	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	fileReader := bufio.NewReader(f)
	fileReader.Read(buffer)

	e.Inline_image[filename] = base64.StdEncoding.EncodeToString(buffer)

	return filename
}

func (e *EmailOptions) AddAttachment(f *os.File) error {

	filename := f.Name()

	if _, ok := e.Attachment[filename]; ok == true {
		err := fmt.Errorf("That file is already attached to the email.")
		return err
	}

	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	fileReader := bufio.NewReader(f)
	fileReader.Read(buffer)

	e.Attachment[filename] = base64.StdEncoding.EncodeToString(buffer)

	return nil
}
