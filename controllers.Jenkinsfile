pipeline {
    agent { label 'kthw-agent' }
    environment {
        JENKINS_NODE_COOKIE = 'do_not_kill'
        ANSIBLE_PRIVATE_KEY = credentials('admin-agent')
    }
    stages {
        stage('Run k8s-control-bootstrap.yml') {
            steps {
                script {
                    // Run the second playbook
                    def result = sh(script: "ansible-playbook -i hosts.hosts --private-key=$ANSIBLE_PRIVATE_KEY k8s-control-bootstrap.yml", returnStatus: true)

                    // Check if the second playbook ran successfully
                    if (result != 0) {
                        currentBuild.result = 'FAILURE'
                        error "k8s-control-bootstrap.yml failed, skipping subsequent playbooks."
                    }
                }
            }
        }
    }

    post {
        success {
            echo "All playbooks executed successfully!"
        }
        failure {
            echo "One or more playbooks failed."
        }
    }
}
