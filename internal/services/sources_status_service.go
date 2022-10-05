package services

import (
	"errors"
	stdhttp "net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
)

var UnknownProviderFromSourcesErr = errors.New("unknown provider returned from sources")

// SourcesStatus fetches information from sources and then performs a smallest possible
// request on the cloud provider (list keys or similar). Reports an error if sources configuration
// is no longer valid.
func SourcesStatus(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	sourceId := chi.URLParam(r, "ID")

	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "can't init sources client", err))
		return
	}

	auth, err := sourcesClient.GetAuthentication(r.Context(), sourceId)
	if err != nil {
		if errors.Is(err, http.ApplicationNotFoundErr) {
			renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources: application not found", err, 404))
			return
		}
		if errors.Is(err, http.AuthenticationForSourcesNotFoundErr) {
			renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources: authentication not found", err, 404))
			return
		}
		renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources", err, 500))
		return
	}

	var statusClient clients.ClientStatuser
	switch auth.Type() {
	case models.ProviderTypeAWS:
		statusClient, err = clients.GetCustomerEC2Client(r.Context(), auth, config.AWS.DefaultRegion)
		if err != nil {
			renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get client", err))
			return
		}
	case models.ProviderTypeGCP:
		statusClient, err = clients.GetGCPClient(r.Context(), auth)
		if err != nil {
			renderError(w, r, payloads.NewGCPError(r.Context(), "unable to get client", err))
			return
		}
	case models.ProviderTypeAzure:
		statusClient, err = clients.GetAzureClient(r.Context(), auth)
		if err != nil {
			renderError(w, r, payloads.NewAzureError(r.Context(), "unable to get client", err))
			return
		}
	case models.ProviderTypeNoop:
	case models.ProviderTypeUnknown:
	default:
		renderError(w, r, payloads.NewStatusError(r.Context(), UnknownProviderFromSourcesErr))
		return
	}

	err = statusClient.Status(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewStatusError(r.Context(), err))
		return
	}

	write200(w, r)
}