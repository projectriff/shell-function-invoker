/*
 * Copyright 2018-Present the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/projectriff/shell-function-invoker/pkg/function"
)

const (
	CorrelationId = "correlationId"
)


func invoke(fnUri string, in *function.Message) (*function.Message, error) {
	var outputBuf bytes.Buffer

	cmd := exec.Command(fnUri)

	// Record any correlation id header.
	hasCorrelationId := false
	correlationId := []string{}
	if _, ok := in.Headers[CorrelationId]; ok {
		hasCorrelationId = true
		correlationId = in.Headers[CorrelationId].Values
	}

	cmd.Stdin = bytes.NewReader(in.Payload)
	cmd.Stdout = &outputBuf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	outHeaders := make(map[string]*function.Message_HeaderValue)

	// If there was a correlation id header, add it to the output headers.
	if hasCorrelationId {
		headerValue := function.Message_HeaderValue{
			Values: correlationId,
		}
		outHeaders[CorrelationId] = &headerValue
	}

	out := &function.Message{
		Payload: outputBuf.Bytes(),
		Headers: outHeaders,
	}
	return out, nil
}
