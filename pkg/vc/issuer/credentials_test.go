/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package issuer

import (
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"

	"github.com/stretchr/testify/require"
)

const (
	id              = "e9e0f944-7b74-4298-9f3e-00ca609d6266"
	inviteeDID      = "did:example:7b744298e9e0f"
	inviterDID      = "agc"
	inviterLabel    = "user-agent"
	threadID        = "92d01dp5251sien42ad4dq4q2"
	connectionState = "completed"

	vc = `{
	   "@context":[
		  "https://www.w3.org/2018/credentials/v1",
		  "https://www.w3.org/2018/credentials/examples/v1"
	   ],
	   "id":"http://example.edu/credentials/1872",
	   "type":[
		  "VerifiableCredential",
		  "DIDConnectCredential"
	   ],
	   "credentialSubject":{
		  "id":"` + id + `",
		  "inviteeDID":"` + inviteeDID + `",
		  "inviterDID":"` + inviterDID + `",
		  "inviterLabel":"` + inviterLabel + `",
		  "threadID":"` + threadID + `",
		  "connectionState": "` + connectionState + `"
	   },
	   "issuer":"did:example:76e12ec712ebc6f1c221ebfeb1f",
	   "issuanceDate":"2010-01-01T19:23:24Z"
	}`
)

func TestParseWalletResponse(t *testing.T) {
	t.Run("test parse wallet - success", func(t *testing.T) {
		conn, err := ParseWalletResponse(getTestVP(t))
		require.NoError(t, err)
		require.NotNil(t, conn)

		require.Equal(t, conn.ID, id)
		require.Equal(t, conn.InviteeDID, inviteeDID)
		require.Equal(t, conn.InviterDID, inviterDID)
		require.Equal(t, conn.InviterLabel, inviterLabel)
		require.Equal(t, conn.ConnectionState, connectionState)
		require.Equal(t, conn.ThreadID, threadID)
	})

	t.Run("test parse wallet - invalid vp", func(t *testing.T) {
		conn, err := ParseWalletResponse([]byte("invalid json"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid presentation")
		require.Nil(t, conn)
	})

	t.Run("test parse wallet - no credentials inside vp", func(t *testing.T) {
		vp := verifiable.Presentation{
			Context: []string{"https://www.w3.org/2018/credentials/v1"},
			Type:    []string{"VerifiablePresentation"},
		}
		vpJSON, err := vp.MarshalJSON()
		require.NoError(t, err)

		conn, err := ParseWalletResponse(vpJSON)
		require.Error(t, err)
		require.Contains(t, err.Error(), "there should be only one credential")
		require.Nil(t, conn)
	})

	t.Run("test parse wallet - invalid credential inside vp", func(t *testing.T) {
		vc := verifiable.Credential{
			Context: []string{"https://www.w3.org/2018/credentials/v1"},
			Types:   []string{"VerifiablePresentation"},
		}

		vp, err := vc.Presentation()
		require.NoError(t, err)

		vpJSON, err := vp.MarshalJSON()
		require.NoError(t, err)

		conn, err := ParseWalletResponse(vpJSON)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse credential")
		require.Nil(t, conn)
	})

	t.Run("test parse wallet - no credential of DIDConnectCredential type inside vp", func(t *testing.T) {
		vc, err := verifiable.ParseCredential([]byte(vc))
		require.NoError(t, err)

		vc.Types = []string{"VerifiableCredential"}

		vp, err := vc.Presentation()
		require.NoError(t, err)

		vpJSON, err := vp.MarshalJSON()
		require.NoError(t, err)

		conn, err := ParseWalletResponse(vpJSON)
		require.Error(t, err)
		require.Contains(t, err.Error(), "vc doesn't contain DIDConnectCredential type")
		require.Nil(t, conn)
	})
}

func getTestVP(t *testing.T) []byte {
	vc, err := verifiable.ParseCredential([]byte(vc))
	require.NoError(t, err)

	vp, err := vc.Presentation()
	require.NoError(t, err)

	vpJSON, err := vp.MarshalJSON()
	require.NoError(t, err)

	return vpJSON
}