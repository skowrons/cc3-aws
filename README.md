# Cloud Computing Teil 3

Es werden drei verschiedene Usecases betrachtet und als Prototyp umgesetzt.

## DevOps Benachrichtigung

Lambda Funktionen können gut für verschiedene Aufgaben im DevOps Bereich genutzt werden.
So könnte wenn eine CI/CD Pipeline fehlschlägt mittels Webhooks eine SMS an entsprechendes Personal versendet werden.

In diesem Beispiel wird bei jedem Commit ein Webhook getriggert, welcher dann eine SMS versendet mittels eines Lambdas.
 
 ┌──────────────────────────────────────┐
 │GitHub / GitLab                       │
 │                                      │
 │ ┌──────────┐         ┌───────────────┤
 │ │          │Git Event│               │
 │ │ Git Repo ├────────►│ GitHub Webhook├──┐
 │ │          │         │               │  │
 │ └──────────┘         └───────────────┤  │
 │                                      │  │ HTTP-Request
 └──────────────────────────────────────┘  │
                                           │
          ┌────────────────────────────────┘
          │
 ┌────────┼───────────────────────────┐
 │AWS     │                           │
 │        ▼                           │
 │ ┌─────────────┐ ┌────────┐ ┌─────┐ │
 │ │             │ │        │ │     │ │
 │ │ API Gateway │►│ Lambda │►│ SNS │ │
 │ │             │ │        │ │     │ │
 │ └─────────────┘ └────────┘ └──┬──┘ │
 │                               │    │
 └───────────────────────────────┼────┘
                                 │
                   ┌─────────────┘SMS Benachrichtigung
                   ▼
            ┌─────────────┐
            │             │
            │ Smartphone  │
            │             │
            └─────────────┘

## Datalake

Ein Datalake dient zum speichern von strukturierten und unstrukturierten Datensätzen.
So können Daten z.B. immer als eine Datei in einem Dateisystem oder S3 abgelegt werden.
Später könnten mehrere Lambdas bzw. Logik hinzukommen, die Daten je nach Typ in die entsprechenden Datenbanken schreiben.
RDS, DynamoDB, EFS, S3 und co. könnten dann zum speichern von Daten dienen.
Das gesamte Konstrukt könnte dann noch mit Cognito abgesichert werden, um den Zugriff zu beschränken.

                                            ┌───────────────┐
                                            │               │
 speichern/     ┌─────────────┐    ┌────────┤ eingehangenes │
 laden von Daten│             ├───►│        │               │
 ──────────────►│ API-Gateway │    │ Lambda │ Dateisystem   │
                │             │◄───┤        │               │
                └─────────────┘    └────────┤ (EFS)         │
                                            │               │
                                            └───────────────┘

## Fargate PaaS mit Copilot

                          Request
                             │
                             ▼
                          ┌─────┐
                          │     │
                          │ ELB │
                          │     │
                          └──┬──┘
                             │
 ┌───────────────────────────┼────────────────────────┐
 │ Fargat/ECS                │                        │
 │                           ▼                        │
 │                        ┌─────┐                     │
 │                        │     │                     │
 │                        │ API ├────────┐            │
 │                        │     │        │            │
 │                        └──┬──┘        ▼            │
 │                           │           ECS Services │
 │                           │           ▲            │
 │  ┌───────────┐   ┌────────┴────────┐  │            │
 │  │           │   │                 │  │            │
 │  │ Another   ◄───┘ some non public │  │            │
 │  │ Service   ┌───► backend service ├──┘            │
 │  │           │   │                 │               │
 │  └───────────┘   └─────────────────┘               │
 │                                                    │
 └────────────────────────────────────────────────────┘

