pipeline {
    agent none // We don't use the master node for heavy lifting

    stages {
        // --- STAGE 1: BACKEND ---
        stage('Build Backend (Go)') {
            agent {
                docker { image 'golang:1.25-alpine' } // Ephemeral Agent
            }
            steps {
                dir('backend') {
                    sh 'go version'
                    sh 'go mod download'
                    sh 'go build -o main ./cmd/main.go'
                }
            }
        }

        // --- STAGE 2: FRONTEND ---
        stage('Build Frontend (React)') {
            agent {
                docker { image 'node:22-alpine' } // Ephemeral Agent
            }
            steps {
                dir('frontend') {
                    sh 'npm install'
                    sh 'npm run build'
                }
            }
        }
    }
}