// SPDX-License-Identifier: Apache-2.0

package ocpi_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpi"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"net/http"
	"testing"
)

func TestGetVersions(t *testing.T) {
	engine := inmemory.NewStore()
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")

	want := []ocpi.Version{
		{
			Version: "2.2",
			Url:     "/ocpi/2.2",
		},
	}

	got, err := ocpiApi.GetVersions(context.Background())
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetVersionDetails(t *testing.T) {
	engine := inmemory.NewStore()
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")

	want := ocpi.VersionDetail{
		Version: "2.2",
		Endpoints: []ocpi.Endpoint{
			{
				Identifier: "credentials",
				Role:       ocpi.RECEIVER,
				Url:        "/ocpi/2.2/credentials",
			},
			{
				Identifier: "tokens",
				Role:       ocpi.RECEIVER,
				Url:        "/ocpi/receiver/2.2/tokens/",
			},
		},
	}

	got, err := ocpiApi.GetVersion(context.Background())
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetToken(t *testing.T) {
	engine := inmemory.NewStore()
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")
	err := engine.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		ContractId:  "GBTWKTWTW000018",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "ALWAYS",
	})
	require.NoError(t, err)

	want := &ocpi.Token{
		ContractId:  "GBTWKTWTW000018",
		CountryCode: "GB",
		Issuer:      "Thoughtworks",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		Valid:       true,
		Whitelist:   "ALWAYS",
	}

	got, err := ocpiApi.GetToken(context.Background(), "GB", "TWK", "DEADBEEF")
	require.NoError(t, err)

	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`, got.LastUpdated)
	got.LastUpdated = ""
	assert.Equal(t, want, got)
}

func TestSetToken(t *testing.T) {
	engine := inmemory.NewStore()
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")

	err := ocpiApi.SetToken(context.Background(), ocpi.Token{
		ContractId:  "GBTWKTWTW000018",
		CountryCode: "GB",
		Issuer:      "Thoughtworks",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		Valid:       true,
		Whitelist:   "ALWAYS",
	})
	require.NoError(t, err)

	want := &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		ContractId:  "GBTWKTWTW000018",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "ALWAYS",
	}

	got, err := engine.LookupToken(context.Background(), "DEADBEEF")
	require.NoError(t, err)
	got.LastUpdated = ""
	assert.Equal(t, want, got)
}
