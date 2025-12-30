# Jenkins Pipeline from Git: An ELI5 Guide

## What's Actually Happening Here?

Think of Jenkins as a **robot chef** in a kitchen. The **Jenkinsfile** is the recipe, and **Git** is the cookbook on the shelf. When you want to cook (build), the chef grabs the recipe from the cookbook and follows the steps.

---

## The Big Picture

```
┌─────────────┐      pulls code       ┌─────────────┐
│   Git Repo  │ ──────────────────▶   │   Jenkins   │
│ (has code + │                       │  (runs the  │
│ Jenkinsfile)│                       │    build)   │
└─────────────┘                       └─────────────┘
```

---

## Key Concepts

### Jenkinsfile

A text file (named `Jenkinsfile`, no extension) that lives **in your Git repo** alongside your code. It tells Jenkins: "Here's how to build, test, and deploy this project."

**ELI5:** It's like IKEA instructions that come *inside* the furniture box—the build instructions travel with the thing being built.

### Pipeline

The entire automated workflow: checkout → build → test → deploy. Each step happens in order.

**ELI5:** An assembly line where your code moves from station to station until it becomes a working application.

### SCM (Source Control Management)

How Jenkins talks to Git. When you configure a "Pipeline from SCM," you're telling Jenkins: "Go look in this Git repo for instructions."

---

## Two Ways to Set This Up

### Option 1: Pipeline Script from SCM (Recommended)

Jenkins pulls the Jenkinsfile directly from your repo. The recipe lives with the ingredients.

**Pros:**
- Jenkinsfile is version-controlled with your code
- Changes to the pipeline go through code review
- Different branches can have different build steps

### Option 2: Pipeline Script (Inline)

You paste the Jenkinsfile contents directly into Jenkins UI.

**Pros:**
- Quick for testing
- No Git setup needed

**Cons:**
- Not version-controlled
- Easy to lose changes
- Doesn't scale

---

## Minimal Jenkinsfile Example

```groovy
pipeline {
    agent any                          // Run on any available Jenkins node
    
    stages {
        stage('Checkout') {
            steps {
                // Jenkins does this automatically when using "Pipeline from SCM"
                // but you can be explicit:
                checkout scm
            }
        }
        
        stage('Build') {
            steps {
                sh 'echo "Building the application..."'
                sh 'make build'        // or: npm install, pip install, go build, etc.
            }
        }
        
        stage('Test') {
            steps {
                sh 'make test'
            }
        }
        
        stage('Deploy') {
            steps {
                sh 'echo "Deploying..."'
            }
        }
    }
}
```

**ELI5 Breakdown:**
- `pipeline {}` — "Here's my recipe"
- `agent any` — "Any chef (worker node) can make this"
- `stages {}` — "Here are the steps, in order"
- `stage('Build')` — "This step is called Build"
- `steps {}` — "Do these specific things"
- `sh 'command'` — "Run this shell command"

---

## Setting Up "Pipeline from SCM" in Jenkins

### Step 1: Create a New Pipeline Job

1. Jenkins Dashboard → **New Item**
2. Enter a name → Select **Pipeline** → OK

### Step 2: Configure Git Source

Scroll to **Pipeline** section:

| Field | Value |
|-------|-------|
| Definition | Pipeline script from SCM |
| SCM | Git |
| Repository URL | `https://github.com/yourname/yourrepo.git` or `git@github.com:yourname/yourrepo.git` |
| Credentials | Select or add Git credentials |
| Branch | `*/main` or `*/master` or `*/${BRANCH_NAME}` for multibranch |
| Script Path | `Jenkinsfile` (default, or `ci/Jenkinsfile` if nested) |

### Step 3: Save and Build

Click **Save** → **Build Now**

---

## What Happens When You Click "Build Now"

