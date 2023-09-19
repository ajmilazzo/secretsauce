FROM ghcr.io/cse-labs/dind:latest@sha256:0c578efa8cfbb600a6bcbe66aa9135b427f12af8efaa9fc40e724a96ce656404

ARG USERNAME=vscode

COPY library-scripts/*.sh /scripts/

# install kind / k3d
RUN /bin/bash /scripts/kind-k3d-debian.sh

# change ownership of the home directory
RUN chown -R ${USERNAME}:${USERNAME} /home/${USERNAME} && \
    echo "👋 Welcome to the k3d Codespaces image\n" >> /usr/local/etc/vscode-dev-containers/first-run-notice.txt

# update the container
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get autoremove -y && \
    apt-get clean -y
