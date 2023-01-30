ARG USER_REPO=""


FROM fedora:35

LABEL name="VF dockerfile"

RUN groupadd -g 1001 -o validation
RUN useradd -m -u 1001 -g 1001 -o -s /bin/bash validation

RUN dnf install -y python3-pip gcc python3-devel libffi-devel jq openssh openssh-clients

# Clone the Framework and common Validations
RUN python3 -m pip install validations-libs validations-common

# Clone user repository if provided
RUN if [ "$USER_REPO" != "" ] ; then \
      git clone $USER_REPO $HOME; \
      cd $HOME/$USER_REPO && python3 -m pip install . ; \
    fi

#Setting up the default directory structure for both ansible,
#and the VF
RUN ln -s /usr/local/share/ansible  /usr/share/ansible

ENV ANSIBLE_HOST_KEY_CHECKING false
ENV ANSIBLE_RETRY_FILES_ENABLED false
ENV ANSIBLE_KEEP_REMOTE_FILES 1
ENV ANSIBLE_REMOTE_USER validation
ENV ANSIBLE_PRIVATE_KEY_FILE /home/validation/containerhost_private_key

USER validation
