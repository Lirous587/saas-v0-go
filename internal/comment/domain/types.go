package domain

type TenantID string

func (t TenantID) IsZero() bool {
	return t == ""
}

func (t TenantID) String() string {
	return string(t)
}

type UserID string

func (u UserID) IsZero() bool {
	return u == ""
}

func (u UserID) String() string {
	return string(u)
}

type UserIDs []UserID

func NewUserIDsFromStrings(ids []string) UserIDs {
	result := make(UserIDs, len(ids))
	for i, id := range ids {
		result[i] = UserID(id)
	}
	return result
}

func (ids UserIDs) ToStringSlice() []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = string(id)
	}
	return result
}

type CommentID string

func (c CommentID) IsZero() bool {
	return c == ""
}

func (c CommentID) String() string {
	return string(c)
}

type CommentIDs []CommentID

func (ids CommentIDs) ToStringSlice() []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = string(id)
	}
	return result
}

func (ids CommentIDs) ToMap() map[CommentID]struct{} {
	if len(ids) == 0 {
		return nil
	}

	cMap := make(map[CommentID]struct{}, len(ids))

	for i := range ids {
		cMap[ids[i]] = struct{}{}
	}

	return cMap
}

type PlateID string

func (p PlateID) IsZero() bool {
	return p == ""
}

func (p PlateID) String() string {
	return string(p)
}
