# Fargate Beispiel

## Setup 

### Vorraussetzungen

- ein aws Account
- aws-cli installiert
- copilot-cli installiert
- docker lokal auf dem System vorhanden

### Docker installieren

Alle wichtigen Schritte sind in der [Docker Dokumentation beschrieben.](https://docs.docker.com/get-docker/)

## Verwendung

Copilot benutzt Fargate, ECR, ELB und Cloud Formation, um die als Dockerfile vorhandenen Microservices automatisch bereitzustellen.
Am Ende entsteht folgendes Konstrukt, insofern alle Befehle entsprechend ausgeführt wurden:

```
                Request
                   |
                   v
                +-----+
                |     |
                | ELB |
                |     |
                +--+--+
                   |
+------------------+------------------------+
| Fargat/ECS       |                        |
|                  v                        |
|               +-----+                     |
|               |     |                     |
|               | API +--------+            |
|               |     |        |            |
|               +--+--+        v            |
|                  |           ECS Services |
|                  v           ^            |
|         +-----------------+  |            |
|         |                 |  |            |
|         | some non public |  |            |
|         | backend service +--+            |
|         |                 |               |
|         +-----------------+               |
|                                           |
+-------------------------------------------+
```

Um Copilot und Fargate zu initialisieren muss der folgende Befehl ausgeführt werden:

```bash
copilot app init NAMEDERAPP
```

Dadurch wird ein `copilot` Verzeichnis angelegt, welches entsprechende Konfigurationen hält.
Nun müssen alle Microservices initalisiert werden.
Dabei unterscheiden wir öffentliche Microservices, wie der API-Service, und backend services, welche am ende die Logik halten.

Der folgende Befehl stellt den API Service bereit. 
Dafür muss der Pfad zur Dockerfile, ein Name und der Typ angegeben werden.
Mittels `--deploy` wird der Service in einem Testsystem bereitgestellt und nicht produktive.
Der Typ `Load Balanced Web Service` sorgt für die Erstellung eines öffentlich zugänglichen ELBs.

```bash
copilot init -n api -t "Load Balanced Web Service" -d ./api/Dockerfile --deploy
```

Nachdem der Service initialisiert wurde kann unter `./copilot/api/manifest` die genaue Konfiguration gefunden werden.
Es muss nun noch unter `http:` der Healthcheckendpoint hinzugefügt werden.
Ungefähr so:

```yaml
http:
  path: '/'
  healthcheck: '/_healthcheck'
```

Das entsprechende Container Image wird nun mittel dem lokalen Docker gebaut und in das ECR gepusht, von wo es dann innerhalb der AWS verwendet werden kann.
Verwendung findet es dann im ECS Cluster, wo es als Task mit einem dazugehörigen Service angelegt wird.

```bash
copilot app delete NAMEDERAPP
```