```
1. Jenkins receives trigger (manual click, webhook, schedule)
         │
         ▼
2. Jenkins connects to Git repo using stored credentials
         │
         ▼
3. Jenkins clones/fetches the repo to a workspace directory
         │
         ▼
4. Jenkins reads the Jenkinsfile from the repo
         │
         ▼
5. Jenkins executes each stage in order
         │
         ▼
6. Build succeeds ✓ or fails ✗ (check console output)
```

**ELI5:** Jenkins goes to the library (Git), checks out your book (repo), reads chapter 1 (Jenkinsfile), then follows the instructions step by step.

---

## Real-World Jenkinsfile: Python App

```groovy
pipeline {
    agent any
    
    environment {
        VENV = 'venv'
    }
    
    stages {
        stage('Setup') {
            steps {
                sh '''
                    python3 -m venv ${VENV}
                    . ${VENV}/bin/activate
                    pip install --upgrade pip
                    pip install -r requirements.txt
                '''
            }
        }
        
        stage('Lint') {
            steps {
                sh '''
                    . ${VENV}/bin/activate
                    flake8 src/ --max-line-length=120
                '''
            }
        }
        
        stage('Test') {
            steps {
                sh '''
                    . ${VENV}/bin/activate
                    pytest tests/ -v --junitxml=test-results.xml
                '''
            }
            post {
                always {
                    junit 'test-results.xml'
                }
            }
        }
        
        stage('Build') {
            steps {
                sh '''
                    . ${VENV}/bin/activate
                    python setup.py sdist bdist_wheel
                '''
            }
        }
    }
    
    post {
        success {
            echo 'Build succeeded!'
        }
        failure {
            echo 'Build failed!'
        }
        cleanup {
            cleanWs()    // Clean workspace after build
        }
    }
}
```

---

## Real-World Jenkinsfile: Docker Build

```groovy
pipeline {
    agent any
    
    environment {
        REGISTRY = 'registry.example.com'
        IMAGE_NAME = 'myapp'
        IMAGE_TAG = "${env.BUILD_NUMBER}"
    }
    
    stages {
        stage('Build Image') {
            steps {
                sh "docker build -t ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} ."
            }
        }
        
        stage('Test Image') {
            steps {
                sh "docker run --rm ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} /app/run-tests.sh"
            }
        }
        
        stage('Push Image') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'docker-registry-creds',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh '''
                        echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin ${REGISTRY}
                        docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
                    '''
                }
            }
        }
    }
}
```

---

## Common Gotchas

### "Jenkinsfile not found"

- Check `Script Path` in job config (is it `Jenkinsfile` or `ci/Jenkinsfile`?)
- Is the file actually committed and pushed?
- Is the branch correct?

### "Permission denied" on Git clone

- Credentials not configured or expired
- SSH key not added to Jenkins credentials store
- For HTTPS, use a personal access token, not your password

### "sh: command not found"

The Jenkins agent doesn't have that tool installed. Either:
- Install it on the agent
- Use a Docker agent with the tools pre-installed
- Use a tool installer in the pipeline

### Workspace issues

Each build gets a workspace like `/var/jenkins_home/workspace/job-name/`. If builds interfere with each other, add:

```groovy
options {
    disableConcurrentBuilds()
}
```

---

## Triggering Builds Automatically

### Webhook (Push to Git = Build)

In your job config, under **Build Triggers**:
- Check **GitHub hook trigger for GITScm polling**

In GitHub/GitLab, add a webhook pointing to:
```
https://your-jenkins-server/github-webhook/
```

### Poll SCM (Check Git every X minutes)

```groovy
triggers {
    pollSCM('H/5 * * * *')    // Every 5 minutes
}
```

**ELI5:** Webhook = Git calls Jenkins and says "I changed!" Poll = Jenkins keeps asking Git "Did you change yet? How about now?"

### Scheduled

```groovy
triggers {
    cron('0 2 * * *')    // Every day at 2 AM
}
```

---

## Quick Reference: Declarative vs Scripted

