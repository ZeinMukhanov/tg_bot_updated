package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Call struct {
	Type                 string
	Status               string
	Time                 string
	CallSchemeID         string
	Scheme               string
	OutgoingLine         string
	From                 string
	To                   string
	Manager              string
	CallDuration         string
	TalkDuration         string
	AnswerTime           string
	Rating               string
	RecordID             string
	Label                string
	Tags                 string
	HangupInitiator      string
	CallbackOrderID      string
	DoesRecordExist      string
	IsNewClient          string
	CallbackState        string
	CallbackTime         string
	CRMInformation       string
	CRMResponsiblePerson string
}

type SipuniRequest struct {
	User            string
	FromDate        string
	ToDate          string
	TimeFrom        string
	TimeTo          string
	CallType        string
	State           string
	Tree            string
	ShowTreeID      string
	ToNumber        string
	FromNumber      string
	NumbersRinged   string
	NumbersInvolved string
	Names           string
	OutgoingLine    string
	ToAnswer        string
	Anonymous       string
	FirstTime       string
	DtmfUserAnswer  string
	Rating          string
	HangupInitor    string
	IgnoreSpecChar  string
	Secret          string
}

func buildSipuniRequest() SipuniRequest {
	return SipuniRequest{
		User:            "USERNAME",
		FromDate:        time.Now().Format("02.01.2006"),
		ToDate:          time.Now().Format("02.01.2006"),
		TimeFrom:        "00:00",
		TimeTo:          "23:59",
		CallType:        "0",
		State:           "0",
		Tree:            "",
		ShowTreeID:      "1",
		ToNumber:        "",
		FromNumber:      "",
		NumbersRinged:   "0",
		NumbersInvolved: "0",
		Names:           "1",
		OutgoingLine:    "1",
		ToAnswer:        "",
		Anonymous:       "0",
		FirstTime:       "0",
		DtmfUserAnswer:  "0",
		Rating:          "0",
		HangupInitor:    "1",
		IgnoreSpecChar:  "1",
		Secret:          "SECRETKEY",
	}
}

func buildHash(request SipuniRequest) string {
	hashString := strings.Join([]string{
		request.Anonymous, request.DtmfUserAnswer, request.FirstTime, request.FromDate,
		request.FromNumber, request.HangupInitor, request.IgnoreSpecChar, request.Names,
		request.NumbersInvolved, request.NumbersRinged, request.OutgoingLine,
		request.Rating, request.ShowTreeID, request.State, request.TimeFrom,
		request.TimeTo, request.ToDate, request.ToAnswer, request.ToNumber,
		request.Tree, request.CallType, request.User, request.Secret,
	}, "+")

	hasher := md5.New()
	hasher.Write([]byte(hashString))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Run() (map[string][]Call, error) {
	request := buildSipuniRequest()
	hash := buildHash(request)

	urlStr := "https://sipuni.com/api/statistic/export"
	data := url.Values{
		"anonymous":       {request.Anonymous},
		"dtmfUserAnswer":  {request.DtmfUserAnswer},
		"firstTime":       {request.FirstTime},
		"from":            {request.FromDate},
		"fromNumber":      {request.FromNumber},
		"hangupinitor":    {request.HangupInitor},
		"ignoreSpecChar":  {request.IgnoreSpecChar},
		"names":           {request.Names},
		"numbersInvolved": {request.NumbersInvolved},
		"numbersRinged":   {request.NumbersRinged},
		"outgoingLine":    {request.OutgoingLine},
		"rating":          {request.Rating},
		"showTreeId":      {request.ShowTreeID},
		"state":           {request.State},
		"timeFrom":        {request.TimeFrom},
		"timeTo":          {request.TimeTo},
		"to":              {request.ToDate},
		"toAnswer":        {request.ToAnswer},
		"toNumber":        {request.ToNumber},
		"tree":            {request.Tree},
		"type":            {request.CallType},
		"user":            {request.User},
		"hash":            {hash},
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	//fmt.Println("API Response:", string(body))
	response := string(body)
	if strings.HasPrefix(response, "\ufeff") {
		response = strings.TrimPrefix(response, "\ufeff")
	}
	lines := strings.Split(response, "\n")

	if len(lines) < 2 {
		fmt.Println("Invalid API response: Less than 2 lines.")
		return nil, nil
	}

	callsByManager := make(map[string][]Call)

	for _, line := range lines[1:] {
		cells := strings.Split(line, ";")
		if len(cells) < 24 { // Проверяем, что в строке достаточно ячеек
			continue
		}
		call := Call{
			Type:                 cells[0],
			Status:               cells[1],
			Time:                 cells[2],
			CallSchemeID:         cells[3],
			Scheme:               cells[4],
			OutgoingLine:         cells[5],
			From:                 cells[6],
			To:                   cells[7],
			Manager:              cells[8],
			CallDuration:         cells[9],
			TalkDuration:         cells[10],
			AnswerTime:           cells[11],
			Rating:               cells[12],
			RecordID:             cells[13],
			Label:                cells[14],
			Tags:                 cells[15],
			HangupInitiator:      cells[16],
			CallbackOrderID:      cells[17],
			DoesRecordExist:      cells[18],
			IsNewClient:          cells[19],
			CallbackState:        cells[20],
			CallbackTime:         cells[21],
			CRMInformation:       cells[22],
			CRMResponsiblePerson: cells[23],
		}
		callsByManager[call.Manager] = append(callsByManager[call.Manager], call)
	}

	return callsByManager, nil
}

func sendSipuniData(name string) (string, error) {
	data, err := Run()
	if err != nil {
		return "", err
	}

	desiredCRMResponsible := name
	var incomingCalls strings.Builder
	var outgoingCalls strings.Builder
	var analytics strings.Builder

	totalCalls := 0
	answeredCalls := 0
	unansweredCalls := 0

	today := time.Now().Format("2006-01-02")

	incomingCalls.WriteString(fmt.Sprintf("%s: Входящие %s\n", name, today))
	outgoingCalls.WriteString(fmt.Sprintf("%s: Исходящие %s\n", name, today))

	for _, calls := range data {
		for _, call := range calls {
			if call.CRMResponsiblePerson == desiredCRMResponsible {
				totalCalls++
				if call.Status == "Отвечен" {
					answeredCalls++
				} else if call.Status == "Не отвечен" {
					unansweredCalls++
				}
				text := fmt.Sprintf("%s, %s",
					getStatusIcon(call.Status), call.To)

				if call.Type == "Входящий" {
					incomingCalls.WriteString(text + "\n")
				} else if call.Type == "Исходящий" {
					outgoingCalls.WriteString(text + "\n")
				}
			}
		}
	}

	percentageUnanswered := float64(unansweredCalls) / float64(totalCalls) * 100
	analytics.WriteString(fmt.Sprintf("Всего звонков: %d\nОтвеченные звонки: %d\nНеотвеченные звонки: %d\nПроцент неотвеченных звонков: %.2f%%\n",
		totalCalls, answeredCalls, unansweredCalls, percentageUnanswered))

	result := fmt.Sprintf("%s\n------------------\n%s\n------------------\n%s", incomingCalls.String(), outgoingCalls.String(), analytics.String())

	return result, nil
}
