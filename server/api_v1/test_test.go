package api_v1

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/pkg/storage_mock"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/policy_mock"
	"google.golang.org/protobuf/proto"
)

const OtpStub = "123456"

type Setup struct {
	store           *storage_mock.MockStore
	engine          *engine.Engine
	opts            engine.Opts
	userPolicy      *policy_mock.MockUserPolicy
	workspacePolicy *policy_mock.MockWorkspacePolicy
	ctrl            *gomock.Controller
	res             *httptest.ResponseRecorder
}

func MustSetup(t *testing.T) *Setup {
	ctrl := gomock.NewController(t)
	store := storage_mock.NewMockStore(ctrl)
	userPolicy := policy_mock.NewMockUserPolicy(ctrl)
	workspacePolicy := policy_mock.NewMockWorkspacePolicy(ctrl)
	res := httptest.NewRecorder()
	opts := engine.Opts{
		AllowOrigin:     []string{"*"},
		UserPolicy:      userPolicy,
		WorkspacePolicy: workspacePolicy,
		OtpStub:         OtpStub,
		Production:      true,
	}

	app, err := engine.New(store, opts)
	require.NoError(t, err)
	router := Extend(app)

	userPolicy.EXPECT().Wrap(gomock.Any()).Return(userPolicy).AnyTimes()
	workspacePolicy.EXPECT().Wrap(gomock.Any(), gomock.Any()).Return(workspacePolicy).AnyTimes()

	return &Setup{
		ctrl:            ctrl,
		opts:            opts,
		store:           store,
		userPolicy:      userPolicy,
		workspacePolicy: workspacePolicy,
		res:             res,
		engine:          router,
	}
}

func (s *Setup) MustPOST(t *testing.T, path string, message proto.Message) {
	bytea, err := proto.Marshal(message)
	require.NoError(t, err, "marshalling message")

	req, err := http.NewRequest("POST", path, bytes.NewReader(bytea))
	req.Header.Add("Content-Type", "application/x-protobuf")
	req.Header.Add("Cookie", "jwt=foo")
	require.NoError(t, err, fmt.Sprintf("POST TO %s", path))

	s.engine.ServeHTTP(s.res, req)
}

func (s *Setup) MustBindResponse(t *testing.T, message proto.Message) {
	body, err := ioutil.ReadAll(s.res.Body)
	require.NoError(t, err, "body read")
	err = proto.Unmarshal(body, message)
	require.NoError(t, err, "unmarshall error")
}

func (s *Setup) LoginAs(t *testing.T, session domain.Session, user domain.User) {
	if !session.OtpSucceed.Valid {
		session.OtpSucceed = domain.NewEmptyBool(true)
	}

	s.store.EXPECT().RetrieveSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(session, nil)
	s.store.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(user, nil)
}