| Declarative (Modern) | Scripted (Legacy) |
|---------------------|-------------------|
| `pipeline { }` wrapper | `node { }` wrapper |
| Structured, opinionated | Flexible, Groovy-heavy |
| Easier to read | More powerful |
| Use this for new pipelines | You'll see this in old projects |

---

## Checklist: From Zero to Building

- [ ] Jenkinsfile exists in repo root (or known path)
- [ ] Jenkinsfile syntax is valid (use Jenkins → Pipeline Syntax helper)
- [ ] Git credentials stored in Jenkins
- [ ] Pipeline job created with "Pipeline script from SCM"
- [ ] Repository URL and branch are correct
- [ ] Script Path matches actual file location
- [ ] Build tools exist on Jenkins agent (or use Docker)
- [ ] First manual build succeeds
- [ ] Webhook configured for auto-builds

---

---

## Copy-Paste Setup Instructions

### 1. Create Your Jenkinsfile

Create a file named `Jenkinsfile` (no extension) in your repo root:

```bash
# Navigate to your project
cd /path/to/your/project

# Create the Jenkinsfile
cat > Jenkinsfile << 'EOF'
pipeline {
    agent any
    
    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timestamps()
        disableConcurrentBuilds()
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
                sh 'echo "Checked out branch: ${GIT_BRANCH}"'
            }
        }
        
        stage('Build') {
            steps {
                sh 'echo "Building..."'
                // Replace with your build command:
                // sh 'make build'
                // sh 'npm install && npm run build'
                // sh 'go build -o app .'
                // sh 'pip install -r requirements.txt'
            }
        }
        
        stage('Test') {
            steps {
                sh 'echo "Testing..."'
                // Replace with your test command:
                // sh 'make test'
                // sh 'npm test'
                // sh 'go test ./...'
                // sh 'pytest'
            }
        }
        
        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                sh 'echo "Deploying..."'
                // Add deployment steps here
            }
        }
    }
    
    post {
        success {
            echo '✅ Build succeeded!'
        }
        failure {
            echo '❌ Build failed!'
        }
        always {
            cleanWs()
        }
    }
}
EOF

# Commit and push
git add Jenkinsfile
git commit -m "Add Jenkinsfile for CI/CD pipeline"
git push origin main
```

---

### 2. Add Git Credentials to Jenkins (CLI)

```bash
# Set your Jenkins URL
JENKINS_URL="http://localhost:8080"

# --- Option A: SSH Key Credential ---

# First, generate an SSH key if you don't have one
ssh-keygen -t ed25519 -C "jenkins@yourcompany.com" -f ~/.ssh/jenkins_git_key -N ""

# Add the public key to GitHub/GitLab (copy this output)
cat ~/.ssh/jenkins_git_key.pub

# Create credentials XML file
cat > /tmp/ssh-credential.xml << 'EOF'
<com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey plugin="ssh-credentials">
  <scope>GLOBAL</scope>
  <id>git-ssh-key</id>
  <description>Git SSH Key for Pipeline</description>
  <username>git</username>
  <privateKeySource class="com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey$DirectEntryPrivateKeySource">
    <privateKey>PASTE_YOUR_PRIVATE_KEY_HERE</privateKey>
  </privateKeySource>
</com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey>
EOF

# Insert your private key into the XML
PRIVATE_KEY=$(cat ~/.ssh/jenkins_git_key)
# Then manually edit /tmp/ssh-credential.xml and paste the key

# Upload to Jenkins
curl -X POST "${JENKINS_URL}/credentials/store/system/domain/_/createCredentials" \
  --user admin:YOUR_API_TOKEN \
  --data-urlencode "json=$(cat /tmp/ssh-credential.xml)"


# --- Option B: Username + Token (HTTPS) ---

cat > /tmp/token-credential.xml << 'EOF'
<com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl>
  <scope>GLOBAL</scope>
  <id>git-https-token</id>
  <description>Git HTTPS Token</description>
  <username>YOUR_GIT_USERNAME</username>
  <password>YOUR_PERSONAL_ACCESS_TOKEN</password>
</com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl>
EOF

curl -X POST "${JENKINS_URL}/credentials/store/system/domain/_/createCredentials" \
  --user admin:YOUR_API_TOKEN \
  -H "Content-Type: application/xml" \
  -d @/tmp/token-credential.xml
```

