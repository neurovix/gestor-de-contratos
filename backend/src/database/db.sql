-- Tabla de plantas
CREATE TABLE plantas (
    id_planta SERIAL PRIMARY KEY,
    nombre TEXT
);

-- Tabla de usuarios
CREATE TABLE usuarios (
    id_usuario SERIAL PRIMARY KEY,
    nombre TEXT,
    id_planta INT NOT NULL REFERENCES plantas(id_planta) ON DELETE CASCADE,
    email TEXT UNIQUE,
    password TEXT
);

-- Tabla de trámites
CREATE TABLE tramites (
    id_tramite SERIAL PRIMARY KEY,
    no_contrato TEXT,
    creado_por INT NOT NULL REFERENCES usuarios(id_usuario) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    nombre_tramite TEXT
);

-- Tabla de estatus (aprobado, en espera, rechazado, etc.)
CREATE TABLE estatus (
    id_estatus SERIAL PRIMARY KEY,
    nombre_estatus TEXT NOT NULL
);

-- Tabla pivote: relación N:N entre trámites y verificadores, con estatus
CREATE TABLE verificaciones (
    id_verificacion SERIAL PRIMARY KEY,
    id_tramite INT NOT NULL REFERENCES tramites(id_tramite) ON DELETE CASCADE,
    id_usuario INT NOT NULL REFERENCES usuarios(id_usuario) ON DELETE CASCADE,
    id_estatus INT NOT NULL REFERENCES estatus(id_estatus),
    verificado_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (id_tramite, id_usuario) -- evita registros duplicados
);

-- Tabla de archivos (para manejar S3)
CREATE TABLE archivos (
    id_archivo SERIAL PRIMARY KEY,
    id_tramite INT NOT NULL REFERENCES tramites(id_tramite) ON DELETE CASCADE,
    s3_key TEXT NOT NULL,               -- clave interna del bucket
    s3_url TEXT,                        -- URL accesible al documento
    mime_type TEXT,                     -- tipo de archivo
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
