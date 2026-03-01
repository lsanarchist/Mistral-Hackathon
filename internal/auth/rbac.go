package auth

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// RBACManager handles role-based access control for enterprise features
type RBACManager struct {
	mu            sync.RWMutex
	roles         map[string]model.Role
	permissions   map[string]model.Permission
	userRoles     map[string][]string // userID -> roleIDs
	teamMembers   map[string][]string // teamID -> userIDs
}

// NewRBACManager creates a new RBAC manager with default roles and permissions
func NewRBACManager() *RBACManager {
	r := &RBACManager{
		roles:         make(map[string]model.Role),
		permissions:   make(map[string]model.Permission),
		userRoles:     make(map[string][]string),
		teamMembers:   make(map[string][]string),
	}

	// Initialize default permissions
	r.AddPermission(model.Permission{
		ID:          "run_analysis",
		Name:        "Run Performance Analysis",
		Description: "Allow running performance analysis commands",
	})
	r.AddPermission(model.Permission{
		ID:          "view_reports",
		Name:        "View Reports",
		Description: "Allow viewing performance reports and findings",
	})
	r.AddPermission(model.Permission{
		ID:          "manage_users",
		Name:        "Manage Users",
		Description: "Allow adding/removing users and assigning roles",
	})
	r.AddPermission(model.Permission{
		ID:          "manage_teams",
		Name:        "Manage Teams",
		Description: "Allow creating/managing teams",
	})
	r.AddPermission(model.Permission{
		ID:          "configure_system",
		Name:        "Configure System",
		Description: "Allow modifying system configuration",
	})

	// Initialize default roles
	r.AddRole(model.Role{
		ID:          "admin",
		Name:        "Administrator",
		Permissions: []string{"run_analysis", "view_reports", "manage_users", "manage_teams", "configure_system"},
	})
	r.AddRole(model.Role{
		ID:          "analyst",
		Name:        "Performance Analyst",
		Permissions: []string{"run_analysis", "view_reports"},
	})
	r.AddRole(model.Role{
		ID:          "viewer",
		Name:        "Viewer",
		Permissions: []string{"view_reports"},
	})

	return r
}

// AddPermission adds a new permission to the system
func (r *RBACManager) AddPermission(perm model.Permission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.permissions[perm.ID]; exists {
		return fmt.Errorf("permission %s already exists", perm.ID)
	}

	r.permissions[perm.ID] = perm
	return nil
}

// AddRole adds a new role to the system
func (r *RBACManager) AddRole(role model.Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.roles[role.ID]; exists {
		return fmt.Errorf("role %s already exists", role.ID)
	}

	// Validate permissions exist
	for _, permID := range role.Permissions {
		if _, exists := r.permissions[permID]; !exists {
			return fmt.Errorf("permission %s does not exist", permID)
		}
	}

	r.roles[role.ID] = role
	return nil
}

// AssignRoleToUser assigns a role to a user
func (r *RBACManager) AssignRoleToUser(userID, roleID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.roles[roleID]; !exists {
		return fmt.Errorf("role %s does not exist", roleID)
	}

	r.userRoles[userID] = append(r.userRoles[userID], roleID)
	return nil
}

// RemoveRoleFromUser removes a role from a user
func (r *RBACManager) RemoveRoleFromUser(userID, roleID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if roles, exists := r.userRoles[userID]; exists {
		for i, rID := range roles {
			if rID == roleID {
				r.userRoles[userID] = append(roles[:i], roles[i+1:]...)
				return nil
			}
		}
	}

	return fmt.Errorf("user %s does not have role %s", userID, roleID)
}

// CreateTeam creates a new team
func (r *RBACManager) CreateTeam(team model.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.teamMembers[team.ID]; exists {
		return fmt.Errorf("team %s already exists", team.ID)
	}

	r.teamMembers[team.ID] = team.Members
	return nil
}

// AddUserToTeam adds a user to a team
func (r *RBACManager) AddUserToTeam(teamID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.teamMembers[teamID]; !exists {
		return fmt.Errorf("team %s does not exist", teamID)
	}

	r.teamMembers[teamID] = append(r.teamMembers[teamID], userID)
	return nil
}

// RemoveUserFromTeam removes a user from a team
func (r *RBACManager) RemoveUserFromTeam(teamID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if members, exists := r.teamMembers[teamID]; exists {
		for i, mID := range members {
			if mID == userID {
				r.teamMembers[teamID] = append(members[:i], members[i+1:]...)
				return nil
			}
		}
	}

	return fmt.Errorf("user %s is not in team %s", userID, teamID)
}