---

### 3. Create Pipeline Job (CLI with jenkins-cli.jar)

```bash
# Download Jenkins CLI
JENKINS_URL="http://localhost:8080"
curl -o jenkins-cli.jar "${JENKINS_URL}/jnlpJars/jenkins-cli.jar"

# Create job config XML
cat > /tmp/pipeline-job.xml << 'EOF'
<?xml version='1.1' encoding='UTF-8'?>
<flow-definition plugin="workflow-job">
  <description>Pipeline job pulling Jenkinsfile from Git</description>
  <keepDependencies>false</keepDependencies>
  <properties>
    <org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
      <triggers>
        <com.cloudbees.jenkins.GitHubPushTrigger plugin="github">
          <spec></spec>
        </com.cloudbees.jenkins.GitHubPushTrigger>
      </triggers>
    </org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
  </properties>
  <definition class="org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition" plugin="workflow-cps">
    <scm class="hudson.plugins.git.GitSCM" plugin="git">
      <configVersion>2</configVersion>
      <userRemoteConfigs>
        <hudson.plugins.git.UserRemoteConfig>
          <url>git@github.com:YOUR_ORG/YOUR_REPO.git</url>
          <credentialsId>git-ssh-key</credentialsId>
        </hudson.plugins.git.UserRemoteConfig>
      </userRemoteConfigs>
      <branches>
        <hudson.plugins.git.BranchSpec>
          <name>*/main</name>
        </hudson.plugins.git.BranchSpec>
      </branches>
      <doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>
      <submoduleCfg class="empty-list"/>
      <extensions/>
    </scm>
    <scriptPath>Jenkinsfile</scriptPath>
    <lightweight>true</lightweight>
  </definition>
  <triggers/>
  <disabled>false</disabled>
</flow-definition>
EOF

# Create the job
java -jar jenkins-cli.jar -s "${JENKINS_URL}" -auth admin:YOUR_API_TOKEN \
  create-job my-pipeline-job < /tmp/pipeline-job.xml

# Trigger a build
java -jar jenkins-cli.jar -s "${JENKINS_URL}" -auth admin:YOUR_API_TOKEN \
  build my-pipeline-job
```

---

### 4. Quick Variable Replacement Script

Save this and run it to customize the templates above:

```bash
#!/bin/bash
# save as: setup-jenkins-pipeline.sh

read -p "Jenkins URL [http://localhost:8080]: " JENKINS_URL
JENKINS_URL=${JENKINS_URL:-http://localhost:8080}

read -p "Git repo URL (SSH format): " GIT_REPO_URL
read -p "Git branch [main]: " GIT_BRANCH
GIT_BRANCH=${GIT_BRANCH:-main}

read -p "Jenkins job name: " JOB_NAME
read -p "Credential ID [git-ssh-key]: " CRED_ID
CRED_ID=${CRED_ID:-git-ssh-key}

read -p "Jenkinsfile path [Jenkinsfile]: " JENKINSFILE_PATH
JENKINSFILE_PATH=${JENKINSFILE_PATH:-Jenkinsfile}

cat << EOF

============================================
Your Jenkins Pipeline Configuration
============================================

Jenkins URL:     ${JENKINS_URL}
Git Repo:        ${GIT_REPO_URL}
Branch:          ${GIT_BRANCH}
Job Name:        ${JOB_NAME}
Credential ID:   ${CRED_ID}
Jenkinsfile:     ${JENKINSFILE_PATH}

Next steps:
1. Add credential '${CRED_ID}' to Jenkins
2. Create pipeline job with these settings
3. Add webhook: ${JENKINS_URL}/github-webhook/

EOF
```

