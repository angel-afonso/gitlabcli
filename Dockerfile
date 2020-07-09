FROM scratch
COPY gitlabcli /
ENTRYPOINT ["/gitlabcli"]
