/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package operation

import (
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/client/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/pkg/errors"

	"github.com/trustbloc/edge-adapter/pkg/internal/common/adapterutil"
	"github.com/trustbloc/edge-adapter/pkg/presentationex"
	"github.com/trustbloc/edge-adapter/pkg/vc"
	"github.com/trustbloc/edge-adapter/pkg/vc/rp"
)

var errInvalidCredential = errors.New("malformed credential")

func parseWalletResponse(
	definitions *presentationex.PresentationDefinitions, vpBytes []byte) (*vc.UserConsentCredential, error) {
	vp, err := verifiable.ParsePresentation(vpBytes)
	if err != nil {
		return nil, errors.Wrapf(
			errInvalidCredential, fmt.Sprintf("error parsing a verifiable presentation : %s", err))
	}

	err = evaluatePresentationSubmission(definitions, vp)
	if err != nil {
		return nil, errors.Wrapf(errInvalidCredential, "invalid presentation submission : %s", err)
	}

	rawCreds, err := vp.MarshalledCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credentials from vp : %w", err)
	}

	var base *verifiable.Credential

	for i := range rawCreds {
		raw := rawCreds[i]

		cred, parseErr := verifiable.ParseCredential(raw)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse raw credential %s : %w", string(raw), parseErr)
		}

		if adapterutil.StringsContains(vc.UserConsentCredentialType, cred.Types) {
			base = cred
			break
		}

		logger.Warnf("ignoring credential with unrecognized types: %+v", cred.Types)
	}

	if base == nil {
		return nil, errors.Wrapf(
			errInvalidCredential, "no suitable credential of type %s found", vc.UserConsentCredentialType)
	}

	consentVC := &vc.UserConsentCredential{Base: base}

	err = adapterutil.DecodeJSONMarshaller(base, consentVC)
	if err != nil {
		return nil, fmt.Errorf("unable to decode user consent credential : %w", err)
	}

	return consentVC, nil
}

func parseIssuerResponse(def *presentationex.PresentationDefinitions,
	pres *presentproof.Presentation) (*rp.PresentationSubmissionPresentation, error) {
	if len(pres.PresentationsAttach) == 0 {
		return nil, fmt.Errorf("%w : expected at least 1 attachment but got 0", errInvalidCredential)
	}

	attachment := pres.PresentationsAttach[0]

	vpBytes, err := attachment.Data.Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contents of attachment with id %s : %w", attachment.ID, err)
	}

	vp, err := verifiable.ParsePresentation(vpBytes)
	if err != nil {
		return nil,
			errors.Wrapf(errInvalidCredential, fmt.Sprintf("failed to parse a verifiable presentation : %s", err))
	}

	err = evaluatePresentationSubmission(def, vp)
	if err != nil {
		return nil, errors.Wrapf(errInvalidCredential, "invalid presentation submission : %s", err)
	}

	presentationSubmissionVP := &rp.PresentationSubmissionPresentation{
		Base: vp,
	}

	err = adapterutil.DecodeJSONMarshaller(vp, presentationSubmissionVP)
	if err != nil {
		return nil, fmt.Errorf("failed to decode presentation_submission VP : %w", err)
	}

	return presentationSubmissionVP, nil
}

// TODO validate presentation_submission against presentation_definitions
//  https://github.com/trustbloc/edge-adapter/issues/108
func evaluatePresentationSubmission(_ *presentationex.PresentationDefinitions, vp *verifiable.Presentation) error {
	submission := &rp.PresentationSubmissionPresentation{
		Base: vp,
	}

	return adapterutil.DecodeJSONMarshaller(vp, submission)
}