---

### 5. GitHub Webhook Setup (curl)

```bash
# Variables
GITHUB_TOKEN="ghp_your_personal_access_token"
GITHUB_ORG="your-org"
GITHUB_REPO="your-repo"
JENKINS_URL="https://jenkins.yourcompany.com"

# Create webhook
curl -X POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/${GITHUB_ORG}/${GITHUB_REPO}/hooks" \
  -d '{
    "name": "web",
    "active": true,
    "events": ["push", "pull_request"],
    "config": {
      "url": "'"${JENKINS_URL}"'/github-webhook/",
      "content_type": "json",
      "insecure_ssl": "0"
    }
  }'
```

---

### 6. GitLab Webhook Setup (curl)

```bash
# Variables
GITLAB_TOKEN="glpat-your_token"
GITLAB_PROJECT_ID="12345"  # or "group/project" URL-encoded
JENKINS_URL="https://jenkins.yourcompany.com"

# Create webhook
curl -X POST \
  -H "PRIVATE-TOKEN: ${GITLAB_TOKEN}" \
  "https://gitlab.com/api/v4/projects/${GITLAB_PROJECT_ID}/hooks" \
  -d "url=${JENKINS_URL}/project/YOUR_JOB_NAME" \
  -d "push_events=true" \
  -d "merge_requests_events=true" \
  -d "enable_ssl_verification=true"
```

---

### 7. Validate Jenkinsfile Syntax (Local)

```bash
# Option A: Use Jenkins CLI
JENKINS_URL="http://localhost:8080"
curl -o jenkins-cli.jar "${JENKINS_URL}/jnlpJars/jenkins-cli.jar"

java -jar jenkins-cli.jar -s "${JENKINS_URL}" -auth admin:YOUR_API_TOKEN \
  declarative-linter < Jenkinsfile


# Option B: Use npm package (no Jenkins needed)
npm install -g jenkins-pipeline-linter-connector

jplc --url "${JENKINS_URL}" --user admin --token YOUR_API_TOKEN Jenkinsfile


# Option C: Basic Groovy syntax check (catches obvious errors)
# Requires groovy installed
groovy -e "new GroovyShell().parse(new File('Jenkinsfile'))"
```

---

### 8. Starter Jenkinsfiles by Language

#### Python

```bash
cat > Jenkinsfile << 'EOF'
pipeline {
    agent any
    
    environment {
        PYTHONDONTWRITEBYTECODE = '1'
    }
    
    stages {
        stage('Setup') {
            steps {
                sh '''
                    python3 -m venv venv
                    . venv/bin/activate
                    pip install --upgrade pip
                    pip install -r requirements.txt
                    pip install pytest flake8
                '''
            }
        }
        stage('Lint') {
            steps {
                sh '''
                    . venv/bin/activate
                    flake8 . --count --select=E9,F63,F7,F82 --show-source --statistics
                '''
            }
        }
        stage('Test') {
            steps {
                sh '''
                    . venv/bin/activate
                    pytest --junitxml=results.xml
                '''
            }
            post {
                always {
                    junit 'results.xml'
                }
            }
        }
    }
}
EOF
```

#### Node.js

```bash
cat > Jenkinsfile << 'EOF'
pipeline {
    agent any
    
    tools {
        nodejs 'NodeJS-18'  // Must match Jenkins Global Tool name
    }
    
    stages {
        stage('Install') {
            steps {
                sh 'npm ci'
            }
        }
        stage('Lint') {
            steps {
                sh 'npm run lint'
            }
        }
        stage('Test') {
            steps {
                sh 'npm test -- --coverage'
            }
        }
        stage('Build') {
            steps {
                sh 'npm run build'
            }
        }
    }
}
EOF
```

#### Go

