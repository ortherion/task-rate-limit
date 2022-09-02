module gitlab.com/g6834/team17/task-service

go 1.18

replace gitlab.com/g6834/team17/task-service/pkg/auth_service => ./pkg/auth_service

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi/v5 v5.0.7
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/jackc/pgx/v4 v4.16.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/rs/zerolog v1.27.0
	github.com/satori/go.uuid v1.2.0
	gitlab.com/g6834/team17/task-service/pkg/auth_service v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/grpc v1.48.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/getsentry/sentry-go v0.13.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.12.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lib/pq v1.10.2 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.34.0 // indirect
	go.opentelemetry.io/otel v1.9.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.9.0 // indirect
	go.opentelemetry.io/otel/metric v0.31.0 // indirect
	go.opentelemetry.io/otel/sdk v1.9.0 // indirect
	go.opentelemetry.io/otel/trace v1.9.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.0.0-20220706163947-c90051bbdb60 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220503193339-ba3ae3f07e29 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
