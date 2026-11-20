
CREATE EXTENSION IF NOT EXISTS vector;


CREATE TABLE tenants (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name       text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE documents (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    filename     text NOT NULL,
    content_hash text NOT NULL,
    status       text NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending','processing','done','failed')),
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, content_hash),
    UNIQUE (id, tenant_id)
);

CREATE TABLE chunks (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL,
    document_id uuid NOT NULL,
    chunk_index int  NOT NULL,
    content     text NOT NULL,
    embedding   vector(768) NOT NULL,
    tsv         tsvector GENERATED ALWAYS AS (to_tsvector('english', content)) STORED,
    UNIQUE (document_id, chunk_index),
    FOREIGN KEY (document_id, tenant_id)
        REFERENCES documents(id, tenant_id) ON DELETE CASCADE
);

CREATE TABLE ingest_jobs (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL,
    document_id uuid NOT NULL,
    status      text NOT NULL DEFAULT 'pending'
                CHECK (status IN ('pending','processing','done','failed')),
    attempts    int  NOT NULL DEFAULT 0,
    error       text,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (document_id, tenant_id)
        REFERENCES documents(id, tenant_id) ON DELETE CASCADE
);


CREATE INDEX chunks_tenant_idx ON chunks (tenant_id);
CREATE INDEX chunks_embedding_hnsw ON chunks USING hnsw (embedding vector_cosine_ops);
CREATE INDEX chunks_tsv_gin ON chunks USING gin (tsv);
CREATE INDEX ingest_jobs_status_idx ON ingest_jobs (status);