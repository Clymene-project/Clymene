pipeline {
    environment {
        registry = "bourbonkk/clymene/agent"
        registryCredential = 'bourbonkk'
        BRANCH_NAME = "${GIT_BRANCH.split("/")[1]}"
    }
    agent any
    stages {
        stage('docker build') {
            steps {
                script {

                    dockerImage = docker.build registry + ":$BBRANCH_NAME"

                }
                bat 'docker build -t '+ registry+':'+BRANCH_NAME+' -f=agent.Dockerfile .'
            }
        }
        stage('docker deploy') {
            steps {
                bat 'docker push '+ registry+':'+BRANCH_NAME

            }
        }
        stage('Clean docker image') {
            steps{
                bat "docker rmi "+ registry+':'+BRANCH_NAME
            }
        }
    }
}
