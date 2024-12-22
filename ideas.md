# Server

## global

1) au demarrage de serveur
- Appel API en mode batch (configurable: premier appel ON/OFF)
- Génération de fichier JSON avec les données de API
- Récup des blagues perso
- fusion blagues perso/API 
- go run main.go --no-api (pas de requetes API)

# Client

- /getJoke                -> une blague au pif
- /getJoke/{id}           -> blague à l'id
- /getJokes               -> toutes les blagues

- /initJokes              -> fait appel API et reset fichier JSON de blagues + régénére map de blagues fusionnées