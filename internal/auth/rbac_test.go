package auth

import (
	"testing"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRBACManager(t *testing.T) {
	r := NewRBACManager()

	// Test default permissions and roles
	assert.Equal(t, 5, len(r.GetAllPermissions()))
	assert.Equal(t, 3, len(r.GetAllRoles()))

	// Test adding a new permission
	err := r.AddPermission(model.Permission{
		ID:          "test_permission",
		Name:        "Test Permission",
		Description: "A test permission",
	})
	assert.NoError(t, err)
	assert.Equal(t, 6, len(r.GetAllPermissions()))

	// Test adding a duplicate permission
	err = r.AddPermission(model.Permission{
		ID:          "test_permission",
		Name:        "Duplicate Test Permission",
		Description: "A duplicate test permission",
	})
	assert.Error(t, err)

	// Test adding a new role
	err = r.AddRole(model.Role{
		ID:          "test_role",
		Name:        "Test Role",
		Permissions: []string{"test_permission"},
	})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(r.GetAllRoles()))

	// Test assigning role to user
	err = r.AssignRoleToUser("user1", "test_role")
	assert.NoError(t, err)

	// Test user has permission
	assert.True(t, r.HasPermission("user1", "test_permission"))
	assert.False(t, r.HasPermission("user1", "run_analysis"))

	// Test getting user permissions
	permissions := r.GetUserPermissions("user1")
	assert.Equal(t, 1, len(permissions))
	assert.Equal(t, "test_permission", permissions[0])

	// Test team functionality
	err = r.CreateTeam(model.Team{
		ID:       "team1",
		Name:     "Test Team",
		Members:  []string{"user1"},
	})
	assert.NoError(t, err)

	// Test adding user to team
	err = r.AddUserToTeam("team1", "user2")
	assert.NoError(t, err)

	// Test team members
	members := r.GetTeamMembers("team1")
	assert.Equal(t, 2, len(members))

	// Test user teams
	teams := r.GetUserTeams("user1")
	assert.Equal(t, 1, len(teams))
	assert.Equal(t, "team1", teams[0])

	// Test removing role from user
	err = r.RemoveRoleFromUser("user1", "test_role")
	assert.NoError(t, err)
	assert.False(t, r.HasPermission("user1", "test_permission"))

	// Test removing user from team
	err = r.RemoveUserFromTeam("team1", "user1")
	assert.NoError(t, err)
	members = r.GetTeamMembers("team1")
	assert.Equal(t, 1, len(members))
}

func TestRBACPermissionChecking(t *testing.T) {
	r := NewRBACManager()

	// Test admin role has all permissions
	err := r.AssignRoleToUser("admin_user", "admin")
	assert.NoError(t, err)

	permissions := r.GetUserPermissions("admin_user")
	assert.True(t, len(permissions) >= 5) // Should have all default permissions

	// Test analyst role has limited permissions
	err = r.AssignRoleToUser("analyst_user", "analyst")
	assert.NoError(t, err)

	analystPermissions := r.GetUserPermissions("analyst_user")
	assert.True(t, len(analystPermissions) >= 2) // Should have run_analysis and view_reports
	assert.False(t, r.HasPermission("analyst_user", "manage_users"))

	// Test viewer role has only view permissions
	err = r.AssignRoleToUser("viewer_user", "viewer")
	assert.NoError(t, err)

	viewerPermissions := r.GetUserPermissions("viewer_user")
	assert.Equal(t, 1, len(viewerPermissions))
	assert.True(t, r.HasPermission("viewer_user", "view_reports"))
	assert.False(t, r.HasPermission("viewer_user", "run_analysis"))
}

func TestRBACContextChecking(t *testing.T) {
	r := NewRBACManager()

	// Create a user with no direct permissions
	err := r.AssignRoleToUser("team_user", "viewer")
	assert.NoError(t, err)

	// Create a team with an admin user
	err = r.CreateTeam(model.Team{
		ID:      "admin_team",
		Members: []string{"admin_user"},
	})
	assert.NoError(t, err)

	err = r.AssignRoleToUser("admin_user", "admin")
	assert.NoError(t, err)

	// Test context-based permission checking
	context := map[string]interface{}{
		"team": "admin_team",
	}

	// Team user should not have admin permissions directly
	assert.False(t, r.HasPermission("team_user", "manage_users"))

	// But with team context, they should have team permissions
	assert.True(t, r.CheckTeamPermission("admin_team", "manage_users"))

	// Context checking should work
	assert.True(t, r.CheckPermissionWithContext("admin_user", "manage_users", context))
	assert.False(t, r.CheckPermissionWithContext("team_user", "manage_users", context))
}