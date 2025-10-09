-- 安装扩展
CREATE
    EXTENSION IF NOT EXISTS pg_trgm;


-- casbin_rule
CREATE TABLE public.casbin_rule
(
    id    serial4      NOT NULL,
    ptype varchar(100) NOT NULL,
    v0    varchar(100) NULL,
    v1    varchar(100) NULL,
    v2    varchar(100) NULL,
    v3    varchar(100) NULL,
    v4    varchar(100) NULL,
    v5    varchar(100) NULL,
    CONSTRAINT casbin_rule_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_casbin_rule_ptype ON public.casbin_rule USING btree (ptype);
CREATE INDEX idx_casbin_rule_v0 ON public.casbin_rule USING btree (v0);
CREATE INDEX idx_casbin_rule_v1 ON public.casbin_rule USING btree (v1);
CREATE INDEX IF NOT EXISTS idx_casbin_rule_ptype_v0_v1_v2 ON public.casbin_rule (ptype, v0, v1, v2);


-- API表
-- 定义项目所有的api 便于后续sys为权限分配api
-- CREATE TABLE public.apis
-- (
--     id          bigserial PRIMARY KEY,
--     method      varchar(10) NOT NULL,
--     path        varchar(60) NOT NULL,
--     description varchar(60) NULL,
--     CONSTRAINT api_check CHECK (
--         method IN ('get', 'post', 'put', 'delete', 'patch')
--         ),
--     CONSTRAINT apis_method_path_unique UNIQUE (method, path)
-- );
-- CREATE INDEX IF NOT EXISTS idx_apis_path ON public.apis (path);
-- CREATE INDEX IF NOT EXISTS idx_apis_path_trgm ON public.apis USING gin (path gin_trgm_ops);


-- 计划表
CREATE TABLE public.plans
(
    id          bigserial PRIMARY KEY,
    name        varchar(50)    NOT NULL UNIQUE,
    price       numeric(12, 2) NOT NULL DEFAULT 0,
    tenant_count int2 NOT NULL,
    number_count int4 NOT NULL,
    description text NOT NULL,
    created_at  timestamptz(6) NOT NULL DEFAULT now(),
    updated_at  timestamptz(6) NOT NULL DEFAULT now()
);


-- 租户表
CREATE TABLE public.tenants
(
    id          bigserial PRIMARY KEY,
    name        varchar(50)    NOT NULL,
    created_at  timestamptz(6) NOT NULL DEFAULT now(),
    updated_at  timestamptz(6) NOT NULL DEFAULT now(),
    description varchar(120)
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_tenants_name ON public.tenants (name);


-- 租户计划关联表
CREATE TABLE public.tenant_plan
(
    tenant_id bigint         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    plan_id   bigint         NOT NULL REFERENCES public.plans (id),
    start_at  timestamptz(6) NOT NULL DEFAULT now(),
    end_at    timestamptz(6) NULL,
    CONSTRAINT ux_tenant_plan_tenant UNIQUE (tenant_id),
    PRIMARY KEY (tenant_id, plan_id)
);


-- 角色表 
-- 与租户无关 应预定义 避免去实现租户层面的权限分配 因为这对于非超大型SaaS系统而言是一个灾难设计
CREATE TABLE public.roles
(
    id          bigserial PRIMARY KEY,
    name        varchar(30) NOT NULL,
    description text,
);
CREATE UNIQUE INDEX ux_roles_default_name ON public.roles (name);


-- 菜单表
-- CREATE TABLE public.menus
-- (
--     id        bigserial PRIMARY KEY,
--     tenant_id bigint      NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
--     name      varchar(50) NOT NULL,
--     path      varchar(100),
--     parent_id bigint,
--     sort      int,
--     UNIQUE (tenant_id, name)
-- );
-- -- 按 tenant 查询、按 parent 查子节点并按 sort 排序、path 模糊搜索
-- CREATE INDEX IF NOT EXISTS idx_menus_tenant_id ON public.menus (tenant_id);
-- CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON public.menus (parent_id);
-- CREATE INDEX IF NOT EXISTS idx_menus_tenant_parent_sort ON public.menus (tenant_id, parent_id, sort);
-- CREATE INDEX IF NOT EXISTS idx_menus_path_trgm ON public.menus USING gin (path gin_trgm_ops);


-- -- 角色菜单关联
-- CREATE TABLE public.role_menu
-- (
--     id      bigserial PRIMARY KEY,
--     role_id bigint NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
--     menu_id bigint NOT NULL REFERENCES menus (id) ON DELETE CASCADE,
--     UNIQUE (role_id, menu_id)
-- );


-- -- 按钮表
-- CREATE TABLE public.buttons
-- (
--     id        bigserial PRIMARY KEY,
--     tenant_id bigint      NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
--     menu_id   bigint      NOT NULL REFERENCES menus (id) ON DELETE CASCADE,
--     name      varchar(50) NOT NULL,
--     code      varchar(50) NOT NULL,
--     UNIQUE (tenant_id, menu_id, code)
-- );
-- -- 按 tenant/menu 查、按 menu_id 加速关联查询
-- CREATE INDEX IF NOT EXISTS idx_buttons_tenant_id ON public.buttons (tenant_id);
-- CREATE INDEX IF NOT EXISTS idx_buttons_menu_id ON public.buttons (menu_id);


-- -- 角色按钮关联
-- CREATE TABLE public.role_button
-- (
--     id        bigserial PRIMARY KEY,
--     role_id   bigint NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
--     button_id bigint NOT NULL REFERENCES buttons (id) ON DELETE CASCADE,
--     UNIQUE (role_id, button_id)
-- );


-- 用户表
CREATE TABLE public.users
(
    id            bigserial      NOT NULL PRIMARY KEY,
    nickname      varchar(20)    NOT NULL,
    email         varchar(80)    NOT NULL UNIQUE,
    github_id     varchar(60)    NULL UNIQUE,
    google_id     varchar(60)    NULL UNIQUE,
    password_hash text           NULL,
    last_login_at timestamptz(6) NOT NULL,
    created_at    timestamptz(6) NOT NULL DEFAULT now(),
    updated_at    timestamptz(6) NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON public.users (created_at);
CREATE INDEX IF NOT EXISTS idx_users_updated_at ON public.users (updated_at);
CREATE INDEX IF NOT EXISTS idx_users_nickname ON public.users (nickname);


-- 租户-用户-角色 关联表
-- 一个用户在一个租户下只能有一个角色
CREATE TABLE public.tenant_user_role
(
    user_id   bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    tenant_id bigint NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    role_id   bigint NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tenant_id)
);
-- 常按 user 查租户、按 tenant 查用户；避免重复记录
CREATE INDEX IF NOT EXISTS idx_tenant_user_role_user_id ON public.tenant_user_role (user_id);
CREATE INDEX IF NOT EXISTS idx_tenant_user_role_tenant_id ON public.tenant_user_role (tenant_id);
CREATE INDEX IF NOT EXISTS idx_tenant_user_role_role_id ON public.tenant_user_role (role_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_tenant_user_role_user_tenant ON public.tenant_user_role (user_id, tenant_id);


-- img_categories
CREATE TABLE public.img_categories
(
    id         bigserial PRIMARY KEY,
    tenant_id  bigint         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    title      varchar(10)    NOT NULL,
    prefix     varchar(20)    NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, title)
);


-- img 表
CREATE TABLE public.imgs
(
    id          bigserial PRIMARY KEY,
    tenant_id   bigint         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    path        varchar(120)   NOT NULL,
    description varchar(60),
    created_at  timestamptz(6) NOT NULL DEFAULT now(),
    updated_at  timestamptz(6) NOT NULL,
    deleted_at  timestamptz(6),
    category_id bigint REFERENCES img_categories (id),
    UNIQUE (tenant_id, path)
);
CREATE INDEX idx_img_deleted_at ON public.imgs (deleted_at);
CREATE INDEX idx_img_description_trgm ON public.imgs USING gin (description gin_trgm_ops);


-- 租户r2配置表
CREATE TABLE public.tenant_r2_configs (
    tenant_id bigint NOT NULL REFERENCES public.tenants(id) ON DELETE CASCADE PRIMARY KEY,  
    account_id varchar(255) NOT NULL,
    access_key_id varchar(255) NOT NULL,
    secret_access_key varchar(255) NOT NULL,  -- 加密存储
    public_bucket varchar(255) NOT NULL,
    public_url_prefix varchar(500) NOT NULL,
    delete_bucket varchar(255) NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT now(),
    updated_at timestamptz(6) NOT NULL DEFAULT now()
);


-- 评论表
CREATE TABLE public.comments
(
    id         bigserial PRIMARY KEY,
    belong_key varchar(50)    NOT NULL,  -- 资源标识，如 "article:123"
    tenant_id  bigint         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    user_id    bigint         NOT NULL REFERENCES public.users (id) ON DELETE CASCADE,
    parent_id  bigint         NULL REFERENCES public.comments (id) ON DELETE CASCADE,
    root_id    bigint         NULL REFERENCES public.comments (id) ON DELETE CASCADE,
    content    text           NOT NULL,
    status     varchar(10)    NOT NULL DEFAULT 'pending',
    like_count int4           NOT NULL DEFAULT 0,
    created_at timestamptz(6) NOT NULL DEFAULT now(),
    CONSTRAINT comments_status_check CHECK (
        status IN ('pending', 'approved', 'rejected')
    )
);
-- 按租户查询
CREATE INDEX IF NOT EXISTS idx_comments_tenant_id ON public.comments (tenant_id);
-- 查询资源的所有评论
CREATE INDEX IF NOT EXISTS idx_comments_tenant_belong_key ON public.comments (tenant_id, belong_key);  
-- 查询子评论
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON public.comments (parent_id);
-- 查询整个评论树
CREATE INDEX IF NOT EXISTS idx_comments_root_id ON public.comments (root_id, created_at);
-- 按状态查询
CREATE INDEX IF NOT EXISTS idx_comments_status ON public.comments (status);
-- 全文搜索内容
CREATE INDEX IF NOT EXISTS idx_comments_content_trgm ON public.comments USING gin (content gin_trgm_ops);
-- 按时间排序查询
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON public.comments (created_at DESC);
-- 热门评论查询（按点赞数）
CREATE INDEX IF NOT EXISTS idx_comments_like_count ON public.comments (like_count DESC);


-- 评论点赞表
CREATE TABLE public.comment_likes
(
    comment_id bigint         NOT NULL REFERENCES public.comments (id) ON DELETE CASCADE,
    user_id    bigint         NOT NULL REFERENCES public.users (id) ON DELETE CASCADE,
    created_at timestamptz(6) NOT NULL DEFAULT now(),
    PRIMARY KEY (comment_id, user_id)  -- 确保一个用户对一个评论只能点赞一次
);
-- 索引优化
CREATE INDEX IF NOT EXISTS idx_comment_likes_user_id ON public.comment_likes (user_id);  -- 用户的点赞列表
CREATE INDEX IF NOT EXISTS idx_comment_likes_comment_id ON public.comment_likes (comment_id);  -- 评论的点赞列表


-- 评论租户全局配置（默认配置）
CREATE TABLE public.comment_tenant_configs
(
    tenant_id    bigint      NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE PRIMARY KEY,
    client_token text        NOT NULL,  -- 客户端令牌，用于 API 访问控制，防止接口被刷
    if_audit     boolean     NOT NULL DEFAULT true,   -- 默认是否开启审核
    created_at   timestamptz(6) NOT NULL DEFAULT now(),
    updated_at   timestamptz(6) NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_comment_tenant_configs_if_audit ON public.comment_tenant_configs (if_audit);

-- 评论板块配置（资源级别，高精细度）
CREATE TABLE public.comment_configs
(
    tenant_id     bigint      NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    belong_key    varchar(50) NOT NULL,  -- 与 comments.belong_key 对应
    client_token text         NOT NULL,  -- 客户端令牌，用于 API 访问控制，防止接口被刷
    if_audit      boolean     NOT NULL DEFAULT true,  -- 是否开启审核
    created_at    timestamptz(6) NOT NULL DEFAULT now(),
    updated_at    timestamptz(6) NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, belong_key)
);
CREATE INDEX IF NOT EXISTS idx_comment_configs_tenant_id ON public.comment_configs (tenant_id);
CREATE INDEX IF NOT EXISTS idx_comment_configs_if_audit ON public.comment_configs (if_audit);  
