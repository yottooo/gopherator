package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handleWord(body []byte) ([]byte, *serverError) {

	var english, gopher string
	err := json.Unmarshal(body, &struct {
		Word *string `json:"english-word"`
	}{
		Word: &english,
	})
	if err != nil {
		return nil, &serverError{Error: err, Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	if len(english) == 0 {
		return nil, &serverError{Error: nil, Msg: "No english word provided.", Code: http.StatusBadRequest}
	}

	var ok bool
	if gopher, ok = history.Load(english); !ok {
		gopher, err = TranslateWord(english)
		if err != nil {
			return nil, &serverError{Error: err, Msg: "Translation failed", Code: http.StatusBadRequest}
		}
		history.Store(english, gopher)
	}

	data, err := json.Marshal(struct {
		Word *string `json:"gopher-word"`
	}{
		Word: &gopher,
	})
	if err != nil {
		return nil, &serverError{Error: err, Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return data, nil
}

func handleSentence(body []byte) ([]byte, *serverError) {

	var english, gopher string
	err := json.Unmarshal(body, &struct {
		Word *string `json:"english-sentence"`
	}{
		Word: &english,
	})
	if err != nil {
		return nil, &serverError{Error: err, Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	if len(english) == 0 {
		return nil, &serverError{Error: nil, Msg: "No english sentence provided.", Code: http.StatusBadRequest}
	}
	var ok bool
	if gopher, ok = history.Load(english); !ok {
		gopher, err = TranslateSentence(english)
		if err != nil {
			return nil, &serverError{Error: err, Msg: "Translation failed", Code: http.StatusBadRequest}
		}
		history.Store(english, gopher)
	}

	data, err := json.Marshal(struct {
		Word *string `json:"gopher-sentence"`
	}{
		Word: &gopher,
	})
	if err != nil {
		return nil, &serverError{Error: err, Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return data, nil
}

var history History

func handleHistory(body []byte) ([]byte, *serverError) {

	data, err := history.ToJSON()
	if err != nil {
		return nil, &serverError{Error: err, Msg: err.Error(), Code: http.StatusInternalServerError}

	}
	return data, nil
}

type serverError struct {
	Error error
	Msg   string
	Code  int
}

type serverHandler struct {
	Handle func(body []byte) ([]byte, *serverError)
	Method string
}

func (sh serverHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if req.Method != sh.Method {
		http.Error(rw, fmt.Sprintf("%s method not supported", req.Method), http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	data, srvErr := sh.Handle(body)
	if srvErr != nil {
		http.Error(rw, srvErr.Msg, srvErr.Code)
		return
	}

	_, err = rw.Write(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

func initServer(port int) error {
	http.Handle("/word", serverHandler{Handle: handleWord, Method: http.MethodPost})
	http.Handle("/sentence", serverHandler{Handle: handleSentence, Method: http.MethodPost})
	http.Handle("/history", serverHandler{Handle: handleHistory, Method: http.MethodGet})

	log.Println("Start server on ", port)
	strPort := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(strPort, nil)
	if err != nil {
		return err
	}
	return nil
}
