package api_v1

//import (
//	"errors"
//	"testing"
//
//	"github.com/oka-is/alice/desc/alice_v1"
//	"github.com/oka-is/alice/pkg/validator"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/require"
//)
//
//func TestCreateUser_Errors(t *testing.T) {
//	type args struct {
//		desc       string
//		storeError error
//		wantCode   int
//	}
//
//	tests := []args{
//		{
//			desc:       "when success",
//			storeError: nil,
//			wantCode:   200,
//		}, {
//			desc:       "when validation error",
//			storeError: validator.NewError("foo"),
//			wantCode:   422,
//		}, {
//			desc:       "when internal error",
//			storeError: errors.New("foo"),
//			wantCode:   500,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.desc, func(t *testing.T) {
//			setup := MustSetup(t)
//			defer setup.Finish()
//
//			setup.store.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.storeError)
//			setup.MustPOST(t, "/v1/users", &alice_v1.CreateUserRequest{})
//			require.Equal(t, tt.wantCode, setup.res.Code)
//		})
//	}
//}
//
//func TestCreateUser_Render(t *testing.T) {
//	t.Skip()
//
//	t.Run("it works", func(t *testing.T) {
//		setup := MustSetup(t)
//		defer setup.Finish()
//
//		setup.store.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
//		setup.MustPOST(t, "/v1/users", &alice_v1.CreateUserRequest{
//			Username: "foo",
//		})
//
//		proto := new(alice_v1.CreateUserResponse)
//		require.Equal(t, 200, setup.res.Code)
//		setup.ShouldBind(t, proto)
//		require.Equal(t, "foo", proto.GetUser().GetUsername())
//	})
//}
