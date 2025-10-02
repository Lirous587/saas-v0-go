package handler

import (
    "saas/internal/role/domain"
)

type RoleResponse struct {
    ID          int64  `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
    CreatedAt   int64  `json:"created_at"`
    UpdatedAt   int64  `json:"updated_at"`
}

func domainRoleToResponse(role *domain.Role) *RoleResponse {
    if role == nil {
        return nil
    }

    return &RoleResponse{
        ID:          role.ID,
        Title:       role.Title,
        Description: role.Description,
        CreatedAt:   role.CreatedAt.Unix(),
        UpdatedAt:   role.UpdatedAt.Unix(),
    }
}

func domainRolesToResponse(roles []*domain.Role) []*RoleResponse {
    if len(roles) == 0 {
        return nil
    }

    ret := make([]*RoleResponse, 0, len(roles))

    for _, role := range roles {
        if role != nil {
            ret = append(ret, domainRoleToResponse(role))
        }
    }
    return ret
}

type CreateRequest struct {
    Title       string  `json:"title" binding:"required,max=30"`
    Description string  `json:"description" binding:"max=60"`
}

type UpdateRequest struct {
    Title       string  `json:"title" binding:"required,max=30"`
    Description string  `json:"description" binding:"max=60"`
}

type ListRequest struct {
    Page     int    `form:"page,default=1" binding:"min=1"`
    PageSize int    `form:"page_size,default=5" binding:"min=5,max=20"`
    KeyWord  string `form:"keyword" binding:"max=20"`
}

type RoleListResponse struct {
    Total int64                         `json:"total"`
    List  []*RoleResponse   `json:"list"`
}

func domainRoleListToResponse(data *domain.RoleList) *RoleListResponse {
    if data == nil {
        return nil
    }

    return &RoleListResponse{
        Total: data.Total,
        List:  domainRolesToResponse(data.List),
    }
}
