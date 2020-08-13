pipeline {
    agent none

    stages {
        stage('Build potty-server') {
            agent {
                label 'potty-go'
            }
            steps {
                bat 'go build cmd/potty-server'
            }
        }
        stage('Run tests') {
            agent {
                label 'potty-go'
            }
            steps {
                bat 'go test cmd/potty-server'
            }
        }
        stage('Build docker image') {
            agent {
                label 'potty-go-oc-docker'
            }
            steps {
                bat 'docker build --tag gyorireka/potty -f Dockerfile .'
                bat 'docker push gyorireka/potty'
            }
        }
        stage('Deploy to OC') {
            agent {
                label 'potty-go-oc-docker'
            }
            steps {
                bat 'oc rollout latest deployment.apps/potty-server'
            }
        }
    }
}