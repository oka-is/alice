package storage

import (
	"context"
	"testing"

	"github.com/oka-is/alice/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestStorage_CloneCard(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	t.Run("it works", func(t *testing.T) {
		ctx := context.Background()
		card := mustCreateCard(t, store, &domain.Card{
			Archived: domain.NewEmptyBool(true),
		})

		mustCreateCardItem(t, store, &domain.CardItem{})
		mustCreateCardItem(t, store, &domain.CardItem{CardID: card.ID})

		clone, err := store.CloneCard(ctx, card.ID.String, []byte("foo"))
		require.NoError(t, err)

		items, err := store.ListCardItems(ctx, clone.ID.String)
		require.NoError(t, err)

		require.NotEmpty(t, clone.ID.String)
		require.NotEqual(t, clone.ID.String, card.ID.String)

		require.Equal(t, []byte("foo"), clone.TitleEnc.Bytea)
		require.Equal(t, card.WorkspaceID.String, clone.WorkspaceID.String)
		require.Equal(t, card.TagsEnc.Slice, clone.TagsEnc.Slice)
		require.Equal(t, card.Archived.Bool, clone.Archived.Bool)
		require.Len(t, items, 1)
	})
}

func mustBuildCard(t *testing.T, storage *Storage, input *domain.Card) *domain.Card {
	out := &domain.Card{
		TitleEnc: domain.NewEmptyBytes([]byte("TitleEnc")),
		TagsEnc:  domain.NewEmptyByteSlice([][]byte{[]byte("tag")}),
		Archived: domain.NewEmptyBool(false),
	}

	if input.TitleEnc.Valid {
		out.TitleEnc = input.TitleEnc
	}

	if input.WorkspaceID.Valid {
		out.WorkspaceID = input.WorkspaceID
	} else {
		out.WorkspaceID = mustCreateWorkspace(t, storage, &domain.Workspace{}).ID
	}

	if input.TagsEnc.Valid {
		out.TagsEnc = input.TagsEnc
	}

	if input.TagsEnc.Valid {
		out.TagsEnc = input.TagsEnc
	}

	return out
}

func mustCreateCard(t *testing.T, storage *Storage, input *domain.Card) *domain.Card {
	ctx := context.Background()
	output := mustBuildCard(t, storage, input)
	if err := storage.insertCard(ctx, storage.db, output); err != nil {
		t.Errorf("cant create factory card: %s", err)
		t.FailNow()
	}
	return output
}
