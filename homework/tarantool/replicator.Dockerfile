FROM centos:7

# RUN yum install -y epel-release; yum clean all

# dependencies
RUN yum install -y git ncurses-devel cmake gcc-c++ boost boost-devel wget unzip nano bzip2 mysql-devel mysql-lib
RUN yum install -y make
#RUN yum install -y expat-devel zlib-devel bzip2-devel lua-devel which openssh-server

# install replicator
RUN git clone https://github.com/tarantool/mysql-tarantool-replication.git
WORKDIR /mysql-tarantool-replication
RUN git submodule update --init --recursive
RUN cmake . && make
RUN cp replicatord /usr/local/sbin/replicatord
# default config
RUN cp replicatord.yml /usr/local/etc/replicatord.yml

CMD ["/usr/local/sbin/replicatord","-c","/usr/local/etc/replicatord.yml"]
