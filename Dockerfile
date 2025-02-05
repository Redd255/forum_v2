# Étape de construction
FROM golang:latest AS builder

# Définir le répertoire de travail
WORKDIR /app

# Copier tous les fichiers du projet en une seule commande
COPY . .

# Télécharger les dépendances
RUN go mod tidy

# Compiler l'application
RUN go build -o forum

# Étape d'exécution (image plus légère)
FROM debian:bookworm-slim

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers nécessaires depuis l'étape de build
COPY --from=builder /app/forum .
COPY --from=builder /app/templates ./templates/
COPY --from=builder /app/static ./static/
COPY --from=builder /app/databases/data.db ./databases/data.db
COPY --from=builder /app/databases/mine.sql ./databases/mine.sql

# Assurer les bonnes permissions pour la base SQLite (si nécessaire)
RUN chmod 777 /app/databases/data.db || true

# Exposer le bon port
EXPOSE 9090

# Lancer l'application
CMD ["./forum"]
