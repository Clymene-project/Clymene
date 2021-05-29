pipeline {
    environment {
        registry = "registry.localhub:5000/clymene/agent"
        registryCredential = 'bourbonkk'
        BRANCH_NAME = "${GIT_BRANCH.split("/")[1]}"
        DOCKER_SCAN_SUGGEST = false
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
//                 withDockerRegistry([ credentialsId: "docker-hub-credentials", url: "" ]) {
                    bat 'docker push '+ registry+':'+BRANCH_NAME
//                 }
            }
        }
        stage('Clean docker image') {
            steps{
                bat "docker rmi "+ registry+':'+BRANCH_NAME
            }
        }
    }
}