// HasPermission checks if a user has a specific permission
func (r *RBACManager) HasPermission(userID, permissionID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get user's roles
	roleIDs, exists := r.userRoles[userID]
	if !exists {
		return false
	}

	// Check each role for the permission
	for _, roleID := range roleIDs {
		if role, exists := r.roles[roleID]; exists {
			for _, permID := range role.Permissions {
				if permID == permissionID {
					return true
				}
			}
		}
	}

	return false
}

// GetUserPermissions returns all permissions for a user
func (r *RBACManager) GetUserPermissions(userID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	permissions := make(map[string]bool)
	roleIDs, exists := r.userRoles[userID]
	if !exists {
		return []string{}
	}

	for _, roleID := range roleIDs {
		if role, exists := r.roles[roleID]; exists {
			for _, permID := range role.Permissions {
				permissions[permID] = true
			}
		}
	}

	result := make([]string, 0, len(permissions))
	for permID := range permissions {
		result = append(result, permID)
	}

	return result
}

// GetUserRoles returns all roles for a user
func (r *RBACManager) GetUserRoles(userID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if roles, exists := r.userRoles[userID]; exists {
		return roles
	}

	return []string{}
}

// GetTeamMembers returns all members of a team
func (r *RBACManager) GetTeamMembers(teamID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if members, exists := r.teamMembers[teamID]; exists {
		return members
	}

	return []string{}
}

// GetUserTeams returns all teams a user belongs to
func (r *RBACManager) GetUserTeams(userID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var teams []string
	for teamID, members := range r.teamMembers {
		for _, memberID := range members {
			if memberID == userID {
				teams = append(teams, teamID)
				break
			}
		}
	}

	return teams
}

// CheckTeamPermission checks if any member of a team has a specific permission
func (r *RBACManager) CheckTeamPermission(teamID, permissionID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if members, exists := r.teamMembers[teamID]; exists {
		for _, memberID := range members {
			if r.HasPermission(memberID, permissionID) {
				return true
			}
		}
	}

	return false
}

// GetAllRoles returns all roles in the system
func (r *RBACManager) GetAllRoles() []model.Role {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roles := make([]model.Role, 0, len(r.roles))
	for _, role := range r.roles {
		roles = append(roles, role)
	}

	return roles
}

// GetAllPermissions returns all permissions in the system
func (r *RBACManager) GetAllPermissions() []model.Permission {
	r.mu.RLock()
	defer r.mu.RUnlock()

	perms := make([]model.Permission, 0, len(r.permissions))
	for _, perm := range r.permissions {
		perms = append(perms, perm)
	}

	return perms
}

// GetAllTeams returns all teams in the system
func (r *RBACManager) GetAllTeams() []model.Team {
	r.mu.RLock()
	defer r.mu.RUnlock()

	teams := make([]model.Team, 0, len(r.teamMembers))
	for teamID, members := range r.teamMembers {
		teams = append(teams, model.Team{
			ID:       teamID,
			Members:  members,
		})
	}

	return teams
}

// GetAllUsers returns all users in the system
func (r *RBACManager) GetAllUsers() []model.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.User, 0, len(r.userRoles))
	for userID, roleIDs := range r.userRoles {
		users = append(users, model.User{
			ID:    userID,
			Roles: roleIDs,
		})
	}

	return users
}

// CheckPermissionWithContext checks permission with additional context (team membership, etc.)
func (r *RBACManager) CheckPermissionWithContext(userID, permissionID string, context map[string]interface{}) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Direct user permission check
	if r.HasPermission(userID, permissionID) {
		return true
	}

	// Check if user is member of the team in context
	if teamID, ok := context["team"].(string); ok {
		// Check if user is actually a member of this team
		if members, exists := r.teamMembers[teamID]; exists {
			for _, memberID := range members {
				if memberID == userID {
					// User is a member of this team, check team permissions
					if r.CheckTeamPermission(teamID, permissionID) {
						return true
					}
					break
				}
			}
		}
	}

	return false
}

// String representation of RBAC configuration
func (r *RBACManager) String() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString("RBAC Configuration:\n")
	sb.WriteString(fmt.Sprintf("  Roles: %d\n", len(r.roles)))
	sb.WriteString(fmt.Sprintf("  Permissions: %d\n", len(r.permissions)))
	sb.WriteString(fmt.Sprintf("  Users: %d\n", len(r.userRoles)))
	sb.WriteString(fmt.Sprintf("  Teams: %d\n", len(r.teamMembers)))

	return sb.String()
}