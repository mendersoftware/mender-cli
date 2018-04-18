// Copyright 2018 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package deployments

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/pkg/errors"

	"github.com/mendersoftware/mender-cli/log"
)

const (
	artifactUploadUrl = "/api/management/v1/deployments/artifacts"
)

type Client struct {
	url               string
	artifactUploadUrl string
	client            *http.Client
}

func NewClient(url string, skipVerify bool) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	return &Client{
		url:               url,
		artifactUploadUrl: JoinURL(url, artifactUploadUrl),
		client: &http.Client{
			Transport: tr,
		},
	}
}

func (c *Client) UploadArtifact(description, artifactPath, tokenPath string, noProgress bool) error {

	var bar *pb.ProgressBar

	artifact, err := os.Open(artifactPath)
	if err != nil {
		return errors.Wrap(err, "Cannot read artifact file")
	}

	artifactStats, err := artifact.Stat()
	if err != nil {
		return errors.Wrap(err, "Cannot read artifact file stats")
	}

	// create pipe
	pR, pW := io.Pipe()

	// create multipart writer
	writer := multipart.NewWriter(pW)

	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return errors.Wrap(err, "Please Login first")
	}

	req, err := http.NewRequest(http.MethodPost, c.artifactUploadUrl, pR)
	if err != nil {
		return errors.Wrap(err, "Cannot create request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+string(token))

	reqDump, _ := httputil.DumpRequest(req, false)
	log.Verbf("sending request: \n%v", string(reqDump))

	if !noProgress {
		// create progress bar
		bar = pb.New64(artifactStats.Size()).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 100)
		bar.Start()
	}

	go func() {
		var part io.Writer
		defer pW.Close()
		defer artifact.Close()

		writer.WriteField("size", strconv.FormatInt(artifactStats.Size(), 10))
		writer.WriteField("description", description)
		part, _ = writer.CreateFormFile("artifact", artifactStats.Name())

		if !noProgress {
			part = io.MultiWriter(part, bar)
		}

		if _, err := io.Copy(part, artifact); err != nil {
			writer.Close()
			pR.CloseWithError(err)
			return
		}

		writer.Close()
		if !noProgress {
			bar.FinishPrint("Processing uploaded file. This may take around one minute.\n")
		}
		return
	}()

	rsp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "POST /artifacts request failed")
	}
	defer rsp.Body.Close()
	pR.Close()

	rspDump, _ := httputil.DumpResponse(rsp, true)
	log.Verbf("response: \n%v\n", string(rspDump))

	if rsp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return errors.Wrap(err, "can't read request body")
		}
		if rsp.StatusCode == http.StatusUnauthorized {
			log.Verbf("artifact upload failed wih status %d, reason: %s", rsp.StatusCode, body)
			return errors.New("Unauthorized. Please Login first")
		}
		return errors.New(fmt.Sprintf("artifact upload failed wih status %d, reason: %s", rsp.StatusCode, body))
	}

	return nil
}

func JoinURL(base, url string) string {
	if strings.HasPrefix(url, "/") {
		url = url[1:]
	}
	if !strings.HasSuffix(base, "/") {
		base = base + "/"
	}
	return base + url
}
