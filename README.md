# Simple link Interceptor

## Utilisation

### Build

Avec docker d'installé, il suffit de lancer la commande suivante: `docker build -t link-interceptor .`

Une version de l'image est disponible sur le docker hub: `docker pull ery4z/link-interceptor:latest`

### Run

Il est nécessaire de définir les variables d'environnement suivantes:

-   SQL_USER= # Utilisateur de la base de données
-   SQL_USER_PASSWORD= # Mot de passe de l'utilisateur de la base de données
-   SQL_URL= # Url de la base de données
-   SQL_PORT= # Port de la base de données
-   SQL_DATABASE= # Nom de la base de données
-   SQL_SCHEMA= # Schéma de la base de données
-   SQL_TABLE= # Table de la base de données
-   SELF_URL= # Url du service
-   URL_REDIRECT= # Url de redirection

Le service a été testé avec une base de données Azure SQL. Il est possible qu'une autre base de données fonctionne mais cela n'a pas été testé.

Le schéma et la table seront créés automatiquement si ils n'existent pas.

### Utilisation

Plusieurs routes sont disponibles:

-   POST /create : Ajoute un couple email lien dans la base de données et renvoie un lien de redirection. Le body doit contenir ```json
    {
    "email": "email",
    "link": "link"
    }

```
- GET /key/{key} : Redirige vers le lien défini par URL_REDIRECT dans les variables d'environnement. Si la clé n'existe pas, redirige quand même. Lors de l'utilisation de cette route la base de données est mise à jour pour stocker l'heure d'utilisation.



```
