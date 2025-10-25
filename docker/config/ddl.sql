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


-- 用户表
CREATE TABLE public.users
(
    id            bigserial      NOT NULL PRIMARY KEY,
    nickname      varchar(20)    NOT NULL,
    email         varchar(80)    NOT NULL UNIQUE,
    avatar_url    varchar(255)   NOT NULL,
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


-- 计划表
CREATE TABLE public.plans
(
    id          bigserial PRIMARY KEY,
    name        varchar(50)    NOT NULL UNIQUE,
    price       numeric(12, 2) NOT NULL DEFAULT 0,
    description text NOT NULL,
    created_at  timestamptz(6) NOT NULL DEFAULT now(),
    updated_at  timestamptz(6) NOT NULL DEFAULT now()
);


-- 租户表
CREATE TABLE public.tenants
(
    id          bigserial PRIMARY KEY,
    creator_id  bigint         NOT NULL REFERENCES public.users (id) ON DELETE RESTRICT,
    name        varchar(20)    NOT NULL,
    created_at  timestamptz(6) NOT NULL DEFAULT now(),
    updated_at  timestamptz(6) NOT NULL DEFAULT now(),
    description varchar(120)
);
CREATE INDEX IF NOT EXISTS idx_tenants_creator_id ON public.tenants (creator_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_tenants_creator_name ON public.tenants (creator_id, name);



-- 租户计划关联表
CREATE TABLE public.tenant_plan
(
    tenant_id   bigint         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    plan_id     bigint         NOT NULL REFERENCES public.plans (id),
    creator_id  bigint       NOT NULL REFERENCES public.users (id) ON DELETE RESTRICT,
    start_at    timestamptz(6) NOT NULL DEFAULT now(),
    end_at      timestamptz(6) NULL,
    CONSTRAINT ux_tenant_plan_tenant UNIQUE (tenant_id),
    PRIMARY KEY (tenant_id, plan_id)
);
-- 确保每个用户只能有一个Free和一个Caring
CREATE UNIQUE INDEX IF NOT EXISTS ux_user_one_free_plan ON public.tenant_plan (creator_id) WHERE plan_id = 1;
CREATE UNIQUE INDEX IF NOT EXISTS ux_user_one_caring_plan ON public.tenant_plan (creator_id) WHERE plan_id = 2;



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


-- 评论板块表
CREATE TABLE public.comment_plates
(
    id          bigserial      PRIMARY KEY,
    summary     varchar(60)    NOT NULL,
    belong_key  varchar(50)    NOT NULL,  -- 资源标识，如 "article:123"
    tenant_id   bigint         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    UNIQUE (tenant_id, belong_key)
);
-- 索引优化
CREATE INDEX IF NOT EXISTS idx_comment_plates_tenant_id ON public.comment_plates (tenant_id);
CREATE INDEX IF NOT EXISTS idx_comment_plates_belong_key ON public.comment_plates (belong_key);
CREATE INDEX idx_comment_plates_description_trgm ON public.comment_plates USING gin (description gin_trgm_ops);


-- 评论表
CREATE TYPE comment_status AS ENUM (
    'pending', -- 待审核
    'approved' -- 已通过
);
CREATE TABLE public.comments
(
    id         bigserial      PRIMARY KEY,
    tenant_id  bigint         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    plate_id   bigint         NOT NULL REFERENCES public.comment_plates (id) ON DELETE CASCADE,
    user_id    bigint         NOT NULL REFERENCES public.users (id) ON DELETE CASCADE,
    parent_id  bigint         NULL REFERENCES public.comments (id) ON DELETE CASCADE,
    root_id    bigint         NULL REFERENCES public.comments (id) ON DELETE CASCADE,
    content    text           NOT NULL,
    status     comment_status NOT NULL DEFAULT 'pending',
    like_count int8           NOT NULL DEFAULT 0,
    created_at timestamptz(6) NOT NULL DEFAULT now()
);
-- 按租户查询
CREATE INDEX IF NOT EXISTS idx_comments_tenant_id ON public.comments (tenant_id);
-- 查询板块的所有评论
CREATE INDEX IF NOT EXISTS idx_comments_tenant_plate_id ON public.comments (tenant_id, plate_id);  
-- 查询子评论
CREATE INDEX IF NOT EXISTS idx_comments_tenant_parent_id ON public.comments (tenant_id,parent_id);
-- 查询整个评论树
CREATE INDEX IF NOT EXISTS idx_comments_tenant_root_created_at ON public.comments (tenant_id, root_id, created_at);
-- 按状态查询
CREATE INDEX IF NOT EXISTS idx_comments_status ON public.comments (status);
-- 全文搜索内容
CREATE INDEX IF NOT EXISTS idx_comments_content_trgm ON public.comments USING gin (content gin_trgm_ops);
-- 热门评论查询（按点赞数）
CREATE INDEX IF NOT EXISTS idx_comments_like_count ON public.comments (like_count DESC);

CREATE INDEX idx_comments_tenant_plate_status_root_parent ON public.comments (tenant_id, plate_id, status, root_id, parent_id);


-- 评论租户全局配置（默认配置）
CREATE TABLE public.comment_tenant_configs
(
    tenant_id    bigint      NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE PRIMARY KEY,
    client_token text        NOT NULL,  -- 客户端令牌，用于 API 访问控制，防止接口被刷
    if_audit     boolean     NOT NULL DEFAULT true,   -- 是否开启审核
    created_at   timestamptz(6) NOT NULL DEFAULT now(),
    updated_at   timestamptz(6) NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_comment_tenant_configs_if_audit ON public.comment_tenant_configs (if_audit);

-- 评论板块配置（资源级别，高精细度）
CREATE TABLE public.comment_plate_configs
(
    tenant_id     bigint      NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    plate_id      bigint         NOT NULL REFERENCES public.comment_plates (id) ON DELETE CASCADE,
    if_audit      boolean     NOT NULL DEFAULT true,  -- 是否开启审核
    created_at    timestamptz(6) NOT NULL DEFAULT now(),
    updated_at    timestamptz(6) NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, plate_id)
);
CREATE INDEX IF NOT EXISTS idx_comment_plate_configs_tenant_id ON public.comment_plate_configs (tenant_id);
CREATE INDEX IF NOT EXISTS idx_comment_plate_configs_if_audit ON public.comment_plate_configs (if_audit);  
