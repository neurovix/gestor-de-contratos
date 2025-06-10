CREATE TABLE plantas (
    id_planta SERIAL PRIMARY KEY,
    nombre TEXT
);

INSERT INTO plantas(nombre)VALUES('Planta Bajio'),('Planta Prosede'),('Planta Ramos Arizpe'),('Planta Tlaxcala'),('Planta Ecatepec'),('Planta Morelia'),('Planta San Rafael'),('Planta Texmelucan'),('Planta Toluca');

CREATE TABLE usuarios (
    id_usuario SERIAL PRIMARY KEY,
    nombre TEXT,
    id_planta INT NOT NULL REFERENCES plantas(id_planta) ON DELETE CASCADE,
    email TEXT UNIQUE,
    password TEXT,
    cargo TEXT
);

CREATE TABLE tramites (
    id_tramite SERIAL PRIMARY KEY,
    no_contrato TEXT,
    nombre_tramite TEXT,
    archivo_pdf_url TEXT,
    creador_id INT NOT NULL REFERENCES usuarios(id_usuario) ON DELETE CASCADE,
    fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE verificaciones (
    id_verificacion SERIAL PRIMARY KEY,
    id_tramite INT NOT NULL REFERENCES tramites(id_tramite) ON DELETE CASCADE,
    id_verificador INT NOT NULL REFERENCES usuarios(id_usuario) ON DELETE CASCADE,
    orden INT NOT NULL,
    verificado BOOLEAN DEFAULT FALSE,
    fecha_verificacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (id_tramite, id_verificador)
);
