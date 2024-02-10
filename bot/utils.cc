#include <linux/prctl.h>
#include <sys/prctl.h>
#include <unistd.h>

#include <cstring>
#include <ctime>

#include "headers/include.hpp"

string GetArch() {
#ifdef ARCH
  return string(ARCH);
#endif
  return "UNKNOWN";
}

unsigned int RandomInt(int min, int max) {
  srand(static_cast<unsigned int>(time(nullptr)));
  return rand() % (max - min + 1) + min;
}

void RandomBuf(char* buffer, const int buffer_size) {
  const char charSet[] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  srand(static_cast<unsigned int>(time(nullptr)));

  int csize = sizeof(charSet) - 1;

  for (int i = 0; i < buffer_size - 1; ++i) {
    int randomIndex = rand() % csize;
    buffer[i] = charSet[randomIndex];
  }

  buffer[buffer_size - 1] = '\0';
}

void DynamicName(char* argv[]) {
  char buffer[10];
  while (1) {
    RandomBuf(buffer, sizeof(buffer) / sizeof(buffer[0]));
    strcpy(argv[0], buffer);
    // printf("Argv0: %s\n", argv[0]);
    prctl(PR_SET_NAME, argv[0]);
    sleep(1);
  }
}

#ifdef KILLER
bool IsInteger(const string str) {
  if (str.length() < 1) {
    return false;
  }

  for (char c : str) {
    if (!isdigit(c)) {
      return false;
    }
  }

  return true;
}
#endif

#ifdef DEVMODE
void printBytes(const string& str) {
  const unsigned char* bytes = reinterpret_cast<const unsigned char*>(str.c_str());

  for (unsigned long i = 0; i < str.length(); i++) {
    printf("BYTE %lu: %d\n", i, static_cast<int>(bytes[i]));
  }
}

#endif
