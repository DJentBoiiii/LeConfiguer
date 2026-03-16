pipeline {
  // Node is allocated by the first stage; all Docker agents reuse it via reuseNode true
  agent any

  options {
    timestamps()
  }

  environment {
    REGISTRY              = 'djentboiii'
    DOCKER_CREDENTIALS_ID = 'docker-registry-credentials'
    AWS_CREDENTIALS_ID    = 'aws-jenkins-credentials'
    AWS_REGION            = 'eu-central-1'
    EKS_CLUSTER_NAME      = 'leconfiguer-eks'
    K8S_NAMESPACE         = 'default'
    IMAGE_TAG             = "${env.BUILD_NUMBER}"
  }

  stages {
    stage('Prepare') {
      agent {
        docker {
          image 'golang:1.25'
          reuseNode true
        }
      }
      steps {
        checkout scm
        sh '''
          set -eux
          go version
          for svc in api-gateway config-storage indexing; do
            (cd "$svc" && go mod download)
          done
        '''
      }
    }

    stage('Test') {
      agent {
        docker {
          // golang:1.25 (debian-based) includes gcc for CGO required by go-sqlite3
          image 'golang:1.25'
          reuseNode true
        }
      }
      steps {
        sh '''
          set -eux
          for svc in api-gateway config-storage indexing; do
            (cd "$svc" && go test ./...)
          done
        '''
      }
    }

    stage('Static code analysis') {
      agent {
        docker {
          image 'golang:1.25'
          reuseNode true
        }
      }
      steps {
        sh '''
          set -eux
          for svc in api-gateway config-storage indexing; do
            (
              cd "$svc"
              UNFORMATTED="$(find . -name '*.go' -not -path './vendor/*' -print0 | xargs -0 gofmt -l)"
              if [ -n "$UNFORMATTED" ]; then
                echo "Unformatted files in $svc:"
                echo "$UNFORMATTED"
                exit 1
              fi
              go vet ./...
            )
          done
        '''
      }
    }

    stage('Build Docker images') {
      agent {
        docker {
          image 'docker:27-cli'
          // Mount the host Docker socket so we can build images
          args  '-v /var/run/docker.sock:/var/run/docker.sock'
          reuseNode true
        }
      }
      steps {
        sh '''
          set -eux
          for svc in api-gateway config-storage indexing; do
            docker build \
              -t "$REGISTRY/$svc:$IMAGE_TAG" \
              -t "$REGISTRY/$svc:latest" \
              "$svc"
          done
        '''
      }
    }

    stage('Push images to registry') {
      agent {
        docker {
          image 'docker:27-cli'
          args  '-v /var/run/docker.sock:/var/run/docker.sock'
          reuseNode true
        }
      }
      steps {
        withCredentials([usernamePassword(
          credentialsId: env.DOCKER_CREDENTIALS_ID,
          usernameVariable: 'REG_USER',
          passwordVariable: 'REG_PASS'
        )]) {
          sh '''
            set -eux
            echo "$REG_PASS" | docker login -u "$REG_USER" --password-stdin

            for svc in api-gateway config-storage indexing; do
              docker push "$REGISTRY/$svc:$IMAGE_TAG"
              docker push "$REGISTRY/$svc:latest"
            done

            docker logout
          '''
        }
      }
    }

    stage('Deploy to EKS') {
      when {
        anyOf {
          branch 'main'
          branch 'master'
        }
      }
      agent {
        docker {
          // alpine/k8s bundles kubectl + aws-cli + helm in a single image
          image 'alpine/k8s:1.29.2'
          reuseNode true
        }
      }
      steps {
        withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', credentialsId: env.AWS_CREDENTIALS_ID]]) {
          sh '''
            set -eux

            aws eks update-kubeconfig \
              --region "$AWS_REGION" \
              --name "$EKS_CLUSTER_NAME"

            kubectl get ns "$K8S_NAMESPACE" >/dev/null 2>&1 || kubectl create ns "$K8S_NAMESPACE"

            kubectl -n "$K8S_NAMESPACE" apply -f kubernetes/secrets.yaml
            kubectl -n "$K8S_NAMESPACE" apply -f kubernetes/services.yaml
            kubectl -n "$K8S_NAMESPACE" apply -f kubernetes/statefulsets.yaml
            kubectl -n "$K8S_NAMESPACE" apply -f kubernetes/deployment.yaml

            kubectl -n "$K8S_NAMESPACE" set image deployment/api-gateway     api-gateway="$REGISTRY/api-gateway:$IMAGE_TAG"
            kubectl -n "$K8S_NAMESPACE" set image deployment/config-storage   config-storage="$REGISTRY/config-storage:$IMAGE_TAG"
            kubectl -n "$K8S_NAMESPACE" set image deployment/indexing         indexing="$REGISTRY/indexing:$IMAGE_TAG"

            kubectl -n "$K8S_NAMESPACE" rollout status deployment/api-gateway     --timeout=180s
            kubectl -n "$K8S_NAMESPACE" rollout status deployment/config-storage  --timeout=180s
            kubectl -n "$K8S_NAMESPACE" rollout status deployment/indexing        --timeout=180s
          '''
        }
      }
    }
  }

  post {
    always {
      cleanWs()
    }
  }
}
