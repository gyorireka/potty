pipeline {
    agent none

    stages {
        stage('Build') {
            agent {
                label 'potty-go'
            }
            steps {
                bat 'go version'
            }
        }
    }
}