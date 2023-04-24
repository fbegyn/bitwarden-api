package bitwarden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type URI struct {
	Match int    `json:"match,omitempty"`
	URI   string `json:"uri,omitempty"`
}

type Login struct {
	URIs     []URI  `json:"uris,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	TOTP     string `json:"totp,omitempty"`
}

type Field struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	Type  int    `json:"type,omitempty"`
}

type Item struct {
	ID             string   `json:"id,omitempty"`
	OrganizationID string   `json:"organizationid,omitempty"`
	CollectionIDs  []string `json:"collectionid,omitempty"`
	FolderID       string   `json:"folderid,omitempty"`
	Type           int      `json:"type,omitempty"`
	Name           string   `json:"name,omitempty"`
	Notes          string   `json:"notes,omitempty"`
	Favorite       bool     `json:"favorite,omitempty"`
	Fields         []Field  `json:"fields,omitempty"`
	Login          Login    `json:"login,omitempty"`
	Reprompt       int      `json:"reprompt,omitempty"`
}

type ItemCreateResp struct {
	Success      bool
	Data         map[string]interface{}
	RevisionDate time.Time
	DeleteDate   time.Time
}

func (bw *BitwardenClient) CreateItem(item Item) error {
	jsonItem, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal Item object: %w", err)
	}
	req, err := http.NewRequest("POST", bw.BaseURL+"/object/item", bytes.NewBuffer(jsonItem))
	if err != nil {
		return fmt.Errorf("failed to create create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := bw.Client.Do(req)
	if err != nil {
		return fmt.Errorf("create request caused an error: %w", err)
	}
	defer resp.Body.Close()
	var createResp ItemCreateResp
	json.NewDecoder(resp.Body).Decode(&createResp)
	if !createResp.Success {
		return fmt.Errorf("create operation was unsuccesful: %v", resp)
	}
	return nil
}

type ItemGetResp struct {
	Success      bool
	Data         Item
	RevisionDate time.Time
	DeleteDate   time.Time
}

func (bw *BitwardenClient) GetItem(item Item) (Item, error) {
	req, err := http.NewRequest("GET", bw.BaseURL+"/object/item/"+item.ID, nil)
	if err != nil {
		return Item{}, fmt.Errorf("failed to create get request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := bw.Client.Do(req)
	if err != nil {
		return Item{}, fmt.Errorf("get request caused an error: %w", err)
	}
	defer resp.Body.Close()
	var getResp ItemGetResp
	json.NewDecoder(resp.Body).Decode(&getResp)
	if !getResp.Success {
		return Item{}, fmt.Errorf("get operation was unsuccessful: %v", resp)
	}
	return getResp.Data, nil
}

type ItemUpdateResp struct {
	Success      bool
	Data         Item
	RevisionDate time.Time
	DeleteDate   time.Time
}

func (bw *BitwardenClient) UpdateItem(item Item) error {
	jsonItem, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal Item object: %w", err)
	}
	req, err := http.NewRequest("PUT", bw.BaseURL+"/object/item/"+item.ID, bytes.NewBuffer(jsonItem))
	if err != nil {
		return fmt.Errorf("failed to create update request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := bw.Client.Do(req)
	if err != nil {
		return fmt.Errorf("update request caused an error: %w", err)
	}
	defer resp.Body.Close()
	var updateResp ItemUpdateResp
	json.NewDecoder(resp.Body).Decode(&updateResp)
	if !updateResp.Success {
		return fmt.Errorf("delete operation was unsuccessful: %v", resp)
	}
	return nil
}

type ItemDeleteResp struct {
	Success bool
}

func (bw *BitwardenClient) DeleteItem(item Item) error {
	jsonItem, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal Item object: %w", err)
	}
	req, err := http.NewRequest("DELETE", bw.BaseURL+"/object/item/"+item.ID, bytes.NewBuffer(jsonItem))
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := bw.Client.Do(req)
	if err != nil {
		return fmt.Errorf("delete request caused an error: %w", err)
	}
	defer resp.Body.Close()
	var deleteResp ItemDeleteResp
	json.NewDecoder(resp.Body).Decode(&deleteResp)
	if !deleteResp.Success {
		return fmt.Errorf("delete operation was unsuccessful: %v", resp)
	}
	return nil
}

type UnlockResp struct {
	Success bool
	Data    map[string]interface{}
}

type LockResp struct {
	Success bool
	Data    map[string]interface{}
}

func (bw *BitwardenClient) Lock() error {
	resp, err := http.Post(bw.BaseURL+"/lock", "application/json", &bytes.Buffer{})
	if err != nil {
		return fmt.Errorf("failed to create lock request: %w", err)
	}
	defer resp.Body.Close()
	var lockResp LockResp
	json.NewDecoder(resp.Body).Decode(&lockResp)
	if !lockResp.Success {
		return fmt.Errorf("lock operation was unsuccessful: %w", err)
	}
	return nil
}

type ItemListResp struct {
	Success bool
	Data    struct {
		Object string
		Data   []Item
	}
}

// List implements the objects items list API functionality. It retturns a
// ListResponse according to the Bitwarden API.
func (bw *BitwardenClient) ListItems(search string) ([]Item, error) {
	req, err := http.NewRequest("GET", bw.BaseURL+"/list/object/items", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list request: %w", err)
	}
	q := req.URL.Query()
	q.Add("search", search)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Content-Type", "application/json")
	resp, err := bw.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list request caused an error: %w", err)
	}
	defer resp.Body.Close()
	var listResp ItemListResp
	json.NewDecoder(resp.Body).Decode(&listResp)

	if !listResp.Success {
		return nil, fmt.Errorf("list operation was unsuccesful: %w", err)
	}
	return listResp.Data.Data, nil
}
