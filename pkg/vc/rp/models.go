/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package rp

import (
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"

	"github.com/trustbloc/edge-adapter/pkg/presentationex"
)

const (
	// PresentationSubmissionPresentationType is the PresentationSubmissionPresentation's JSON-LD type.
	PresentationSubmissionPresentationType = "PresentationSubmission"
)

// PresentationSubmissionPresentation is the PresentationSubmission VerifiablePresentation.
// https://identity.foundation/presentation-exchange/#presentation-submission.
type PresentationSubmissionPresentation struct {
	Base                   *verifiable.Presentation               `json:"-"`
	PresentationSubmission *presentationex.PresentationSubmission `json:"presentation_submission,omitempty"`
}
