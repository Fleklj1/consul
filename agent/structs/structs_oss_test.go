package structs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceID_String(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		sid := NewServiceID("the-id", &EnterpriseMeta{})
		require.Equal(t, "the-id", fmt.Sprintf("%v", sid))
	})
	t.Run("pointer", func(t *testing.T) {
		sid := NewServiceID("the-id", &EnterpriseMeta{})
		require.Equal(t, "the-id", fmt.Sprintf("%v", &sid))
	})
}

func TestCheckID_String(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		cid := NewCheckID("the-id", &EnterpriseMeta{})
		require.Equal(t, "the-id", fmt.Sprintf("%v", cid))
	})
	t.Run("pointer", func(t *testing.T) {
		cid := NewCheckID("the-id", &EnterpriseMeta{})
		require.Equal(t, "the-id", fmt.Sprintf("%v", &cid))
	})
}

func TestServiceName_String(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		sn := NewServiceName("the-id", &EnterpriseMeta{})
		require.Equal(t, "the-id", fmt.Sprintf("%v", sn))
	})
	t.Run("pointer", func(t *testing.T) {
		sn := NewServiceName("the-id", &EnterpriseMeta{})
		require.Equal(t, "the-id", fmt.Sprintf("%v", &sn))
	})
}
