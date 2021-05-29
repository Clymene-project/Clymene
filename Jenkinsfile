pipeline {
    environment {
        registry = "registry.localhub:5000/clymene/agent"
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
                withDockerRegistry([ credentialsId: registryCredential, url: "https://registry.hub.docker.com" ]) {
                    bat 'docker push '+ registry+':'+BRANCH_NAME
                }
            }
        }
        stage('Clean docker image') {
            steps{
                bat "docker rmi $registry"
            }
        }
    }
}

// def checkOs(){
//     if (isUnix()) {
//         def uname = sh script: 'uname', returnStdout: true
//         if (uname.startsWith("Darwin")) {
//             return "Macos"
//         }
//         // Optionally add 'else if' for other Unix OS
//         else {
//             return "Linux"
//         }
//     }
//     else {
//         return "Windows"
//     }
// }