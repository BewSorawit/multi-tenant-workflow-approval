pipeline {
    agent any

    options {
        timestamps()
        ansiColor('xterm')
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Tool Versions') {
            steps {
                sh '''
                    git --version
                    go version
                    make --version
                    docker version
                    docker compose version
                    curl --version
                '''
            }
        }

        stage('Format Check') {
            steps {
                sh 'make format-check'
            }
        }

        stage('Vet') {
            steps {
                sh 'make vet'
            }
        }

        stage('Test') {
            steps {
                sh 'make test'
            }
        }

        stage('Build') {
            steps {
                sh 'make build'
            }
        }

        stage('Docker Build') {
            steps {
                sh 'make docker-build'
            }
        }

        stage('Smoke Test') {
            steps {
                sh 'make smoke-test'
            }
        }
    }

    post {
        always {
            sh 'docker compose down || true'
        }

        success {
            echo 'CI pipeline completed successfully'
        }

        failure {
            echo 'CI pipeline failed'
        }
    }
}
