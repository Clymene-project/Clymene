pipeline {
    environment {
        registry = "bourbonkk/clymene/agent"
        registryCredential = 'docker-hub'
        BRANCH_NAME = "${GIT_BRANCH.split("/")[1]}"
    }
    agent any
    stages {
        stage('docker build') {
            steps {
                bat 'docker build -t '+ registry+':'+BRANCH_NAME+' -f=agent.Dockerfile .'
            }
        }
        stage('docker deploy') {
            steps {
                 docker.withRegistry('https://registry.hub.docker.com', 'docker-hub'){
                    bat 'docker push '+ registry+':'+BRANCH_NAME
                 }
            }
        }
        stage('Clean docker image') {
            steps{
                bat "docker rmi "+ registry+':'+BRANCH_NAME
            }
        }
    }
}
