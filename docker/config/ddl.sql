-- 安装扩展
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 用户表
CREATE TABLE public.users
(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
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



-- 租户表
CREATE TYPE tenant_plan_type AS ENUM ('free', 'care','pro');
CREATE TYPE tenant_status AS ENUM ('active', 'inactive');
CREATE TYPE tnant_plan_billing_cycle AS ENUM ('monthly', 'yearly', 'lifetime');
CREATE TABLE public.tenants
(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    creator_id     UUID               NOT NULL REFERENCES public.users (id) ON DELETE RESTRICT,
    plan_type      tenant_plan_type     NOT NULL DEFAULT 'free',
    name           varchar(20)          NOT NULL,
    status         tenant_status        NOT NULL DEFAULT 'active',
    billing_cycle  tnant_plan_billing_cycle  NOT NULL DEFAULT 'monthly',
    description    text         NULL,
    start_at       timestamptz(6)       NOT NULL DEFAULT now(),
    end_at         timestamptz(6)       NULL,
    created_at     timestamptz(6)       NOT NULL DEFAULT now(),
    updated_at     timestamptz(6)       NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_tenants_creator_id ON public.tenants (creator_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_tenants_creator_name ON public.tenants (creator_id, name);
-- 确保每个用户只能有一个Free和一个Care
CREATE UNIQUE INDEX IF NOT EXISTS ux_user_one_free_plan ON public.tenants (creator_id) WHERE plan_type = 'free';
CREATE UNIQUE INDEX IF NOT EXISTS ux_user_one_care_plan ON public.tenants (creator_id) WHERE plan_type = 'care';



-- -- 租户计划历史表
-- CREATE TABLE public.tenant_plan_history
-- (
--     id            bigserial PRIMARY KEY,
--     tenant_id     UUID         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
--     old_plan_type tenant_plan_type,
--     new_plan_type tenant_plan_type NOT NULL,
--     upgraded_at   timestamptz(6) NOT NULL DEFAULT now()
-- );

-- 租户限制配置表
-- CREATE TABLE public.tenant_limits
-- (
--     plan_id    UUID NOT NULL REFERENCES public.plans (id) ON DELETE CASCADE PRIMARY KEY,
--     api_calls  int    NOT NULL DEFAULT 1000,
--     plates     int    NOT NULL DEFAULT 5,
--     created_at timestamptz(6) NOT NULL DEFAULT now(),
--     updated_at timestamptz(6) NOT NULL DEFAULT now()
-- );



-- img_categories
CREATE TABLE public.img_categories
(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id  UUID         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    title      varchar(10)    NOT NULL,
    prefix     varchar(20)    NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, title)
);



-- img 表
CREATE TABLE public.imgs
(
    id          UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id   UUID         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    category_id UUID REFERENCES img_categories (id),
    path        text   NOT NULL,
    description varchar(60),
    created_at  timestamptz(6) NOT NULL DEFAULT now(),
    updated_at  timestamptz(6) NOT NULL DEFAULT now(),
    deleted_at  timestamptz(6),
    UNIQUE (tenant_id, path)
);
CREATE INDEX idx_img_deleted_at ON public.imgs (deleted_at);
CREATE INDEX idx_img_description_trgm ON public.imgs USING gin (description gin_trgm_ops);



-- 租户r2配置表
CREATE TABLE public.tenant_r2_configs (
    tenant_id UUID NOT NULL REFERENCES public.tenants(id) ON DELETE CASCADE PRIMARY KEY,  
    account_id varchar(32) NOT NULL,
    access_key_id varchar(32) NOT NULL,
    secret_access_key text  NULL,  -- 加密存储
    public_bucket varchar(32) NOT NULL,
    public_url_prefix varchar(128) NOT NULL,
    delete_bucket varchar(32) NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT now(),
    updated_at timestamptz(6) NOT NULL DEFAULT now()
);



-- 评论板块表
CREATE TABLE public.comment_plates
(
    id          UUID PRIMARY KEY DEFAULT uuidv7(),
    belong_key  varchar(50)    NOT NULL,  -- 资源标识，如 "article:123"
    related_url varchar(255) NOT NULL,
    summary     text    NOT NULL,
    tenant_id   UUID         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    UNIQUE (tenant_id, belong_key)
);
-- 索引优化
CREATE INDEX IF NOT EXISTS idx_comment_plates_tenant_id ON public.comment_plates (tenant_id);
CREATE INDEX IF NOT EXISTS idx_comment_plates_belong_key ON public.comment_plates (belong_key);
CREATE INDEX idx_comment_plates_summary_trgm ON public.comment_plates USING gin (summary gin_trgm_ops);



-- 评论表
CREATE TYPE comment_status AS ENUM (
    'pending', -- 待审核
    'approved' -- 已通过
);
CREATE TABLE public.comments
(
    id         UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id  UUID         NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    plate_id   UUID         NOT NULL REFERENCES public.comment_plates (id) ON DELETE CASCADE,
    user_id    UUID         NOT NULL REFERENCES public.users (id) ON DELETE CASCADE,
    parent_id  UUID         NULL REFERENCES public.comments (id) ON DELETE CASCADE,
    root_id    UUID         NULL REFERENCES public.comments (id) ON DELETE CASCADE,
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


-- 评论点赞表
CREATE TABLE public.comment_likes (
    tenant_id UUID NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    comment_id UUID NOT NULL REFERENCES public.comments (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES public.users (id) ON DELETE CASCADE,
    created_at timestamptz(6) NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, user_id, comment_id)
);

-- 索引优化
CREATE INDEX IF NOT EXISTS idx_comment_likes_user_like ON public.comment_likes (tenant_id, user_id);
CREATE INDEX IF NOT EXISTS idx_comment_likes_comment_like ON public.comment_likes (tenant_id, comment_id);
CREATE INDEX IF NOT EXISTS idx_comment_likes_created_at ON public.comment_likes (created_at DESC);


-- 评论租户全局配置（默认配置）
CREATE TABLE public.comment_tenant_configs
(
    tenant_id    UUID      NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE PRIMARY KEY,
    if_audit     boolean     NOT NULL DEFAULT true,   -- 是否开启审核
    created_at   timestamptz(6) NOT NULL DEFAULT now(),
    updated_at   timestamptz(6) NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_comment_tenant_configs_if_audit ON public.comment_tenant_configs (if_audit);

-- 评论板块配置（资源级别，高精细度）
CREATE TABLE public.comment_plate_configs
(
    plate_id      UUID           NOT NULL REFERENCES public.comment_plates (id) ON DELETE CASCADE,
    tenant_id     UUID           NOT NULL REFERENCES public.tenants (id) ON DELETE CASCADE,
    if_audit      boolean        NOT NULL DEFAULT true,  -- 是否开启审核
    created_at    timestamptz(6) NOT NULL DEFAULT now(),
    updated_at    timestamptz(6) NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, plate_id)
);
CREATE INDEX IF NOT EXISTS idx_comment_plate_configs_tenant_id ON public.comment_plate_configs (tenant_id);
CREATE INDEX IF NOT EXISTS idx_comment_plate_configs_if_audit ON public.comment_plate_configs (if_audit);  
