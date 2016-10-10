package azure

import (
	"fmt"
	"strings"

	"koding/db/mongodb/modelhelper"
	"koding/kites/kloud/stack"

	"github.com/Azure/azure-sdk-for-go/management"
	"golang.org/x/net/context"
)

func (s *Stack) Authenticate(ctx context.Context) (interface{}, error) {
	var arg stack.AuthenticateRequest
	if err := s.Req.Args.One().Unmarshal(&arg); err != nil {
		return nil, err
	}

	if err := arg.Valid(); err != nil {
		return nil, err
	}

	if err := s.Builder.BuildCredentials(s.Req.Method, s.Req.Username, arg.GroupName, arg.Identifiers); err != nil {
		return nil, err
	}

	s.Log.Debug("Fetched terraform data: koding=%+v, template=%+v", s.Builder.Koding, s.Builder.Template)

	resp := make(stack.AuthenticateResponse)

	for _, cred := range s.Builder.Credentials {
		res := &stack.AuthenticateResult{}
		resp[cred.Identifier] = res

		if cred.Provider != "azure" {
			res.Message = "unable to authenticate non-azure credential: " + cred.Provider
			continue
		}

		meta := cred.Meta.(*Cred)

		if err := meta.Valid(); err != nil {
			res.Message = fmt.Sprintf("validation error: %s", cred.Identifier, err)
			continue
		}

		c, err := management.ClientFromPublishSettingsData(meta.PublishSettings, meta.SubscriptionID)
		if err != nil {
			res.Message = err.Error()
			continue
		}

		err = c.WaitForOperation("invalid", nil)
		if err != nil && !strings.Contains(err.Error(), "The operation request ID was not found") {
			res.Message = err.Error()
			continue
		}

		if err := modelhelper.SetCredentialVerified(cred.Identifier, true); err != nil {
			res.Message = err.Error()
			continue
		}

		res.Verified = true
	}

	s.Log.Debug("Authenticate credentials result: %+v", resp)

	return resp, nil
}
