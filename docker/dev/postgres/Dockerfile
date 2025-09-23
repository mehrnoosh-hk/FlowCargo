FROM postgres:15.4

# Install dependencies for pgaudit
RUN apt-get update && apt-get install -y \
    postgresql-15-pgaudit \
    && rm -rf /var/lib/apt/lists/*

# Enable pgaudit in postgresql.conf
RUN echo "shared_preload_libraries = 'pgaudit'" >> /usr/share/postgresql/postgresql.conf.sample
RUN echo "pgaudit.log = 'write,ddl'" >> /usr/share/postgresql/postgresql.conf.sample

# Create initialization script
RUN echo "CREATE EXTENSION IF NOT EXISTS pgaudit;" > /docker-entrypoint-initdb.d/init-pgaudit.sql