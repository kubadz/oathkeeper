/*
 * Copyright © 2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author       Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright  2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license  	   Apache-2.0
 */

package mutate

import (
	"bytes"
	"encoding/json"
	"github.com/ory/x/httpx"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/ory/oathkeeper/driver/configuration"
	"github.com/ory/oathkeeper/pipeline"
	"github.com/ory/oathkeeper/pipeline/authn"
)

var ErrMalformedResponseFromUpstreamAPI = "The call to an external API returned an invalid JSON object"
var ErrMissingAPIURL = "Missing URL in mutator configuration"
var ErrInvalidAPIURL = "Invalid URL in mutator configuration"
var ErrNon200ResponseFromAPI = "The call to an external API returned a non-200 HTTP response"
var ErrInvalidCredentials = "Invalid credentials were provided in mutator configuration"
var ErrNoCredentialsProvided = "No credentials were provided in mutator configuration"

type MutatorEnhancer struct {
	c      configuration.Provider
	client *http.Client
}

type BasicAuthn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Authentication struct {
	Basic BasicAuthn `json:"basic"`
}

type RetryConfig struct {
	Number         int `json:"number"`
	DelayInSeconds int `json:"delayInSeconds"`
}

type externalAPIConfig struct {
	Url   string          `json:"url"`
	Authn *Authentication `json:"authn,omitempty"`
	Retry RetryConfig     `json:"retry"`
}

type MutatorEnhancerConfig struct {
	Api externalAPIConfig `json:"api"`
}

func NewMutatorEnhancer(c configuration.Provider) *MutatorEnhancer {
	return &MutatorEnhancer{c: c, client: httpx.NewResilientClientLatencyToleranceSmall(nil)}
}

func (a *MutatorEnhancer) GetID() string {
	return "enhancer"
}

func (a *MutatorEnhancer) Mutate(r *http.Request, session *authn.AuthenticationSession, config json.RawMessage, _ pipeline.Rule) (http.Header, error) {
	if len(config) == 0 {
		config = []byte("{}")
	}
	var cfg MutatorEnhancerConfig
	d := json.NewDecoder(bytes.NewBuffer(config))
	d.DisallowUnknownFields()
	if err := d.Decode(&cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(session)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if cfg.Api.Url == "" {
		return nil, errors.New(ErrMissingAPIURL)
	} else if _, err := url.ParseRequestURI(cfg.Api.Url); err != nil {
		return nil, errors.New(ErrInvalidAPIURL)
	}
	req, err := http.NewRequest("POST", cfg.Api.Url, &b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	if cfg.Api.Authn != nil {
		credentials := cfg.Api.Authn.Basic
		req.SetBasicAuth(credentials.Username, credentials.Password)
	}

	res, err := a.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	switch res.StatusCode {
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		if cfg.Api.Authn != nil {
			return nil, errors.New(ErrInvalidCredentials)
		} else {
			return nil, errors.New(ErrNoCredentialsProvided)
		}
	default:
		return nil, errors.New(ErrNon200ResponseFromAPI)
	}
	sessionFromUpstream := authn.AuthenticationSession{}
	err = json.NewDecoder(res.Body).Decode(&sessionFromUpstream)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if sessionFromUpstream.Extra == nil || sessionFromUpstream.Subject != session.Subject { // TODO: should API be able to modify subject?
		return nil, errors.New(ErrMalformedResponseFromUpstreamAPI)
	}
	*session = sessionFromUpstream

	return nil, nil
}

func (a *MutatorEnhancer) Validate() error {
	if !a.c.MutatorEnhancerIsEnabled() {
		return errors.WithStack(ErrMutatorNotEnabled.WithReasonf(`Mutator "%s" is disabled per configuration.`, a.GetID()))
	}
	return nil
}
