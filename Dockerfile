FROM golang:1.21 as BUILDER

# Build binary
COPY . /go/src/github.com/openmerlin/cronjob
RUN cd /go/src/github.com/openmerlin/cronjob && GO111MODULE=on CGO_ENABLED=0 go build -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'"

# Copy binary, config, and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update --repo OS --repo update && \
    dnf in -y shadow --repo OS --repo update && \
    dnf remove -y gdb-gdbserver && \
    groupadd -g 1000 cronjob && \
    useradd -u 1000 -g cronjob -s /sbin/nologin -m cronjob && \
    echo > /etc/issue && echo > /etc/issue.net && echo > /etc/motd && \
    sed -i 's/^PASS_MAX_DAYS.*/PASS_MAX_DAYS   90/' /etc/login.defs && \ # Set password expiration policies
    echo "umask 027" >> /root/.bashrc && \
    echo 'set +o history' >> /root/.bashrc

# Set additional security settings for the cronjob user's environment
USER cronjob
WORKDIR /home/cronjob
COPY --chown=cronjob --from=BUILDER /go/src/github.com/openmerlin/cronjob/cronjob /home/cronjob
RUN chmod 550 /home/cronjob/cronjob && \  # Restrict permissions for the cronjob binary
    echo "umask 027" >> /home/cronjob/.bashrc && \
    echo 'set +o history' >> /home/cronjob/.bashrc

ENTRYPOINT ["/home/cronjob/cronjob"]
