package has2be

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"testing"
)

var calledTimes int

type dummyEvCertificateProvider struct{}

func (d dummyEvCertificateProvider) ProvideCertificate(_ context.Context, exiRequest string) (services.EvCertificate15118Response, error) {
	calledTimes++
	if exiRequest == "success" {
		return services.EvCertificate15118Response{
			Status:                     ocpp201.Iso15118EVCertificateStatusEnumTypeAccepted,
			CertificateInstallationRes: "dummy exi",
		}, nil
	} else {
		return services.EvCertificate15118Response{}, errors.New("failure, try again")
	}

}

func TestGet15118EvCertificate(t *testing.T) {
	schemaVersion := "urn:iso:15118:2:2013:MsgDef"
	req := &types.Get15118EVCertificateRequestJson{
		A15118SchemaVersion: &schemaVersion,
		ExiRequest:          "success",
	}

	h := Get15118EvCertificateHandler{
		EvCertificateProvider: dummyEvCertificateProvider{},
	}

	got, err := h.HandleCall(context.Background(), "cs001", req)
	want := &types.Get15118EVCertificateResponseJson{
		Status:      types.Iso15118EVCertificateStatusEnumTypeAccepted,
		ExiResponse: "dummy exi",
	}

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
