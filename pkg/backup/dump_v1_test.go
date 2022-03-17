package backup

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/lib/encoder"
	"github.com/wault-pw/alice/pkg/backup_mock"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/pkg/storage_mock"
	"google.golang.org/protobuf/proto"
)

func Test_Whoami(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		ctrl, buf, backup, store, flush := MustSetup(t)
		defer ctrl.Finish()

		user := domain.User{PasswdSalt: domain.NewEmptyBytes([]byte{1, 2, 3})}

		flush.EXPECT().Flush().Times(1)
		store.EXPECT().
			FindUser(gomock.Any(), gomock.Any()).
			Return(user, nil).
			Times(1)

		err := Whoami(backup, "1")
		require.NoError(t, err)

		res := new(alice_v1.WhoAmIResponse)
		marker := MustParse(t, buf, res)
		require.Equal(t, MarkerWhoAmI, marker)

		require.Equal(t, user.PasswdSalt.Bytea, res.User.PasswdSalt)
	})
}

func Test_ListWorkspace(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		ctrl, buf, backup, _, flush := MustSetup(t)
		defer ctrl.Finish()

		workspace := domain.UserWithWorkspace{WorkspaceID: domain.NewEmptyString("foo")}

		flush.EXPECT().Flush().Times(1)

		err := ListWorkspace(backup, workspace)
		require.NoError(t, err)

		res := new(alice_v1.UserWithWorkspace)
		marker := MustParse(t, buf, res)
		require.Equal(t, MarkerWorkspace, marker)

		require.Equal(t, res.WorkspaceId, workspace.WorkspaceID.String)
	})
}

func Test_ListCard(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		ctrl, buf, backup, _, flush := MustSetup(t)
		defer ctrl.Finish()

		flush.EXPECT().Flush().Times(1)

		card := domain.Card{ID: domain.NewEmptyString("foo")}
		err := ListCard(backup, card)
		require.NoError(t, err)

		res := new(alice_v1.Card)
		marker := MustParse(t, buf, res)
		require.Equal(t, MarkerCard, marker)

		require.Equal(t, res.Id, card.ID.String)
	})
}

func Test_ListCardItems(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		ctrl, buf, backup, store, flush := MustSetup(t)
		defer ctrl.Finish()

		item := domain.CardItem{ID: domain.NewEmptyString("foo")}

		flush.EXPECT().Flush().Times(1)
		store.EXPECT().ListCardItems(gomock.Any(), gomock.Any()).Return([]domain.CardItem{item}, nil)

		err := ListCardItems(backup, domain.Card{})
		require.NoError(t, err)

		res := new(alice_v1.CardItem)
		marker := MustParse(t, buf, res)
		require.Equal(t, MarkerCardItem, marker)
		require.Equal(t, res.Id, item.ID.String)
	})
}

func MustSetup(t *testing.T) (*gomock.Controller, *bytes.Buffer, *Backup, *storage_mock.MockStore, *backup_mock.MockFlush) {
	buf := bytes.NewBuffer(nil)
	ctrl := gomock.NewController(t)
	store := storage_mock.NewMockStore(ctrl)
	flush := backup_mock.NewMockFlush(ctrl)
	backup := NewBackup(store, buf, flush)
	return ctrl, buf, backup, store, flush
}

func MustDecode(t *testing.T, input []byte, message proto.Message) {
	err := proto.Unmarshal(input, message)
	require.NoError(t, err, "proto decoding")
}

func MustParse(t *testing.T, reader *bytes.Buffer, message proto.Message) byte {
	marker, err := reader.ReadByte()
	require.NoError(t, err)

	sized := encoder.MakeUint32()
	_, err = reader.Read(sized)
	require.NoError(t, err)

	size := int(encoder.Uint32Binary(sized))
	body := make([]byte, size)
	_, err = reader.Read(body)
	require.NoError(t, err)
	require.Equal(t, len(body), size, "wrong body size")

	MustDecode(t, body, message)
	return marker
}
