package state

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/consul/agent/connect"
	"github.com/hashicorp/consul/agent/consul/stream"
	"github.com/hashicorp/consul/agent/structs"
)

func TestCARootsEvents(t *testing.T) {
	store := testStateStore(t)
	rootA := connect.TestCA(t, nil)

	_, err := store.CARootSetCAS(1, 0, structs.CARoots{rootA})
	require.NoError(t, err)

	t.Run("roots changed", func(t *testing.T) {
		tx := store.db.WriteTxn(2)
		defer tx.Abort()

		rootB := connect.TestCA(t, nil)
		err = caRootSetCASTxn(tx, 2, 1, structs.CARoots{rootB})
		require.NoError(t, err)

		events, err := caRootsChangeEvents(tx, Changes{Index: 2, Changes: tx.Changes()})
		require.NoError(t, err)
		require.Equal(t, []stream.Event{
			{
				Topic: EventTopicCARoots,
				Index: 2,
				Payload: EventPayloadCARoots{
					CARoots: structs.CARoots{rootB},
				},
			},
		}, events)
	})

	t.Run("no change", func(t *testing.T) {
		tx := store.db.ReadTxn()
		defer tx.Abort()

		events, err := caRootsChangeEvents(tx, Changes{Index: 2, Changes: tx.Changes()})
		require.NoError(t, err)
		require.Empty(t, events)
	})
}

func TestCARootsSnapshot(t *testing.T) {
	store := testStateStore(t)
	fn := caRootsSnapshot((*readDB)(store.db.db))

	var req stream.SubscribeRequest

	t.Run("no roots", func(t *testing.T) {
		buf := &snapshotAppender{}

		idx, err := fn(req, buf)
		require.NoError(t, err)
		require.Equal(t, uint64(0), idx)

		require.Len(t, buf.events, 1)
		require.Len(t, buf.events[0], 1)

		payload := buf.events[0][0].Payload.(EventPayloadCARoots)
		require.Empty(t, payload.CARoots)
	})

	t.Run("with roots", func(t *testing.T) {
		buf := &snapshotAppender{}

		root := connect.TestCA(t, nil)

		_, err := store.CARootSetCAS(1, 0, structs.CARoots{root})
		require.NoError(t, err)

		idx, err := fn(req, buf)
		require.NoError(t, err)
		require.Equal(t, uint64(1), idx)

		require.Equal(t, buf.events, [][]stream.Event{
			{
				{
					Topic: EventTopicCARoots,
					Index: 1,
					Payload: EventPayloadCARoots{
						CARoots: structs.CARoots{root},
					},
				},
			},
		})
	})
}
