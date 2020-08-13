pipeline {
    agent none

    stages {
        stage('Build') {
            agent {
                label 'potty-go'
            }
            steps {
                sh 'ls'
            }
        }
    }
}