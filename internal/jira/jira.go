package jira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL  string
	Email    string
	APIToken string
	Client   *http.Client
}

func NewClient(baseURL, email, apiToken string) *Client {
	return &Client{
		BaseURL:  strings.TrimRight(baseURL, "/"),
		Email:    email,
		APIToken: apiToken,
		Client:   &http.Client{},
	}
}

func (c *Client) request(method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.Email + ":" + c.APIToken))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("jira api error: %s %s", resp.Status, string(respBody))
	}

	return respBody, nil
}

// Minimal ADF structure
type ADF struct {
	Version int       `json:"version"`
	Type    string    `json:"type"`
	Content []ADFNode `json:"content"`
}

type ADFNode struct {
	Type    string    `json:"type"`
	Content []ADFNode `json:"content,omitempty"`
	Text    string    `json:"text,omitempty"`
}

func TextToADFMultiline(text string) ADF {
	lines := strings.Split(text, "\n")
	var content []ADFNode
	for _, line := range lines {
		// if line == "" {
		// 	continue
		// }
		// Preserve empty lines as empty paragraphs?
		// Empty paragraph needs content?
		// Let's just create a paragraph with text.
		
		node := ADFNode{
			Type: "paragraph",
		}
		if line != "" {
			node.Content = []ADFNode{
				{
					Type: "text",
					Text: line,
				},
			}
		}
		content = append(content, node)
	}
	return ADF{
		Version: 1,
		Type:    "doc",
		Content: content,
	}
}

type Project struct {
	Key string `json:"key"`
}

type IssueType struct {
	Name string `json:"name"`
}

type Priority struct {
	Name string `json:"name"`
}

func (c *Client) CreateTicket(projectKey, summary, description, issueTypeName, priorityName string) (string, error) {
	type CreateIssueFields struct {
		Project   Project   `json:"project"`
		Summary   string    `json:"summary"`
		Description ADF     `json:"description"`
		IssueType IssueType `json:"issuetype"`
		Labels    []string  `json:"labels"`
		Priority  *Priority `json:"priority,omitempty"`
	}
	
	fields := CreateIssueFields{
		Project: Project{Key: projectKey},
		Summary: summary,
		Description: TextToADFMultiline(description),
		IssueType: IssueType{Name: issueTypeName},
		Labels: []string{"debtbomb", "expired"},
	}
	
	if priorityName != "" {
		fields.Priority = &Priority{Name: priorityName}
	}

	body := map[string]interface{}{
		"fields": fields,
	}

	resp, err := c.request("POST", "/rest/api/3/issue", body)
	if err != nil {
		return "", err
	}

	var result struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	return result.Key, nil
}

func (c *Client) AddComment(issueKey, comment string) error {
	body := map[string]interface{}{
		"body": TextToADFMultiline(comment),
	}
	_, err := c.request("POST", fmt.Sprintf("/rest/api/3/issue/%s/comment", issueKey), body)
	return err
}

func (c *Client) UpdatePriority(issueKey, priorityName string) error {
	body := map[string]interface{}{
		"fields": map[string]interface{}{
			"priority": map[string]string{
				"name": priorityName,
			},
		},
	}
	_, err := c.request("PUT", fmt.Sprintf("/rest/api/3/issue/%s", issueKey), body)
	return err
}

func (c *Client) CloseTicket(issueKey string) error {
	transitionsResp, err := c.request("GET", fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueKey), nil)
	if err != nil {
		return err
	}
	
	var transitionsResult struct {
		Transitions []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"transitions"`
	}
	if err := json.Unmarshal(transitionsResp, &transitionsResult); err != nil {
		return err
	}
	
	var transitionID string
	for _, t := range transitionsResult.Transitions {
		lowerName := strings.ToLower(t.Name)
		if lowerName == "done" || lowerName == "closed" || lowerName == "resolve" || lowerName == "resolved" {
			transitionID = t.ID
			break
		}
	}
	
	if transitionID == "" {
		return fmt.Errorf("could not find close transition for issue %s", issueKey)
	}
	
	body := map[string]interface{}{
		"transition": map[string]string{
			"id": transitionID,
		},
	}
	
	_, err = c.request("POST", fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueKey), body)
	return err
}
