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
CREATE TABLE public.apis
(
    id          bigserial PRIMARY KEY,
    method      varchar(10) NOT NULL,
    path        varchar(60) NOT NULL,
    description varchar(60) NULL,
    CONSTRAINT api_check CHECK (
        method IN ('get', 'post', 'put', 'delete', 'patch')
        ),
    CONSTRAINT apis_method_path_unique UNIQUE (method, path)
);
CREATE INDEX IF NOT EXISTS idx_apis_path ON public.apis (path);
CREATE INDEX IF NOT EXISTS idx_apis_path_trgm ON public.apis USING gin (path gin_trgm_ops);


-- 计划表
CREATE TABLE public.plans
(
    id          bigserial PRIMARY KEY,
    name        varchar(50)    NOT NULL UNIQUE,
    price       numeric(12, 2) NOT NULL DEFAULT 0,
    description text,
    created_at  timestamptz(6) NOT NULL DEFAULT now()
);

-- 租户表
CREATE TABLE public.tenants
(
    id   bigserial PRIMARY KEY,
    name varchar(50) NOT NULL
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
CREATE TABLE public.roles
(
    id          bigserial PRIMARY KEY,
    tenant_id   bigint      NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    name        varchar(30) NOT NULL,
    description text,
    CONSTRAINT ux_roles_tenant_name UNIQUE (tenant_id, name)
);
CREATE INDEX IF NOT EXISTS idx_roles_tenant_id ON public.roles (tenant_id);


-- 菜单表
CREATE TABLE public.menus
(
    id        bigserial PRIMARY KEY,
    tenant_id bigint      NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    name      varchar(50) NOT NULL,
    path      varchar(100),
    parent_id bigint,
    sort      int,
    UNIQUE (tenant_id, name)
);
-- 按 tenant 查询、按 parent 查子节点并按 sort 排序、path 模糊搜索
CREATE INDEX IF NOT EXISTS idx_menus_tenant_id ON public.menus (tenant_id);
CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON public.menus (parent_id);
CREATE INDEX IF NOT EXISTS idx_menus_tenant_parent_sort ON public.menus (tenant_id, parent_id, sort);
CREATE INDEX IF NOT EXISTS idx_menus_path_trgm ON public.menus USING gin (path gin_trgm_ops);


-- 角色菜单关联
CREATE TABLE public.role_menu
(
    id      bigserial PRIMARY KEY,
    role_id bigint NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    menu_id bigint NOT NULL REFERENCES menus (id) ON DELETE CASCADE,
    UNIQUE (role_id, menu_id)
);


-- 按钮表
CREATE TABLE public.buttons
(
    id        bigserial PRIMARY KEY,
    tenant_id bigint      NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    menu_id   bigint      NOT NULL REFERENCES menus (id) ON DELETE CASCADE,
    name      varchar(50) NOT NULL,
    code      varchar(50) NOT NULL,
    UNIQUE (tenant_id, menu_id, code)
);
-- 按 tenant/menu 查、按 menu_id 加速关联查询
CREATE INDEX IF NOT EXISTS idx_buttons_tenant_id ON public.buttons (tenant_id);
CREATE INDEX IF NOT EXISTS idx_buttons_menu_id ON public.buttons (menu_id);

-- 角色按钮关联
CREATE TABLE public.role_button
(
    id        bigserial PRIMARY KEY,
    role_id   bigint NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    button_id bigint NOT NULL REFERENCES buttons (id) ON DELETE CASCADE,
    UNIQUE (role_id, button_id)
);

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

-- 用户-租户关联表
-- 一个用户在一个租户下只能有一个角色
CREATE TABLE public.user_tenants
(
    user_id   bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    tenant_id bigint NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    role_id   bigint NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tenant_id)
);
-- 常按 user 查租户、按 tenant 查用户；避免重复记录
CREATE INDEX IF NOT EXISTS idx_user_tenants_user_id ON public.user_tenants (user_id);
CREATE INDEX IF NOT EXISTS idx_user_tenants_tenant_id ON public.user_tenants (tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_tenants_role_id ON public.user_tenants (role_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_user_tenants_user_tenant ON public.user_tenants (user_id, tenant_id);

-- img_categories
CREATE TABLE public.img_categories
(
    id         bigserial PRIMARY KEY,
    tenant_id  bigint         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    title      varchar(10)    NOT NULL UNIQUE,
    prefix     varchar(20)    NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT now()
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
    category_id bigint REFERENCES img_categories (id)
);
CREATE UNIQUE INDEX idx_img_path ON public.imgs (path);
CREATE INDEX idx_img_deleted_at ON public.imgs (deleted_at);
CREATE INDEX idx_img_description_trgm ON public.imgs USING gin (description gin_trgm_ops);