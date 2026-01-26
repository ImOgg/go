FROM ubuntu:latest

# 自動化你剛才的手動步驟
RUN apt update && apt install -y wget git tar && \
    wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /app
CMD ["/bin/bash"]