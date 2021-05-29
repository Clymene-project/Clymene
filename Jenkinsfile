pipeline {
    environment {
        registry = "clymene/agent"
        registryCredential = 'dockerhub'
        BRANCH_NAME = "${GIT_BRANCH.split("/")[1]}"
    }
    agent any
    stages {
        stage('docker build') {
            steps {
                bat 'docker build -t '+ registry+':'+BRANCH_NAME+' cmd/agent/.'
            }
        }
        stage('docker deploy') {
            steps {
                withDockerRegistry([ credentialsId: registryCredential, url: "" ]) {
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