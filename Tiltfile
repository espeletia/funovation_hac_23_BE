load('ext://restart_process', 'docker_build_with_restart')
load_dynamic('./ci/tilt/postgres.Tiltfile')
load_dynamic('./ci/tilt/minio.Tiltfile')
k8s_yaml("./ci/funovation.yaml")

s3Url = "http://minio:9000"

local_resource(
    'regenerate-funovation',
    'go generate cmd/main.go',
    deps=[
    './graph/'
    ],
    ignore=[
    './graph/*.go',
    './graph/generated',
    './graph/model',
    ],
    resource_deps=['postgresql'],
    labels=["compile"],
)

local_resource(
      'compile funovation',
      'bash ./ci/build.sh',
      deps=[
      './',
      ],
      ignore=[
      'tilt_modules',
      'Tiltfile',
      'graph/schema.graphqls',
      'build',
      'dep',
      'ci/docker-compose.yaml',
      'swagger.yaml',
      'internal/handlers/swagger.yaml',
      'internal/handlers/generated.go',
      '**/testdata'
      ],
      labels=["compile"],
  )
  
docker_build_with_restart('funovation',
    '.',
    dockerfile='ci/Dockerfile',
    entrypoint='/app/start_server',
    only=[
        './build',
        './configurations',
        './certs',
        './migrations',
        './cmd/migrations'
    ],
    live_update=[
        sync('./configurations', '/app/configurations'),
        sync('./build', '/app')
    ]
)

docker_build_with_restart('funovation-migrations',
    '.',
    dockerfile='ci/Dockerfile',
    entrypoint='/app/run_migrations',
    only=[
        './build',                    # Ensure the binary "run_migrations" is here
        './configurations',
        './certs',
        './migrations'
    ],
    live_update=[
        sync('./build', '/app'),       # Sync the "build" directory with /app inside the container
        sync('./configurations', '/app/configurations')
    ],
    build_args={"app": "funovation","s3Url": s3Url}
)
     

k8s_resource("funovation", port_forwards=["0.0.0.0:8080:8080"], resource_deps=['postgresql', 'minio'], labels=["BE"])