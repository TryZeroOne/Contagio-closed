# cross compile

FROM ubuntu:latest

RUN apt-get update && apt update && apt full-upgrade -y

RUN apt-get install --no-install-recommends -y \
    clang \
    sudo \
    pkg-config \ 
    file \
    nano 


COPY bot /home/root/bot/

ENV CC=/usr/bin/clang
ENV CXX=/usr/bin/clang++

USER root
WORKDIR /home/root/bot
ENV HOME /home/root



RUN sudo apt-get install -y build-essential 
RUN sudo apt-get install -y g++-mips-linux-gnu
RUN sudo apt-get install -y g++-mipsel-linux-gnu
RUN sudo apt-get install -y g++-mips64-linux-gnuabi64

RUN sudo apt-get update; sudo apt-get install -y g++-multilib -y


CMD ["tail", "-f", "/dev/null"]