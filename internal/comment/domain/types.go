package domain

type TenantID string

func (t TenantID) IsZero() bool {
	return t == ""
}

type UserID string

func (u UserID) IsZero() bool {
	return u == ""
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

type CommentIDs []CommentID

func (ids CommentIDs) ToStringSlice() []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = string(id)
	}
	return result
}

type PlateID string

func (p PlateID) IsZero() bool {
	return p == ""
}
