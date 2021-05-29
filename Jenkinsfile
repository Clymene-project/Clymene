pipeline {
    environment {
        registry = "clymene/agent"
        registryCredential = 'docker.io'
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
                bat 'docker push '+ registry+':'+BRANCH_NAME
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