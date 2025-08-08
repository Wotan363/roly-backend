TODO: Adjust and polish the README file

# roly-backend
The Backend for roly. 

You have to create a file called .env.db 
You have to create a file called .env.secrets and fill out evything like .env.secretsTemplate
There you should put the database credentials like in .env.dbTemplate

0. install Docker
-> Start Docker

1. Start Database with docker
irgendwas mit dockercompose

2. Go Starten
alle commands für windows und linux

3. Go bauen
alle commands für windows und linux

4. Go testen
alle commands für windows und linux

Um das Go Backend zu starten in dev umgebung:
$env:APP_ENV = "development"
go run ./cmd/roly-backend

für docker muss man das eingeben:
docker compose -f dockerDatabase/docker-compose.yml up -d