```bash
cat > Jenkinsfile << 'EOF'
pipeline {
    agent any
    
    environment {
        GO111MODULE = 'on'
        GOPATH = "${WORKSPACE}/go"
    }
    
    stages {
        stage('Build') {
            steps {
                sh 'go build -v ./...'
            }
        }
        stage('Test') {
            steps {
                sh 'go test -v -race -coverprofile=coverage.out ./...'
            }
        }
        stage('Vet') {
            steps {
                sh 'go vet ./...'
            }
        }
    }
}
EOF
```

#### Docker Build + Push

```bash
cat > Jenkinsfile << 'EOF'
pipeline {
    agent any
    
    environment {
        REGISTRY = 'docker.io'
        IMAGE = 'youruser/yourapp'
        TAG = "${env.GIT_COMMIT?.take(7) ?: 'latest'}"
    }
    
    stages {
        stage('Build') {
            steps {
                sh "docker build -t ${REGISTRY}/${IMAGE}:${TAG} ."
            }
        }
        stage('Push') {
            when {
                branch 'main'
            }
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-creds',
                    usernameVariable: 'USER',
                    passwordVariable: 'PASS'
                )]) {
                    sh '''
                        echo "$PASS" | docker login -u "$USER" --password-stdin ${REGISTRY}
                        docker push ${REGISTRY}/${IMAGE}:${TAG}
                        docker tag ${REGISTRY}/${IMAGE}:${TAG} ${REGISTRY}/${IMAGE}:latest
                        docker push ${REGISTRY}/${IMAGE}:latest
                    '''
                }
            }
        }
    }
    
    post {
        always {
            sh "docker rmi ${REGISTRY}/${IMAGE}:${TAG} || true"
        }
    }
}
EOF
```

---

### 9. Multibranch Pipeline Job (CLI)

```bash
cat > /tmp/multibranch-job.xml << 'EOF'
<?xml version='1.1' encoding='UTF-8'?>
<org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject plugin="workflow-multibranch">
  <description>Multibranch pipeline - auto-discovers branches</description>
  <properties/>
  <folderViews class="jenkins.branch.MultiBranchProjectViewHolder" plugin="branch-api"/>
  <healthMetrics/>
  <orphanedItemStrategy class="com.cloudbees.hudson.plugins.folder.computed.DefaultOrphanedItemStrategy" plugin="cloudbees-folder">
    <pruneDeadBranches>true</pruneDeadBranches>
    <daysToKeep>7</daysToKeep>
    <numToKeep>5</numToKeep>
  </orphanedItemStrategy>
  <sources class="jenkins.branch.MultiBranchProject$BranchSourceList" plugin="branch-api">
    <data>
      <jenkins.branch.BranchSource>
        <source class="jenkins.plugins.git.GitSCMSource" plugin="git">
          <id>git-source</id>
          <remote>git@github.com:YOUR_ORG/YOUR_REPO.git</remote>
          <credentialsId>git-ssh-key</credentialsId>
          <traits>
            <jenkins.plugins.git.traits.BranchDiscoveryTrait/>
            <jenkins.plugins.git.traits.TagDiscoveryTrait/>
          </traits>
        </source>
      </jenkins.branch.BranchSource>
    </data>
  </sources>
  <factory class="org.jenkinsci.plugins.workflow.multibranch.WorkflowBranchProjectFactory">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
    <scriptPath>Jenkinsfile</scriptPath>
  </factory>
</org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject>
EOF

java -jar jenkins-cli.jar -s "${JENKINS_URL}" -auth admin:YOUR_API_TOKEN \
  create-job my-multibranch-pipeline < /tmp/multibranch-job.xml
```

---

## Further Reading

- [Jenkins Pipeline Syntax](https://www.jenkins.io/doc/book/pipeline/syntax/)
- [Pipeline Steps Reference](https://www.jenkins.io/doc/pipeline/steps/)
- [Blue Ocean](https://www.jenkins.io/doc/book/blueocean/) — Modern Jenkins UI
