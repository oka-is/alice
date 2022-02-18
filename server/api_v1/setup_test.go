package api_v1

//import (
//	"bytes"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/oka-is/alice/server/engine"
//	"github.com/oka-is/alice/server/engine_mock"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/require"
//	"google.golang.org/protobuf/proto"
//)
//
//type Setup struct {
//	store  *engine_mock.MockSotre
//	engine *engine.Engine
//	ctrl   *gomock.Controller
//	res    *httptest.ResponseRecorder
//}
//
//func (s *Setup) Finish() {
//	s.ctrl.Finish()
//}
//
//func (s *Setup) MustPOST(t *testing.T, path string, message proto.Message) {
//	bytea, err := proto.Marshal(message)
//	require.NoError(t, err, "marshalling message")
//
//	req, err := http.NewRequest("POST", path, bytes.NewReader(bytea))
//	req.Header.Add("Content-Type", "application/x-protobuf")
//	require.NoError(t, err, fmt.Sprintf("POST TO %s", path))
//
//	s.engine.ServeHTTP(s.res, req)
//}
//
//func (s *Setup) ShouldBind(t *testing.T, message proto.Message) {
//	body, err := ioutil.ReadAll(s.res.Body)
//	require.NoError(t, err, "body read")
//	err = proto.Unmarshal(body, message)
//	require.NoError(t, err, "unmarshall error")
//}
//
//func MustSetup(t *testing.T) *Setup {
//	ctrl := gomock.NewController(t)
//	store := engine_mock.NewMockSotre(ctrl)
//	res := httptest.NewRecorder()
//	router := Extend(engine.New(store))
//
//	return &Setup{
//		ctrl:   ctrl,
//		store:  store,
//		res:    res,
//		engine: router,
//	}
//}
