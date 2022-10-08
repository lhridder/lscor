package corero

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"lscor/config"
	"net/http"
	"time"
)

type Recent struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Results []APIAttack `json:"results"`
}

type APIAttack struct {
	TenantID     int    `json:"tenantId"`
	TenantName   string `json:"tenant_name"`
	IPaddress    string `json:"ip_address"`
	Status       string `json:"status"`
	StartTime    int    `json:"start_time"`
	Duration     int    `json:"duration"`
	MaxPeak      int    `json:"max_peak"`
	ServiceLevel string `json:"service_level"`
	AttackID     string `json:"attack_id"`
	Description  string `json:"description"`
	Volume       int    `json:"volume"`
}

type Top struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Results []APIIP `json:"results"`
}

type APIIP struct {
	IPaddress   string `json:"ip_address"`
	AttackCount int    `json:"attack_count"`
	TenantName  string `json:"tenant_name"`
}

type Attack struct {
	Zone        string `json:"zone"`
	IPaddress   string `json:"ip_address"`
	Status      string `json:"status"`
	Start       int    `json:"start"`
	Duration    int    `json:"duration"`
	End         int    `json:"end"`
	Max         int    `json:"max_megabits"`
	AttackID    string `json:"attack_id"`
	Description string `json:"description"`
	Traffic     int    `json:"traffic_megabytes"`
}

type IP struct {
	IPaddress   string `json:"ip_address"`
	AttackCount int    `json:"attack_count"`
	Zone        string `json:"zone"`
}

func GetRecent(duration time.Duration, ip string) ([]Attack, error) {
	apiurl := config.GlobalConfig.Corero.URL + "api/"
	t := time.Now().Add(-duration)
	starttime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	req, err := http.NewRequest("GET", apiurl+"get_recent_attacks", nil)
	if err != nil {
		return []Attack{}, fmt.Errorf("failed to format get: %s", err)
	}

	q := req.URL.Query()
	q.Add("start", starttime)
	if ip != "" {
		q.Add("dip", ip)
	}
	req.URL.RawQuery = q.Encode()
	req.AddCookie(Cookie)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []Attack{}, fmt.Errorf("failed to get: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return []Attack{}, fmt.Errorf("failed to get, non 200 status code: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Attack{}, fmt.Errorf("failed to read body: %s", err)
	}

	var recent Recent
	err = json.Unmarshal(body, &recent)
	if err != nil {
		return []Attack{}, fmt.Errorf("failed to unmarshal: %s", err)
	}

	if recent.Message != "Success" {
		return []Attack{}, fmt.Errorf("failed to get, no Success: %s", recent.Message)
	}

	if len(recent.Results) == 0 {
		log.Println("no attacks")
		return []Attack{}, nil
	}

	var attacks []Attack
	for _, attack := range recent.Results {
		attacks = append(attacks, Attack{
			Zone:        attack.TenantName,
			IPaddress:   attack.IPaddress,
			Status:      attack.Status,
			Start:       attack.StartTime,
			Duration:    attack.Duration,
			End:         attack.StartTime + attack.Duration,
			Max:         attack.MaxPeak,
			AttackID:    attack.AttackID,
			Description: attack.Description,
			Traffic:     attack.Volume,
		})
	}

	return attacks, nil
}

func GetTop(duration time.Duration) ([]IP, error) {
	apiurl := config.GlobalConfig.Corero.URL + "api/"
	t := time.Now().Add(-duration)
	starttime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	req, err := http.NewRequest("GET", apiurl+"get_top_attacked_ip", nil)
	if err != nil {
		return []IP{}, fmt.Errorf("failed to format get: %s", err)
	}

	q := req.URL.Query()
	q.Add("start", starttime)
	req.URL.RawQuery = q.Encode()
	req.AddCookie(Cookie)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []IP{}, fmt.Errorf("failed to get: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return []IP{}, fmt.Errorf("failed to get, non 200 status code: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []IP{}, fmt.Errorf("failed to read body: %s", err)
	}

	var top Top
	err = json.Unmarshal(body, &top)
	if err != nil {
		return []IP{}, fmt.Errorf("failed to unmarshal: %s", err)
	}

	if top.Message != "Success" {
		return []IP{}, fmt.Errorf("failed to get, no Success: %s", top.Message)
	}

	if len(top.Results) == 0 {
		log.Println("no ips")
		return []IP{}, nil
	}

	var ips []IP
	for _, ip := range top.Results {
		ips = append(ips, IP{
			IPaddress:   ip.IPaddress,
			AttackCount: ip.AttackCount,
			Zone:        ip.TenantName,
		})
	}

	return ips, nil
}
