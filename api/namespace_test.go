// +build consulent

package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPI_Namespaces(t *testing.T) {
	t.Parallel()
	c, s := makeACLClient(t)
	defer s.Stop()

	namespaces := c.Namespaces()
	acl := c.ACL()

	nsPolicy, _, err := acl.PolicyCreate(&ACLPolicy{
		Name:  "ns-policy",
		Rules: `operator = "write"`,
	}, nil)
	require.NoError(t, err)

	nsRole, _, err := acl.RoleCreate(&ACLRole{
		Name: "ns-role",
		Policies: []*ACLRolePolicyLink{
			&ACLRolePolicyLink{
				ID: nsPolicy.ID,
			},
		},
	}, nil)

	require.NoError(t, err)

	t.Run("Create Nameless", func(t *testing.T) {
		ns := Namespace{
			Description: "foo",
		}

		_, _, err := namespaces.Create(&ns, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Must specify a Name for Namespace creation")
	})

	t.Run("Create", func(t *testing.T) {
		ns, _, err := namespaces.Create(&Namespace{
			Name: "foo",
		}, nil)
		require.NoError(t, err)
		require.NotNil(t, ns)
		require.Equal(t, "foo", ns.Name)
		require.Nil(t, ns.ACLs)

		ns, _, err = namespaces.Create(&Namespace{
			Name:        "acls",
			Description: "This namespace has ACL config attached",
			ACLs: &NamespaceACLConfig{
				PolicyDefaults: []ACLLink{
					ACLLink{ID: nsPolicy.ID},
				},
				RoleDefaults: []ACLLink{
					ACLLink{ID: nsRole.ID},
				},
			},
		}, nil)

		require.NoError(t, err)
		require.NotNil(t, ns)
		require.NotNil(t, ns.ACLs)
		require.Nil(t, ns.DeletedAt)
	})

	t.Run("Update Nameless", func(t *testing.T) {
		ns := Namespace{
			Description: "foo",
		}

		_, _, err := namespaces.Update(&ns, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Must specify a Name for Namespace updating")
	})

	t.Run("Update", func(t *testing.T) {
		ns, _, err := namespaces.Update(&Namespace{
			Name:        "foo",
			Description: "updated description",
		}, nil)

		require.NoError(t, err)
		require.NotNil(t, ns)
		require.Equal(t, "updated description", ns.Description)
	})

	t.Run("List", func(t *testing.T) {
		nsList, _, err := namespaces.List(nil)

		require.NoError(t, err)
		require.Len(t, nsList, 3)

		found := make(map[string]struct{})
		for _, ns := range nsList {
			found[ns.Name] = struct{}{}
		}

		require.Contains(t, found, "default")
		require.Contains(t, found, "foo")
		require.Contains(t, found, "acls")
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := namespaces.Delete("foo", nil)
		require.NoError(t, err)

		// due to deferred deletion the namespace might still exist
		// this checks that either it is in fact gone and we get a 404 or
		// that the namespace is still there but marked for deletion
		ns, _, err := namespaces.Read("foo", nil)
		if err != nil {
			require.Contains(t, err.Error(), "Unexpected response code: 404")
			require.Nil(t, ns)
		} else {
			require.NotNil(t, ns)
			require.NotNil(t, ns.DeletedAt)
			require.False(t, ns.DeletedAt.IsZero())
		}
	})
}
