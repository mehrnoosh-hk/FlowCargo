FROM postgres:15.4

# Install dependencies for pgaudit
RUN apt-get update && apt-get install -y \
    postgresql-15-pgaudit \
    && rm -rf /var/lib/apt/lists/*