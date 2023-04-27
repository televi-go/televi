package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/televi-go/televi/telegram"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"runtime"
)

type Api struct {
	Token   string
	client  http.Client
	Address string
}

func (api *Api) LogError(err error, request telegram.Request) {
	var buf = make([]byte, 2048)
	runtime.Stack(buf, true)
	log.Printf("error with request %#v: %v\n%s\n", request, err, buf)
}

func buildParams(in telegram.Params, files []telegram.File) (*bytes.Buffer, string) {

	if in == nil {
		return nil, ""
	}

	out := &bytes.Buffer{}
	w := multipart.NewWriter(out)
	defer w.Close()
	for key, value := range in {
		w.WriteField(key, value)
	}

	for _, file := range files {
		if file.FileId != "" {
			w.WriteField(file.FieldName, file.FileId)
			continue
		}

		fileBytes, err := io.ReadAll(file.Reader)
		if err != nil {
			panic(err)
		}
		mimeType := mimetype.Detect(fileBytes)
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				file.FieldName, file.Name))
		header.Set("Content-Type", mimeType.String())
		part, err := w.CreatePart(header)
		_, err = io.Copy(part, bytes.NewReader(fileBytes))
		if err != nil {
			panic(err)
		}
	}

	return out, w.FormDataContentType()
}

func (api *Api) getHttpRequest(request telegram.Request, ctx context.Context) (*http.Request, error) {
	endPoint := fmt.Sprintf("%s/bot%s/%s", api.Address, api.Token, request.Method())

	params, err := request.Params()
	if err != nil {
		return nil, err
	}

	var files []telegram.File

	fileRequest, isFileRequest := request.(telegram.RequestWithFiles)
	if isFileRequest {
		files = fileRequest.Files()
	}

	urlParameters, cType := buildParams(params, files)
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", endPoint, urlParameters)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("Content-Type", cType)

	return httpRequest, nil
}

func (api *Api) RequestContext(request telegram.Request, ctx context.Context) (telegram.Response, error) {
	httpRequest, err := api.getHttpRequest(request, ctx)

	if err != nil {
		return telegram.Response{}, err
	}

	response, err := api.client.Do(httpRequest)
	if err != nil {
		return telegram.Response{}, err
	}

	defer response.Body.Close()

	text, err := io.ReadAll(response.Body)

	decoder := json.NewDecoder(bytes.NewReader(text))
	var responseObj telegram.Response
	err = decoder.Decode(&responseObj)
	if !responseObj.Ok {
		err = fmt.Errorf("%s\n\n%s", responseObj.Result, text)
	}
	return responseObj, err
}

func (api *Api) Request(request telegram.Request) (telegram.Response, error) {
	return api.RequestContext(request, context.Background())
}

func (api *Api) LaunchRequest(request telegram.Request) {
	go func() {
		_, err := api.Request(request)
		if err != nil {
			var buf = make([]byte, 1024)
			runtime.Stack(buf, false)
			log.Printf("error with request %#v %v\n%s\n", request, err, buf)
		}
	}()
}

func NewApi(token string, address string) *Api {
	return &Api{
		Token: token,
		client: http.Client{
			Transport: &http.Transport{MaxConnsPerHost: 50},
		},
		Address: address,
	}
}
