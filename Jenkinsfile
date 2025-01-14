pipeline {
    agent { label 'kthw-agent' }
    environment {
        JENKINS_NODE_COOKIE = 'do_not_kill'
        ANSIBLE_PRIVATE_KEY = credentials('admin-agent')
    }
    stages {
        stage('Run client_tools.yml') {
            steps {
                script {
                    // Run the first playbook
                    def result = sh(script: "ansible-playbook -i hosts.hosts --private-key=$ANSIBLE_PRIVATE_KEY client_tools.yml", returnStatus: true)

                    // Check if the first playbook ran successfully
                    if (result != 0) {
                        currentBuild.result = 'FAILURE'
                        error "client_tools.yml failed, skipping the next playbook."
                    }
                }
            }
        }

        stage('Run tls_cert_gen.yml') {
            when {
                expression {
                    // Run only if the previous stage was successful
                    return currentBuild.result == 'SUCCESS'
                }
            }
            steps {
                script {
                    // Run the second playbook
                    def result = sh(script: "ansible-playbook -i hosts.hosts --private-key=$ANSIBLE_PRIVATE_KEY tls_cert_gen.yml", returnStatus: true)

                    // Check if the second playbook ran successfully
                    if (result != 0) {
                        currentBuild.result = 'FAILURE'
                        error "tls_cert_gen.yml failed, skipping subsequent playbooks."
                    }
                }
            }
        }

        stage('Run etcd-bootstrap.yml') {
            when {
                expression {
                    // Run only if the previous stage was successful
                    return currentBuild.result == 'SUCCESS'
                }
            }
            steps {
                script {
                    // Run another playbook
                    def result = sh(script: "ansible-playbook -i hosts.hosts --private-key=$ANSIBLE_PRIVATE_KEY etcd-bootstrap.yml", returnStatus: true)

                    // Check if the third playbook ran successfully
                    if (result != 0) {
                        currentBuild.result = 'FAILURE'
                        error "etcd-bootstrap.yml failed."
                    }
                }
            }
        }


        stage('Run k8s-control-bootstrap.yml') {
            when {
                expression {
                    // Run only if the previous stage was successful
                    return currentBuild.result == 'SUCCESS'
                }
            }
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

        stage('Run k8s-worker-bootstrap.yml') {
            when {
                expression {
                    // Run only if the previous stage was successful
                    return currentBuild.result == 'SUCCESS'
                }
            }
            steps {
                script {
                    // Run another playbook
                    def result = sh(script: "ansible-playbook -i hosts.hosts --private-key=$ANSIBLE_PRIVATE_KEY k8s-worker-bootstrap.yml", returnStatus: true)

                    // Check if the third playbook ran successfully
                    if (result != 0) {
                        currentBuild.result = 'FAILURE'
                        error "k8s-worker-bootstrap.yml failed."
                    }
                }
            }
        }

        stage('Run kubectl-accessyml') {
            when {
                expression {
                    // Run only if the previous stage was successful
                    return currentBuild.result == 'SUCCESS'
                }
            }
            steps {
                script {
                    // Run another playbook
                    def result = sh(script: "ansible-playbook -i hosts.hosts --private-key=$ANSIBLE_PRIVATE_KEY kubectl-access.yml", returnStatus: true)

                    // Check if the third playbook ran successfully
                    if (result != 0) {
                        currentBuild.result = 'FAILURE'
                        error "kubectl-access.yml failed."
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
