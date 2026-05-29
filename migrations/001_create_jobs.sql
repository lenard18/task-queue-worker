-- ============================================================
--  Cola de Tareas — Schema principal
-- ============================================================

CREATE TABLE IF NOT EXISTS jobs (
    id            BIGSERIAL    PRIMARY KEY,
    type          VARCHAR(50)  NOT NULL,                  -- send_email | send_report | send_welcome
    payload       JSONB        NOT NULL DEFAULT '{}',     -- datos del job en JSON
    status        VARCHAR(20)  NOT NULL DEFAULT 'pending',-- pending | processing | done | failed
    priority      INT          NOT NULL DEFAULT 0,        -- mayor número = mayor prioridad
    attempts      INT          NOT NULL DEFAULT 0,        -- intentos realizados
    max_attempts  INT          NOT NULL DEFAULT 3,        -- máximo de reintentos
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    scheduled_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),    -- no procesar antes de esta fecha
    processed_at  TIMESTAMPTZ,                            -- cuándo se completó
    failed_at     TIMESTAMPTZ,                            -- cuándo falló definitivamente
    error_msg     TEXT                                    -- último mensaje de error
);

-- Índice para que el worker encuentre jobs rápido
CREATE INDEX IF NOT EXISTS idx_jobs_status_scheduled
    ON jobs (status, scheduled_at, priority DESC)
    WHERE status IN ('pending', 'failed');

-- ============================================================
--  Vista para el dashboard
-- ============================================================
CREATE OR REPLACE VIEW jobs_stats AS
SELECT
    status,
    COUNT(*)                                      AS total,
    MIN(created_at)                               AS oldest,
    MAX(created_at)                               AS newest,
    AVG(EXTRACT(EPOCH FROM (processed_at - created_at)))::INT AS avg_seconds
FROM jobs
GROUP BY status;
