package domain

type TenantService interface {
	Create(tenant *Tenant, planID int64, userID int64) (*Tenant, error)
	Read(id int64) (*Tenant, error)
	Update(tenant *Tenant) (*Tenant, error)
	Delete(id int64) error
	List(query *TenantQuery) (*TenantList, error)

	GenInviteToken(payload *GenInviteTokenPayload) (string, error)
	Invite(payload *InvitePayload) error
	Enter(paylod *EnterPayload) error

	CheckRoleValidity(roleID int64) error
	ListUsersWithRole(query *UserWithRoleQuery) (*UserWithRoleList, error)
}
