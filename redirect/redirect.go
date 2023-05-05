package redirect

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
)

type Redirect interface {
	Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, path string, errHandlingResponse func(w http.ResponseWriter, err error))
}

type redirect struct {
	baseUrl string
	client  *http.Client
}

func Init(ctx context.Context, baseUrl string, client *http.Client) Redirect {
	return &redirect{
		baseUrl: baseUrl,
		client:  client,
	}
}

func (redirect *redirect) Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, path string, errHandlingResponse func(w http.ResponseWriter, err error)) {
	if errHandlingResponse == nil {
		errHandlingResponse = func(w http.ResponseWriter, err error) {
			log.Println(err)
		}
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errHandlingResponse(w, err)
		return
	}

	request, err := http.NewRequest(r.Method, redirect.baseUrl+path, bytes.NewBuffer(reqBody))
	if err != nil {
		errHandlingResponse(w, err)
		return
	}

	response, err := redirect.client.Do(request)
	if err != nil {
		errHandlingResponse(w, err)
		return
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		errHandlingResponse(w, err)
		return
	}

	w.Write(resBody)
	for k, v := range response.Header {
		head := ""
		for _, s := range v {
			head += s + ";"
		}
		w.Header().Set(k, head)
	}
	return
}
