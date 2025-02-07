# FORUM WEB

## Description
Ce projet est un forum web développé en **Go** avec **SQLite** comme base de données. Il permet aux utilisateurs de s'inscrire, se connecter, poster des publications, liker/disliker des contenus et filtrer les posts par catégories.

## Fonctionnalités
- Authentification des utilisateurs (inscription, connexion, sessions).
- Création de publications.
- Système de "like" et "dislike".
- Association de catégories aux publications.
- Filtrage des publications par catégories.

## Technologies Utilisées
- **Back-end :** Go
- **Base de données :** SQLite
- **Serveur Web :** `net/http`
- **Gestion des sessions :** Cookies
- **Conteneurisation :** Docker

---

## Structure du Projet
```
forum_v2/
│── databases/
│   ├── mine.sql         # Fichier de base de données SQLite
│── src/                 # Code source principal
│── static/              # Fichiers CSS  
│── templates/           # Fichiers HTML pour le rendu
│── Dockerfile           # Fichier de configuration Docker
│── go.mod               # Fichier des dépendances Go
│── main.go              # Point d'entrée du projet
│── README.md            # Documentation du projet
```

---

## Structure de la Base de Données

**Aperçu de la table de données :**

![Aperçu de la table](static/dia.png)

---

## Installation et Exécution sans Docker
### 1. Cloner le projet
```sh
git clone <URL_DU_REPO>
cd forum_v2
code forum_v2
```

### 2. Installer Go
Assurez-vous d'avoir **Go 1.22.2 ou plus** installé.

### 3. Installer les dépendances
```sh
go mod tidy
```

### 4. Exécuter le projet
```sh
go run main.go
```

Le serveur démarrera sur `http://localhost:9090`.

---

## Exécution avec Docker

### a) Construction de l'image Docker
```sh
docker build -t <nom_de_image> .
```

### b) Exécution du conteneur
```sh
docker run -p 9090:9090 <nom_de_image>
```

> 💡 Accédez au site via **http://localhost:9090**

---

## Auteurs
- **Youssef HAYYANI**
- **Hassane OUHAMOU**
- **Ibrahim FARES**
- **Mohammed AADOU**
- **Agiel OTCHOUN**

📌 **Licence :** Ce projet est sous licence MIT.